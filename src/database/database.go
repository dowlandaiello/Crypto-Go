package database

import (
	"github.com/mitsukomegumi/FakeCrypto/src/accounts"
	"gopkg.in/mgo.v2"
)

// ReadDatabase - attempt to fetch database from specified address
func ReadDatabase(address string) (*mgo.Database, error) {
	session, err := mgo.Dial(address)
	if err != nil {
		return nil, err
	}

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	c := session.DB("crypto")

	return c, nil
}

// AddAccount - add account to database
func addAccount(database *mgo.Database, account accounts.Account) error {
	c := database.C("accounts")

	err := c.Insert(account)

	if err != nil {
		return err
	}

	return nil
}
