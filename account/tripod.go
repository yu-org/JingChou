package account

import (
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
	//TODO: verify owner script
	return nil
}

func (a *Account) InitChain(block *types.Block) {

}

type ClaimAccountRequest struct {
	ID    AccountID `json:"id"`
	Owner *Script   `json:"owner"`
}

func (a *Account) ClaimAccount(ctx *context.WriteContext) error {
	req := new(ClaimAccountRequest)
	if err := ctx.BindJson(req); err != nil {
		return err
	}
	if req.Owner == nil {
		return errors.New("owner-script is nil")
	}

	claimed := &CommonAccount{
		Id:    req.ID,
		Owner: req.Owner,
		UDTs:  nil, //TODO: get the old UDTs from the state
	}
	return nil
}

type TransferRequest struct {
	FromID    AccountID           `json:"from_id"`
	OwnerArgs []byte              `json:"owner_args"`
	To        AccountID           `json:"to"`
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
	FromID    AccountID `json:"from_id"`
	OwnerArgs []byte    `json:"owner_args"`
	ScriptID  string    `json:"script_id"`
	Args      []byte    `json:"args"`
}

func (a *Account) InvokeScript(ctx *context.WriteContext) error {
	req := new(InvokeScriptRequest)
	if err := ctx.BindJson(req); err != nil {
		return err
	}

	return nil
}
