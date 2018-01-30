package cli

import "context"

// Instance abstraction layer above Application
type Instance interface {
	Run() error
	Shutdown(ctx context.Context) error
}
