package account

import (
	"encoding/json"
	"errors"
	"github.com/yu-org/JingChou/script"
	"github.com/yu-org/JingChou/udt"
	"github.com/yu-org/yu/core/context"
	"github.com/yu-org/yu/core/tripod"
	"github.com/yu-org/yu/core/types"
	"math/big"
)

type AccountTripod struct {
	*tripod.Tripod
}

func NewAccountTripod() *AccountTripod {
	a := &AccountTripod{
		Tripod: tripod.NewTripodWithName("account"),
	}
	a.SetInit(a)
	a.SetTxnChecker(a)
	a.SetWritings(a.ClaimAccount, a.Transfer, a.InvokeScript)

	return a
}

func (a *AccountTripod) CheckTxn(tx *types.SignedTxn) error {
	//TODO: 1. verify owner script
	//TODO: 2. verify UDTs is enough
	//TODO: 3. verify UDTs include native token
	return nil
}

func (a *AccountTripod) InitChain(block *types.Block) {

}

type ClaimAccountRequest struct {
	Owner       script.ScriptID `json:"owner"`
	OwnerScript *script.Script  `json:"owner_script"`
	Args        []byte          `json:"args,omitempty"`
}

func (a *AccountTripod) ClaimAccount(ctx *context.WriteContext) error {
	req := new(ClaimAccountRequest)
	if err := ctx.BindJson(req); err != nil {
		return err
	}
	if req.Owner == "" {
		return errors.New("owner-script is nil")
	}

	oldAccountByt, err := a.Get([]byte(req.Owner))
	if err != nil {
		return err
	}
	oldAccount := new(Account)
	if err = json.Unmarshal(oldAccountByt, oldAccount); err != nil {
		return err
	}

	if oldAccount.Owner != req.Owner {
		return errors.New("owner-script is not the same")
	}
	claimed := &Account{
		Owner:   req.Owner,
		UDTs:    oldAccount.UDTs,
		Scripts: []script.ScriptID{req.Owner},
	}
	byt, err := json.Marshal(claimed)
	if err != nil {
		return err
	}
	a.Set([]byte(req.Owner), byt)
	return nil
}

type TransferRequest struct {
	FromID    script.ScriptID          `json:"from_id"`
	OwnerArgs []byte                   `json:"owner_args"`
	To        script.ScriptID          `json:"to"`
	UDTs      map[udt.TokenID]*big.Int `json:"udts"`
}

func (a *AccountTripod) Transfer(ctx *context.WriteContext) error {
	req := new(TransferRequest)
	if err := ctx.BindJson(req); err != nil {
		return err
	}
	return nil
}

type InvokeScriptRequest struct {
	FromID    script.ScriptID `json:"from_id"`
	OwnerArgs []byte          `json:"owner_args"`
	ScriptID  string          `json:"script_id"`
	Args      []byte          `json:"args"`
}

func (a *AccountTripod) InvokeScript(ctx *context.WriteContext) error {
	req := new(InvokeScriptRequest)
	if err := ctx.BindJson(req); err != nil {
		return err
	}

	return nil
}
