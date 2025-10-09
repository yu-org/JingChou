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
// ⚠️ 重要：这个函数必须返回与你的 OpenVM 程序中 reveal_public_values() 输出完全一致的数据
func (z *ZkRollup) extractPublicValues(proofResult *prover.ProofResult) []byte {
	// TODO: 根据你的 OpenVM Rust 程序实际输出来实现
	//
	// 理想情况下，publicValues 应该从 Axiom API 的响应中直接获取
	// 而不是在这里重新构造
	//
	// 如果 Axiom API 返回了 publicValues:
	// return proofResult.PublicValues
	//
	// 如果需要从其他地方提取，确保格式与你的 OpenVM 程序匹配：
	//
	// 示例 OpenVM Rust 代码：
	// ```rust
	// #[app]
	// fn main() {
	//     let from_block = ...;
	//     let to_block = ...;
	//     let pre_state = ...;
	//     let new_state = ...;
	//
	//     reveal_public_values(&[
	//         from_block.to_le_bytes(),  // 注意：little-endian 还是 big-endian
	//         to_block.to_le_bytes(),
	//         pre_state.as_bytes(),
	//         new_state.as_bytes(),
	//     ]);
	// }
	// ```

	var publicValues []byte

	if proofResult.Proof != nil {
		// 当前示例实现（需要根据实际情况调整）
		// FromBlockNum (8 bytes, big-endian)
		fromBlockBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(fromBlockBytes, proofResult.Proof.FromBlockNum)
		publicValues = append(publicValues, fromBlockBytes...)

		// ToBlockNum (8 bytes, big-endian)
		toBlockBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(toBlockBytes, proofResult.Proof.ToBlockNum)
		publicValues = append(publicValues, toBlockBytes...)

		// PreStateRoot (32 bytes)
		publicValues = append(publicValues, proofResult.Proof.PreStateRoot.Bytes()...)

		// NewStateRoot (32 bytes)
		publicValues = append(publicValues, proofResult.Proof.NewStateRoot.Bytes()...)

		// 总共 80 bytes
	}

	return publicValues
}

// calculateAppExeCommit 获取应用执行承诺
// appExeCommit = Hash(ELF文件)，标识特定版本的程序
// ⚠️ 这个值是固定的，不随每次执行而变化！
func (z *ZkRollup) calculateAppExeCommit(proofResult *prover.ProofResult) [32]byte {
	var commit [32]byte

	// 方式1：从配置中获取（如果配置了固定值）
	if z.cfg.AppExeCommit != "" {
		copy(commit[:], ethcommon.HexToHash(z.cfg.AppExeCommit).Bytes())
		return commit
	}

	// 方式2：从 AxiomProver 中获取（推荐）
	// AxiomProver 在初始化时已经计算并存储了 appExeCommit
	if axiomProver, ok := z.prover.(*prover.AxiomProver); ok {
		commit = axiomProver.GetAppExeCommit()
		if commit != [32]byte{} {
			return commit
		}
	}

	// 如果都没有，警告并返回零值
	logrus.Warn("AppExeCommit not available, using zero value - verification will likely fail!")
	return commit
}

// calculateAppVmCommit 计算应用 VM 承诺
// appVmCommit = Hash(VM配置)，标识使用的 VM 版本和配置
// ⚠️ 这个值应该是固定的，不随每次执行而变化！
func (z *ZkRollup) calculateAppVmCommit(proofResult *prover.ProofResult) [32]byte {
	var commit [32]byte

	// 方式1：使用配置中的固定值（强烈推荐）
	// 这个值应该与部署的 Verifier 合约使用的 VM 配置匹配
	if z.cfg.AppVmCommit != "" {
		copy(commit[:], ethcommon.HexToHash(z.cfg.AppVmCommit).Bytes())
		return commit
	}

	// 方式2：从 Axiom API 响应中获取（如果可用）
	// TODO: 添加字段 proofResult.AppVmCommit
	// if proofResult.AppVmCommit != [32]byte{} {
	//     return proofResult.AppVmCommit
	// }

	// 方式3：使用 OpenVM 默认 VM 配置的承诺
	// 需要从 OpenVM v1.4 文档获取默认配置的哈希值
	// 例如：
	// defaultVmCommit := "0x1234567890abcdef..." // OpenVM v1.4 default config hash
	// copy(commit[:], ethcommon.HexToHash(defaultVmCommit).Bytes())
	// return commit

	logrus.Warn("AppVmCommit not configured, using zero value - verification may fail!")
	return commit
}
