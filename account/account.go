package account

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"

	"github.com/exgolang/go-nano/common"
	"github.com/mr-tron/base58"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

var (
	// ErrInvalidMnemonic is returned when trying to use a malformed mnemonic.
	ErrInvalidMnemonic = errors.New("invalid mnenomic")
)

// Generate new mnemonic text.
func Create() (string, error) {

	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return "", err
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}

	return mnemonic, nil

}

// Struct account struct.
type Account struct {
	_PrivateKey string
	_PublicKey  string
	_Address    string
}

// Master account.
func Master(mnemonic, passphrase string) (*Account, error) {

	// Is returned when trying to use a malformed mnemonic.
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, ErrInvalidMnemonic
	}

	// Generate a Bip32 HD wallet for the mnemonic and a user supplied passphrase.
	seed := bip39.NewSeed(mnemonic, passphrase)
	master, _ := bip32.NewMasterKey(seed)

	// m/44'
	key, err := master.NewChildKey(2147483648 + 44)
	if err != nil {
		return nil, err
	}

	// Serialize base58 decode.
	serialize, err := base58.Decode(key.B58Serialize())
	if err != nil {
		return nil, err
	}

	// Private bytes key.
	private, err := common.ToPrivate(serialize[46:78])
	if err != nil {
		return nil, err
	}

	// Here we start with a new hash.
	h := sha1.New()

	// Write expects bytes. If you have a string `s`, use []byte(s) to coerce it to bytes.
	h.Write([]byte(common.FromPublic(&private.PublicKey)))

	return &Account{
		_PrivateKey: "0x" + hex.EncodeToString(common.FromPrivate(private)),
		_PublicKey:  "0x" + hex.EncodeToString(common.FromPublic(&private.PublicKey)),
		_Address:    "0x" + hex.EncodeToString(h.Sum(nil)),
	}, nil

}

// Params result.
func (p *Account) Address() string { return p._Address }    // Account address.
func (p *Account) Private() string { return p._PrivateKey } // Account private key
func (p *Account) Public() string  { return p._PublicKey }  // Account public key.

func (p *Account) Value() {

}
