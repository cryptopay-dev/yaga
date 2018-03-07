package nop

import (
	"io"

	"github.com/cryptopay-dev/yaga/logger"
	"github.com/labstack/gommon/log"
)

// Logger struct
type Logger struct{}

// New creates new nop logger
func New() logger.Logger {
	return new(Logger)
}

// Output not implemented
func (l *Logger) Output() io.Writer { return nil }

// SetOutput not implemented
func (l *Logger) SetOutput(w io.Writer) {}

// Prefix not implemented
func (l *Logger) Prefix() string { return "" }

// SetPrefix not implemented
func (l *Logger) SetPrefix(p string) {}

// Level of logger
func (l *Logger) Level() log.Lvl { return log.Level() }

// SetLevel of logger
func (l *Logger) SetLevel(v log.Lvl) { log.SetLevel(v) }

// Print not implemented
func (l *Logger) Print(i ...interface{}) {}

// Printf not implemented
func (l *Logger) Printf(format string, args ...interface{}) {}

// Printj not implemented
func (l *Logger) Printj(j log.JSON) {}

// Debug not implemented
func (l *Logger) Debug(i ...interface{}) {}

// Debugf not implemented
func (l *Logger) Debugf(format string, args ...interface{}) {}

// Debugj not implemented
func (l *Logger) Debugj(j log.JSON) {}

// Debugw not implemented
func (l *Logger) Debugw(message string, args ...interface{}) {}

// Info not implemented
func (l *Logger) Info(i ...interface{}) {}

// Infof not implemented
func (l *Logger) Infof(format string, args ...interface{}) {}

// Infoj not implemented
func (l *Logger) Infoj(j log.JSON) {}

// Infow not implemented
func (l *Logger) Infow(message string, args ...interface{}) {}

// Warn not implemented
func (l *Logger) Warn(i ...interface{}) {}

// Warnf not implemented
func (l *Logger) Warnf(format string, args ...interface{}) {}

// Warnj not implemented
func (l *Logger) Warnj(j log.JSON) {}

// Warnw not implemented
func (l *Logger) Warnw(message string, args ...interface{}) {}

// Error not implemented
func (l *Logger) Error(i ...interface{}) {}

// Errorf not implemented
func (l *Logger) Errorf(format string, args ...interface{}) {}

// Errorj not implemented
func (l *Logger) Errorj(j log.JSON) {}

// Errorw not implemented
func (l *Logger) Errorw(message string, args ...interface{}) {}

// Fatal not implemented
func (l *Logger) Fatal(i ...interface{}) {}

// Fatalf not implemented
func (l *Logger) Fatalf(format string, args ...interface{}) {}

// Fatalj not implemented
func (l *Logger) Fatalj(j log.JSON) {}

// Fatalw not implemented
func (l *Logger) Fatalw(message string, args ...interface{}) {}

// Panic not implemented
func (l *Logger) Panic(i ...interface{}) {}

// Panicf not implemented
func (l *Logger) Panicf(format string, args ...interface{}) {}

// Panicj not implemented
func (l *Logger) Panicj(j log.JSON) {}

// Panicw not implemented
func (l *Logger) Panicw(message string, args ...interface{}) {}

// WithContext not implemented
func (l *Logger) WithContext(fields map[string]interface{}) logger.Logger { return l }

// Named not implemented
func (l *Logger) Named(name string) logger.Logger { return l }
