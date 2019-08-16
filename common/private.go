package common

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"errors"
	"fmt"
	"math/big"
)

// Exports a private key into a binary dump.
func FromPrivate(priv *ecdsa.PrivateKey) []byte {
	if priv == nil {
		return nil
	}

	bitSize := priv.Params().BitSize / 8

	if priv.D.BitLen()/8 >= bitSize {
		return priv.D.Bytes()
	}
	ret := make([]byte, bitSize)

	i := len(ret)
	for _, d := range priv.D.Bits() {
		for j := 0; j < 32<<(uint64(^big.Word(0))>>63)/8 && i > 0; j++ {
			i--
			ret[i] = byte(d)
			d >>= 8
		}
	}

	return ret
}

// Private key with the given D value.
func ToPrivate(d []byte) (*ecdsa.PrivateKey, error) {

	secp256k1N, _ := new(big.Int).SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)
	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = elliptic.P256()

	if 8*len(d) != priv.Params().BitSize {
		return nil, fmt.Errorf("invalid length, need %d bits", priv.Params().BitSize)
	}
	priv.D = new(big.Int).SetBytes(d)

	// The priv.D must < N
	if priv.D.Cmp(secp256k1N) >= 0 {
		return nil, errors.New("invalid private key, >=N")
	}

	// The priv.D must not be zero or negative.
	if priv.D.Sign() <= 0 {
		return nil, errors.New("invalid private key, zero or negative")
	}

	priv.PublicKey.X, priv.PublicKey.Y = priv.PublicKey.Curve.ScalarBaseMult(d)
	if priv.PublicKey.X == nil {
		return nil, errors.New("invalid private key")
	}

	return priv, nil

}

// Parses hex a private key.
func HexToPrivate(hexkey string) (*ecdsa.PrivateKey, error) {

	b, err := HexByte(hexkey)
	if err != nil {
		return nil, err
	}

	return ToPrivate(b)

}
