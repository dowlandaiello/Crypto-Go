package api

import (
	"github.com/buaazp/fasthttprouter"
	"gopkg.in/mgo.v2"
)

// SetupMarketRoutes - setup all required routes for market operation
func SetupMarketRoutes(router *fasthttprouter.Router, db *mgo.Database) (*fasthttprouter.Router, error) {
	_, err := setPriceGets(router, db)

	if err != nil {
		return router, err
	}

	return router, nil
}

func setPriceGets(initRouter *fasthttprouter.Router, db *mgo.Database) (*fasthttprouter.Router, error) {
	req, err := NewRequestServer("?pair", "/api/markets/price", "GET", db, db, "?pair")
	if err != nil {
		return nil, err
	}

	_, err = req.AttemptToServeRequestsWithRouter(initRouter)

	if err != nil {
		return nil, err
	}

	return initRouter, nil
}
