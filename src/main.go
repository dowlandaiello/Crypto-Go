package main

import (
	"fmt"
	"strings"

	"github.com/mitsukomegumi/Crypto-Go/src/accounts"
	"github.com/mitsukomegumi/Crypto-Go/src/wallets"

	"github.com/mitsukomegumi/Crypto-Go/src/database"
)

func main() {
	db, err := database.ReadDatabase("127.0.0.1")

	if err != nil {
		panic(err)
	}

	acc, err := database.FindAccount(db, "test")

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			pub, _, err := wallets.NewWallets()

			if err != nil {
				panic(err)
			}

			tempAccount := accounts.NewAccount("test", "test@test.com", "test", pub)
			acc = &tempAccount
			database.AddAccount(db, acc)
		}
	}

	update := accounts.NewAccount("test", "test@test.com", "mongo is amazing", acc.WalletAddresses)

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
		- How would wallet private keys be stored?
*/
