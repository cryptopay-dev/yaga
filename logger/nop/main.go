package nop

import (
	"io"

	"github.com/cryptopay-dev/yaga/logger"
	"github.com/labstack/gommon/log"
)

type Logger struct{}

func New() logger.Logger {
	return new(Logger)
}

func (l *Logger) Output() io.Writer                                       { return logger.Null }
func (l *Logger) SetOutput(w io.Writer)                                   {}
func (l *Logger) Prefix() string                                          { return "" }
func (l *Logger) SetPrefix(p string)                                      {}
func (l *Logger) Level() log.Lvl                                          { return log.Level() }
func (l *Logger) SetLevel(v log.Lvl)                                      { log.SetLevel(v) }
func (l *Logger) Print(i ...interface{})                                  {}
func (l *Logger) Printf(format string, args ...interface{})               {}
func (l *Logger) Printj(j log.JSON)                                       {}
func (l *Logger) Debug(i ...interface{})                                  {}
func (l *Logger) Debugf(format string, args ...interface{})               {}
func (l *Logger) Debugj(j log.JSON)                                       {}
func (l *Logger) Debugw(message string, args ...interface{})              {}
func (l *Logger) Info(i ...interface{})                                   {}
func (l *Logger) Infof(format string, args ...interface{})                {}
func (l *Logger) Infoj(j log.JSON)                                        {}
func (l *Logger) Infow(message string, args ...interface{})               {}
func (l *Logger) Warn(i ...interface{})                                   {}
func (l *Logger) Warnf(format string, args ...interface{})                {}
func (l *Logger) Warnj(j log.JSON)                                        {}
func (l *Logger) Warnw(message string, args ...interface{})               {}
func (l *Logger) Error(i ...interface{})                                  {}
func (l *Logger) Errorf(format string, args ...interface{})               {}
func (l *Logger) Errorj(j log.JSON)                                       {}
func (l *Logger) Errorw(message string, args ...interface{})              {}
func (l *Logger) Fatal(i ...interface{})                                  {}
func (l *Logger) Fatalf(format string, args ...interface{})               {}
func (l *Logger) Fatalj(j log.JSON)                                       {}
func (l *Logger) Fatalw(message string, args ...interface{})              {}
func (l *Logger) Panic(i ...interface{})                                  {}
func (l *Logger) Panicf(format string, args ...interface{})               {}
func (l *Logger) Panicj(j log.JSON)                                       {}
func (l *Logger) Panicw(message string, args ...interface{})              {}
func (l *Logger) WithContext(fields map[string]interface{}) logger.Logger { return l }
func (l *Logger) Named(name string) logger.Logger                         { return l }
