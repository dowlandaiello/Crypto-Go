package accounts

import (
	"encoding/base64"

	"github.com/mitsukomegumi/Crypto-Go/src/wallets"

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
	WalletHashedKeys []string `json:"hashedkeys"`
}

// NewAccount - create, return new account
func NewAccount(username string, email string, pass string) Account {
	pub, priv, _ := wallets.NewWallets()
	encrypted := encryptPrivateKeys(priv, pass)
	pass = common.HashAndSalt([]byte(pass))
	rAccount := Account{Username: username, Email: email, PassHash: pass, WalletAddresses: pub, WalletHashedKeys: encrypted}
	return rAccount
}

func decryptPrivateKeys(encryptedKeys []string, key string) []string {
	decrypted := []string{}

	x := 0

	for x != len(encryptedKeys)-1 {
		singleDecrypted, _ := common.Decrypt([]byte(encryptedKeys[x]), []byte(key))
		decrypted = append(decrypted, base64.URLEncoding.EncodeToString(singleDecrypted))
		x++
	}

	return decrypted
}

func encryptPrivateKeys(privatekeys []string, key string) []string {
	encrypted := []string{}

	x := 0

	for x != len(privatekeys) {
		singleEncrypted, _ := common.Encrypt([]byte(key), []byte(privatekeys[x]))

		encrypted = append(encrypted, base64.URLEncoding.EncodeToString(singleEncrypted))
		x++
	}

	return encrypted
}
