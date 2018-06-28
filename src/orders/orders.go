package orders

import (
	"time"

	"github.com/mitsukomegumi/FakeCrypto/src/accounts"
	"github.com/mitsukomegumi/FakeCrypto/src/common"
	"github.com/mitsukomegumi/FakeCrypto/src/pairs"
)

// Order - definition of order, fields attributed to a single order
type Order struct {
	Filled bool `json:"filled"`

	IssuanceTime time.Time `json:"issuancetime"`
	FillTime     time.Time `json:"filletime"`

	Amount int `json:"amount"`

	OrderType string     `json:"ordertype"`
	OrderFee  int        `json:"orderfee"`
	OrderPair pairs.Pair `json:"tradingpair"`

	Issuer *accounts.Account `json:"issuer"`

	ID string `json:"order"`
}

// NewOrder - creates, retursn new instance of order struct
func NewOrder(account *accounts.Account, ordertype string, tradingpair pairs.Pair, amount int) (Order, error) {
	rOrder := Order{Filled: false, IssuanceTime: time.Now().UTC(), Amount: amount, OrderType: ordertype, OrderPair: tradingpair, Issuer: account, ID: ""}

	hash, err := common.Hash(rOrder)

	if err != nil {
		return rOrder, err
	}

	rOrder.ID = hash

	return rOrder, nil
}
