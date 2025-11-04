package utxo

import (
	"github.com/yu-org/yu/common"
	"github.com/yu-org/yu/core/tripod"
	"math/big"
)

type UTXO struct {
	Capacity   uint64  `json:"capacity"`
	Owner      *Script `json:"owner"`
	Transition *Script `json:"transition"`

	UdtName string   `json:"udt_name"`
	Amount  *big.Int `json:"amount"`
}

type Script struct {
	Name     string `json:"name"`
	Code     []byte `json:"code"`
	Args     []byte `json:"args"`
	Capacity uint64 `json:"capacity"`
}

func (s *Script) ID() string {
	return common.Bytes2Hex(common.Sha256(s.Code))
}

type Account struct {
	*tripod.Tripod
	UtxoLib   []*UTXO
	ScriptLib []*Script
}