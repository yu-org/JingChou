package orderbook

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/yu-org/yu/core/context"
	"github.com/yu-org/yu/core/tripod"
	"github.com/yu-org/yu/core/types"
	"sort"
)

type Orderbook struct {
	*tripod.Tripod

	BuyOrders  map[OrderPair]Orders
	SellOrders map[OrderPair]Orders
}

func (ob *Orderbook) StartBlock(block *types.Block) {
	//TODO implement me
	panic("implement me")
}

func (ob *Orderbook) EndBlock(block *types.Block) {
	//TODO implement me
	panic("implement me")
}

func (ob *Orderbook) FinalizeBlock(block *types.Block) {
	//TODO implement me
	panic("implement me")
}

func NewOrderbook() *Orderbook {
	ob := &Orderbook{
		Tripod:     tripod.NewTripod(),
		BuyOrders:  make(map[OrderPair]Orders),
		SellOrders: make(map[OrderPair]Orders),
	}
	ob.SetWritings(ob.AddOrder, ob.CancelOrder)
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
	pair := req.Order.Pair()
	switch req.Order.Type {
	case Buy:
		ob.BuyOrders[pair] = append(ob.BuyOrders[pair], req.Order)
		sort.Sort(ob.BuyOrders[pair])
	case Sell:
		ob.SellOrders[pair] = append(ob.SellOrders[pair], req.Order)
		sort.Sort(ob.SellOrders[pair])
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
