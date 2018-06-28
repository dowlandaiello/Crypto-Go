package main

import (
	"fmt"

	"github.com/mitsukomegumi/Crypto-Go/src/accounts"

	"github.com/mitsukomegumi/Crypto-Go/src/database"
)

func main() {
	db, err := database.ReadDatabase("127.0.0.1")

	if err != nil {
		panic(err)
	}

	acc, err := database.FindAccount(db, "test")

	update := accounts.NewAccount("test", "test@test.com", "mongo is amazing")

	err = database.UpdateAccount(db, *acc, &update)

	nAcc, err := database.FindAccount(db, "test")

	if err != nil {
		panic(err)
	}

	fmt.Println(nAcc)
}

/*
	Questions to ask:
		- Should orders be stored in the account struct?
	TODO:
		- Generate random wallet addresses on account creation
*/
