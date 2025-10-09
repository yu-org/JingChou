package zkrollup

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"github.com/yu-org/JingChou/zkrollup/config"
	"github.com/yu-org/JingChou/zkrollup/prover"
	"github.com/yu-org/yu/common"
	"github.com/yu-org/yu/core/tripod"
	"github.com/yu-org/yu/core/types"
)

type ZkRollup struct {
	*tripod.Tripod

	cfg       *config.Config
	ethCli    *ethclient.Client
	prover    prover.Prover
	proofChan chan *prover.ProofResult
}

func NewZkRollup(cfg *config.Config) (*ZkRollup, error) {
	ethCli, err := ethclient.Dial(cfg.L1ChainAddr)
	if err != nil {
		return nil, err
	}
	prv, err := prover.NewAxiomProver(&cfg.Prover)
	if err != nil {
		return nil, err
	}
	proofChan := make(chan *prover.ProofResult, 10)
	return &ZkRollup{
		Tripod:    tripod.NewTripod(),
		cfg:       cfg,
		ethCli:    ethCli,
		prover:    prv,
		proofChan: proofChan,
	}, nil
}

func (z *ZkRollup) StartBlock(block *types.Block) {
	//TODO implement me
	panic("implement me")
}

func (z *ZkRollup) EndBlock(block *types.Block) {

}

func (z *ZkRollup) FinalizeBlock(block *types.Block) {
	if uint(block.Height)%z.cfg.BlockBatchSizeForProve != 0 {
		return
	}

	startProveBlockHeight := block.Height - common.BlockNum(z.cfg.BlockBatchSizeForProve) + 1

	blocks, err := z.Chain.GetRangeBlocks(startProveBlockHeight, block.Height)
	if err != nil {
		logrus.Errorf("get range blocks failed: %v", err)
		return
	}
	// 证明
	proofID, err := z.prover.GenerateProof(blocks, z.proofChan)
	if err != nil {
		logrus.Errorf("start to prove blocks from %d to %d failed: %v", startProveBlockHeight, block.Height, err)
		return
	}
	logrus.Infof("start to prove blocks from %d to %d, proofID: %s", startProveBlockHeight, block.Height, proofID)
}

func (z *ZkRollup) GetProof() {
	for {
		select {
		case proofResult := <-z.proofChan:
			logrus.Infof("get proof: %s", proofResult.ProofID)
			// TODO: Verify proof
			// 发送到 L1
			err := z.SendProofToL1(proofResult)
			if err != nil {
				logrus.Errorf("send proof to L1 failed: %v", err)
			}
		}
	}
}

func (z *ZkRollup) SendProofToL1(proofResult *prover.ProofResult) error {

}
