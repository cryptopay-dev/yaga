package cli

import "context"

type Instance interface {
	Run() error
	Shutdown(ctx context.Context) error
}
