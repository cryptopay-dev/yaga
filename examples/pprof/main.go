package main

import (
	"github.com/cryptopay-dev/yaga/errors"
	"github.com/cryptopay-dev/yaga/logger/nop"
	"github.com/cryptopay-dev/yaga/pprof"
	"github.com/cryptopay-dev/yaga/web"
)

func main() {
	e := web.New(web.Options{
		Logger: nop.New(),
		Error: errors.Logic{errors.Options{
			Logger: nop.New(),
		}},
	})

	pprof.Wrap(e)

	e.Start(":8080")
}
