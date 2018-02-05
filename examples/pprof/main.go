package main

import (
	"github.com/cryptopay-dev/yaga/pprof"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	pprof.Wrap(e)
	e.Start(":8080")
}
