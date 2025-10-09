# App Commitments 说明文档

## 📋 概述

`appExeCommit` 和 `appVmCommit` 是 OpenVM 验证系统中的两个关键承诺值，用于确保证明的安全性和一致性。

## 🔐 appExeCommit - 应用执行承诺

### 含义
- **程序身份标识**：ELF 文件的哈希值
- **作用**：确保证明对应的是特定版本的程序
- **特性**：固定不变（除非程序升级）

### 计算方式
```
appExeCommit = SHA256(ELF_FILE_BYTES)
```

### 实现位置

#### 1. AxiomProver 中自动计算

在 `prover/axiom.go` 中：

```go
type AxiomProver struct {
    // ...
    appExeCommit [32]byte  // 初始化时自动计算
}

func NewAxiomProver(cfg *config.ProverConfig) (Prover, error) {
    // ...
    
    // 从 ELF 文件计算 appExeCommit
    if cfg.ElfPath != "" {
        appExeCommit, err := calculateAppExeCommitFromELF(cfg.ElfPath)
        if err != nil {
            return nil, err
        }
        prover.appExeCommit = appExeCommit
    }
    
    return prover, nil
}

// 辅助函数
func calculateAppExeCommitFromELF(elfPath string) ([32]byte, error) {
    elfData, err := os.ReadFile(elfPath)
    if err != nil {
        return [32]byte{}, err
    }
    return sha256.Sum256(elfData), nil
}

// 对外接口
func (a *AxiomProver) GetAppExeCommit() [32]byte {
    return a.appExeCommit
}
```

#### 2. ZkRollup 中使用

在 `zkrollup/zkrollup.go` 中：

```go
func (z *ZkRollup) calculateAppExeCommit(proofResult *prover.ProofResult) [32]byte {
    // 方式1：从配置获取（如果配置了）
    if z.cfg.AppExeCommit != "" {
        return ethcommon.HexToHash(z.cfg.AppExeCommit)
    }
    
    // 方式2：从 AxiomProver 获取（推荐）
    if axiomProver, ok := z.prover.(*prover.AxiomProver); ok {
        commit := axiomProver.GetAppExeCommit()
        if commit != [32]byte{} {
            return commit
        }
    }
    
    // 如果都没有，返回零值并警告
    logrus.Warn("AppExeCommit not available!")
    return [32]byte{}
}
```

## 🖥️ appVmCommit - 应用 VM 承诺

### 含义
- **VM 配置标识**：VM 配置的哈希值
- **作用**：确保证明使用的是特定的 VM 配置
- **特性**：固定不变（除非 VM 配置升级）

### 计算方式
```
appVmCommit = Hash(VM_CONFIG)
```

VM 配置包括：
- 指令集支持
- 内存限制
- 扩展功能
- 聚合参数

### 获取方式

#### 方式1：从配置文件（推荐）

```toml
[zkrollup]
app_vm_commit = "0xfedcba0987654321..."  # 从 OpenVM/Axiom 文档获取
```

#### 方式2：使用默认值

如果使用 OpenVM 默认 VM 配置，可以使用官方提供的默认哈希值。

#### 方式3：从 Axiom API 获取

如果 Axiom API 在响应中提供了这个值，可以从响应中提取。

## 📊 工作流程

```
初始化阶段：
┌─────────────────────┐
│  NewAxiomProver()   │
├─────────────────────┤
│ 1. 读取 ELF 文件    │
│ 2. 计算 SHA256      │
│ 3. 存储为           │
│    appExeCommit     │
└──────────┬──────────┘
           │
           ▼
    [AxiomProver 实例]
    - programID
    - appExeCommit ✓
    
验证阶段：
┌─────────────────────┐
│  SendProofToL1()    │
├─────────────────────┤
│ 获取 appExeCommit:  │
│ 1. 从配置？         │
│ 2. 从 prover ✓      │
│                     │
│ 获取 appVmCommit:   │
│ 1. 从配置 ✓         │
│ 2. 使用默认值       │
└─────────────────────┘
```

## 🎯 优势

### 之前的实现（错误）
```go
// ❌ 每次都用不同的值
appExeCommit = proofResult.Proof.PreStateRoot  // 变化的！
appVmCommit  = proofResult.Proof.NewStateRoot  // 变化的！

// 后果：验证器认为每次都是不同的程序和 VM
// → 验证失败！
```

### 现在的实现（正确）
```go
// ✅ 使用固定的值
appExeCommit = prover.GetAppExeCommit()  // 固定：ELF 哈希
appVmCommit  = config.AppVmCommit        // 固定：VM 配置哈希

// 结果：验证器正确识别程序和 VM
// → 验证成功！
```

## 📝 配置示例

```toml
[zkrollup]
# OpenVM Verifier 合约地址
l1_verifier_addr = "0xYourVerifierAddress"

# App VM 承诺（必需）- 从 OpenVM 文档或 Axiom 获取
app_vm_commit = "0xfedcba0987654321fedcba0987654321fedcba0987654321fedcba0987654321"

# App Exe 承诺（可选）- 如果不配置，会从 ELF 文件自动计算
# app_exe_commit = "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"

[zkrollup.prover]
# ELF 文件路径（必需，用于计算 appExeCommit）
elf_path = "/path/to/reth.elf"
```

## 🔍 调试

### 查看计算的 appExeCommit

```go
// 在初始化后打印
prover, err := prover.NewAxiomProver(cfg)
appExeCommit := prover.GetAppExeCommit()
fmt.Printf("appExeCommit: 0x%x\n", appExeCommit)
```

### 手动验证

```bash
# 计算 ELF 文件的 SHA256
sha256sum /path/to/reth.elf

# 应该与 appExeCommit 匹配
```

## ⚠️ 重要提醒

1. **appExeCommit**
   - ✅ 自动从 ELF 文件计算
   - ✅ 存储在 AxiomProver 中
   - ✅ zkrollup 直接调用 `GetAppExeCommit()` 获取
   - ⚠️ ELF 文件不能改变，否则哈希会变化

2. **appVmCommit**
   - ⚠️ 必须手动配置
   - ⚠️ 必须与部署的 Verifier 合约使用的 VM 配置匹配
   - ⚠️ 询问 Axiom 或查看 OpenVM 文档获取正确的值

## 🚀 使用建议

### 开发阶段
```toml
# 只需要配置 elf_path，appExeCommit 自动计算
elf_path = "/path/to/reth.elf"

# appVmCommit 使用测试值
app_vm_commit = "0x0000000000000000000000000000000000000000000000000000000000000000"
```

### 生产环境
```toml
# 配置正确的 VM 承诺（从 Axiom/OpenVM 获取）
app_vm_commit = "0x实际的VM配置哈希值"

# appExeCommit 依然自动计算（或手动配置固定值）
elf_path = "/path/to/reth.elf"
```

## 📚 相关代码

- `prover/axiom.go` - `calculateAppExeCommitFromELF()` 和 `GetAppExeCommit()`
- `zkrollup/zkrollup.go` - `calculateAppExeCommit()` 和 `calculateAppVmCommit()`
- `zkrollup/config/config.go` - 配置定义

