package eth

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/yu-org/yu/apps/eth/evm"
	"github.com/yu-org/yu/core/tripod"
	"github.com/yu-org/yu/core/types"
)

type EthRelayer struct {
	*tripod.Tripod
	cfg      *Config
	Solidity evm.Solidity `tripod:"solidity"`
	ethCli   *ethclient.Client
}

func NewETHRelayer(cfg *Config) (*EthRelayer, error) {
	ethCli, err := ethclient.Dial(cfg.L1ClientAddress)
	if err != nil {
		return nil, err
	}
	return &EthRelayer{
		Tripod: tripod.NewTripod(),
		cfg:    cfg,
		ethCli: ethCli,
	}, nil
}

func (eth *EthRelayer) StartBlock(block *types.Block) {

}

func (eth *EthRelayer) EndBlock(block *types.Block) {

}

func (eth *EthRelayer) FinalizeBlock(block *types.Block) {

}
