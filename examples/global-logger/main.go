package main

import (
	"github.com/cryptopay-dev/yaga/config"
	"github.com/cryptopay-dev/yaga/logger/log"
)

func main() {
	if err := config.ReadInConfig(); err != nil {
		panic(err)
	}

	// Setup logger..
	log.New()

	log.Info("all fine!")
	log.Print("all fine!")
	log.Debug("all fine!")
	log.Warn("all fine!")
	log.Error("all fine!")
	log.Panic("all fine!")
}
