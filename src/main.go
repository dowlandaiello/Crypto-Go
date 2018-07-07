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
		- solution to no main wallet, use 0 confs
*/
