package accounts

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mitsukomegumi/Crypto-Go/src/common"
	"github.com/mitsukomegumi/Crypto-Go/src/wallets"
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

	WalletAddresses  []string `json:"walletaddresses"`
	WalletBalances   []float64
	WalletHashedKeys []string `json:"hashedkeys"`
}

// NewAccount - create, return new account
func NewAccount(username string, email string, pass string) Account {
	pub, priv, _ := wallets.NewWallets()
	encrypted := encryptPrivateKeys(priv, pass)
	pass = common.HashAndSalt([]byte(pass))
	rAccount := Account{Balance: 0, Username: username, Email: email, PassHash: pass, WalletAddresses: pub, WalletBalances: []float64{float64(0), float64(0), float64(0)}, WalletHashedKeys: encrypted}
	return rAccount
}

// Deposit - wait for deposit into account
func (acc Account) Deposit(symbol string, db *mgo.Database) error {
	if common.StringInSlice(symbol, common.AvailableSymbols) {
		received := false

		startTime := time.Now()

		prevBalance := acc.Balance

		for received != true {

			if time.Since(startTime) > common.TxTimeout*time.Second {
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
			response, err := http.Get("https://blockchain.info/balance?active=" + acc.WalletAddresses[2])
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
		} else if strings.ToLower(symbol) == "ltc" {
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
		} else if strings.ToLower(symbol) == "eth" {
			response, err := http.Get("https://api.etherscan.io/api?module=account&action=balance&address=" + acc.WalletAddresses[0] + "&tag=latest&apikey=" + common.EtherscanToken)
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
	}
	return 0, errors.New("invalid symbol")
}

func decryptPrivateKeys(encryptedKeys []string, key string) []string {
	decrypted := []string{}

	x := 0

	for x != len(encryptedKeys)-1 {
		singleDecrypted, _ := common.Decrypt([]byte(encryptedKeys[x]), []byte(key))
		decrypted = append(decrypted, base64.URLEncoding.EncodeToString(singleDecrypted))
		x++
	}

	return decrypted
}

func encryptPrivateKeys(privatekeys []string, key string) []string {
	encrypted := []string{}

	x := 0

	for x != len(privatekeys) {
		singleEncrypted, _ := common.Encrypt([]byte(key), []byte(privatekeys[x]))

		encrypted = append(encrypted, base64.URLEncoding.EncodeToString(singleEncrypted))
		x++
	}

	return encrypted
}

func updateAccount(database *mgo.Database, account Account, update *Account) error {
	c := database.C("accounts")

	err := c.Update(account, update)

	if err != nil {
		return err
	}

	return nil
}
