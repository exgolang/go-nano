package block

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/exgolang/go-nano/types"
	//"github.com/davecgh/go-spew/spew"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"

	log "github.com/sirupsen/logrus"
)

const (
	GenesisPrevHash = "0x0000000000000000000000000000000000000000000000000000000000000000"
)

// Block interface.
type Interface interface {
	Hash(block types.Block) (string, error)
}

var (
	// Block index already exists.
	//ErrIndexAlready = errors.New("block index already exists: block.go")

	// The original hash of the previous block matches the hash of the new block.
	ErrPreviousHash = errors.New("the original hash of the previous block matches the hash of the new block: block.go")

	// Invalid block indexing.
	ErrInvalidIndexing = errors.New("invalid block indexing: block.go")

	// Invalid block hash.
	ErrInvalidBlockHash = errors.New("invalid block hash: block.go")
)

// Components struct.
type Components struct {
	Db      *leveldb.DB
	Collect types.Block
}

// Master create block.
func Master(db *leveldb.DB) (*Components, error) {

	var (
		prev, current types.Block
	)

	iterator := db.NewIterator(util.BytesPrefix([]byte("block-")), nil)
	if iterator.Next() {

		if err := json.Unmarshal(iterator.Value(), &prev); err != nil {
			return nil, err
		}

		current.Index, current.Prev = prev.Index+1, prev.Current
	} else {
		current.Index, current.Prev = 0, GenesisPrevHash
	}
	iterator.Release() // Note: you should first get data and then release iterator

	if err := iterator.Error(); err != nil {
		return nil, err
	}

	return &Components{
		Db:      db,
		Collect: current,
	}, nil

}

// Commit new block.
func (c *Components) Commit() error {

	//spew.Dump(c.Collect.Index)

	//index := append([]byte("block-"), []byte(strconv.Itoa(c.Collect.Index))...)
	//if b, err := c.Db.Has(index, nil); !b {
	//	if err != nil {
	//		return err
	//	}

	if _, err := c.Hash(types.Block{}); err != nil {
		return err
	}

	//if err := c.isValidate(); err != nil {
	//	return err
	//}

	block, err := json.Marshal(c.Collect)
	if err != nil {
		return err
	}

	if err = c.Db.Put(append([]byte("block-"), []byte(strconv.Itoa(c.Collect.Index))...), block, nil); err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"index":     c.Collect.Index,
		"prev":      c.Collect.Prev,
		"current":   c.Collect.Current,
		"trx-count": len(c.Collect.Transactions),
	}).Info("Commit block successful")

	return nil

	//} else {
	//	return ErrIndexAlready
	//}

}

// Validate new block.
func (c *Components) isValidate() error {

	prev, index := new(types.Block), c.Collect.Index
	if index > 0 {

		block, err := c.Db.Get(append([]byte("block-"), []byte(strconv.Itoa(index))...), nil)
		if err != nil {
			return err
		}

		if err = json.Unmarshal(block, prev); err != nil {
			return err
		}

		if prev.Current != c.Collect.Prev {
			return ErrPreviousHash
		}

		if prev.Index+1 != c.Collect.Index {
			return ErrInvalidIndexing
		}

		if current, _ := c.Hash(types.Block{}); c.Collect.Current != current {
			return ErrInvalidBlockHash
		}

	}

	return nil

}

// Hashed block.
func (c *Components) Hash(block types.Block) (string, error) {

	var (
		appends []byte
	)

	current := new(types.Block)
	if block.Index > 0 {
		current = &block
	} else {
		current = &c.Collect
	}

	trx, err := json.Marshal(current.Transactions)
	if err != nil {
		return "", err
	}

	//c.Mutex.Lock()

	appends = append(appends, trx...)
	appends = append(appends, []byte(strconv.Itoa(current.Index))...)
	appends = append(appends, []byte(strconv.FormatInt(current.Timestamp, 10))...)
	appends = append(appends, []byte(strconv.Itoa(current.Fees))...)
	appends = append(appends, []byte(current.Prev)...)

	//c.Mutex.Unlock()

	h := sha256.New()
	h.Write(appends)

	c.Collect.Current = "0x" + hex.EncodeToString(h.Sum(nil))

	return c.Collect.Current, nil

}
