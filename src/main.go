package main

import (
	"github.com/mitsukomegumi/Crypto-Go/src/api"
	"github.com/mitsukomegumi/Crypto-Go/src/database"
)

func main() {
	db, err := database.ReadDatabase("127.0.0.1")

	if err != nil {
		panic(err)
	}

	api.SetupRoutes(db)
}

/*
	TODO:
		- Handle UPDATE requests
		- Handle onchain wallet deposits
		- On deployment, check for non-testnet
		- Solution to no main wallet, use 0 confs
		- Add API methods to get current price of certain assets
		- Check that an order's amount does not exceed the amount of coins that can be bought at a price
*/
