package orderbook

import (
	"encoding/json"
	"github.com/yu-org/JingChou/udt"
	"github.com/yu-org/yu/common"
	"math/big"
)

type Order struct {
	Type OrderType `json:"type"`

	OrderToken udt.TokenID `json:"order_token"`
	Amount     *big.Int    `json:"amount"`

	PricingToken udt.TokenID `json:"pricing_token"`
	Price        *big.Int    `json:"price"`

	Owner *OrderScript `json:"owner"`
	// Taker *OrderScript `json:"taker"`

	Capacity uint64 `json:"capacity"`
}

func (o *Order) ID() (string, error) {
	byt, err := json.Marshal(o)
	if err != nil {
		return "", err
	}
	return common.Bytes2Hex(common.Sha256(byt)), nil
}

func (o *Order) Pair() OrderPair {
	return OrderPair{
		OrderToken:   o.OrderToken,
		PricingToken: o.PricingToken,
	}
}

type OrderPair struct {
	OrderToken   udt.TokenID `json:"order_token"`
	PricingToken udt.TokenID `json:"pricing_token"`
}

type Orders []*Order

func (o Orders) Len() int {
	return len(o)
}

func (o Orders) Less(i, j int) bool {
	return o[i].Price.Cmp(o[j].Price) < 0
}

func (o Orders) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

type OrderType uint8

const (
	Buy OrderType = iota
	Sell
)

func MatchOrders(buys, sells []*Order) {
	// TODO: Optimize the matching algorithm
	for _, buy := range buys {
		for _, sell := range sells {
			if buy.PricingToken != sell.PricingToken {
				continue
			}
			if buy.OrderToken != sell.OrderToken {
				continue
			}
			if buy.Price.Cmp(sell.Price) < 0 {
				continue
			}

		}
	}
}

type OrderScript struct {
	//Amount       *big.Int    `json:"amount"`
	//PricingToken udt.TokenID `json:"pricing_token"`
	//Code         []byte      `json:"code"`
	Args []byte `json:"args"`
}
