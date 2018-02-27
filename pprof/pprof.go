package pprof

import (
	"net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/cryptopay-dev/yaga/logger"
	"github.com/cryptopay-dev/yaga/web"
)

const (
	pprofPortEnv    = "PPROF_PORT"
	tplInfoPprof    = "Pprof start on port: %s"
	errNilWebEngine = "web.Engine is nil, can't add pprof"
)

// Wrap adds several routes from package `net/http/pprof` to *echo.Echo object.
func Wrap(logger logger.Logger, e *web.Engine) {
	port := os.Getenv(pprofPortEnv)
	if len(port) == 0 {
		if e == nil {
			logger.Error(errNilWebEngine)
			return
		}
		WrapGroup("", e.Group("/debug"))
		return
	}

	pprofWeb := web.New(web.Options{})
	WrapGroup("", pprofWeb.Group("/debug"))

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGABRT)
	go func() {
		logger.Infof(tplInfoPprof, port)
		if err := web.StartServer(pprofWeb, port); err != nil {
			logger.Error(err)
			ch <- syscall.SIGABRT
		}
	}()
}

// Wrapper make sure we are backward compatible.
var Wrapper = Wrap

// WrapGroup adds several routes from package `net/http/pprof` to *echo.Group object.
func WrapGroup(prefix string, g *web.Group) {
	routers := []struct {
		Method  string
		Path    string
		Handler web.HandlerFunc
	}{
		{"GET", "/pprof/", IndexHandler()},
		{"GET", "/pprof/heap", HeapHandler()},
		{"GET", "/pprof/goroutine", GoroutineHandler()},
		{"GET", "/pprof/block", BlockHandler()},
		{"GET", "/pprof/threadcreate", ThreadCreateHandler()},
		{"GET", "/pprof/cmdline", CmdlineHandler()},
		{"GET", "/pprof/profile", ProfileHandler()},
		{"GET", "/pprof/symbol", SymbolHandler()},
		{"POST", "/pprof/symbol", SymbolHandler()},
		{"GET", "/pprof/trace", TraceHandler()},
		{"GET", "/pprof/mutex", MutexHandler()},
	}

	for _, r := range routers {
		switch r.Method {
		case "GET":
			g.GET(strings.TrimPrefix(r.Path, prefix), r.Handler)
		case "POST":
			g.POST(strings.TrimPrefix(r.Path, prefix), r.Handler)
		}
	}
}

// IndexHandler will pass the call from /debug/pprof to pprof.
func IndexHandler() web.HandlerFunc {
	return func(ctx web.Context) error {
		pprof.Index(ctx.Response().Writer, ctx.Request())
		return nil
	}
}

// HeapHandler will pass the call from /debug/pprof/heap to pprof.
func HeapHandler() web.HandlerFunc {
	return func(ctx web.Context) error {
		pprof.Handler("heap").ServeHTTP(ctx.Response(), ctx.Request())
		return nil
	}
}

// GoroutineHandler will pass the call from /debug/pprof/goroutine to pprof.
func GoroutineHandler() web.HandlerFunc {
	return func(ctx web.Context) error {
		pprof.Handler("goroutine").ServeHTTP(ctx.Response().Writer, ctx.Request())
		return nil
	}
}

// BlockHandler will pass the call from /debug/pprof/block to pprof.
func BlockHandler() web.HandlerFunc {
	return func(ctx web.Context) error {
		pprof.Handler("block").ServeHTTP(ctx.Response().Writer, ctx.Request())
		return nil
	}
}

// ThreadCreateHandler will pass the call from /debug/pprof/threadcreate to pprof.
func ThreadCreateHandler() web.HandlerFunc {
	return func(ctx web.Context) error {
		pprof.Handler("threadcreate").ServeHTTP(ctx.Response().Writer, ctx.Request())
		return nil
	}
}

// CmdlineHandler will pass the call from /debug/pprof/cmdline to pprof.
func CmdlineHandler() web.HandlerFunc {
	return func(ctx web.Context) error {
		pprof.Cmdline(ctx.Response().Writer, ctx.Request())
		return nil
	}
}

// ProfileHandler will pass the call from /debug/pprof/profile to pprof.
func ProfileHandler() web.HandlerFunc {
	return func(ctx web.Context) error {
		pprof.Profile(ctx.Response().Writer, ctx.Request())
		return nil
	}
}

// SymbolHandler will pass the call from /debug/pprof/symbol to pprof.
func SymbolHandler() web.HandlerFunc {
	return func(ctx web.Context) error {
		pprof.Symbol(ctx.Response().Writer, ctx.Request())
		return nil
	}
}

// TraceHandler will pass the call from /debug/pprof/trace to pprof.
func TraceHandler() web.HandlerFunc {
	return func(ctx web.Context) error {
		pprof.Trace(ctx.Response().Writer, ctx.Request())
		return nil
	}
}

// MutexHandler will pass the call from /debug/pprof/mutex to pprof.
func MutexHandler() web.HandlerFunc {
	return func(ctx web.Context) error {
		pprof.Handler("mutex").ServeHTTP(ctx.Response().Writer, ctx.Request())
		return nil
	}
}
