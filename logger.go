package logger

// Import packages
import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"time"
)

var (
	// GlobalLevel specifies the default logging level that each logger will get assigned if not explicitly defined.
	GlobalLevel = LevelInfo

	// DefaultLoggerName specifies the default logger name that each logger will get assigned if not explicitly defined.
	DefaultLoggerName = "default"
)

var (
	// defaultFormatter specifies the default formatter to use for each logger.
	defaultFormatter = NewTextFormatter()
)

// Level determines the level for which a message will be sent to the logger.
type Level int

const (
	LevelFatal Level = iota + 1
	LevelError
	LevelWarning
	LevelInfo
	LevelDebug
	LevelTrace

	LevelFatalValue   = "fatal"
	LevelErrorValue   = "error"
	LevelWarningValue = "warn"
	LevelInfoValue    = "info"
	LevelDebugValue   = "debug"
	LevelTraceValue   = "trace"
)

// String will print the pretty name of the logging level.
func (lvl Level) String() string {
	switch lvl {
	case LevelFatal:
		return LevelFatalValue
	case LevelError:
		return LevelErrorValue
	case LevelWarning:
		return LevelWarningValue
	case LevelInfo:
		return LevelInfoValue
	case LevelDebug:
		return LevelDebugValue
	case LevelTrace:
		return LevelTraceValue
	default:
		return ""
	}
}

// Logger defines the structure that is able to transmit a log line to the writer.
type Logger struct {
	formatter Formatter
	name      string
	level     Level
	w         io.Writer

	// only used for testing ...
	ignoreExit bool
}

// Options is used when creating a more extensive logger with the need of customization.
type Options struct {
	Name      string
	Formatter Formatter
	Level     Level
	Writer    io.Writer
}

// New returns a new logger instance. It will create a logger with optimistic defaults for ease of use.
func New(writer io.Writer) *Logger {
	return NewWithOptions(Options{
		Writer: writer,
		Name:   DefaultLoggerName,
	})
}

// NewWithName returns a new logger instance. It will create a logger with optimistic defaults for ease of use.
// It will write to os.Stdout by default and use the default assigned formatter.
func NewWithName(name string) *Logger {
	return NewWithOptions(Options{
		Name:   name,
		Writer: os.Stdout,
	})
}

// NewWithOptions returns a new logger instance according to the configuration of the Options.
// If no writer is specified it will use io.Discard. If no Formatter is specified it will use the default formatter.
// If no name is assigned it will use the DefaultLoggerName. And if no level is assigned it will use the GlobalLevel.
func NewWithOptions(opts Options) *Logger {
	if opts.Writer == nil {
		opts.Writer = io.Discard
	}

	if opts.Formatter == nil {
		opts.Formatter = defaultFormatter
	}

	return &Logger{
		w:         opts.Writer,
		formatter: opts.Formatter,
		name:      opts.Name,
		level:     opts.Level,
	}
}

// SetLogLevel will assign a new log level to the logger instance.
func (l *Logger) SetLogLevel(lvl Level) {
	l.level = lvl
}

// should returns true if the log event should be logged.
func (l *Logger) should(lvl Level) bool {
	if l.w == nil {
		return false
	}
	if l.level <= 0 {
		// use global level -> level not specified by logger
		return lvl <= GlobalLevel
	} else {
		// use logger level -> specific logger level
		return lvl <= l.level
	}
}

// log is the function available to user to log message, lvl specifies the severity of the message
// whilst message contains the actual information.
func (l *Logger) log(lvl Level, message string) {
	enabled := l.should(lvl)
	if !enabled {
		return // skip log line
	}

	// create new event
	e := getEvent()
	e.Time = time.Now()
	e.Module = l.name
	e.Level = lvl
	e.Message = message

	if pf, ok := l.formatter.(*PrettyFormatter); ok && pf.AppendSource {
		// append caller information for pretty formatter
		_, filename, line, _ := runtime.Caller(2)
		e.Filename = path.Base(filename)
		e.Line = line
	}

	// format using logger formatter -> this will update internal buffer of event
	l.formatter.Format(e)
	if _, err := l.w.Write(e.buf); err != nil {
		fmt.Fprintf(os.Stderr, "logger: could not write event: %v\n", err)
	}

	// put event back in event pool
	putEvent(e)
}

// WithName clones the logger instance and changes the name of the logger.
func (l *Logger) WithName(name string) *Logger {
	clone := &Logger{
		formatter: l.formatter,
		name:      name,
		w:         l.w,
		level:     l.level,
	}
	return clone
}

// Panic is just like Fatal except that it is followed by a call to panic.
func (l *Logger) Panic(message string) {
	l.log(LevelFatal, message)
	panic(message)
}

// Panicf is just like Fatalf except that it is followed by a call to panic.
func (l *Logger) Panicf(format string, a ...interface{}) {
	l.log(LevelFatal, fmt.Sprintf(format, a...))
	panic(fmt.Sprintf(format, a...))
}

// Fatal logs a message at a Fatal Level that is followed by an OS exit code.
func (l *Logger) Fatal(message string) {
	l.log(LevelFatal, message)
	if !l.ignoreExit {
		os.Exit(1)
	}
}

// Fatalf logs a message at Fatal level that is followed by an OS exit code.
func (l *Logger) Fatalf(format string, a ...interface{}) {
	l.log(LevelFatal, fmt.Sprintf(format, a...))
	if !l.ignoreExit {
		os.Exit(1)
	}
}

// Error logs a message at Error level.
func (l *Logger) Error(message string) {
	l.log(LevelError, message)
}

// Errorf logs a message at Error level.
func (l *Logger) Errorf(format string, a ...interface{}) {
	l.log(LevelError, fmt.Sprintf(format, a...))
}

// Warning logs a message at Warning level
func (l *Logger) Warning(message string) {
	l.log(LevelWarning, message)
}

// Warningf logs a message at Warning level.
func (l *Logger) Warningf(format string, a ...interface{}) {
	l.log(LevelWarning, fmt.Sprintf(format, a...))
}

// Info logs a message at Info level.
func (l *Logger) Info(message string) {
	l.log(LevelInfo, message)
}

// Infof logs a message at Info level.
func (l *Logger) Infof(format string, a ...interface{}) {
	l.log(LevelInfo, fmt.Sprintf(format, a...))
}

// Debug logs a message at Debug level.
func (l *Logger) Debug(message string) {
	l.log(LevelDebug, message)
}

// Debugf logs a message at Debug level.
func (l *Logger) Debugf(format string, a ...interface{}) {
	l.log(LevelDebug, fmt.Sprintf(format, a...))
}

// Trace logs a message at Debug level.
func (l *Logger) Trace(message string) {
	l.log(LevelTrace, message)
}

// Tracef logs a message at Debug level.
func (l *Logger) Tracef(format string, a ...interface{}) {
	l.log(LevelTrace, fmt.Sprintf(format, a...))
}
