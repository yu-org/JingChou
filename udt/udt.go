package udt

import (
	"math/big"
)

type UDT struct {
	Name          string      `json:"name"` // global unique name
	Description   string      `json:"description"`
	OriginalToken *ChainToken `json:"original_token,omitempty"`
	Total         *big.Int    `json:"total"`
	Locked        *big.Int    `json:"locked"`
	Issued        *big.Int    `json:"issued"`
}

type ChainToken struct {
	ChainURL     string `json:"chain_url"`
	TokenAddress []byte `json:"token_address"`
}

var NativeToken = UDT{
	Name:          "JingChou",
	Description:   "JingChou Chain Native Token",
	OriginalToken: nil,
	Total:         new(big.Int).SetUint64(1000000000000000000),
	Locked:        new(big.Int).SetUint64(0),
	Issued:        new(big.Int).SetUint64(1000000000000000),
}
