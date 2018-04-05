package zap

import (
	"io"
	"strconv"
	"time"

	"github.com/cryptopay-dev/yaga/logger"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// Development mode
	Development = "development"
	// Production mode
	Production = "production"
)

// Logger struct
type Logger struct {
	logger *zap.Logger
	sugar  *zap.SugaredLogger
}

// New creates new logger
func New(platform string, callerSkip int) logger.Logger {
	var l *zap.Logger

	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	switch platform {
	case Production:
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case Development:
		config.Encoding = "console"
		config.EncoderConfig = zap.NewDevelopmentEncoderConfig()
	default:
		// not implemented
	}

	l, err := config.Build()
	if err != nil {
		panic(errors.Wrap(err, "can't create logger"))
	}

	l = l.WithOptions(zap.AddCallerSkip(callerSkip))

	sugar := l.Sugar()

	sugaredLogger := &Logger{
		logger: l,
		sugar:  sugar,
	}

	return sugaredLogger
}

// TimeFromStringField creates new zapcore.Field
func TimeFromStringField(key string, val string) zapcore.Field {
	bt, _ := strconv.ParseInt(val, 10, 64)
	dt := time.Unix(bt, 0)
	return zap.Time(key, dt)
}

// StringField creates new zapcore.Field
func StringField(key, val string) zapcore.Field {
	return zap.String(key, val)
}

// Named adds a new path segment to the logger's name. Segments are joined by
// periods. By default, Loggers are unnamed.
func (l *Logger) Named(name string) logger.Logger {
	lg := l.logger.Named(name)
	sugar := lg.Sugar()
	return &Logger{
		logger: lg,
		sugar:  sugar,
	}
}

// WithContext creates a child logger and adds structured context to it. Fields added
// to the child don't affect the parent, and vice versa.
func (l *Logger) WithContext(fields map[string]interface{}) logger.Logger {
	zapFields := make([]zapcore.Field, len(fields))
	i := 0
	for key, field := range fields {
		zapFields[i] = zap.Any(key, field)
		i++
	}
	lg := l.logger.With(zapFields...)
	sugar := lg.Sugar()
	return &Logger{
		logger: lg,
		sugar:  sugar,
	}
}

// Output not implemented
func (l *Logger) Output() io.Writer { return logger.Null }

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

// Print uses fmt.Sprint to construct and log a message.
func (l *Logger) Print(i ...interface{}) { l.sugar.Debug(i...) }

// Printf uses fmt.Sprintf to log a templated message.
func (l *Logger) Printf(format string, args ...interface{}) { l.sugar.Debugf(format, args...) }

// Printj not implemented
func (l *Logger) Printj(j log.JSON) {}

// Debug uses fmt.Sprint to construct and log a message.
func (l *Logger) Debug(i ...interface{}) { l.sugar.Debug(i...) }

// Debugf uses fmt.Sprintf to log a templated message.
func (l *Logger) Debugf(format string, args ...interface{}) { l.sugar.Debugf(format, args...) }

// Debugj not implemented
func (l *Logger) Debugj(j log.JSON) {}

// Debugw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
//
// When debug-level logging is disabled, this is much faster than
//  s.With(keysAndValues).Debug(msg)
func (l *Logger) Debugw(message string, args ...interface{}) { l.sugar.Debugw(message, args...) }

// Info uses fmt.Sprint to construct and log a message.
func (l *Logger) Info(i ...interface{}) { l.sugar.Info(i...) }

// Infof uses fmt.Sprintf to log a templated message.
func (l *Logger) Infof(format string, args ...interface{}) { l.sugar.Infof(format, args...) }

// Infoj not implemented
func (l *Logger) Infoj(j log.JSON) {}

// Infow logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (l *Logger) Infow(message string, args ...interface{}) { l.sugar.Infow(message, args...) }

// Warn uses fmt.Sprint to construct and log a message.
func (l *Logger) Warn(i ...interface{}) { l.sugar.Warn(i...) }

// Warnf uses fmt.Sprintf to log a templated message.
func (l *Logger) Warnf(format string, args ...interface{}) { l.sugar.Warnf(format, args...) }

// Warnj not implemented
func (l *Logger) Warnj(j log.JSON) {}

// Warnw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (l *Logger) Warnw(message string, args ...interface{}) { l.sugar.Warnw(message, args...) }

// Error uses fmt.Sprint to construct and log a message.
func (l *Logger) Error(i ...interface{}) { l.sugar.Error(zap.Any("error", i)) }

// Errorf uses fmt.Sprintf to log a templated message.
func (l *Logger) Errorf(format string, args ...interface{}) { l.sugar.Errorf(format, args...) }

// Errorj not implemented
func (l *Logger) Errorj(j log.JSON) {}

// Errorw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (l *Logger) Errorw(message string, args ...interface{}) { l.sugar.Errorw(message, args...) }

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func (l *Logger) Fatal(i ...interface{}) { l.sugar.Fatal(i...) }

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func (l *Logger) Fatalf(format string, args ...interface{}) { l.sugar.Fatalf(format, args...) }

// Fatalj not implemented
func (l *Logger) Fatalj(j log.JSON) {}

// Fatalw logs a message with some additional context, then calls os.Exit. The
// variadic key-value pairs are treated as they are in With.
func (l *Logger) Fatalw(message string, args ...interface{}) { l.sugar.Fatalw(message, args...) }

// Panic uses fmt.Sprint to construct and log a message, then panics.
func (l *Logger) Panic(i ...interface{}) { l.sugar.Panic(i...) }

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func (l *Logger) Panicf(format string, args ...interface{}) { l.sugar.Panicf(format, args...) }

// Panicj not implemented
func (l *Logger) Panicj(j log.JSON) {}

// Panicw logs a message with some additional context, then panics. The
// variadic key-value pairs are treated as they are in With.
func (l *Logger) Panicw(message string, args ...interface{}) { l.sugar.Panicw(message, args...) }
