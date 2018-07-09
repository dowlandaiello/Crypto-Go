package common

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	cryptorand "crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"math/rand"
	"reflect"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/golangcrypto/bcrypt"
)

// AvailableSymbols - acceptable trading symbols
var AvailableSymbols = []string{"BTC", "LTC", "ETH"}

// AvailableOrderType - acceptable trading order types
var AvailableOrderType = []string{"BUY", "SELL"}

// EtherscanToken - global reference to blockcypher api token
var EtherscanToken = "M91WQ2WASNXW6QP3PBZMG7K9FWHYS9GYU3"

// TxTimeout - global expiry time for tx
const TxTimeout = 10

// FeeRate - global exchange fee
const FeeRate = 0.1

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// BlockchainRequest - request struct for blockchain.com
type BlockchainRequest struct {
	Balance  float64 `json:"final_balance"`
	TxCount  int     `json:"n_tx"`
	Received float64 `json:"total_received"`
}

// EtherscanRequest - request struct for etherscan.io
type EtherscanRequest struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

// BlockcypherRequest - request struct for blockcypher.com
// NOTE: all balances are in satoshis
type BlockcypherRequest struct {
	Address          string `json:"address"`
	Received         int    `json:"total_received"`
	Sent             int    `json:"total_sent"`
	SatBalance       int    `json:"balance"`
	UncomfSatBalance int    `json:"unconfirmed_balance"`
	FinalSatBalance  int    `json:"final_balance"`
	TxCount          int    `json:"n_tx"`
	UncomfTxCount    int    `json:"unconfirmed_n_tx"`
	FinalTxCount     int    `json:"final_n_tx"`
}

// Hash - hash specified interface, return string
func Hash(obj interface{}) (string, error) {
	b, err := GetBytes(obj)

	if err != nil {
		return "", err
	}

	hasher := sha256.New()
	hasher.Write(b)

	return base64.URLEncoding.EncodeToString(hasher.Sum(nil)), nil
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

// IndexInSlice - attempts to retrieve position of item in slice
func IndexInSlice(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}

// CheckSafe - check that value of specified interface is not nil
func CheckSafe(val interface{}) bool {
	if reflect.ValueOf(val).IsNil() {
		return false
	}
	return true
}

// CheckSafeSlice - check that value of specified slice is not nil
func CheckSafeSlice(val interface{}) bool {
	switch reflect.TypeOf(val).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(val)

		if s.Len() == 0 {
			return false
		}
	}
	return true
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

	networkParams.Name = network

	if network == "bitcoin" {
		networkParams.PubKeyHashAddrID = 0x00
		networkParams.PrivateKeyID = 0x80
		networkParams.Net = 0xf9beb4d9
	} else if network == "litecoin" {
		networkParams.PubKeyHashAddrID = 0x30
		networkParams.PrivateKeyID = 0xb0
		networkParams.Net = 0xfbc0b6db
	}

	return networkParams
}

// CreateWIF - creates WIF
func CreateWIF(network string) (*btcutil.WIF, error) {
	secret, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return nil, err
	}

	return btcutil.NewWIF(secret, GetNetworkParams(network), false)
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

// HashSlice - hashes all elements in slice
func HashSlice(s [][]byte) []string {
	hashed := []string{}

	x := 0

	for x != len(s) {
		hashed = append(hashed, hex.EncodeToString(s[x]))
		x++
	}

	return hashed
}

// ComparePasswords - compare specified passwords (hash, actual), to verify correct
func ComparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

// BytesToKey - converts standard byte array to byte array with key length
func BytesToKey(raw []byte) []byte {
	fixed := make([]byte, 32)
	copy(fixed, raw)

	return fixed
}

// Encrypt - encrypt specified byte array with key
func Encrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(cryptorand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return ciphertext, nil
}

// Decrypt - decrypt specified byte array with key
func Decrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return nil, err
	}
	return data, nil
}

// CheckPrice - checks price of asset
func CheckPrice(symbol string) (float64, error) {
	return float64(0), nil
}
