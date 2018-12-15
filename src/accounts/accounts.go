package accounts

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dowlandaiello/Crypto-Go/src/common"
	"github.com/dowlandaiello/Crypto-Go/src/wallets"
	"github.com/dowlandaiello/Crypto-Go/src/wallets/ethwallets"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Account - exchange account
type Account struct {
	Balance float64 `json:"balance"`

	Username string `json:"username"`
	Email    string `json:"email"`
	PassHash string `json:"passwordhash"`

	Orders []string

	WalletAddresses     []string `json:"walletaddresses"`
	WalletBalances      []float64
	WalletHashedKeys    []string `json:"hashedkeys"`
	WalletRawHashedKeys [][]byte `json:"rawkeys"`
}

// NewAccount - create, return new account
func NewAccount(username string, email string, password string) Account {
	pub, priv, _ := wallets.NewWallets()
	encrypted, err := encryptPrivateKeys(priv, password)

	if err != nil {
		return Account{}
	}

	password = common.HashAndSalt([]byte(password))
	rAccount := Account{Balance: 0, Username: username, Email: email, PassHash: password, WalletAddresses: pub, WalletBalances: []float64{float64(0), float64(0), float64(0)}, WalletRawHashedKeys: encrypted, WalletHashedKeys: common.HashSlice(encrypted)}
	return rAccount
}

// Deposit - wait for deposit into account
func (acc Account) Deposit(symbol string, db *mgo.Database) error {
	if common.StringInSlice(symbol, common.AvailableSymbols) {
		received := false

		startTime := time.Now()

		prevBalance := acc.Balance

		for received != true {
			if time.Since(startTime) > common.TxTimeout*time.Minute {
				break
			}

			balance, err := acc.checkBalance(symbol)

			if err != nil && !strings.Contains(err.Error(), "invalid character") {
				return err
			}

			if balance >= prevBalance {
				acc.WalletBalances[common.IndexInSlice(strings.ToUpper(symbol), []string{"BTC", "LTC", "ETH"})] = balance

				fAcc, err := findAccount(db, acc.Username)

				if err != nil {
					break
				}

				updateAccount(db, fAcc, &acc)

				received = true
			}
		}

		if received != true {
			return errors.New("tx timed out")
		}
		return nil
	}
	return errors.New("invalid symbol")
}

func findAccount(database *mgo.Database, username string) (Account, error) {
	c := database.C("accounts")

	result := Account{}

	err := c.Find(bson.M{"username": username}).One(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (acc *Account) checkBalance(symbol string) (float64, error) {
	if common.StringInSlice(symbol, common.AvailableSymbols) {
		if strings.ToLower(symbol) == "btc" {
			return acc.handleBtc()
		} else if strings.ToLower(symbol) == "ltc" {
			return acc.handleLtc()
		} else if strings.ToLower(symbol) == "eth" {
			return acc.handleEth()
		}
	}
	return 0, errors.New("invalid symbol")
}

func (acc *Account) handleBtc() (float64, error) {
	response, err := http.Get("https://blockchain.info/balance?active=" + acc.WalletAddresses[0])

	if err != nil {
		return float64(0), err
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return float64(0), err
	}

	formatted := common.BlockchainRequest{}

	err = json.Unmarshal(contents, &formatted)
	if err != nil {
		return float64(0), err
	}

	return formatted.Balance / 100000000, nil
}

func (acc *Account) handleEth() (float64, error) {
	response, err := http.Get("https://api.etherscan.io/api?module=account&action=balance&address=" + acc.WalletAddresses[2] + "&tag=latest&apikey=" + common.EtherscanToken)
	if err != nil {
		return float64(0), err
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return float64(0), err
	}

	formatted := common.EtherscanRequest{}
	err = json.Unmarshal(contents, &formatted)

	if err != nil {
		return float64(0), err
	}

	val, err := strconv.ParseFloat(formatted.Result, 64)

	if err != nil {
		return float64(0), err
	}

	return val / 1000000000000000000, nil
}

func (acc *Account) handleLtc() (float64, error) {
	response, err := http.Get("http://api.blockcypher.com/v1/ltc/main/addrs/" + acc.WalletAddresses[1] + "/balance")
	if err != nil {
		return float64(0), err
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return float64(0), err
	}

	formatted := common.EtherscanRequest{}
	err = json.Unmarshal(contents, &formatted)

	if err != nil {
		return float64(0), err
	}

	val, err := strconv.ParseFloat(formatted.Result, 64)

	if err != nil {
		return float64(0), err
	}

	return val / 100000000, nil
}

// MoveAssets - moves specified amount of coins from one wallet to another
func (acc Account) MoveAssets(symbol string, sending string, destination string, privatekey string, amount float64) error {
	if symbol == "ETH" {
		err := ethwallets.SendCoins(sending, privatekey, amount)

		if err != nil {
			return err
		}
	}
	return nil
}

// DecryptPrivateKeys - decrypts private keys
func DecryptPrivateKeys(encryptedKeys [][]byte, key string) ([]string, error) {
	decrypted := []string{}

	x := 0

	for x != len(encryptedKeys) {
		singleDecrypted, err := common.Decrypt(common.BytesToKey([]byte(key)), encryptedKeys[x])

		if err != nil {
			return []string{}, err
		}

		decrypted = append(decrypted, hex.EncodeToString(singleDecrypted))
		x++
	}

	return decrypted, nil
}

func encryptPrivateKeys(privatekeys [][]byte, key string) ([][]byte, error) {
	encrypted := [][]byte{}

	x := 0

	for x != len(privatekeys) {
		singleEncrypted, err := common.Encrypt(common.BytesToKey([]byte(key)), privatekeys[x])

		if err != nil {
			return [][]byte{}, err
		}

		encrypted = append(encrypted, singleEncrypted)
		x++
	}

	return encrypted, nil
}

func updateAccount(database *mgo.Database, account Account, update *Account) error {
	c := database.C("accounts")

	err := c.Update(account, update)

	if err != nil {
		return err
	}

	return nil
}
