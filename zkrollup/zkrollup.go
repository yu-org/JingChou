package zkrollup

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/yu-org/yu/core/tripod"
	"github.com/yu-org/yu/core/types"
)

type ZkRollup struct {
	*tripod.Tripod

	cfg    *Config
	ethCli *ethclient.Client
}

func NewZkRollup(cfg *Config) (*ZkRollup, error) {
	ethCli, err := ethclient.Dial(cfg.L1ChainAddr)
	if err != nil {
		return nil, err
	}
	return &ZkRollup{
		Tripod: tripod.NewTripod(),
		cfg:    cfg,
		ethCli: ethCli,
	}, nil
}

func (z *ZkRollup) StartBlock(block *types.Block) {
	//TODO implement me
	panic("implement me")
}

func (z *ZkRollup) EndBlock(block *types.Block) {
	//TODO implement me
	panic("implement me")
}

func (z *ZkRollup) FinalizeBlock(block *types.Block) {
	//TODO implement me
	panic("implement me")
}
