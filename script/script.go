package script

import (
	"encoding/json"

	"github.com/yu-org/JingChou/udt"
	"github.com/yu-org/yu/common"
)

type ScriptType int

const (
	Once ScriptType = iota
	Permanent
)

type Script struct {
	Type     ScriptType  `json:"type"`
	Code     []byte      `json:"code"`
	GasToken udt.TokenID `json:"gas_token,omitempty"`
}

func (s *Script) Id() (string, error) {
	byt, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return common.Bytes2Hex(common.Sha256(byt)), nil
}
