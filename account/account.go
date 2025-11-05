package account

import (
	"github.com/yu-org/JingChou/udt"
)

type Account struct {
	Owner   string                   `json:"owner"`
	UDTs    map[udt.TokenID]*udt.UDT `json:"udts"`
	Scripts []string                 `json:"scripts"`
}

func (a *Account) VerifyOwner(args []byte) error {
	return nil
}
