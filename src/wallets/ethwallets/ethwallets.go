package ethwallets

import (
	"math/big"

	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	cryptocommon "github.com/dowlandaiello/Crypto-Go/src/common"
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

// SendCoins - sends coins from specified Ether address to specified destination address
func SendCoins(address string, privatekey string, amount float64) error {
	key, err := crypto.HexToECDSA(privatekey)

	if err != nil {
		return err
	}

	tx := types.NewTransaction(uint64(0), common.BytesToAddress([]byte(address)), big.NewInt(int64(amount*1000000000000000000)), 38, big.NewInt(21000), []byte(cryptocommon.Tag))
	tx, _ = types.SignTx(tx, types.HomesteadSigner{}, key)

	client, err := ethclient.Dial("http://localhost:8545")

	context := context.Background()

	err = client.SendTransaction(context, tx)

	if err != nil {
		return err
	}

	return nil
}
