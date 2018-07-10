package pairs

import (
	"strings"

	"github.com/mitsukomegumi/Crypto-Go/src/common"
)

// Pair - trading pair definition
type Pair struct {
	StartingSymbol string `json:"startingsymbol"`
	EndingSymbol   string `json:"endingsymbol"`
}

// NewPair - returns pair, checks if valid
func NewPair(startingSymbol string, endingSymbol string) Pair {
	startingSymbol = strings.ToUpper(startingSymbol)
	endingSymbol = strings.ToUpper(endingSymbol)

	if startingSymbol != endingSymbol && common.StringInSlice(startingSymbol, common.AvailableSymbols) && common.StringInSlice(endingSymbol, common.AvailableSymbols) && startingSymbol != endingSymbol {
		return Pair{StartingSymbol: startingSymbol, EndingSymbol: endingSymbol}
	}
	return Pair{}
}

// ToString - converts trading pair to string
func (pair Pair) ToString() string {
	return pair.StartingSymbol + "-" + pair.EndingSymbol
}
