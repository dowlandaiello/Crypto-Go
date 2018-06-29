package main

import (
	"github.com/mitsukomegumi/Crypto-Go/src/accounts"
	"github.com/mitsukomegumi/Crypto-Go/src/api"
	"github.com/mitsukomegumi/Crypto-Go/src/database"
	"github.com/mitsukomegumi/Crypto-Go/src/wallets"
)

func main() {
	db, err := database.ReadDatabase("127.0.0.1")

	if err != nil {
		panic(err)
	}

	pub, _, _ := wallets.NewWallets()

	fAcc, _ := database.FindAccount(db, "mitsukom")
	database.RemoveAccount(db, fAcc)

	acc := accounts.NewAccount("mitsukom", "mitsukomegumii@gmail.com", "dnalwod080304", pub)

	database.AddAccount(db, &acc)

	api.SetupAccountRoutes(db)
}

/*
	Questions to ask:
		- Should orders be stored in the account struct?
		- How would wallet private keys be stored?
*/
