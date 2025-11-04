package orderbook

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/yu-org/yu/core/context"
	"github.com/yu-org/yu/core/tripod"
)

type Orderbook struct {
	*tripod.Tripod

	BuyOrders  []*Order
	SellOrders []*Order
}

func NewOrderbook() *Orderbook {
	ob := &Orderbook{
		Tripod: tripod.NewTripod(),
	}
	ob.SetWritings(ob.AddOrder, ob.EatOrder, ob.CancelOrder)
	ob.SetReadings(ob.QueryOrder)
	return ob
}

type AddOrderRequest struct {
	Order       *Order   `json:"order"`
	FromUtxoIDs []string `json:"from_utxo_ids"`
	Args        []byte   `json:"args"`
}

func (ob *Orderbook) AddOrder(ctx *context.WriteContext) error {
	req := new(AddOrderRequest)
	if err := ctx.BindJson(req); err != nil {
		return err
	}
	return nil
}

type EatOrderRequest struct {
	OrderID string `json:"order_id"`
	EatArgs []byte `json:"eat_args"`
}

func (ob *Orderbook) EatOrder(ctx *context.WriteContext) error {
	req := new(EatOrderRequest)
	if err := ctx.BindJson(req); err != nil {
		return err
	}
	return nil
}

type CancelOrderRequest struct {
	OrderID    string `json:"order_id"`
	CancelArgs []byte `json:"cancel_args"`
}

func (ob *Orderbook) CancelOrder(ctx *context.WriteContext) error {
	req := new(CancelOrderRequest)
	if err := ctx.BindJson(req); err != nil {
		return err
	}
	return nil
}

func (ob *Orderbook) QueryOrder(ctx *context.ReadContext) {

}
