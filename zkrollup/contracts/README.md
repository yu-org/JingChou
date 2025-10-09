# OpenVM Halo2 Verifier - Golang ABI

这个目录包含了 OpenVM Halo2 Verifier 合约的 Golang ABI 绑定，用于在链上验证 OpenVM 生成的零知识证明。

## 文件说明

- `IOpenVmHalo2Verifier.sol` - OpenVM Halo2 Verifier 合约的 Solidity 接口
- `IOpenVmHalo2Verifier.abi` - 合约的 ABI JSON 文件
- `openvm_halo2_verifier.go` - 自动生成的 Golang 绑定（由 abigen 生成）
- `example_usage.go` - 使用示例
- `README.md` - 本文件

## 基于的版本

- **OpenVM Version**: v1.4
- **Solidity Version**: 0.8.19
- **默认配置**: 使用默认的 aggregation VM config，public values 固定为 32 个

## 安装依赖

```bash
go get github.com/ethereum/go-ethereum
```

## 合约接口

```solidity
interface IOpenVmHalo2Verifier {
    function verify(
        bytes calldata publicValues,
        bytes calldata proofData,
        bytes32 appExeCommit,
        bytes32 appVmCommit
    ) external view;
}
```

**重要：该函数没有返回值，验证失败会直接 revert！**

### 参数说明

1. **publicValues** - `bytes`
   - 公开值的字节数组
   - 格式需要与你的 OpenVM 程序输出匹配

2. **proofData** - `bytes`
   - 证明数据的字节数组
   - 从 Axiom API 获取的 proof 数据

3. **appExeCommit** - `bytes32`
   - 应用执行承诺（Application Execution Commitment）

4. **appVmCommit** - `bytes32`
   - 应用 VM 承诺（Application VM Commitment）

## 使用方法

### 1. 基本使用

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/ethereum/go-ethereum/accounts/abi/bind"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/ethclient"
    "github.com/yu-org/JingChou/zkrollup/contracts"
)

func main() {
    // 连接以太坊节点
    client, err := ethclient.Dial("https://eth-mainnet.alchemyapi.io/v2/YOUR_API_KEY")
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // Verifier 合约地址
    verifierAddr := common.HexToAddress("0xYourVerifierAddress")

    // 创建合约实例
    verifier, err := contracts.NewOpenVmHalo2Verifier(verifierAddr, client)
    if err != nil {
        log.Fatal(err)
    }

    // 准备验证参数
    publicValues := []byte{/* 你的 public values 数据 */}
    proofData := []byte{/* 你的证明数据 */}
    appExeCommit := [32]byte{/* app exe commit */}
    appVmCommit := [32]byte{/* app vm commit */}

    // 调用验证（无返回值，失败会 revert）
    callOpts := &bind.CallOpts{Context: context.Background()}
    err = verifier.Verify(callOpts, publicValues, proofData, appExeCommit, appVmCommit)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("✓ 证明验证成功！")
}
```

### 2. 使用 Session

```go
session := &contracts.OpenVmHalo2VerifierSession{
    Contract: verifier,
    CallOpts: bind.CallOpts{
        Context: context.Background(),
    },
}

err := session.Verify(publicValues, proofData, appExeCommit, appVmCommit)
```

### 3. 封装的便捷函数

```go
import "github.com/yu-org/JingChou/zkrollup/contracts"

err := contracts.VerifyOpenVMProof(
    "https://eth-mainnet.alchemyapi.io/v2/YOUR_API_KEY",
    verifierAddress,
    publicValues,
    proofData,
    appExeCommit,
    appVmCommit,
)
```

## 与 Axiom Prover 集成

在你的 zkrollup 项目中，可以这样使用：

```go
// 1. 使用 Axiom Prover 生成证明
proofChan := make(chan *ProofResult, 10)
proofID, err := axiomProver.GenerateProof(blockBatch, proofChan)

// 2. 等待证明完成
result := <-proofChan

if result.StatusCode == ProveSuccess {
    // 3. 在链上验证证明
    err := contracts.VerifyOpenVMProof(
        ethClientURL,
        verifierAddress,
        publicValues,      // 需要从 proof 中提取
        result.Proof.ZKProof,
        appExeCommit,
        appVmCommit,
    )
    
    if err != nil {
        fmt.Printf("链上验证失败: %v\n", err)
    } else {
        fmt.Println("链上验证成功！")
    }
}
```

## 部署 Verifier 合约

### 使用 OpenVM SDK 部署

```bash
# 克隆 OpenVM Solidity SDK
git clone --recursive https://github.com/openvm-org/openvm-solidity-sdk.git
cd openvm-solidity-sdk

# 部署
forge create src/v1.4/OpenVmHalo2Verifier.sol:OpenVmHalo2Verifier \
    --rpc-url $RPC_URL \
    --private-key $PRIVATE_KEY \
    --broadcast
```

### 使用已部署的合约

如果使用已部署的 Verifier 合约，只需要知道合约地址即可：

```go
verifierAddr := common.HexToAddress("0xDeployedVerifierAddress")
```

## 注意事项

1. **验证失败会 Revert**
   - `verify()` 函数没有返回值
   - 验证失败会直接 revert，抛出异常
   - 在 Go 中表现为返回 error

2. **Public Values 格式**
   - publicValues 是字节数组，格式需要与你的 OpenVM 程序输出匹配
   - 不是固定大小，长度可变

3. **Gas 消耗**
   - 验证证明是一个复杂的链上操作，会消耗较多 Gas
   - 建议先在测试网测试

4. **证明格式**
   - 证明数据格式必须与 OpenVM 生成的格式匹配
   - 确保使用相同版本的 OpenVM

5. **合约版本**
   - 合约使用 Solidity 0.8.19
   - OpenVM 保证 patch 版本的向后兼容性

## 参考链接

- [OpenVM Documentation](https://docs.openvm.dev/)
- [OpenVM Solidity SDK](https://github.com/openvm-org/openvm-solidity-sdk)
- [Axiom API Documentation](https://docs.axiom.xyz/)

## 重新生成 ABI

如果需要重新生成 Golang 绑定：

```bash
# 确保已安装 abigen
# go install github.com/ethereum/go-ethereum/cmd/abigen@latest

# 生成绑定
abigen --abi zkrollup/contracts/IOpenVmHalo2Verifier.abi \
    --pkg contracts \
    --type OpenVmHalo2Verifier \
    --out zkrollup/contracts/openvm_halo2_verifier.go
```

## 许可证

根据 OpenVM Solidity SDK，代码采用 Apache-2.0 和 MIT 双许可证。

