package logger

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestParseLevel(t *testing.T) {
	equals := func(expected, actual interface{}) {
		t.Helper()
		if expected != actual {
			t.Errorf("expected: %s, got: %s", expected, actual)
		}
	}

	equals(LevelFatal, ParseLevel("fatal"))
	equals(LevelError, ParseLevel("error"))
	equals(LevelWarning, ParseLevel("warn"))
	equals(LevelInfo, ParseLevel("info"))
	equals(LevelDebug, ParseLevel("debug"))
	equals(LevelTrace, ParseLevel("trace"))
	equals(Level(0), ParseLevel("invalid")) // empty log level
}

func TestLoggerNew(t *testing.T) {
	log := New(io.Discard)
	if log == nil {
		t.Errorf("unable to create new logger")
		t.FailNow()
	}

	// writing to discard writer -> no output expected
	log.Infof("hello %s", "world")
}

func TestLoggerNewWithName(t *testing.T) {
	var buf bytes.Buffer
	log := NewWithOptions(Options{Name: "testing", Writer: &buf})

	log.Info("hello world!")

	got := buf.String()
	if !strings.Contains(got, "logger=testing") {
		t.Errorf("expected log message to contain logger=testing")
	}
}

func TestLogLevel(t *testing.T) {
	var tests = []struct {
		level    Level
		message  string
		expected int
	}{
		{
			LevelFatal,
			"Fatal Logging",
			1, // only expecting fatal output
		},
		{
			LevelError,
			"Error logging",
			2, // expecting 'fatal' + 'error'
		},
		{
			LevelWarning,
			"Warning logging",
			3, // expecting 'fatal' + 'error' + 'warning'
		},
		{
			LevelInfo,
			"Info Logging",
			4, // expecting 'fatal' + 'error' + 'warning' + 'info'
		},
		{
			LevelDebug,
			"Debug logging",
			5, // expecting 'fatal' + 'error' + 'warning' + 'info' + 'debug'
		},
		{
			LevelTrace,
			"Trace logging",
			6, // expecting 'fatal' + 'error' + 'warning' + 'info' + 'debug' + 'trace'
		},
	}

	var buf bytes.Buffer
	logger := New(&buf)
	logger.ignoreExit = true

	for _, test := range tests {
		logger.SetLogLevel(test.level)

		logger.Fatal("logging a fatal message")
		logger.Error("logging an error message")
		logger.Warning("logging a warning message")
		logger.Info("logging an info message")
		logger.Debug("logging a debug message")
		logger.Trace("logging a trace message")

		// Count output lines from logger
		lines := strings.Count(buf.String(), "\n")
		if lines != test.expected {
			t.Errorf("expected %d lines bot got %d lines when logging for lvl=%s", test.expected, lines, test.level.String())
		}

		// reset for next test
		buf.Reset()
	}
}

func TestLogLevelGlobal(t *testing.T) {
	var tests = []struct {
		level    Level
		message  string
		expected int
	}{
		{
			LevelFatal,
			"Fatal Logging",
			1, // only expecting fatal output
		},
		{
			LevelError,
			"Error logging",
			2, // expecting 'fatal' + 'error'
		},
		{
			LevelWarning,
			"Warning logging",
			3, // expecting 'fatal' + 'error' + 'warning'
		},
		{
			LevelInfo,
			"Info Logging",
			4, // expecting 'fatal' + 'error' + 'warning' + 'info'
		},
		{
			LevelDebug,
			"Debug logging",
			5, // expecting 'fatal' + 'error' + 'warning' + 'info' + 'debug'
		},
		{
			LevelTrace,
			"Trace logging",
			6, // expecting 'fatal' + 'error' + 'warning' + 'info' + 'debug' + 'trace'
		},
	}

	var buf bytes.Buffer
	logger := New(&buf)
	logger.ignoreExit = true

	for _, test := range tests {
		GlobalLevel = test.level

		logger.Fatal("logging a fatal message")
		logger.Error("logging an error message")
		logger.Warning("logging a warning message")
		logger.Info("logging an info message")
		logger.Debug("logging a debug message")
		logger.Trace("logging a trace message")

		// Count output lines from logger
		lines := strings.Count(buf.String(), "\n")
		if lines != test.expected {
			t.Errorf("expected %d lines bot got %d lines when logging for lvl=%s", test.expected, lines, test.level.String())
		}

		// reset for next test
		buf.Reset()
	}
}
