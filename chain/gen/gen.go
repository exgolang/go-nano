package gen

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/exgolang/go-nano/account"
	"github.com/exgolang/go-nano/chain/block"
	"github.com/exgolang/go-nano/common"
	"github.com/exgolang/go-nano/types"
	"github.com/mitchellh/mapstructure"
	"github.com/syndtr/goleveldb/leveldb"

	log "github.com/sirupsen/logrus"
)

var (
	// Account limit should not exceed more than 10.
	ErrAccountLimit = errors.New("account limit should not exceed more than 10")

	// Genesis block is already recorded in the repository.
	ErrAlreadyRecorded = errors.New("genesis block is already recorded in the repository, we recommend deleting the data folder")
)

// Genesis accounts struct.
type GenesisAccounts struct {
	Mnemonic   string `json:"mnemonic"`
	Passphrase string `json:"passphrase"`
	Private    string `json:"private"`
	Public     string `json:"public"`
	Address    string `json:"address"`
}

// Create new genesis and accounts genesis list.
func New(coint int, db *leveldb.DB) {

	var (
		i          int
		blockchain types.Block
		transfer   types.Transfer
	)

	// Check if there is a genesis file.
	if _, err := os.Stat("genesis.json"); os.IsNotExist(err) {

		var (
			passphrase    string
			accountsLimit int
			amountsLimit  = coint
			accounts      []GenesisAccounts
		)

		log.Print("Enter the number of accounts (1-10):")

		count := bufio.NewScanner(os.Stdin)
		count.Scan()

		accountsLimit, err := strconv.Atoi(count.Text())
		if err != nil {
			log.Fatal(err)
		}

		if accountsLimit > 10 {
			log.Fatal(ErrAccountLimit)
		}

		for i = 0; accountsLimit > i; i++ {

			genaccount := GenesisAccounts{}

			// From account.
			fromMnemonic, err := account.Create()
			if err != nil {
				log.Fatal(err)
			}

			hasher := md5.New()
			hasher.Write([]byte(fromMnemonic))
			passphrase = hex.EncodeToString(hasher.Sum(nil))

			fromAccount, err := account.Master(fromMnemonic, passphrase)
			if err != nil {
				log.Fatal(err)
			}

			fromPrivateECDSA, err := common.HexToPrivate(fromAccount.Private())
			if err != nil {
				log.Fatal(err)
			}

			fromPublicECDSA, err := common.HexToPublic(fromAccount.Public())
			if err != nil {
				log.Fatal(err)
			}

			// To account.
			toMnemonic, err := account.Create()
			if err != nil {
				log.Fatal(err)
			}

			toAccount, err := account.Master(toMnemonic, passphrase)
			if err != nil {
				log.Fatal(err)
			}

			// Write new account.
			genaccount.Passphrase = passphrase
			genaccount.Address = toAccount.Address()
			genaccount.Private = toAccount.Private()
			genaccount.Public = toAccount.Public()
			genaccount.Mnemonic = toMnemonic

			// Sender root.
			transfer.Root = fromAccount.Public()

			// Transfer input to account.
			transfer.Input.To = genaccount.Address
			transfer.Input.Value = amountsLimit / accountsLimit
			transfer.Input.Message = "0x" + hex.EncodeToString([]byte("Genesis account #"+strconv.Itoa(i)+" - It was created in Ukraine: "+time.Now().String()))

			// Transfer input to bytes.
			trx, err := json.Marshal(transfer.Input)
			if err != nil {
				log.Fatal(err)
			}

			// Create new signature.
			signature, err := common.Sign(fromPrivateECDSA, trx)
			if err != nil {
				log.Fatal(err)
			}
			transfer.Signature = signature
			transfer.Timestamp = time.Now().UTC().Unix()

			// Verifycations signature.
			verifycations, err := common.Verify(fromPublicECDSA, trx, signature)
			if err != nil {
				log.Fatal(err)
			}

			log.WithFields(log.Fields{
				"signature": signature,
				"verify":    verifycations,
			}).Info("Transaction successful")

			blockchain.Transactions, accounts = append(blockchain.Transactions, transfer), append(accounts, genaccount)
		}

		blockchain.Timestamp, blockchain.Prev = time.Now().UTC().Unix(), block.DefaultPrevHash

		generate := block.Components{}
		blockchain.Current, err = generate.Hash(blockchain, true)
		if err != nil {
			log.Fatal(err)
		}

		// Create new genesis file.
		_, err = os.Create("genesis.json")
		if err != nil {
			log.Fatal(err)
		}

		marshalGenesis, err := json.MarshalIndent(&blockchain, "", "	")
		if err != nil {
			log.Fatal(err)
		}

		// Open genesis file.
		openGenesis, err := os.OpenFile("genesis.json", os.O_WRONLY, 0777)
		if err != nil {
			log.Fatal(err)
		}

		// Write to genesis file.
		_, err = openGenesis.Write(marshalGenesis)
		if err != nil {
			log.Fatal(err)
		}

		// Create new accounts file.
		_, err = os.Create("accounts.json")
		if err != nil {
			log.Fatal(err)
		}

		marshalAccounts, err := json.MarshalIndent(&accounts, "", "	")
		if err != nil {
			log.Fatal(err)
		}

		// Open accounts file.
		openAccounts, err := os.OpenFile("accounts.json", os.O_WRONLY, 0777)
		if err != nil {
			log.Fatal(err)
		}

		// Write to accounts file.
		_, err = openAccounts.Write(marshalAccounts)
		if err != nil {
			log.Fatal(err)
		}

		log.Info("Generated genesis successful!.")
	} else {

		serialize, err := ioutil.ReadFile("genesis.json")
		if err != nil {
			log.Fatal(err)
		}

		// We unmarshal our byte array which contains our
		// genesis block content into 'serialize' which we defined above
		if err = json.Unmarshal(serialize, &blockchain); err != nil {
			log.Fatal(err)
		}

		// Transfer transactions map to struct.
		for _, trx := range blockchain.Transactions {
			if err = mapstructure.Decode(trx, &transfer); err != nil {
				log.Fatal(err)
			}
			blockchain.Transactions = append(blockchain.Transactions, transfer)
		}

	}

	// Insert to store.
	err := Insert(blockchain, db)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(1)

}

// Insert genesis block to store.
func Insert(block types.Block, db *leveldb.DB) error {

	// Prefix key.
	index := append([]byte("block-"), []byte(strconv.Itoa(0))...)

	// Check if there is a genesis block, then we display an error.
	if b, err := db.Has(index, nil); !b {
		if err != nil {
			return err
		}

		// Interface encoded to byte code.
		serialize, err := json.Marshal(block)
		if err != nil {
			return err
		}

		// Put byte code to store.
		if err = db.Put(index, serialize, nil); err != nil {
			return err
		}

		log.Info("Genesis was successfully recorded in the repository!.")

		return nil
	} else {
		return ErrAlreadyRecorded
	}

}
