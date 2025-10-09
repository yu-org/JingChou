package config

type Config struct {
	BlockBatchSizeForProve uint         `toml:"block_batch_size_for_prove"`
	Prover                 ProverConfig `toml:"prover"`

	L1ChainAddr    string `toml:"l1_chain_addr"`
	L1ContractAddr string `toml:"l1_contract_addr"`
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
