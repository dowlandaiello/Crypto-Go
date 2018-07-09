package orders

import (
	"errors"
	"strings"
	"time"

	"github.com/mitsukomegumi/Crypto-Go/src/accounts"
	"github.com/mitsukomegumi/Crypto-Go/src/common"
	"github.com/mitsukomegumi/Crypto-Go/src/pairs"
)

// Order - definition of order, fields attributed to a single order
type Order struct {
	Filled    bool      `json:"filled"`
	FillTime  time.Time `json:"filltime"`
	FillPrice float64   `json:"fillprice"` // Set price for order to fill at

	IssuanceTime time.Time `json:"issuancetime"` // IssuanceTime - time at which order is created

	Amount float64 `json:"amount"`

	OrderType string     `json:"ordertype"` // OrderType - BUY, SELL
	OrderFee  float64    `json:"orderfee"`
	OrderPair pairs.Pair `json:"tradingpair"` // OrderPair - startingpair (BTC, ETH, LTC), endingpair (BTC, ETH, LTC)

	Issuer *accounts.Account `json:"issuer"` // Account creating order

	OrderID string `json:orderid` // Order's hash
}

// NewOrder - creates, retursn new instance of order struct
func NewOrder(account *accounts.Account, ordertype string, tradingpair pairs.Pair, amount float64, fillprice float64) (Order, error) {
	ordertype = strings.ToUpper(ordertype)                                                                                // Used to check validity of order type
	if amount <= account.WalletBalances[common.IndexInSlice(tradingpair.StartingSymbol, []string{"BTC", "LTC", "ETH"})] { // Checks that amount is not more than account's balance
		rOrder := Order{Filled: false, IssuanceTime: time.Now().UTC(), Amount: (1.0 - common.FeeRate) * amount, OrderType: ordertype, OrderPair: tradingpair, Issuer: account, OrderID: "", OrderFee: common.FeeRate * amount}

		hash, err := common.Hash(rOrder) // Creates order hash

		if err != nil {
			return rOrder, err
		}

		rOrder.OrderID = hash

		account.Orders = append(account.Orders, hash) // Appends

		//account.Balance -= (rOrder.OrderFee + rOrder.Amount) // No clue

		return rOrder, nil
	}
	return Order{}, errors.New("insufficient balance") // Triggered on insufficient balance, nil order
}

// FillOrder - fills order
func FillOrder(order *Order) {
	if order.Issuer.WalletBalances[common.IndexInSlice(order.OrderPair.StartingSymbol, common.AvailableSymbols)] >= (order.Amount + order.OrderFee) { // Checks that order value is not more than account balance
		order.Filled = true
		order.FillTime = time.Now().UTC()
		order.Issuer.WalletBalances[common.IndexInSlice(order.OrderPair.EndingSymbol, common.AvailableSymbols)] += order.Amount   // Adds actual order amount (not including fees) to wallet
		order.Issuer.WalletBalances[common.IndexInSlice(order.OrderPair.StartingSymbol, common.AvailableSymbols)] -= order.Amount // Subtracts order value from wallet

		//TODO: move assets
	}
}
