package database

import (
	"github.com/mitsukomegumi/FakeCrypto/src/accounts"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// AddAccount - add account to database
func AddAccount(database *mgo.Database, account accounts.Account) error {
	c := database.C("accounts")

	err := c.Insert(account)

	if err != nil {
		return err
	}

	return nil
}

// UpdateAccount -
func UpdateAccount(database *mgo.Database, account accounts.Account) error {
	return nil
}

// FindAccount - query specified database, return found account (if found)
func FindAccount(database *mgo.Database, account accounts.Account, username string) (accounts.Account, error) {
	c := database.C("accounts")

	result := accounts.Account{}

	err := c.Find(bson.M{"username": username}).One(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}
