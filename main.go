package main

import (
	"github.com/exgolang/go-nano/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {

	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	cmd.Master()

	<-make(chan bool)
}
