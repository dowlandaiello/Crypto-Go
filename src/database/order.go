package database

import (
	"errors"

	"github.com/mitsukomegumi/Crypto-Go/src/orders"
	mgo "gopkg.in/mgo.v2"
)

// AddOrder - add order to database
func AddOrder(database *mgo.Database, order orders.Order) error {
	c := database.C("orders")

	err := c.Insert(order)

	if err != nil {
		return err
	}

	return nil
}

// UpdateOrder - updates specified order
func UpdateOrder(database *mgo.Database, order orders.Order, update orders.Order) error {
	c := database.C("orders")

	if order.Filled != true {
		err := c.Update(order, update)

		if err != nil {
			return err
		}

		return nil
	}
	return errors.New("order already filled")
}

// CancelOrder - cancel specified order, returns error
func CancelOrder(database *mgo.Database, order orders.Order) error {
	c := database.C("orders")

	if order.Filled != true {
		err := c.Remove(order)

		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("order already filled")
}
