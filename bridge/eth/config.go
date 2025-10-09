package eth

type Config struct {
	L1ClientAddress            string `toml:"l1_client_address"`
	ParentLayerContractAddress string `toml:"parentlayer_contract_address"`
	ChildLayerContractAddress  string `toml:"childlayer_contract_address"`
}
