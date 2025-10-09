package zkrollup

type Config struct {
	BlockBatchSizeForProve uint64       `toml:"block_batch_size_for_prove"`
	Prover                 ProverConfig `toml:"prover"`

	L1ChainAddr    string `toml:"l1_chain_addr"`
	L1ContractAddr string `toml:"l1_contract_addr"`
}

type ProverConfig struct {
	URL        string `toml:"url"`
	ApiKey     string `toml:"api_key"`
	ElfPath    string `toml:"elf_path"`     // reth ELF 文件路径
	ProgramID  string `toml:"program_id"`   // 注册后的 program ID（可选，如果已注册）
	VMConfigID string `toml:"vm_config_id"` // VM 配置 ID
}
