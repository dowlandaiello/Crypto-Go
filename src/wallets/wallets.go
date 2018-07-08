package wallets

import (
	"github.com/mitsukomegumi/Crypto-Go/src/wallets/btcwallets"
	"github.com/mitsukomegumi/Crypto-Go/src/wallets/ethwallets"
	"github.com/mitsukomegumi/Crypto-Go/src/wallets/ltcwallets"
)

// NewWallets - generate pub, private keys for all wallet types
func NewWallets() ([]string, []string, error) {
	btcPub, btcPrivate, err := btcwallets.NewWallet()

	if err != nil {
		return nil, nil, err
	}

	ltcPub, ltcPrivate, err := ltcwallets.NewWallet()

	if err != nil {
		return nil, nil, err
	}

	ethPub, ethPrivate, err := ethwallets.NewWallet()

	if err != nil {
		return nil, nil, err
	}

	return []string{btcPub, ltcPub, ethPub}, []string{ethPrivate, ltcPrivate, btcPrivate}, nil
}
