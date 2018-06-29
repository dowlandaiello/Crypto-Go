package common

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"math/rand"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/golangcrypto/bcrypt"
)

// AvailableSymbols - acceptable trading symbols
var AvailableSymbols = []string{"BTC", "LTC", "ETH"}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Hash - hash specified interface, return string
func Hash(obj interface{}) (string, error) {
	b, err := GetBytes(obj)

	if err != nil {
		return "", err
	}

	h := sha256.Sum256(b)

	return string(h[:]), nil
}

// GetBytes - get bytes of specified interface, return byte array
func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// StringInSlice - checks if specified string is in slice
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// RandStringBytesRmndr - generate random string
func RandStringBytesRmndr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

// TrimLeftChar - trims value of string, removing first character
func TrimLeftChar(s string) string {
	for i := range s {
		if i > 0 {
			return s[i:]
		}
	}
	return s[:0]
}

// GetNetworkParams - returns parameters of selected network
func GetNetworkParams(network string) *chaincfg.Params {
	networkParams := &chaincfg.MainNetParams

	if network == "bitcoin" {
		networkParams.PubKeyHashAddrID = 0x00
		networkParams.PrivateKeyID = 0x80
	} else if network == "litecoin" {
		networkParams.PubKeyHashAddrID = 0x30
		networkParams.PrivateKeyID = 0xb0
	}

	return networkParams
}

// CreateWIF - creates WIF
func CreateWIF(network string) (*btcutil.WIF, error) {
	secret, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return nil, err
	}

	return btcutil.NewWIF(secret, GetNetworkParams(network), true)
}

// GetAddress - get address from specified wif
func GetAddress(network string, wif *btcutil.WIF) (*btcutil.AddressPubKey, error) {
	return btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeCompressed(), GetNetworkParams(network))
}

// HashAndSalt - generate hash for specified byte array
func HashAndSalt(b []byte) string {
	hash, err := bcrypt.GenerateFromPassword(b, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}

	return string(hash)
}
