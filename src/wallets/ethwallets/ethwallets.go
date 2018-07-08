package ethwallets

import (
	"github.com/ethereum/go-ethereum/crypto"
)

// NewWallet - generate pub, private keys for new wallet
func NewWallet() (string, []byte, error) {
	key, err := crypto.GenerateKey()

	if err != nil {
		return "", []byte{}, err
	}

	address := crypto.PubkeyToAddress(key.PublicKey).Hex()
	privateKey := key.D.Bytes()

	return address, privateKey, nil
}
