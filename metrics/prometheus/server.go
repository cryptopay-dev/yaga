package prometheus

import (
	"context"
	"net/http"

	"github.com/cryptopay-dev/yaga/logger/log"
)

// StartWebServer run http server in separate goroutine
func (p *Provider) StartWebServer() {
	go func() {
		if err := p.server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Errorw("start prometheus http server error", "error", err)
			}
		}
	}()
}

// StopWebServer calling http.Server.Shutdown
func (p *Provider) StopWebServer(ctx context.Context) error {
	return p.server.Shutdown(ctx)
}
