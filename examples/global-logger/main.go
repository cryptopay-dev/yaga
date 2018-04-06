package main

import (
	"os"

	"github.com/cryptopay-dev/yaga/logger/log"
)

func main() {
	if err := os.Setenv("LEVEL", "dev"); err != nil {
		panic(err)
	}

	log.Info("all fine!")
	log.Print("all fine!")
	log.Debug("all fine!")
	log.Warn("all fine!")
	log.Error("all fine!")
	log.Panic("all fine!")
}
