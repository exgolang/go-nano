package common

import (
	"crypto/ecdsa"
	"crypto/md5"
	"crypto/rand"
	"encoding/asn1"
	"encoding/hex"
	"hash"
	"math/big"
)

// Signature struct.
type Signature struct {
	R, S *big.Int
}

// Create signature.
func Sign(private *ecdsa.PrivateKey, data []byte) (string, error) {

	var h hash.Hash
	h = md5.New()
	r := big.NewInt(0)
	s := big.NewInt(0)

	// Write expects bytes. If you have a string `s`, use []byte(s) to coerce it to bytes.
	h.Write(data)

	r, s, err := ecdsa.Sign(rand.Reader, private, h.Sum(nil))
	if err != nil {
		return "", err
	}

	signature, err := asn1.Marshal(Signature{
		R: r,
		S: s,
	})
	if err != nil {
		return "", err
	}

	return "0x" + hex.EncodeToString(signature), nil

}

// Verifycation signature.
func Verify(public *ecdsa.PublicKey, data []byte, signature string) (bool, error) {

	var esig Signature

	b, err := HexByte(signature)
	if err != nil {
		return false, err
	}

	var h hash.Hash
	h = md5.New()
	if _, err = asn1.Unmarshal(b, &esig); err != nil {
		return false, err
	}

	// Write expects bytes. If you have a string `s`, use []byte(s) to coerce it to bytes.
	h.Write(data)

	return ecdsa.Verify(public, h.Sum(nil), esig.R, esig.S), nil
}
