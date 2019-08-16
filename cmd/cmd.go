package cmd

import (
	"flag"

	"github.com/exgolang/go-nano/node"
	"github.com/exgolang/go-nano/types"
	log "github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

func Master() {

	var params types.Cmd

	log.WithFields(log.Fields{
		"name": "Catalina",
	}).Info("Master node")

	flag.StringVar(&params.Network, "network", "mainnet", "Type network Default: (mainnet), or testnet.")
	flag.StringVar(&params.Host, "host", "127.0.0.1", "Master node host Default: (127.0.0.1).")
	flag.IntVar(&params.Port, "port", 9005, "Master node port Default: (9005).")
	flag.BoolVar(&params.Genesis, "new-genesis", false, "Created new genesis.")
	flag.IntVar(&params.Coint, "new-coint", 100000000, "Number of coins issued: (100000000 or 100000000000)")
	flag.Parse()

	db, err := leveldb.OpenFile("./store/"+params.Network, &opt.Options{
		OpenFilesCacheCapacity: 512,
		BlockCacheCapacity:     512 / 2 * opt.MiB,
		WriteBuffer:            512 / 4 * opt.MiB, // Two of these are used internally
		Filter:                 filter.NewBloomFilter(10),
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	node.Master(db, &params).Running()

}
