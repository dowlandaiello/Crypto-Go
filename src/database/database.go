package database

import (
	"gopkg.in/mgo.v2"
)

// ReadDatabase - attempt to fetch database from specified address
func ReadDatabase(address string) (*mgo.Database, error) {
	session, err := mgo.Dial(address)
	if err != nil {
		return nil, err
	}

	session.SetMode(mgo.Monotonic, true)

	c := session.DB("crypto")

	return c, nil
}
