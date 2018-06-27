package main

import (
	"github.com/mitsukomegumi/FakeCrypto/src/accounts"
	"gopkg.in/mgo.v2"
)

func main() {
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	c := session.DB("crypto").C("accounts")

	err = c.Insert(accounts.NewAccount("test", "test", "test"))
}
