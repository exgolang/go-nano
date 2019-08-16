package gen

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/exgolang/go-nano/account"
	"github.com/exgolang/go-nano/chain/block"
	"github.com/exgolang/go-nano/common"
	"github.com/exgolang/go-nano/types"

	log "github.com/sirupsen/logrus"
)

var (
	// Account limit should not exceed more than 10.
	ErrAccountLimit = errors.New("account limit should not exceed more than 10")
)

type GenesisAccounts struct {
	Mnemonic   string `json:"mnemonic"`
	Passphrase string `json:"passphrase"`
	Private    string `json:"private"`
	Public     string `json:"public"`
	Address    string `json:"address"`
}

func Master(coint int) {

	var (
		i             int
		transfer      types.Transfer
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

	blockchain := types.Block{}
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

		accounts = append(accounts, genaccount)
		blockchain.Transactions = append(blockchain.Transactions, transfer)
	}

	blockchain.Timestamp = time.Now().UTC().Unix()
	blockchain.Prev = block.GenesisPrevHash

	bl := block.Components{}
	blockchain.Current, err = bl.Hash(blockchain)
	if err != nil {
		log.Fatal(err)
	}

	_, err = os.Create("genesis.json")
	if err != nil {
		log.Fatal(err)
	}

	rawGenesis, err := json.MarshalIndent(&blockchain, "", "	")
	if err != nil {
		log.Fatal(err)
	}

	genesisJson, err := os.OpenFile("genesis.json", os.O_WRONLY, 0777)
	if err != nil {
		log.Fatal(err)
	}

	_, err = genesisJson.Write(rawGenesis)
	if err != nil {
		log.Fatal(err)
	}

	_, err = os.Create("accounts.json")
	if err != nil {
		log.Fatal(err)
	}

	rawAccounts, err := json.MarshalIndent(&accounts, "", "	")
	if err != nil {
		log.Fatal(err)
	}

	accountsJson, err := os.OpenFile("accounts.json", os.O_WRONLY, 0777)
	if err != nil {
		log.Fatal(err)
	}

	_, err = accountsJson.Write(rawAccounts)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Generated genesis successful!...")
	os.Exit(1)

}
