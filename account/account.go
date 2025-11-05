package account

import (
	"github.com/yu-org/JingChou/script"
	"github.com/yu-org/JingChou/udt"
)

type Account struct {
	Owner   script.ScriptID          `json:"owner"`
	UDTs    map[udt.TokenID]*udt.UDT `json:"udts"`
	Scripts []script.ScriptID        `json:"scripts"`
}
