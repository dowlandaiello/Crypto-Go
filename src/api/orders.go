package api

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/mitsukomegumi/Crypto-Go/src/orders"
	"github.com/mitsukomegumi/Crypto-Go/src/pairs"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// SetupOrderRoutes - setup necessary routes for accout database
func SetupOrderRoutes(router *fasthttprouter.Router, db *mgo.Database) (*fasthttprouter.Router, error) {
	_, pErr := setOrderPosts(router, db)

	if pErr != nil {
		return router, pErr
	}

	_, err := setOrderGets(router, db)

	if err != nil {
		return router, err
	}

	return router, nil
}

func setOrderGets(initRouter *fasthttprouter.Router, db *mgo.Database) (*fasthttprouter.Router, error) {
	req, err := NewRequestServer(":pair/:OrderID", "/api/orders", "GET", db, db, "OrderID")
	if err != nil {
		return nil, err
	}

	router, err := req.AttemptToServeRequestsWithRouter(initRouter)

	if err != nil {
		return nil, err
	}

	return router, nil
}

func setOrderPosts(initRouter *fasthttprouter.Router, db *mgo.Database) (*fasthttprouter.Router, error) {
	postReq, rErr := NewRequestServer("POST", "/api/orders", "POST", nil, db, "/:pair/:ordertype/:orderamount/:username/:pass")

	if rErr != nil {
		return nil, rErr
	}

	router, pErr := postReq.AttemptToServeRequestsWithRouter(initRouter)

	if pErr != nil {
		panic(rErr)
	}

	return router, nil
}

func addOrder(database *mgo.Database, order *orders.Order) error {

	_, err := findOrder(database, order.OrderID, order.OrderPair)

	if err != nil {
		c := database.C(order.OrderPair.StartingSymbol + "-" + order.OrderPair.EndingSymbol)

		iErr := c.Insert(order)

		if iErr != nil {
			return iErr
		}

		return nil
	}
	return nil
}

func findOrder(database *mgo.Database, id string, pair pairs.Pair) (*orders.Order, error) {
	c := database.C(pair.StartingSymbol + "-" + pair.EndingSymbol)

	result := orders.Order{}

	err := c.Find(bson.M{"ID": id}).One(&result)
	if err != nil {
		return &result, err
	}

	return &result, nil
}
