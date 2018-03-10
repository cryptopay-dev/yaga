package commands

import (
	"fmt"
	"io"
	"os"

	"github.com/cryptopay-dev/yaga/logger"
	"github.com/labstack/gommon/log"
)

type formatter = func(msg interface{}, styles ...string) string

func output(out io.Writer, format formatter, msg string) {
	clr.SetOutput(out)
	cnt.Add(1)
	clr.Printf("[%03d] %s\n", cnt.Load(), format(msg))
}

// Logger struct
type Logger struct {
	output io.Writer
}

// NewLogger for migrate
func NewLogger() logger.Logger {
	return &Logger{output: os.Stdout}
}

// Output for logger
func (l *Logger) Output() io.Writer { return l.output }

// SetOutput for logger
func (l *Logger) SetOutput(w io.Writer) { l.output = w }

// Prefix for logger
func (l *Logger) Prefix() string { return "" }

// SetPrefix for logger
func (l *Logger) SetPrefix(p string) {}

// Level for logger
func (l *Logger) Level() log.Lvl { return log.Level() }

// SetLevel for logger
func (l *Logger) SetLevel(v log.Lvl) { log.SetLevel(v) }

// Print for logger
func (l *Logger) Print(i ...interface{}) { output(l.output, clr.Blue, fmt.Sprint(i...)) }

// Printf for logger
func (l *Logger) Printf(format string, args ...interface{}) {
	output(l.output, clr.Blue, fmt.Sprintf(format, args...))
}

// Printj for logger
func (l *Logger) Printj(j log.JSON) {}

// Debug for logger
func (l *Logger) Debug(i ...interface{}) {}

// Debugf for logger
func (l *Logger) Debugf(format string, args ...interface{}) {}

// Debugj for logger
func (l *Logger) Debugj(j log.JSON) {}

// Debugw for logger
func (l *Logger) Debugw(message string, args ...interface{}) {}

// Info for logger
func (l *Logger) Info(i ...interface{}) { output(l.output, clr.Green, fmt.Sprint(i...)) }

// Infof for logger
func (l *Logger) Infof(format string, args ...interface{}) {
	output(l.output, clr.Green, fmt.Sprintf(format, args...))
}

// Infoj for logger
func (l *Logger) Infoj(j log.JSON) {}

// Infow for logger
func (l *Logger) Infow(message string, args ...interface{}) {}

// Warn for logger
func (l *Logger) Warn(i ...interface{}) { output(l.output, clr.Yellow, fmt.Sprint(i...)) }

// Warnf for logger
func (l *Logger) Warnf(format string, args ...interface{}) {
	output(l.output, clr.Yellow, fmt.Sprintf(format, args...))
}

// Warnj for logger
func (l *Logger) Warnj(j log.JSON) {}

// Warnw for logger
func (l *Logger) Warnw(message string, args ...interface{}) {}

// Error for logger
func (l *Logger) Error(i ...interface{}) { output(l.output, clr.Red, fmt.Sprint(i...)) }

// Errorf for logger
func (l *Logger) Errorf(format string, args ...interface{}) {
	output(l.output, clr.Red, fmt.Sprintf(format, args...))
}

// Errorj for logger
func (l *Logger) Errorj(j log.JSON) {}

// Errorw for logger
func (l *Logger) Errorw(message string, args ...interface{}) {}

// Fatal for logger
func (l *Logger) Fatal(i ...interface{}) {
	output(l.output, clr.Red, fmt.Sprint(i...))
	os.Exit(1)
}

// Fatalf for logger
func (l *Logger) Fatalf(format string, args ...interface{}) {
	output(l.output, clr.Red, fmt.Sprintf(format, args...))
	os.Exit(1)
}

// Fatalj for logger
func (l *Logger) Fatalj(j log.JSON) {}

// Fatalw for logger
func (l *Logger) Fatalw(message string, args ...interface{}) {}

// Panic for logger
func (l *Logger) Panic(i ...interface{}) {
	output(l.output, clr.Red, fmt.Sprint(i...))
	os.Exit(1)
}

// Panicf for logger
func (l *Logger) Panicf(format string, args ...interface{}) {
	output(l.output, clr.Red, fmt.Sprintf(format, args...))
	os.Exit(1)
}

// Panicj for logger
func (l *Logger) Panicj(j log.JSON) {}

// Panicw for logger
func (l *Logger) Panicw(message string, args ...interface{}) {}

// WithContext for logger
func (l *Logger) WithContext(fields map[string]interface{}) logger.Logger { return l }

// Named for logger
func (l *Logger) Named(name string) logger.Logger { return l }
