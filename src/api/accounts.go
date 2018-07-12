package api

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/mitsukomegumi/Crypto-Go/src/accounts"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// SetupAccountRoutes - setup necessary routes for accout database
func SetupAccountRoutes(db *mgo.Database) (*fasthttprouter.Router, error) {
	router, pErr := setPosts(db)

	if pErr != nil {
		return router, pErr
	}

	_, err := setGets(router, db)

	if err != nil {
		return router, err
	}

	_, err = setUpdates(router, db)

	if err != nil {
		return router, err
	}

	_, prErr := setProtectedGets(router, db)

	if prErr != nil {
		return router, prErr
	}

	_, dErr := setDeletes(router, db)

	if dErr != nil {
		return router, dErr
	}

	_, err = setDeposits(router, db)

	if err != nil {
		return router, err
	}

	_, err = setGeneralAccountRoutes(router, db)

	if err != nil {
		return router, err
	}

	return router, nil
}

func setGets(initRouter *fasthttprouter.Router, db *mgo.Database) (*fasthttprouter.Router, error) {
	req, err := NewRequestServer("?username", "/api/accounts/user", "GET", db, db, "?username")

	if err != nil {
		return nil, err
	}

	router, err := req.AttemptToServeRequestsWithRouter(initRouter)

	if err != nil {
		return nil, err
	}

	return router, nil
}

func setProtectedGets(initRouter *fasthttprouter.Router, db *mgo.Database) (*fasthttprouter.Router, error) {
	req, err := NewRequestServer("?username?password", "/api/accounts/keys", "GET", db, db, "?username?password")

	if err != nil {
		return nil, err
	}

	router, err := req.AttemptToServeRequestsWithRouter(initRouter)

	if err != nil {
		return nil, err
	}

	return router, nil
}

func setPosts(db *mgo.Database) (*fasthttprouter.Router, error) {
	postReq, rErr := NewRequestServer("?username?email?password", "/api/accounts/create", "POST", nil, db, "?username?email?password")

	if rErr != nil {
		return nil, rErr
	}

	router, pErr := postReq.AttemptToServeRequests()

	if pErr != nil {
		return router, pErr
	}

	return router, nil
}

func setUpdates(initrouter *fasthttprouter.Router, db *mgo.Database) (*fasthttprouter.Router, error) {
	updateReq, err := NewRequestServer("?username?email?password", "/api/accounts/update", "POST", nil, db, "?username?email?password")

	if err != nil {
		return initrouter, err
	}

	_, err = updateReq.AttemptToServeRequestsWithRouter(initrouter)

	return initrouter, nil
}

func setDeletes(initRouter *fasthttprouter.Router, db *mgo.Database) (*fasthttprouter.Router, error) {
	delReq, rErr := NewRequestServer("?username?password", "/api/accounts/remove", "DELETE", nil, db, "?username?password")

	if rErr != nil {
		return nil, rErr
	}

	_, dErr := delReq.AttemptToServeRequestsWithRouter(initRouter)

	if dErr != nil {
		return initRouter, dErr
	}

	return initRouter, nil
}

func setDeposits(initRouter *fasthttprouter.Router, db *mgo.Database) (*fasthttprouter.Router, error) {
	depReq, err := NewRequestServer("?username?symbol", "/api/deposit", "POST", nil, db, "?username?symbol")

	if err != nil {
		return nil, err
	}

	_, dErr := depReq.AttemptToServeRequestsWithRouter(initRouter)

	if dErr != nil {
		return initRouter, dErr
	}

	return initRouter, nil
}

func setGeneralAccountRoutes(initRouter *fasthttprouter.Router, db *mgo.Database) (*fasthttprouter.Router, error) {
	getReq, _ := NewRequestServer("GET", "/api/accounts", "GET", nil, db, "")
	initRouter.GET("/api/accounts", getReq.HandleGETCollection)

	return initRouter, nil
}

func addAccount(database *mgo.Database, account *accounts.Account) error {

	_, err := findAccount(database, account.Username)

	if err != nil {
		c := database.C("accounts")

		iErr := c.Insert(account)

		if iErr != nil {
			return iErr
		}

		return nil
	}
	return nil
}

func updateAccount(database *mgo.Database, account accounts.Account, update *accounts.Account) error {
	c := database.C("accounts")

	err := c.Update(account, update)

	if err != nil {
		return err
	}

	return nil
}

func removeAccount(database *mgo.Database, account *accounts.Account) error {
	c := database.C("accounts")

	err := c.Remove(account)

	if err != nil {
		return err
	}

	return nil
}

func findAccount(database *mgo.Database, username string) (accounts.Account, error) {
	c := database.C("accounts")

	result := accounts.Account{}

	err := c.Find(bson.M{"username": username}).One(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}
