package common

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"errors"
)

// Converts ecdsa public key to bytes.
func FromPublic(pub *ecdsa.PublicKey) []byte {

	if pub == nil || pub.X == nil || pub.Y == nil {
		return nil
	}

	return elliptic.Marshal(elliptic.P256(), pub.X, pub.Y)

}

// Converts bytes to a secp256k1 public key.
func ToPublic(pub []byte) (*ecdsa.PublicKey, error) {

	x, y := elliptic.Unmarshal(elliptic.P256(), pub)
	if x == nil {
		return nil, errors.New("invalid secp256k1 public key")
	}

	return &ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}, nil

}

// Parses hex a public key.
func HexToPublic(hexkey string) (*ecdsa.PublicKey, error) {

	b, err := HexByte(hexkey)
	if err != nil {
		return nil, err
	}

	return ToPublic(b)

}
