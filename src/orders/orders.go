package orders

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dowlandaiello/Crypto-Go/src/market"

	"github.com/dowlandaiello/Crypto-Go/src/accounts"
	"github.com/dowlandaiello/Crypto-Go/src/common"
	"github.com/dowlandaiello/Crypto-Go/src/pairs"
)

// Order - definition of order, fields attributed to a single order
type Order struct {
	Filled    bool      `json:"filled"`
	FillTime  time.Time `json:"filltime"`
	FillPrice float64   `json:"fillprice"` // Set price for order to fill at

	IssuanceTime time.Time `json:"issuancetime"` // IssuanceTime - time at which order is created

	Amount float64 `json:"amount"`

	OrderType string      `json:"ordertype"` // OrderType - BUY, SELL
	OrderFee  float64     `json:"orderfee"`
	OrderPair *pairs.Pair `json:"tradingpair"` // OrderPair - startingpair (BTC, ETH, LTC), endingpair (BTC, ETH, LTC)

	Issuer *accounts.Account `json:"issuer"` // Account creating order

	OrderID string `json:"orderid" bson:"orderid"` // Order's hash
}

// NewOrder - creates, retursn new instance of order struct
func NewOrder(account *accounts.Account, ordertype string, tradingpair pairs.Pair, amount float64, fillprice float64) (Order, error) {
	currentPrice, _ := market.CheckPrice(tradingpair)

	ordertype = strings.ToUpper(ordertype) // Used to check validity of order type

	fmt.Println("current " + tradingpair.ToString() + " price: " + common.FloatToString(currentPrice))

	if amount <= account.WalletBalances[common.IndexInSlice(tradingpair.StartingSymbol, []string{"BTC", "LTC", "ETH"})] && tradingpair.StartingSymbol != tradingpair.EndingSymbol && ordertype == "SELL" {
		order, err := newOrder(fillprice, amount, ordertype, tradingpair, account)

		if err != nil {
			return Order{}, err
		}

		return order, nil
	} else if amount*fillprice <= account.WalletBalances[common.IndexInSlice(tradingpair.EndingSymbol, []string{"BTC", "LTC", "ETH"})] && tradingpair.StartingSymbol != tradingpair.EndingSymbol && ordertype == "BUY" { // Checks that amount is not more than account's balance
		order, err := newOrder(fillprice, amount, ordertype, tradingpair, account)

		if err != nil {
			return Order{}, err
		}

		return order, nil
	}

	return Order{}, errors.New("insufficient balance") // Triggered on insufficient balance, nil order
}

func newOrder(fillprice float64, amount float64, ordertype string, tradingpair pairs.Pair, account *accounts.Account) (Order, error) {
	rOrder := Order{Filled: false, FillTime: time.Now().UTC(), FillPrice: fillprice, IssuanceTime: time.Now().UTC(), Amount: (1.0 - common.FeeRate) * amount, OrderType: ordertype, OrderPair: &tradingpair, Issuer: account, OrderID: "", OrderFee: common.FeeRate * amount}

	hash, err := common.Hash(rOrder) // Creates order hash

	if err != nil {
		return rOrder, err
	}

	rOrder.OrderID = hash

	(*account).Orders = append(account.Orders, hash) // Appends

	return rOrder, nil
}

// FillOrder - fills order
func FillOrder(order *Order, password string) error {
	currentPrice, err := market.CheckPrice(*order.OrderPair)

	fmt.Println("\ncurrent " + order.OrderPair.ToString() + " price: " + common.FloatToString(currentPrice))

	if err != nil {
		return err
	}

	if order.Issuer.WalletBalances[common.IndexInSlice(order.OrderPair.EndingSymbol, common.AvailableSymbols)] >= ((order.Amount*currentPrice)+(order.OrderFee*currentPrice)) && currentPrice == order.FillPrice { // Checks that order value is not more than account balance
		fmt.Println("\norder filled at " + common.FloatToString(currentPrice) + " with destination of " + common.FloatToString(order.FillPrice))
		order.Filled = true
		order.FillTime = time.Now().UTC()

		keys, err := accounts.DecryptPrivateKeys(order.Issuer.WalletRawHashedKeys, password)

		if err != nil {
			return err
		}

		if order.OrderType == "BUY" {

			order.Issuer.WalletBalances[common.IndexInSlice(order.OrderPair.StartingSymbol, common.AvailableSymbols)] += order.Amount                                            // Adds actual order amount (not including fees) to wallet
			order.Issuer.WalletBalances[common.IndexInSlice(order.OrderPair.EndingSymbol, common.AvailableSymbols)] -= (order.Amount*currentPrice + order.OrderFee*currentPrice) // Subtracts order value from wallet

			err := order.Issuer.MoveAssets(order.OrderPair.EndingSymbol, order.Issuer.WalletAddresses[common.IndexInSlice(order.OrderPair.EndingSymbol, []string{"BTC", "LTC", "ETH"})], common.ExchangeWallet, keys[common.IndexInSlice(order.OrderPair.EndingSymbol, []string{"BTC", "LTC", "ETH"})], (order.Amount*currentPrice + order.OrderFee*currentPrice))

			if err != nil {
				return err
			}
		} else if order.OrderType == "SELL" {
			order.Issuer.WalletBalances[common.IndexInSlice(order.OrderPair.StartingSymbol, common.AvailableSymbols)] -= order.Amount                                            // Sells actual order amount from wallet
			order.Issuer.WalletBalances[common.IndexInSlice(order.OrderPair.EndingSymbol, common.AvailableSymbols)] += (order.Amount*currentPrice + order.OrderFee*currentPrice) // Adds order value to wallet

			err := order.Issuer.MoveAssets(order.OrderPair.StartingSymbol, order.Issuer.WalletAddresses[common.IndexInSlice(order.OrderPair.StartingSymbol, []string{"BTC", "LTC", "ETH"})], common.ExchangeWallet, keys[common.IndexInSlice(order.OrderPair.StartingSymbol, []string{"BTC", "LTC", "ETH"})], (order.Amount + order.OrderFee))

			if err != nil {
				return err
			}
		}

		return nil

		//TODO: move assets
	}
	return errors.New("invalid request; insufficient balance or invalid fill price")
}
