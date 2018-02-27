package main

import (
	"github.com/cryptopay-dev/yaga/logger/zap"
	"github.com/cryptopay-dev/yaga/pprof"
	"github.com/cryptopay-dev/yaga/web"
)

func main() {
	e := web.New(web.Options{})
	pprof.Wrap(zap.New(zap.Development), e)
	e.Start(":8080")
}
