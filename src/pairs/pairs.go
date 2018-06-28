package pairs

import "github.com/mitsukomegumi/Crypto-Go/src/common"

// Pair - trading pair definition
type Pair struct {
	StartingSymbol string `json:"startingsymbol"`
	EndingSymbol   string `json:"endingsymbol"`
}

// NewPair - returns pair, checks if valid
func NewPair(startingSymbol string, endingSymbol string) Pair {
	if startingSymbol != endingSymbol && common.StringInSlice(startingSymbol, common.AvailableSymbols) && common.StringInSlice(endingSymbol, common.AvailableSymbols) {
		return Pair{StartingSymbol: startingSymbol, EndingSymbol: endingSymbol}
	}
	return Pair{}
}
