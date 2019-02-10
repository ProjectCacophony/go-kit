package discord

import (
	"math/big"
	"strings"
)

// ColorCodeToHex converts a hex string into a Discord Color Code
func HexToColorCode(hex string) int {
	colorInt, ok := new(big.Int).SetString(strings.Replace(hex, "#", "", 1), 16)
	if ok {
		return int(colorInt.Int64())
	}

	return 15957247 // #F37CFF
}

// ColorCodeToHex converts a Discord Color Code into a hex string
func ColorCodeToHex(colour int) (hex string) {
	return strings.ToUpper(big.NewInt(int64(colour)).Text(16))
}
