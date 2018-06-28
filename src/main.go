package main

import (
	"fmt"

	"github.com/mitsukomegumi/FakeCrypto/src/accounts"

	"github.com/mitsukomegumi/FakeCrypto/src/database"
)

func main() {
	db, err := database.ReadDatabase("127.0.0.1")

	if err != nil {
		panic(err)
	}

	acc, err := database.FindAccount(db, "test")

<<<<<<< HEAD
	update := accounts.NewAccount("test", "test@test.com", "mongo is amazing")
=======
	update := accounts.NewAccount("test", "test@test.com", "asuydgfuadskgf")
>>>>>>> 488ab90a1e202036a5422b54673fa83fd7e833d4

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
<<<<<<< HEAD
	TODO:
		- Generate random wallet addresses on account creation
=======
>>>>>>> 488ab90a1e202036a5422b54673fa83fd7e833d4
*/
