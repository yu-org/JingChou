package script

import (
	"encoding/json"
	"errors"
	"github.com/yu-org/yu/core/context"

	"github.com/yu-org/yu/core/tripod"
)

type ScriptTripod struct {
	*tripod.Tripod

	vm VM
}

func NewScriptTripod() *ScriptTripod {
	st := &ScriptTripod{
		Tripod: tripod.NewTripodWithName("script"),
	}
	st.SetReadings(st.GetScript)
	return st
}

func (st *ScriptTripod) GetScript(ctx *context.ReadContext) {
	id := ctx.GetString("script_id")
	script, err := st.GetScriptById(id)
	if err != nil {
		ctx.ErrOk(err)
		return
	}
	ctx.JsonOk(script)
}

func (st *ScriptTripod) InvokeScript(id string, args []byte) ([]byte, error) {
	script, err := st.GetScriptById(id)
	if err != nil {
		return nil, err
	}
	return st.vm.Run(script, args)
}

func (st *ScriptTripod) AddScript(scpt *Script) error {
	if scpt == nil {
		return errors.New("added script is nil")
	}
	scptByt, err := json.Marshal(scpt)
	if err != nil {
		return err
	}
	id, err := scpt.Id()
	if err != nil {
		return err
	}
	st.Set([]byte(id), scptByt)
	return nil
}

func (st *ScriptTripod) GetScriptById(id string) (*Script, error) {
	scptByt, err := st.Get([]byte(id))
	if err != nil {
		return nil, err
	}
	scpt := new(Script)
	if err = json.Unmarshal(scptByt, scpt); err != nil {
		return nil, err
	}
	return scpt, nil
}
