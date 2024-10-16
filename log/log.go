package log

import (
	"os"

	"github.com/jorenkoyen/go-logger"
)

// Logger is the global default logger when no logger instance was defined.
var Logger = logger.NewWithOptions(logger.Options{Name: "", Writer: os.Stderr})

// WithName clones the default logger but changes the name of the logger.
func WithName(name string) *logger.Logger {
	return Logger.WithName(name)
}

// Panic is just like Fatal except that it is followed by a call to panic.
func Panic(message string) {
	Logger.Panic(message)
}

// Panicf is just like Fatalf except that it is followed by a call to panic.
func Panicf(format string, a ...interface{}) {
	Logger.Panicf(format, a...)
}

// Fatal logs a message at a Fatal Level.
func Fatal(message string) {
	Logger.Fatal(message)
}

// Fatalf logs a message at Fatal level.
func Fatalf(format string, a ...interface{}) {
	Logger.Fatalf(format, a...)
}

// Error logs a message at Error level.
func Error(message string) {
	Logger.Error(message)
}

// Errorf logs a message at Error level.
func Errorf(format string, a ...interface{}) {
	Logger.Errorf(format, a...)
}

// Warning logs a message at Warning level
func Warning(message string) {
	Logger.Warning(message)
}

// Warningf logs a message at Warning level.
func Warningf(format string, a ...interface{}) {
	Logger.Warningf(format, a...)
}

// Info logs a message at Info level.
func Info(message string) {
	Logger.Info(message)
}

// Infof logs a message at Info level.
func Infof(format string, a ...interface{}) {
	Logger.Infof(format, a...)
}

// Debug logs a message at Debug level.
func Debug(message string) {
	Logger.Debug(message)
}

// Debugf logs a message at Debug level.
func Debugf(format string, a ...interface{}) {
	Logger.Debugf(format, a...)
}

// Trace logs a message at Debug level.
func Trace(message string) {
	Logger.Trace(message)
}

// Tracef logs a message at Debug level.
func Tracef(format string, a ...interface{}) {
	Logger.Tracef(format, a...)
}
