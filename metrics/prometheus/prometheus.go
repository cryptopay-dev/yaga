package prometheus

import (
	"net/http"
	"os"

	"github.com/cryptopay-dev/yaga/logger/log"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//ErrBadPrometheusBind to serve metrics
var ErrBadPrometheusBind = errors.New("check prometheus bind address")

// Serve prometheus metrics
func Serve() error {
	bind := os.Getenv("PROMETHEUS")

	if len(bind) == 0 {
		return ErrBadPrometheusBind
	}

	go func() {
		log.Debugf("run prometheus metrics on `%s`", bind)

		if err := http.ListenAndServe(bind, promhttp.Handler()); err != nil {
			log.Panic(errors.Wrap(err, "can't serve prometheus metrics"))
		}
	}()

	return nil
}
