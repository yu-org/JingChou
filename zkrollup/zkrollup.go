package zkrollup

import (
	"context"
	"encoding/binary"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"github.com/yu-org/JingChou/zkrollup/config"
	"github.com/yu-org/JingChou/zkrollup/contracts"
	"github.com/yu-org/JingChou/zkrollup/prover"
	"github.com/yu-org/yu/common"
	"github.com/yu-org/yu/core/tripod"
	"github.com/yu-org/yu/core/types"
)

type ZkRollup struct {
	*tripod.Tripod

	cfg       *config.Config
	ethCli    *ethclient.Client
	prover    prover.Prover
	proofChan chan *prover.ProofResult
}

func NewZkRollup(cfg *config.Config) (*ZkRollup, error) {
	ethCli, err := ethclient.Dial(cfg.L1ChainAddr)
	if err != nil {
		return nil, err
	}
	prv, err := prover.NewAxiomProver(&cfg.Prover)
	if err != nil {
		return nil, err
	}
	proofChan := make(chan *prover.ProofResult, 10)
	return &ZkRollup{
		Tripod:    tripod.NewTripod(),
		cfg:       cfg,
		ethCli:    ethCli,
		prover:    prv,
		proofChan: proofChan,
	}, nil
}

func (z *ZkRollup) StartBlock(block *types.Block) {
	//TODO implement me
	panic("implement me")
}

func (z *ZkRollup) EndBlock(block *types.Block) {

}

func (z *ZkRollup) FinalizeBlock(block *types.Block) {
	if uint(block.Height)%z.cfg.BlockBatchSizeForProve != 0 {
		return
	}

	startProveBlockHeight := block.Height - common.BlockNum(z.cfg.BlockBatchSizeForProve) + 1

	blocks, err := z.Chain.GetRangeBlocks(startProveBlockHeight, block.Height)
	if err != nil {
		logrus.Errorf("get range blocks failed: %v", err)
		return
	}
	// 证明
	proofID, err := z.prover.GenerateProof(blocks, z.proofChan)
	if err != nil {
		logrus.Errorf("start to prove blocks from %d to %d failed: %v", startProveBlockHeight, block.Height, err)
		return
	}
	logrus.Infof("start to prove blocks from %d to %d, proofID: %s", startProveBlockHeight, block.Height, proofID)
}

func (z *ZkRollup) GetProof() {
	for {
		select {
		case proofResult := <-z.proofChan:
			logrus.Infof("get proof: %s", proofResult.ProofID)
			// TODO: Verify proof
			// 发送到 L1
			err := z.SendProofToL1(proofResult)
			if err != nil {
				logrus.Errorf("send proof to L1 failed: %v", err)
			}
		}
	}
}

func (z *ZkRollup) SendProofToL1(proofResult *prover.ProofResult) error {
	// 检查证明状态
	if proofResult.StatusCode != prover.ProveSuccess {
		return fmt.Errorf("proof generation failed, status: %s", proofResult.StatusCode.String())
	}

	if proofResult.Proof == nil || len(proofResult.Proof.ZKProof) == 0 {
		return fmt.Errorf("proof data is empty")
	}

	// 1. 创建 OpenVM Halo2 Verifier 合约实例
	verifierAddr := ethcommon.HexToAddress(z.cfg.L1VerifierAddr)
	verifier, err := contracts.NewOpenVmHalo2Verifier(verifierAddr, z.ethCli)
	if err != nil {
		return fmt.Errorf("failed to create verifier contract instance: %w", err)
	}

	// 2. 准备 public values（字节数组）
	// 这里需要根据实际的证明数据来填充
	// 示例：可以从 proof 中提取或者根据区块数据构造
	publicValues := z.extractPublicValues(proofResult)

	// 3. 准备 appExeCommit 和 appVmCommit
	var appExeCommit [32]byte
	var appVmCommit [32]byte

	// 如果配置中有值，使用配置的值
	if z.cfg.AppExeCommit != "" {
		copy(appExeCommit[:], ethcommon.HexToHash(z.cfg.AppExeCommit).Bytes())
	} else {
		// 否则可以从证明中提取或动态计算
		appExeCommit = z.calculateAppExeCommit(proofResult)
	}

	if z.cfg.AppVmCommit != "" {
		copy(appVmCommit[:], ethcommon.HexToHash(z.cfg.AppVmCommit).Bytes())
	} else {
		appVmCommit = z.calculateAppVmCommit(proofResult)
	}

	// 4. 调用链上验证
	callOpts := &bind.CallOpts{
		Context: context.Background(),
		Pending: false,
	}

	logrus.Infof("Verifying proof on L1, proofID: %s, proof size: %d bytes",
		proofResult.ProofID, len(proofResult.Proof.ZKProof))

	err = verifier.Verify(
		callOpts,
		publicValues,
		proofResult.Proof.ZKProof,
		appExeCommit,
		appVmCommit,
	)
	if err != nil {
		return fmt.Errorf("proof verification failed on L1: %w", err)
	}

	logrus.Infof("✓ Proof verified successfully on L1! ProofID: %s, Blocks: %d-%d",
		proofResult.ProofID,
		proofResult.Proof.FromBlockNum,
		proofResult.Proof.ToBlockNum)

	// 5. 可选：如果需要将验证结果写入其他合约或触发其他操作
	// 这里可以调用你的 L1ContractAddr 合约来记录验证结果
	// err = z.submitProofToContract(proofResult)

	return nil
}

// extractPublicValues 从证明结果中提取 public values
// 这个函数需要根据你的具体需求来实现
func (z *ZkRollup) extractPublicValues(proofResult *prover.ProofResult) []byte {
	// 示例实现：
	// 1. 可以从 proof 的 PreStateRoot 和 NewStateRoot 中提取
	// 2. 可以包含区块高度等信息
	// 3. 需要根据你的 OpenVM 程序的输出来决定

	var publicValues []byte

	if proofResult.Proof != nil {
		// 方式1：直接序列化整个 Proof 结构
		// proofBytes, _ := json.Marshal(proofResult.Proof)
		// return proofBytes

		// 方式2：按顺序拼接关键字段（示例）
		// 注意：实际格式需要与你的 OpenVM 程序输出匹配

		// FromBlockNum (8 bytes)
		fromBlockBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(fromBlockBytes, proofResult.Proof.FromBlockNum)
		publicValues = append(publicValues, fromBlockBytes...)

		// ToBlockNum (8 bytes)
		toBlockBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(toBlockBytes, proofResult.Proof.ToBlockNum)
		publicValues = append(publicValues, toBlockBytes...)

		// PreStateRoot (32 bytes)
		publicValues = append(publicValues, proofResult.Proof.PreStateRoot.Bytes()...)

		// NewStateRoot (32 bytes)
		publicValues = append(publicValues, proofResult.Proof.NewStateRoot.Bytes()...)

		// 总共 80 bytes (8 + 8 + 32 + 32)
	}

	return publicValues
}

// calculateAppExeCommit 计算应用执行承诺
// 这个函数需要根据你的具体需求来实现
func (z *ZkRollup) calculateAppExeCommit(proofResult *prover.ProofResult) [32]byte {
	var commit [32]byte

	// 示例实现：
	// 可以根据程序的 ELF 哈希或者其他标识来计算
	// 这里简单地使用 PreStateRoot 作为示例
	if proofResult.Proof != nil {
		copy(commit[:], proofResult.Proof.PreStateRoot.Bytes())
	}

	return commit
}

// calculateAppVmCommit 计算应用 VM 承诺
// 这个函数需要根据你的具体需求来实现
func (z *ZkRollup) calculateAppVmCommit(proofResult *prover.ProofResult) [32]byte {
	var commit [32]byte

	// 示例实现：
	// 可以根据 VM 配置或者其他标识来计算
	// 这里简单地使用 NewStateRoot 作为示例
	if proofResult.Proof != nil {
		copy(commit[:], proofResult.Proof.NewStateRoot.Bytes())
	}

	return commit
}
