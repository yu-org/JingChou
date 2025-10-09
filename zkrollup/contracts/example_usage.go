package contracts

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ExampleVerifyProof 展示如何使用 OpenVM Halo2 Verifier 验证证明
func ExampleVerifyProof() {
	// 1. 连接到以太坊节点
	client, err := ethclient.Dial("https://eth-mainnet.alchemyapi.io/v2/YOUR_API_KEY")
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum client: %v", err)
	}
	defer client.Close()

	// 2. OpenVM Halo2 Verifier 合约地址
	verifierAddress := common.HexToAddress("0xYourVerifierContractAddress")

	// 3. 创建合约实例
	verifier, err := NewOpenVmHalo2Verifier(verifierAddress, client)
	if err != nil {
		log.Fatalf("Failed to instantiate verifier contract: %v", err)
	}

	// 4. 准备验证参数
	// publicValues 现在是字节数组，需要根据你的 OpenVM 程序输出格式来填充
	publicValues := []byte{ /* 你的 public values 字节数据 */ }

	// 证明数据 (从 Axiom API 获取的 proof)
	proofData := []byte{ /* 你的证明数据 */ }

	// Application execution commitment
	appExeCommit := [32]byte{ /* 你的 app exe commit */ }

	// Application VM commitment
	appVmCommit := [32]byte{ /* 你的 app vm commit */ }

	// 5. 调用 verify 函数（无返回值，失败会 revert）
	callOpts := &bind.CallOpts{
		Context: context.Background(),
	}

	err = verifier.Verify(callOpts, publicValues, proofData, appExeCommit, appVmCommit)
	if err != nil {
		log.Fatalf("Proof verification failed: %v", err)
	}

	fmt.Println("✓ 证明验证成功！")
}

// ExampleVerifyProofWithSession 使用 Session 方式验证证明
func ExampleVerifyProofWithSession() {
	client, err := ethclient.Dial("https://eth-mainnet.alchemyapi.io/v2/YOUR_API_KEY")
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum client: %v", err)
	}
	defer client.Close()

	verifierAddress := common.HexToAddress("0xYourVerifierContractAddress")
	verifier, err := NewOpenVmHalo2Verifier(verifierAddress, client)
	if err != nil {
		log.Fatalf("Failed to instantiate verifier contract: %v", err)
	}

	// 创建 Session
	session := &OpenVmHalo2VerifierSession{
		Contract: verifier,
		CallOpts: bind.CallOpts{
			Context: context.Background(),
			Pending: false,
		},
	}

	// 准备参数
	publicValues := []byte{ /* public values bytes */ }
	proofData := []byte{ /* proof data */ }
	appExeCommit := [32]byte{}
	appVmCommit := [32]byte{}

	// 使用 Session 调用
	err = session.Verify(publicValues, proofData, appExeCommit, appVmCommit)
	if err != nil {
		log.Fatalf("Proof verification failed: %v", err)
	}

	fmt.Println("✓ 验证成功！")
}

// VerifyOpenVMProof 封装的验证函数，便于在你的项目中使用
func VerifyOpenVMProof(
	ethClientURL string,
	verifierAddress common.Address,
	publicValues []byte,
	proofData []byte,
	appExeCommit [32]byte,
	appVmCommit [32]byte,
) error {
	// 连接以太坊客户端
	client, err := ethclient.Dial(ethClientURL)
	if err != nil {
		return fmt.Errorf("failed to connect to Ethereum client: %w", err)
	}
	defer client.Close()

	// 创建合约实例
	verifier, err := NewOpenVmHalo2Verifier(verifierAddress, client)
	if err != nil {
		return fmt.Errorf("failed to instantiate verifier contract: %w", err)
	}

	// 调用验证函数
	callOpts := &bind.CallOpts{
		Context: context.Background(),
	}

	err = verifier.Verify(callOpts, publicValues, proofData, appExeCommit, appVmCommit)
	if err != nil {
		return fmt.Errorf("proof verification failed: %w", err)
	}

	return nil
}

// EstimateVerifyGas 估算验证交易所需的 Gas
func EstimateVerifyGas(
	ethClientURL string,
	verifierAddress common.Address,
	publicValues []byte,
	proofData []byte,
	appExeCommit [32]byte,
	appVmCommit [32]byte,
) (*big.Int, error) {
	client, err := ethclient.Dial(ethClientURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum client: %w", err)
	}
	defer client.Close()

	verifier, err := NewOpenVmHalo2Verifier(verifierAddress, client)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate verifier contract: %w", err)
	}

	// 尝试调用以测试是否成功
	callOpts := &bind.CallOpts{Context: context.Background()}
	err = verifier.Verify(callOpts, publicValues, proofData, appExeCommit, appVmCommit)
	if err != nil {
		return nil, fmt.Errorf("failed to estimate gas: %w", err)
	}

	// 这里返回的是一个估算值，实际使用时需要根据具体情况调整
	return big.NewInt(500000), nil // 示例值，实际需要链上估算
}
