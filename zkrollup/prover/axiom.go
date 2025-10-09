package prover

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/yu-org/JingChou/zkrollup/config"

	"github.com/yu-org/yu/core/types"
)

type AxiomProver struct {
	cfg          *config.ProverConfig
	httpClient   *http.Client
	programID    string        // 注册后的 program ID
	pollInterval time.Duration // 轮询间隔
	pollTimeout  time.Duration // 轮询超时时间
	proofType    string        // 证明类型
}

// Axiom API 响应结构（根据文档）

// ProgramUploadResponse 程序上传响应
type ProgramUploadResponse struct {
	ID string `json:"id"`
}

// ProofCreateResponse 创建证明响应
type ProofCreateResponse struct {
	ID string `json:"id"` // 证明 ID
}

// ProofStatusResponse 证明状态响应
type ProofStatusResponse struct {
	ID              string  `json:"id"`
	CreatedAt       string  `json:"created_at"`
	State           string  `json:"state"` // Queued, Executing, Executed, AppProving, AppProvingDone, PostProcessing, Failed, Succeeded
	ProofType       string  `json:"proof_type"`
	ProgramUUID     string  `json:"program_uuid"`
	ErrorMessage    *string `json:"error_message"`
	LaunchedAt      *string `json:"launched_at"`
	TerminatedAt    *string `json:"terminated_at"`
	CreatedBy       string  `json:"created_by"`
	CellsUsed       int     `json:"cells_used"`
	Cost            *int    `json:"cost"`
	MachineType     string  `json:"machine_type"`
	ProofSize       *int    `json:"proof_size"`
	NumInstructions *int    `json:"num_instructions"`
}

// ProofInputData 证明输入数据
type ProofInputData struct {
	Input []string `json:"input"` // 十六进制字符串数组
}

func NewAxiomProver(cfg *config.ProverConfig) (Prover, error) {
	if cfg.URL == "" {
		return nil, fmt.Errorf("axiom URL is required")
	}
	if cfg.ApiKey == "" {
		return nil, fmt.Errorf("axiom API key is required")
	}

	// 设置默认的轮询间隔（5秒）
	pollInterval := 5 * time.Second
	if cfg.PollInterval > 0 {
		pollInterval = time.Duration(cfg.PollInterval) * time.Second
	}

	// 设置默认的轮询超时（2小时）
	pollTimeout := 2 * time.Hour
	if cfg.PollTimeout > 0 {
		pollTimeout = time.Duration(cfg.PollTimeout) * time.Second
	}

	// 设置默认的证明类型（stark）
	proofType := "stark"
	if cfg.ProofType != "" {
		proofType = cfg.ProofType
	}

	prover := &AxiomProver{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		pollInterval: pollInterval,
		pollTimeout:  pollTimeout,
		proofType:    proofType,
	}

	// 如果配置中已有 programID，直接使用
	if cfg.ProgramID != "" {
		prover.programID = cfg.ProgramID
		return prover, nil
	}

	// 否则注册 ELF 文件
	if cfg.ElfPath == "" {
		return nil, fmt.Errorf("elf_path is required when program_id is not provided")
	}

	programID, err := prover.registerProgram(cfg.ElfPath)
	if err != nil {
		return nil, fmt.Errorf("failed to register program: %w", err)
	}

	prover.programID = programID
	return prover, nil
}

// registerProgram 注册 reth ELF 文件到 Axiom
func (a *AxiomProver) registerProgram(elfPath string) (string, error) {
	// 检查 ELF 文件是否存在
	if _, err := os.Stat(elfPath); os.IsNotExist(err) {
		return "", fmt.Errorf("elf file not found: %s", elfPath)
	}

	// 打开 ELF 文件
	file, err := os.Open(elfPath)
	if err != nil {
		return "", fmt.Errorf("failed to open elf file: %w", err)
	}
	defer file.Close()

	// 创建 multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加文件
	part, err := writer.CreateFormFile("program", filepath.Base(elfPath))
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return "", fmt.Errorf("failed to copy file content: %w", err)
	}

	// 添加其他字段
	if a.cfg.VMConfigID != "" {
		_ = writer.WriteField("vm_config_id", a.cfg.VMConfigID)
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}

	// 发送请求
	url := fmt.Sprintf("%s/v1/programs", a.cfg.URL)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Axiom-API-Key", a.cfg.ApiKey)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to register program, status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	// 解析响应（根据文档，直接返回 JSON 对象）
	var uploadResp ProgramUploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&uploadResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return uploadResp.ID, nil
}

// GenerateProof 提交区块批次并生成证明，在后台轮询等待结果
func (a *AxiomProver) GenerateProof(blockBatch []*types.Block, proofChan chan *ProofResult) (string, error) {
	if len(blockBatch) == 0 {
		return "", fmt.Errorf("block batch is empty")
	}

	// 准备输入数据（根据文档，input 是十六进制字符串数组）
	// 这里需要将 blockBatch 序列化为十六进制字符串
	blockBytes, err := json.Marshal(blockBatch)
	if err != nil {
		return "", fmt.Errorf("failed to marshal blocks: %w", err)
	}

	// 添加 0x01 前缀表示这是字节数据
	hexInput := fmt.Sprintf("0x01%x", blockBytes)

	inputData := ProofInputData{
		Input: []string{hexInput},
	}

	bodyBytes, err := json.Marshal(inputData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal input: %w", err)
	}

	// 发送证明生成请求（根据文档，program_id 作为查询参数）
	url := fmt.Sprintf("%s/v1/proofs?program_id=%s&proof_type=%s", a.cfg.URL, a.programID, a.proofType)
	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Axiom-API-Key", a.cfg.ApiKey)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to create proof task, status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	// 解析响应（根据文档，直接返回 {"id": "..."} ）
	var createResp ProofCreateResponse
	if err := json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	proofID := createResp.ID

	// 启动后台 goroutine 轮询获取证明结果
	go a.pollProofResult(proofID, proofChan)

	// 立即返回 proof ID
	return proofID, nil
}

// pollProofResult 在后台轮询获取证明结果，完成后发送到 channel
func (a *AxiomProver) pollProofResult(proofID string, proofChan chan *ProofResult) {
	ticker := time.NewTicker(a.pollInterval)
	defer ticker.Stop()

	// 设置最大轮询时间
	timeout := time.After(a.pollTimeout)

	for {
		select {
		case <-ticker.C:
			result, err := a.GetProof(proofID)
			if err != nil {
				// 出错时发送失败结果
				proofChan <- &ProofResult{
					StatusCode: ProveFailed,
					ProofID:    proofID,
					Proof:      nil,
				}
				return
			}

			// 如果已完成或失败，发送结果并退出
			if result.StatusCode == ProveSuccess || result.StatusCode == ProveFailed {
				proofChan <- result
				return
			}

			// 否则继续轮询

		case <-timeout:
			// 超时，发送超时失败结果
			proofChan <- &ProofResult{
				StatusCode: ProveFailed,
				ProofID:    proofID,
				Proof:      nil,
			}
			return
		}
	}
}

// GetProof 查询证明状态和结果
func (a *AxiomProver) GetProof(proofID string) (*ProofResult, error) {
	if proofID == "" {
		return nil, fmt.Errorf("proof ID is required")
	}

	// 构造请求
	url := fmt.Sprintf("%s/v1/proofs/%s", a.cfg.URL, proofID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Axiom-API-Key", a.cfg.ApiKey)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get proof status, status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	// 解析响应（根据文档）
	var statusResp ProofStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&statusResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// 转换状态码（根据文档的 state 字段）
	proofStatus := &ProofResult{
		ProofID: statusResp.ID,
	}

	switch statusResp.State {
	case "Succeeded":
		proofStatus.StatusCode = ProveSuccess
		// 如果需要获取实际的证明数据，需要调用 GET /v1/proofs/{proof_id}/proof/{proof_type}
		if statusResp.ProofSize != nil && *statusResp.ProofSize > 0 {
			// 这里可以选择获取证明数据
			proofData, err := a.getProofData(proofID, a.proofType)
			if err == nil {
				proofStatus.Proof = &Proof{
					ZKProof: proofData,
				}
			}
		}
	case "Failed":
		proofStatus.StatusCode = ProveFailed
	case "Executing", "AppProving", "PostProcessing":
		proofStatus.StatusCode = Proving
	case "Queued", "Executed", "AppProvingDone":
		proofStatus.StatusCode = ProvePending
	default:
		proofStatus.StatusCode = ProvePending
	}

	return proofStatus, nil
}

// getProofData 获取证明数据（根据文档 GET /v1/proofs/{proof_id}/proof/{proof_type}）
func (a *AxiomProver) getProofData(proofID string, proofType string) ([]byte, error) {
	url := fmt.Sprintf("%s/v1/proofs/%s/proof/%s", a.cfg.URL, proofID, proofType)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Axiom-API-Key", a.cfg.ApiKey)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get proof data, status: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// WaitForProof 等待证明完成（轮询）
func (a *AxiomProver) WaitForProof(proofID string, timeout time.Duration) (*ProofResult, error) {
	ticker := time.NewTicker(a.pollInterval)
	defer ticker.Stop()

	timeoutChan := time.After(timeout)

	for {
		select {
		case <-ticker.C:
			status, err := a.GetProof(proofID)
			if err != nil {
				return nil, err
			}

			// 如果已完成或失败，返回结果
			if status.StatusCode == ProveSuccess || status.StatusCode == ProveFailed {
				return status, nil
			}

		case <-timeoutChan:
			return nil, fmt.Errorf("timeout waiting for proof")
		}
	}
}

// CancelProof 取消证明任务
// 注意：根据当前文档，没有明确的取消 API，这里使用 DELETE 尝试
func (a *AxiomProver) CancelProof(proofID string) (*ProofResult, error) {
	if proofID == "" {
		return nil, fmt.Errorf("proof ID is required")
	}

	// 构造请求 - 使用 DELETE 方法取消任务
	url := fmt.Sprintf("%s/v1/proofs/%s", a.cfg.URL, proofID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Axiom-API-Key", a.cfg.ApiKey)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusAccepted {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to cancel proof, status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	// 返回取消后的状态
	return &ProofResult{
		StatusCode: ProveFailed, // 取消的任务标记为失败
		ProofID:    proofID,
		Proof:      nil,
	}, nil
}
