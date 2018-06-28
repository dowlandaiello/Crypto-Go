package database

import (
	"gopkg.in/mgo.v2"
)

// TODO: determine how to upload private key

// AddWallet - append specified public key to wallet db, returning error
func AddWallet(database *mgo.Database, pub string) error {
	c := database.C("wallets")

	err := c.Insert(pub)

	if err != nil {
		return err
	}

	return nil
}

// RemoveWallet - remove specified wallet from wallet db, returning error
func RemoveWallet(database *mgo.Database, pub string) error {
	c := database.C("wallets")

	err := c.Remove(pub)

	if err != nil {
		return err
	}

	return nil
}
