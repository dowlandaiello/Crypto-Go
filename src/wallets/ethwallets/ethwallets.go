package ethwallets

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/crypto"
)

// NewWallet - generate pub, private keys for new wallet
func NewWallet() (string, string, error) {
	key, err := crypto.GenerateKey()

	if err != nil {
		return "", "", err
	}

	address := crypto.PubkeyToAddress(key.PublicKey).Hex()
	privateKey := hex.EncodeToString(key.D.Bytes())

	return address, privateKey, nil
}
