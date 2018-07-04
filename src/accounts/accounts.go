package accounts

import "github.com/mitsukomegumi/Crypto-Go/src/common"

// Account - exchange account
type Account struct {
	Balance float64 `json:"balance"`

	Username string `json:"username"`
	Email    string `json:"email"`
	PassHash string `json:"passwordhash"`

	WalletAddresses []string `json:"walletaddresses"`
	WalletBalances  []float64
}

// NewAccount - create, return new account
func NewAccount(username string, email string, pass string, walletaddrs []string) Account {
	pass = common.HashAndSalt([]byte(pass))
	rAccount := Account{Username: username, Email: email, PassHash: pass, WalletAddresses: walletaddrs}
	return rAccount
}
