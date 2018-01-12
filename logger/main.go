package logger

import (
	"io"

	"github.com/labstack/gommon/log"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Output() io.Writer
	SetOutput(w io.Writer)

	Prefix() string
	SetPrefix(p string)

	Level() log.Lvl
	SetLevel(v log.Lvl)

	WithContext(fields ...zapcore.Field) Logger
	Named(name string) Logger

	Print(i ...interface{})
	Printf(format string, args ...interface{})
	Printj(j log.JSON)

	Debug(i ...interface{})
	Debugf(format string, args ...interface{})
	Debugj(j log.JSON)
	Debugw(message string, args ...interface{})

	Info(i ...interface{})
	Infof(format string, args ...interface{})
	Infoj(j log.JSON)
	Infow(message string, args ...interface{})

	Warn(i ...interface{})
	Warnf(format string, args ...interface{})
	Warnj(j log.JSON)
	Warnw(message string, args ...interface{})

	Error(i ...interface{})
	Errorf(format string, args ...interface{})
	Errorj(j log.JSON)
	Errorw(message string, args ...interface{})

	Fatal(i ...interface{})
	Fatalj(j log.JSON)
	Fatalf(format string, args ...interface{})
	Fatalw(message string, args ...interface{})

	Panic(i ...interface{})
	Panicj(j log.JSON)
	Panicf(format string, args ...interface{})
	Panicw(message string, args ...interface{})
}
