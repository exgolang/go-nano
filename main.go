package main

import (
	"runtime"

	"github.com/exgolang/go-nano/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	go cmd.Master()

	<-make(chan bool)
}
