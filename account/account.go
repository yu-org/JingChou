package account

import (
	"github.com/yu-org/JingChou/udt"
	"github.com/yu-org/yu/common"
)

//type DefaultAccount struct {
//	Owner   []byte     `json:"owner"`
//	UDTs    []*udt.UDT `json:"udts"`
//	Scripts []*Script  `json:"scripts"`
//}

type CommonAccount struct {
	Id      AccountID           `json:"id"`
	Owner   *Script             `json:"owner"`
	UDTs    map[string]*udt.UDT `json:"udts"`
	Scripts map[string]*Script  `json:"scripts"`
}

type AccountID string

func (a *CommonAccount) ID() AccountID {
	if a.Id != "" {
		return a.Id
	}
	return AccountID(a.Owner.ID())
}

type ScriptType int

const (
	Once ScriptType = iota
	Permanent
)

type Script struct {
	Type     ScriptType `json:"type"`
	Code     []byte     `json:"code"`
	GasToken string     `json:"gas_token,omitempty"`
}

func (s *Script) ID() string {
	return common.Bytes2Hex(common.Sha256(s.Code))
}
