package main

import (
	"github.com/mitsukomegumi/Crypto-Go/src/accounts"
	"github.com/mitsukomegumi/Crypto-Go/src/api"
	"github.com/mitsukomegumi/Crypto-Go/src/database"
	"github.com/mitsukomegumi/Crypto-Go/src/wallets"
)

func main() {
	db, err := database.ReadDatabase("127.0.0.1")

	pub, _, _ := wallets.NewWallets()

	acc := accounts.NewAccount("mitsukom", "mitsukomegumii@gmail.com", "dnalwod080304", pub)

	database.AddAccount(db, &acc)

	if err != nil {
		panic(err)
	}

	api.SetupAccountRoutes(db)
}

/*
	FINDINGS:
		Correct configuration: POST http://localhost:8080/api/accounts/POST username:test
	TODO:
		- Store private keys
*/
