package wallets

import (
	"github.com/dowlandaiello/Crypto-Go/src/wallets/btcwallets"
	"github.com/dowlandaiello/Crypto-Go/src/wallets/ethwallets"
	"github.com/dowlandaiello/Crypto-Go/src/wallets/ltcwallets"
)

// NewWallets - generate pub, private keys for all wallet types
func NewWallets() ([]string, [][]byte, error) {
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

	return []string{btcPub, ltcPub, ethPub}, [][]byte{btcPrivate, ltcPrivate, ethPrivate}, nil
}
