package log

import (
	"io"
	"os"

	"github.com/cryptopay-dev/yaga/logger"
	"github.com/cryptopay-dev/yaga/logger/nop"
	"github.com/cryptopay-dev/yaga/logger/zap"
	"github.com/labstack/gommon/log"
	zaplog "go.uber.org/zap"
)

var defaultLog logger.Logger

// Init setup logger
func init() {
	level := os.Getenv("LEVEL")

	if level == "nop" {
		defaultLog = nop.New()
	} else if level == "prod" {
		defaultLog = zap.New(zap.Production)
	} else {
		defaultLog = zap.New(zap.Development)
	}

	defaultLog.SetOptions(zaplog.AddCallerSkip(2))
}

// Logger getter
func Logger() logger.Logger { return defaultLog }

// SetOptions applies the supplied Options to Logger
func SetOptions(opts ...logger.Option) {
	defaultLog.SetOptions(opts...)
}

// Named adds a new path segment to the logger's name. Segments are joined by
// periods. By default, Loggers are unnamed.
func Named(name string) logger.Logger { return defaultLog.Named(name) }

// WithContext creates a child logger and adds structured context to it. Fields added
// to the child don't affect the parent, and vice versa.
func WithContext(fields map[string]interface{}) logger.Logger { return defaultLog.WithContext(fields) }

// Output not implemented
func Output() io.Writer { return logger.Null }

// SetOutput not implemented
func SetOutput(w io.Writer) {}

// Prefix not implemented
func Prefix() string { return "" }

// SetPrefix not implemented
func SetPrefix(p string) {}

// Level of logger
func Level() log.Lvl { return defaultLog.Level() }

// SetLevel of logger
func SetLevel(v log.Lvl) { defaultLog.SetLevel(v) }

// Print uses fmt.Sprint to construct and log a message.
func Print(i ...interface{}) { defaultLog.Print(i...) }

// Printf uses fmt.Sprintf to log a templated message.
func Printf(format string, args ...interface{}) { defaultLog.Debugf(format, args...) }

// Printj not implemented
func Printj(j log.JSON) {}

// Debug uses fmt.Sprint to construct and log a message.
func Debug(i ...interface{}) { defaultLog.Debug(i...) }

// Debugf uses fmt.Sprintf to log a templated message.
func Debugf(format string, args ...interface{}) { defaultLog.Debugf(format, args...) }

// Debugj not implemented
func Debugj(j log.JSON) {}

// Debugw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
//
// When debug-level logging is disabled, this is much faster than
//  s.With(keysAndValues).Debug(msg)
func Debugw(message string, args ...interface{}) { defaultLog.Debugw(message, args...) }

// Info uses fmt.Sprint to construct and log a message.
func Info(i ...interface{}) { defaultLog.Info(i...) }

// Infof uses fmt.Sprintf to log a templated message.
func Infof(format string, args ...interface{}) { defaultLog.Infof(format, args...) }

// Infoj not implemented
func Infoj(j log.JSON) {}

// Infow logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Infow(message string, args ...interface{}) { defaultLog.Infow(message, args...) }

// Warn uses fmt.Sprint to construct and log a message.
func Warn(i ...interface{}) { defaultLog.Warn(i...) }

// Warnf uses fmt.Sprintf to log a templated message.
func Warnf(format string, args ...interface{}) { defaultLog.Warnf(format, args...) }

// Warnj not implemented
func Warnj(j log.JSON) {}

// Warnw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Warnw(message string, args ...interface{}) { defaultLog.Warnw(message, args...) }

// Error uses fmt.Sprint to construct and log a message.
func Error(i ...interface{}) { defaultLog.Error(zaplog.Any("error", i)) }

// Errorf uses fmt.Sprintf to log a templated message.
func Errorf(format string, args ...interface{}) { defaultLog.Errorf(format, args...) }

// Errorj not implemented
func Errorj(j log.JSON) {}

// Errorw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Errorw(message string, args ...interface{}) { defaultLog.Errorw(message, args...) }

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func Fatal(i ...interface{}) { defaultLog.Fatal(i...) }

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func Fatalf(format string, args ...interface{}) { defaultLog.Fatalf(format, args...) }

// Fatalj not implemented
func Fatalj(j log.JSON) {}

// Fatalw logs a message with some additional context, then calls os.Exit. The
// variadic key-value pairs are treated as they are in With.
func Fatalw(message string, args ...interface{}) { defaultLog.Fatalw(message, args...) }

// Panic uses fmt.Sprint to construct and log a message, then panics.
func Panic(i ...interface{}) { defaultLog.Panic(i...) }

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func Panicf(format string, args ...interface{}) { defaultLog.Panicf(format, args...) }

// Panicj not implemented
func Panicj(j log.JSON) {}

// Panicw logs a message with some additional context, then panics. The
// variadic key-value pairs are treated as they are in With.
func Panicw(message string, args ...interface{}) { defaultLog.Panicw(message, args...) }
