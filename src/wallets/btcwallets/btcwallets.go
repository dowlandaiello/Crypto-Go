package btcwallets

import (
	"github.com/mitsukomegumi/Crypto-Go/src/common"
)

// NewWallet - generate pub, private keys for new wallet
func NewWallet() (string, string, error) {
	priv, err := common.CreateWIF("bitcoin")

	if err != nil {
		return "", "", err
	}

	pub, err := common.GetAddress("bitcoin", priv)

	return pub.String(), priv.String(), nil
}
