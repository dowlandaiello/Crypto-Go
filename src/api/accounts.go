package api

import (
	"github.com/mitsukomegumi/Crypto-Go/src/accounts"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// SetupAccountRoutes - setup necessary routes for accout database
func SetupAccountRoutes(db *mgo.Database) {
	setGets(db)
	setPosts(db)
}

func setGets(db *mgo.Database) {
	req, err := NewRequestServer(":username", "/api/accounts", "GET", db, "username")

	if err != nil {
		panic(err)
	}

	err = req.AttemptToServeRequests()

	if err != nil {
		panic(err)
	}
}

func setPosts(db *mgo.Database) {
	/*
		postReq, rErr := api.NewRequestServer("POST", "/api/accounts", "POST", *nAcc)

		if rErr != nil {
			panic(err)
		}

		rErr = req.AttemptToServeRequests()

		if rErr != nil {
			panic(rErr)
		}

		pErr := postReq.AttemptToServeRequests()

		if pErr != nil {
			panic(pErr)
		}
	*/
}

func findAccount(database *mgo.Database, username string) (*accounts.Account, error) {
	c := database.C("accounts")

	result := accounts.Account{}

	err := c.Find(bson.M{"username": username}).One(&result)
	if err != nil {
		return &result, err
	}

	return &result, nil
}

func findValue(database *mgo.Database, collection string, key string, value string) (interface{}, error) {
	c := database.C(collection)

	var result interface{}

	err := c.Find(bson.M{key: value}).One(&result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}
