package node

import (
	"encoding/json"
	"time"

	"github.com/exgolang/go-nano/account"
	"github.com/exgolang/go-nano/chain/block"
	"github.com/exgolang/go-nano/chain/gen"
	"github.com/exgolang/go-nano/types"

	"github.com/exgolang/go-nano/common"
	"github.com/syndtr/goleveldb/leveldb"

	log "github.com/sirupsen/logrus"
)

// Components struct.
type Components struct {
	Db  *leveldb.DB
	Cmd *types.Cmd
}

func Master(db *leveldb.DB, cmd *types.Cmd) *Components {

	if cmd.Genesis.New {
		gen.New(cmd.Genesis.Coint, db)
	}

	return &Components{
		Db:  db,
		Cmd: cmd,
	}
}

func (n *Components) Running() {

	mnemonic := "rug state yellow climb soul dry unique derive fish reason humor runway pluck rather sight trust soap flower wait toy reform envelope upset street"

	acc, err := account.Master(mnemonic, "01010101")
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Private: ", acc.Private())
	log.Print("Public: ", acc.Public())
	log.Print("Address: ", acc.Address())
	log.Warn("[+++++++++++++++++++++++++++++++++++++++++++++++++]")

	/************************************************/

	privateECDSA, err := common.HexToPrivate(acc.Private())
	if err != nil {
		log.Fatal(err)
	}

	publicECDSA, err := common.HexToPublic(acc.Public())
	if err != nil {
		log.Fatal(err)
	}

	/************************************************/

	signature, err := common.Sign(privateECDSA, []byte(mnemonic))
	if err != nil {
		log.Fatal(err)
	}

	verifycations, err := common.Verify(publicECDSA, []byte(mnemonic), signature)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Data: ", mnemonic)
	log.Print("Signature: ", signature)
	log.Print("Verifycations: ", verifycations)
	log.Warn("[+++++++++++++++++++++++++++++++++++++++++++++++++]")

	/************************************************/

	var (
		i        int
		j        int
		transfer types.Transfer
	)

	for j = 0; j < 20; i++ {

		blockchain, err := block.Master(n.Db)
		if err != nil {
			log.Fatal(err)
		}
		for i = 0; i < 50; i++ {

			transfer.Root = acc.Public()

			transfer.Input.From = "0x13a3495eb612e314db7f4d69ee1dc7b95a151692"
			transfer.Input.To = "0x93a3495eb612e314db7f4d69ee1dc7b95a151691"
			transfer.Input.Value = 1000000
			transfer.Input.Message = "Hello World"

			trx, err := json.Marshal(transfer.Input)
			if err != nil {
				log.Fatal(err)
			}

			signature, err := common.Sign(privateECDSA, trx)
			if err != nil {
				log.Fatal(err)
			}
			transfer.Signature = signature
			transfer.Timestamp = time.Now().UTC().Unix()

			blockchain.Collect.Transactions = append(blockchain.Collect.Transactions, transfer)

		}

		blockchain.Collect.Fees = 100
		blockchain.Collect.Timestamp = time.Now().UTC().Unix()

		if err = blockchain.Commit(); err != nil {
			log.Fatal(err)
		}

	}

}
