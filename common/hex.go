package common

import (
	"encoding/hex"
	"errors"
	"strings"
)

// Hex to bytes.
func HexByte(hexkey string) ([]byte, error) {

	if strings.HasPrefix(hexkey, "0x") {
		hexkey = strings.Split(hexkey, "x")[1]
	}

	b, err := hex.DecodeString(hexkey)
	if err != nil {
		return nil, errors.New("invalid hex string")
	}

	return b, nil
}
