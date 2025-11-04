package orderbook

import (
	"encoding/json"
	"github.com/yu-org/JingChou/utxo"
	"github.com/yu-org/yu/common"
	"math/big"
)

type Order struct {
	Type OrderType `json:"type"`

	UDTs map[string]*big.Int `json:"udts"` // map[UdtName]Amount

	Owner *utxo.Script `json:"owner"`
	Eater *utxo.Script `json:"eater"`

	Capacity uint64 `json:"capacity"`
}

func (o *Order) ID() (string, error) {
	byt, err := json.Marshal(o)
	if err != nil {
		return "", err
	}
	return common.Bytes2Hex(common.Sha256(byt)), nil
}

type OrderType uint8

const (
	Buy OrderType = iota
	Sell
)

type OrderScript struct {
	Args map[string]*big.Int `json:"args"`
	Code []byte              `json:"code"`
}
