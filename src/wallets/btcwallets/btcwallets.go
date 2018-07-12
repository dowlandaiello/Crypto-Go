package btcwallets

import (
	"github.com/mitsukomegumi/Crypto-Go/src/common"
)

// NewWallet - generate pub, private keys for new wallet
func NewWallet() (string, []byte, error) {
	priv, err := common.CreateWIF("bitcoin")

	if err != nil {
		return "", []byte{}, err
	}

	pub, err := common.GetAddress("bitcoin", priv)

	if err != nil {
		return "", []byte{}, err
	}

	return pub.EncodeAddress(), priv.SerializePubKey(), nil
}
