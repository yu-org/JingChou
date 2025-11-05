package udt

import (
	"encoding/json"
	"errors"
	"github.com/yu-org/yu/core/context"
	"github.com/yu-org/yu/core/tripod"
)

type UdtTripod struct {
	*tripod.Tripod
}

func NewUdtTripod() *UdtTripod {
	ut := &UdtTripod{
		Tripod: tripod.NewTripodWithName("udt"),
	}
	return ut
}

func (ut *UdtTripod) GetUDT(ctx *context.ReadContext) {
	tokenID := ctx.GetString("token_id")
	udt, err := ut.GetUdt(TokenID(tokenID))
	if err != nil {
		ctx.ErrOk(err)
		return
	}
	ctx.JsonOk(udt)
}

func (ut *UdtTripod) AddUdt(udt *UDT) error {
	byt, err := json.Marshal(udt)
	if err != nil {
		return err
	}
	ut.Set([]byte(udt.Name), byt)
	return nil
}

func (ut *UdtTripod) GetUdt(id TokenID) (*UDT, error) {
	byt, err := ut.Get([]byte(id))
	if err != nil {
		return nil, err
	}
	udt := new(UDT)
	err = json.Unmarshal(byt, udt)
	return udt, err
}

func (ut *UdtTripod) DeleteUdt(id TokenID) error {
	if id.IsNative() {
		return errors.New("native token cannot be deleted")
	}
	ut.Delete([]byte(id))
	return nil
}
