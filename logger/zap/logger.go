package zap

import (
	"io"
	"strconv"
	"time"

	"github.com/cryptopay-dev/yaga/logger"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	Development             = "development"
	Production              = "production"
	startLoggerWithLevelTpl = "Start logger with '%s' level"
)

type Logger struct {
	core   zapcore.Core
	logger *zap.Logger
	sugar  *zap.SugaredLogger
}

func New(platform string) logger.Logger {
	var l *zap.Logger
	if platform == Development {
		l, _ = zap.NewDevelopment(zap.AddCallerSkip(1))
	} else {
		platform = Production
		l, _ = zap.Config{
			Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
			Encoding:         "json",
			EncoderConfig:    zap.NewProductionEncoderConfig(),
			OutputPaths:      []string{"stderr"},
			ErrorOutputPaths: []string{"stderr"},
		}.Build()
	}

	core := l.Core()
	sugar := l.Sugar()

	sugaredLogger := &Logger{
		core:   core,
		logger: l,
		sugar:  sugar,
	}

	return sugaredLogger
}

func TimeFromStringField(key string, val string) zapcore.Field {
	bt, _ := strconv.ParseInt(val, 10, 64)
	dt := time.Unix(bt, 0)
	return zap.Time(key, dt)
}

func StringField(key, val string) zapcore.Field {
	return zap.String(key, val)
}

func (l *Logger) Named(name string) logger.Logger {
	lg := l.logger.Named(name)
	core := lg.Core()
	sugar := lg.Sugar()
	return &Logger{
		core:   core,
		logger: lg,
		sugar:  sugar,
	}
}

func (l *Logger) WithContext(fields map[string]interface{}) logger.Logger {
	zapFields := make([]zapcore.Field, len(fields))
	i := 0
	for key, field := range fields {
		zapFields[i] = zap.Any(key, field)
		i++
	}
	lg := l.logger.With(zapFields...)
	core := lg.Core()
	sugar := lg.Sugar()
	return &Logger{
		core:   core,
		logger: lg,
		sugar:  sugar,
	}
}

func (l *Logger) Output() io.Writer                          { return logger.Null }
func (l *Logger) SetOutput(w io.Writer)                      {}
func (l *Logger) Prefix() string                             { return "" }
func (l *Logger) SetPrefix(p string)                         {}
func (l *Logger) Level() log.Lvl                             { return log.Level() }
func (l *Logger) SetLevel(v log.Lvl)                         { log.SetLevel(v) }
func (l *Logger) Print(i ...interface{})                     { l.sugar.Debug(i...) }
func (l *Logger) Printf(format string, args ...interface{})  { l.sugar.Debugf(format, args...) }
func (l *Logger) Printj(j log.JSON)                          {}
func (l *Logger) Debug(i ...interface{})                     { l.sugar.Debug(i...) }
func (l *Logger) Debugf(format string, args ...interface{})  { l.sugar.Debugf(format, args...) }
func (l *Logger) Debugj(j log.JSON)                          {}
func (l *Logger) Debugw(message string, args ...interface{}) { l.sugar.Debugw(message, args...) }
func (l *Logger) Info(i ...interface{})                      { l.sugar.Info(i...) }
func (l *Logger) Infof(format string, args ...interface{})   { l.sugar.Infof(format, args...) }
func (l *Logger) Infoj(j log.JSON)                           {}
func (l *Logger) Infow(message string, args ...interface{})  { l.sugar.Infow(message, args...) }
func (l *Logger) Warn(i ...interface{})                      { l.sugar.Warn(i...) }
func (l *Logger) Warnf(format string, args ...interface{})   { l.sugar.Warnf(format, args...) }
func (l *Logger) Warnj(j log.JSON)                           {}
func (l *Logger) Warnw(message string, args ...interface{})  { l.sugar.Warnw(message, args...) }
func (l *Logger) Error(i ...interface{})                     { l.sugar.Error(zap.Any("error", i)) }
func (l *Logger) Errorf(format string, args ...interface{})  { l.sugar.Errorf(format, args...) }
func (l *Logger) Errorj(j log.JSON)                          {}
func (l *Logger) Errorw(message string, args ...interface{}) { l.sugar.Errorw(message, args...) }
func (l *Logger) Fatal(i ...interface{})                     { l.sugar.Fatal(i...) }
func (l *Logger) Fatalf(format string, args ...interface{})  { l.sugar.Fatalf(format, args...) }
func (l *Logger) Fatalj(j log.JSON)                          {}
func (l *Logger) Fatalw(message string, args ...interface{}) { l.sugar.Fatalw(message, args...) }
func (l *Logger) Panic(i ...interface{})                     { l.sugar.Panic(i...) }
func (l *Logger) Panicf(format string, args ...interface{})  { l.sugar.Panicf(format, args...) }
func (l *Logger) Panicj(j log.JSON)                          {}
func (l *Logger) Panicw(message string, args ...interface{}) { l.sugar.Panicw(message, args...) }
