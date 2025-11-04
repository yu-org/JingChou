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
