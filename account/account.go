package account

import (
	"encoding/json"
	"github.com/yu-org/JingChou/udt"
	"github.com/yu-org/yu/common"
)

type CommonAccount struct {
	// Id      AccountID            `json:"id"`
	Owner   ScriptID             `json:"owner"`
	UDTs    map[string]*udt.UDT  `json:"udts"`
	Scripts map[ScriptID]*Script `json:"scripts"`
}

type (
	// AccountID string
	ScriptID string
)

//func (a *CommonAccount) ID() AccountID {
//	if a.Id != "" {
//		return a.Id
//	}
//	return AccountID(a.Owner)
//}

type ScriptType int

const (
	Once ScriptType = iota
	Permanent
)

type Script struct {
	Type     ScriptType `json:"type"`
	Code     []byte     `json:"code"`
	Args     []byte     `json:"args,omitempty"`
	GasToken string     `json:"gas_token,omitempty"`
}

func (s *Script) Id() (ScriptID, error) {
	byt, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return ScriptID(common.Bytes2Hex(common.Sha256(byt))), nil
}
