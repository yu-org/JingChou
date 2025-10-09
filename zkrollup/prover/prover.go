package prover

import (
	"github.com/yu-org/yu/common"
	"github.com/yu-org/yu/core/types"
)

type Prover interface {
	GenerateProof(blockBatch []*types.Block, proofChan chan *ProofResult) (proofID string, err error)
	GetProof(proofID string) (*ProofResult, error)
	CancelProof(proofID string) (*ProofResult, error)
}

type ProofResult struct {
	StatusCode ProofStatusCode `json:"status_code"`
	ProofID    string          `json:"proof_id"`
	Proof      *Proof          `json:"proof"`
}

type ProofStatusCode uint8

const (
	ProveSuccess ProofStatusCode = iota
	ProveFailed
	Proving
	ProvePending
)

func (pc ProofStatusCode) String() string {
	switch pc {
	case ProveSuccess:
		return "GenerateProof Success"
	case ProveFailed:
		return "GenerateProof Failed"
	case Proving:
		return "Proving Now"
	case ProvePending:
		return "GenerateProof Pending"
	default:
		return "Unknown ProofStatusCode"
	}
}

type Proof struct {
	FromBlockNum uint64      `json:"from_block_num"`
	ToBlockNum   uint64      `json:"to_block_num"`
	PreStateRoot common.Hash `json:"pre_state_root"`
	NewStateRoot common.Hash `json:"new_state_root"`
	ZKProof      []byte      `json:"zk_proof"`
}
