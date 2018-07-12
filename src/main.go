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
		- Account for sell orders
		- On deployment, check for non-testnet
		- Add method to clear volume after 24 hours via API
*/
