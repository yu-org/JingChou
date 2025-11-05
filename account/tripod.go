package account

import (
	"encoding/json"
	"errors"
	"github.com/yu-org/yu/core/context"
	"github.com/yu-org/yu/core/tripod"
	"github.com/yu-org/yu/core/types"
	"math/big"
)

type Account struct {
	*tripod.Tripod
}

func NewAccount() *Account {
	a := &Account{
		Tripod: tripod.NewTripod(),
	}
	a.SetInit(a)
	a.SetTxnChecker(a)
	a.SetWritings(a.ClaimAccount, a.Transfer, a.InvokeScript)

	return a
}

func (a *Account) CheckTxn(tx *types.SignedTxn) error {
	//TODO: 1. verify owner script
	//TODO: 2. verify UDTs is enough
	return nil
}

func (a *Account) InitChain(block *types.Block) {

}

type ClaimAccountRequest struct {
	Owner       ScriptID `json:"owner"`
	OwnerScript *Script  `json:"owner_script"`
	Args        []byte   `json:"args,omitempty"`
}

func (a *Account) ClaimAccount(ctx *context.WriteContext) error {
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
	oldAccount := new(CommonAccount)
	if err = json.Unmarshal(oldAccountByt, oldAccount); err != nil {
		return err
	}

	if oldAccount.Owner != req.Owner {
		return errors.New("owner-script is not the same")
	}
	claimed := &CommonAccount{
		Owner:   req.Owner,
		UDTs:    oldAccount.UDTs,
		Scripts: map[ScriptID]*Script{req.Owner: req.OwnerScript},
	}
	byt, err := json.Marshal(claimed)
	if err != nil {
		return err
	}
	a.Set([]byte(req.Owner), byt)
	return nil
}

type TransferRequest struct {
	FromID    ScriptID            `json:"from_id"`
	OwnerArgs []byte              `json:"owner_args"`
	To        ScriptID            `json:"to"`
	UDTs      map[string]*big.Int `json:"udts"`
}

func (a *Account) Transfer(ctx *context.WriteContext) error {
	req := new(TransferRequest)
	if err := ctx.BindJson(req); err != nil {
		return err
	}
	return nil
}

type InvokeScriptRequest struct {
	FromID    ScriptID `json:"from_id"`
	OwnerArgs []byte   `json:"owner_args"`
	ScriptID  string   `json:"script_id"`
	Args      []byte   `json:"args"`
}

func (a *Account) InvokeScript(ctx *context.WriteContext) error {
	req := new(InvokeScriptRequest)
	if err := ctx.BindJson(req); err != nil {
		return err
	}

	return nil
}
