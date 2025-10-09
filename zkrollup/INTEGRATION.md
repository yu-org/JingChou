# ZK Rollup 集成文档

## 📋 概述

`SendProofToL1` 函数实现了将 Axiom 生成的零知识证明提交到 L1 以太坊链上进行验证的完整流程。

## 🔧 实现的功能

### 1. SendProofToL1 函数

位于 `zkrollup/zkrollup.go`，实现了以下流程：

```go
func (z *ZkRollup) SendProofToL1(proofResult *prover.ProofResult) error
```

#### 执行步骤：

1. **验证证明状态**
   - 检查证明是否生成成功
   - 确保证明数据不为空

2. **创建合约实例**
   - 使用配置的 Verifier 合约地址
   - 创建 OpenVM Halo2 Verifier 合约实例

3. **准备验证参数**
   - 提取 public values（32 个 bytes32）
   - 准备 appExeCommit（应用执行承诺）
   - 准备 appVmCommit（应用 VM 承诺）

4. **调用链上验证**
   - 调用 Verifier 合约的 `verify()` 函数
   - 验证证明的有效性

5. **处理验证结果**
   - 如果验证成功，记录日志
   - 如果验证失败，返回错误

## 🔄 完整工作流程

```
区块生成 → 达到批次大小 → 生成证明 → 等待证明完成 → 提交到 L1 验证
    ↓           ↓              ↓            ↓              ↓
FinalizeBlock  判断高度     Axiom API    proofChan    SendProofToL1
                                                           ↓
                                                    OpenVM Verifier
                                                           ↓
                                                      验证成功/失败
```

### 详细流程图

```
┌─────────────────┐
│  FinalizeBlock  │
│   (每个区块)     │
└────────┬────────┘
         │
         ├──> 判断是否达到批次大小
         │    (block.Height % BatchSize == 0)
         │
         ▼
┌─────────────────┐
│  GenerateProof  │
│  (提交到 Axiom)  │
└────────┬────────┘
         │
         ├──> 启动后台轮询
         │
         ▼
┌─────────────────┐
│   proofChan     │
│  (接收证明结果)  │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ SendProofToL1   │
│  (L1 验证)      │
└────────┬────────┘
         │
         ├──> 1. 创建 Verifier 实例
         ├──> 2. 提取 public values
         ├──> 3. 准备 commitments
         ├──> 4. 调用 verify()
         │
         ▼
┌─────────────────┐
│  验证结果       │
│  ✓ Success      │
│  ✗ Failed       │
└─────────────────┘
```

## 📝 配置说明

### zkrollup 配置 (config.toml)

```toml
# 区块批次大小
block_batch_size_for_prove = 10

# L1 链配置
l1_chain_addr = "https://eth-mainnet.alchemyapi.io/v2/YOUR_API_KEY"
l1_verifier_addr = "0xYourVerifierContractAddress"

# 可选：固定的承诺值
app_exe_commit = "0x..."
app_vm_commit = "0x..."

[prover]
url = "https://api.axiom.xyz"
api_key = "your_api_key"
elf_path = "/path/to/reth.elf"
proof_type = "stark"
```

## 🎯 辅助函数

### 1. extractPublicValues

从证明结果中提取 public values（32 个 bytes32）。

```go
func (z *ZkRollup) extractPublicValues(proofResult *prover.ProofResult) [32][32]byte
```

**当前实现：**
- `publicValues[0]` - FromBlockNum (起始区块高度)
- `publicValues[1]` - ToBlockNum (结束区块高度)
- `publicValues[2]` - PreStateRoot (前状态根)
- `publicValues[3]` - NewStateRoot (新状态根)
- `publicValues[4-31]` - 保留（零值）

**需要根据你的 OpenVM 程序输出调整！**

### 2. calculateAppExeCommit

计算应用执行承诺。

```go
func (z *ZkRollup) calculateAppExeCommit(proofResult *prover.ProofResult) [32]byte
```

**可以基于：**
- ELF 文件的哈希
- 程序的标识符
- 执行参数的哈希

### 3. calculateAppVmCommit

计算应用 VM 承诺。

```go
func (z *ZkRollup) calculateAppVmCommit(proofResult *prover.ProofResult) [32]byte
```

**可以基于：**
- VM 配置的哈希
- VM 版本标识
- 其他 VM 相关参数

## 💡 使用示例

### 启动 ZK Rollup

```go
package main

import (
    "log"
    "github.com/yu-org/JingChou/zkrollup"
    "github.com/yu-org/JingChou/zkrollup/config"
)

func main() {
    // 加载配置
    cfg := &config.Config{
        BlockBatchSizeForProve: 10,
        L1ChainAddr:            "https://eth-mainnet.alchemyapi.io/v2/YOUR_API_KEY",
        L1VerifierAddr:         "0xYourVerifierAddress",
        Prover: config.ProverConfig{
            URL:     "https://api.axiom.xyz",
            ApiKey:  "your_api_key",
            ElfPath: "/path/to/reth.elf",
        },
    }

    // 创建 ZK Rollup 实例
    zkRollup, err := zkrollup.NewZkRollup(cfg)
    if err != nil {
        log.Fatal(err)
    }

    // 启动证明监听（在后台运行）
    go zkRollup.GetProof()

    // 你的区块链逻辑...
    // 每个区块会调用 FinalizeBlock
    // 达到批次大小时自动生成证明
    // 证明完成后自动提交到 L1 验证
}
```

### 手动调用验证

```go
// 如果需要手动触发验证
proofResult := &prover.ProofResult{
    StatusCode: prover.ProveSuccess,
    ProofID:    "proof_123",
    Proof: &prover.Proof{
        FromBlockNum: 1,
        ToBlockNum:   10,
        ZKProof:      []byte{/* proof data */},
    },
}

err := zkRollup.SendProofToL1(proofResult)
if err != nil {
    log.Printf("验证失败: %v", err)
}
```

## ⚠️ 注意事项

### 1. Public Values 的匹配

**关键**：`extractPublicValues` 中提取的 public values 必须与你的 OpenVM 程序输出完全匹配！

- 检查你的 OpenVM 程序生成了哪些 public values
- 按照相同的顺序和格式提取
- 确保数量为 32 个（默认配置）

### 2. Commitments 的计算

`appExeCommit` 和 `appVmCommit` 的计算方式取决于：
- 你的 OpenVM 程序如何生成这些值
- 是否使用固定值或动态计算
- 建议在配置文件中使用固定值（如果已知）

### 3. Gas 费用

链上验证会消耗 Gas：
- 预估 Gas: ~500,000+
- 建议在测试网先测试
- 考虑 Gas 价格波动

### 4. 错误处理

当前实现中，如果验证失败：
- 会记录错误日志
- 不会重试
- 可以根据需要添加重试逻辑

## 🔍 调试建议

### 查看日志

```bash
# 关键日志信息
[INFO] start to prove blocks from 1 to 10, proofID: proof_xxx
[INFO] get proof: proof_xxx
[INFO] Verifying proof on L1, proofID: proof_xxx, proof size: 12345 bytes
[INFO] ✓ Proof verified successfully on L1! ProofID: proof_xxx, Blocks: 1-10
```

### 常见问题

1. **验证失败**
   - 检查 public values 是否匹配
   - 检查 commitments 是否正确
   - 确认 proof 数据完整

2. **合约调用失败**
   - 检查 Verifier 合约地址
   - 确认以太坊节点连接正常
   - 检查账户余额（如果需要发送交易）

3. **证明数据为空**
   - 检查 Axiom API 是否正确返回
   - 确认证明已经完全生成

## 📚 相关文档

- [OpenVM Documentation](https://docs.openvm.dev/)
- [Axiom API Documentation](https://docs.axiom.xyz/)
- [zkrollup/contracts/README.md](./contracts/README.md) - 合约 ABI 使用文档

## 🚀 下一步

1. **自定义 public values 提取逻辑**
   - 根据你的 OpenVM 程序调整 `extractPublicValues`

2. **实现 commitments 计算**
   - 实现正确的 `calculateAppExeCommit` 和 `calculateAppVmCommit`

3. **添加额外功能**
   - 将验证结果存储到其他合约
   - 实现重试机制
   - 添加监控和告警

4. **优化性能**
   - 批量验证多个证明
   - 并行处理证明
   - 缓存验证结果

