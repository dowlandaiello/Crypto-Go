package wallets

import (
	"github.com/mitsukomegumi/Crypto-Go/src/wallets/btcwallets"
	"github.com/mitsukomegumi/Crypto-Go/src/wallets/ethwallets"
)

// NewWallets - generate pub, private keys for all wallet types
func NewWallets() ([]string, []string, error) {
	ethPub, ethPrivate, err := ethwallets.NewWallet()

	if err != nil {
		return nil, nil, err
	}

	btcPub, btcPrivate, err := btcwallets.NewWallet()

	if err != nil {
		return nil, nil, err
	}

	ltcPub, ltcPrivate, err := btcwallets.NewWallet()

	if err != nil {
		return nil, nil, err
	}

	return []string{ethPub, ltcPub, btcPub}, []string{ethPrivate, ltcPrivate, btcPrivate}, nil
}
