package config

type Config struct {
	BlockBatchSizeForProve uint         `toml:"block_batch_size_for_prove"`
	Prover                 ProverConfig `toml:"prover"`

	L1ChainAddr    string `toml:"l1_chain_addr"`    // L1 以太坊节点地址
	L1VerifierAddr string `toml:"l1_verifier_addr"` // OpenVM Halo2 Verifier 合约地址
	L1ContractAddr string `toml:"l1_contract_addr"` // 其他 L1 合约地址（如果需要）
	AppExeCommit   string `toml:"app_exe_commit"`   // App 执行承诺（可配置或动态生成）
	AppVmCommit    string `toml:"app_vm_commit"`    // App VM 承诺（可配置或动态生成）
}

type ProverConfig struct {
	URL          string `toml:"url"`
	ApiKey       string `toml:"api_key"`
	ElfPath      string `toml:"elf_path"`      // reth ELF 文件路径
	ProgramID    string `toml:"program_id"`    // 注册后的 program ID（可选，如果已注册）
	VMConfigID   string `toml:"vm_config_id"`  // VM 配置 ID
	PollInterval int    `toml:"poll_interval"` // 轮询间隔（秒），默认 5 秒
	PollTimeout  int    `toml:"poll_timeout"`  // 轮询超时时间（秒），默认 7200 秒（2小时）
	ProofType    string `toml:"proof_type"`    // 证明类型：stark 或 evm，默认 stark
}
