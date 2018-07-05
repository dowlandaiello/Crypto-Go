package accounts

import (
	"crypto/sha256"
	"encoding/base64"

	"github.com/mitsukomegumi/Crypto-Go/src/common"
)

// Account - exchange account
type Account struct {
	Balance float64 `json:"balance"`

	Username string `json:"username"`
	Email    string `json:"email"`
	PassHash string `json:"passwordhash"`

	WalletAddresses  []string `json:"walletaddresses"`
	WalletBalances   []float64
	WalletHashedKeys [][]byte `json:"hashedkeys"`
}

// NewAccount - create, return new account
func NewAccount(username string, email string, pass string, walletaddrs []string, privatekeys []string) Account {
	encrypted := encryptPrivateKeys(privatekeys, pass)
	pass = common.HashAndSalt([]byte(pass))
	rAccount := Account{Username: username, Email: email, PassHash: pass, WalletAddresses: walletaddrs, WalletHashedKeys: encrypted}
	return rAccount
}

func decryptPrivateKeys(encryptedKeys []byte, key string) []string {
	decrypted := []string{}

	x := 0

	for x != len(encryptedKeys)-1 {
		singleDecrypted, _ := common.Decrypt(encryptedKeys[x], []byte(key))
		decrypted = append(singleDecrypted)
		x++
	}
}

func encryptPrivateKeys(privatekeys []string, key string) []string {
	encrypted := []string{}

	x := 0

	for x != len(privatekeys)-1 {
		singleEncrypted, _ := common.Encrypt([]byte(privatekeys[x]), []byte(key))

		hasher := sha256.New()
		hasher.Write(singleEncrypted)

		encrypted = append(encrypted, base64.URLEncoding.EncodeToString(hasher.Sum(nil)))
		x++
	}

	return encrypted
}
