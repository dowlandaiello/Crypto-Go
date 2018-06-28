package database

import (
	"github.com/mitsukomegumi/Crypto-Go/src/accounts"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// AddAccount - add account to database
func AddAccount(database *mgo.Database, account *accounts.Account) error {
	c := database.C("accounts")

	err := c.Insert(account)

	if err != nil {
		return err
	}

	return nil
}

// GetAllAccounts - get collection of accounts in db
func GetAllAccounts(database *mgo.Database) (*mgo.Collection, error) {
	c := database.C("accounts")

	return c, nil
}

// RemoveAccount - remove specified account from database
func RemoveAccount(database *mgo.Database, account *accounts.Account) error {
	c := database.C("accounts")

	err := c.Remove(account)

	if err != nil {
		return err
	}

	return nil
}

// UpdateAccount - update account details in database
func UpdateAccount(database *mgo.Database, account accounts.Account, update *accounts.Account) error {
	c := database.C("accounts")

	err := c.Update(account, update)

	if err != nil {
		return err
	}

	return nil
}

// FindAccount - query specified database, return found account (if found)
func FindAccount(database *mgo.Database, username string) (*accounts.Account, error) {
	c := database.C("accounts")

	result := accounts.Account{}

	err := c.Find(bson.M{"username": username}).One(&result)
	if err != nil {
		return &result, err
	}

	return &result, nil
}
