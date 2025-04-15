package log

import (
	"os"

	"github.com/jorenkoyen/go-logger"
)

// log is the global default logger when no logger instance was defined.
var log = logger.NewWithOptions(logger.Options{Name: "", Writer: os.Stderr})

// CreateDefaultLogger creates a new default logger based on the specified options.
func CreateDefaultLogger(options logger.Options) {
	_logger := logger.NewWithOptions(options)
	SetDefaultLogger(_logger)
}

// SetDefaultLogger will override the default logger for the log package.
func SetDefaultLogger(l *logger.Logger) {
	l.Pos = 3 // we are nested one level deeper in this package
	log = l
}

// WithName clones the default logger but changes the name of the logger.
func WithName(name string) *logger.Logger {
	return log.WithName(name)
}

// Panic is just like Fatal except that it is followed by a call to panic.
func Panic(message string) {
	log.Panic(message)
}

// Panicf is just like Fatalf except that it is followed by a call to panic.
func Panicf(format string, a ...interface{}) {
	log.Panicf(format, a...)
}

// Fatal logs a message at a Fatal Level.
func Fatal(message string) {
	log.Fatal(message)
}

// Fatalf logs a message at Fatal level.
func Fatalf(format string, a ...interface{}) {
	log.Fatalf(format, a...)
}

// Error logs a message at Error level.
func Error(message string) {
	log.Error(message)
}

// Errorf logs a message at Error level.
func Errorf(format string, a ...interface{}) {
	log.Errorf(format, a...)
}

// Warning logs a message at Warning level
func Warning(message string) {
	log.Warning(message)
}

// Warningf logs a message at Warning level.
func Warningf(format string, a ...interface{}) {
	log.Warningf(format, a...)
}

// Info logs a message at Info level.
func Info(message string) {
	log.Info(message)
}

// Infof logs a message at Info level.
func Infof(format string, a ...interface{}) {
	log.Infof(format, a...)
}

// Debug logs a message at Debug level.
func Debug(message string) {
	log.Debug(message)
}

// Debugf logs a message at Debug level.
func Debugf(format string, a ...interface{}) {
	log.Debugf(format, a...)
}

// Trace logs a message at Debug level.
func Trace(message string) {
	log.Trace(message)
}

// Tracef logs a message at Debug level.
func Tracef(format string, a ...interface{}) {
	log.Tracef(format, a...)
}
