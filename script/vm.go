package script

type VM interface {
	Run(script *Script, args []byte) (*VMResult, error)
}

type VMResult struct {
	Output  []byte `json:"output"`
	Error   string `json:"error"`
	GasCost uint64 `json:"gas_cost"`
}
