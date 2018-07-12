package market

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mitsukomegumi/Crypto-Go/src/common"
	"github.com/mitsukomegumi/Crypto-Go/src/pairs"
)

// CheckPrice - checks price of asset
func CheckPrice(tradingpair pairs.Pair) (float64, error) {
	response, err := http.Get("https://min-api.cryptocompare.com/data/price?fsym=" + tradingpair.StartingSymbol + "&tsyms=" + tradingpair.EndingSymbol)

	if err != nil {
		return float64(0), err
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return float64(0), err
	}

	formatted := common.CryptoCompareRequest{}
	err = json.Unmarshal(contents, &formatted)

	if err != nil {
		return float64(0), err
	}

	if formatted.BitcoinPrice != 0 {
		fmt.Println("\ncurrent " + tradingpair.ToString() + " price: " + common.FloatToString(formatted.BitcoinPrice))
		return formatted.BitcoinPrice, nil
	} else if formatted.LitecoinPrice != 0 {
		fmt.Println("\ncurrent " + tradingpair.ToString() + " price: " + common.FloatToString(formatted.LitecoinPrice))
		return formatted.LitecoinPrice, nil
	} else if formatted.EthereumPrice != 0 {
		fmt.Println("\ncurrent " + tradingpair.ToString() + " price: " + common.FloatToString(formatted.EthereumPrice))
		return formatted.EthereumPrice, nil
	}

	return float64(0), errors.New("invalid request")
}
