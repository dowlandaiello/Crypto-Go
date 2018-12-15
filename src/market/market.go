package market

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dowlandaiello/Crypto-Go/src/common"
	"github.com/dowlandaiello/Crypto-Go/src/pairs"
)

// CheckVolume - check volume (float) of specified tradingpair
func CheckVolume(tradingpair pairs.Pair) float64 {
	return tradingpair.Volume
}

// ClearVolume - clear volume of specified trading pair
func ClearVolume(tradingpair pairs.Pair) {
	tradingpair.Volume = 0
}

// CheckPrice - check price of asset pair
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
