package udt

import (
	"encoding/json"
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

func (ut *UdtTripod) AddUDT(udt *UDT) error {
	byt, err := json.Marshal(udt)
	if err != nil {
		return err
	}
	ut.Set([]byte(udt.Name), byt)
	return nil
}
