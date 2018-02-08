package main

import (
	"github.com/cryptopay-dev/yaga/pprof"
	"github.com/cryptopay-dev/yaga/web"
)

func main() {
	e := web.New(web.Options{})
	pprof.Wrap(e)
	e.Start(":8080")
}
