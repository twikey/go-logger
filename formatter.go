package logger

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	// characters
	space        = ' '
	equal        = '='
	quote        = '"'
	newline      = '\n'
	colon        = ':'
	bracketLeft  = "["
	bracketRight = "]"
	hyphen       = "-"
	reset        = "\033[0m"
)

const (
	// colors
	black = iota + 30
	red
	green
	yellow
	blue
	magenta
	cyan
	white

	grey = 90
	bold = 1
)

var (
	// prettyLevelNames are used by PrettyFormatter  for a short level name.
	prettyLevelNames = map[Level]string{
		LevelTrace:   "TRC",
		LevelDebug:   "DBG",
		LevelInfo:    "INF",
		LevelWarning: "WRN",
		LevelError:   "ERR",
		LevelFatal:   "FTL",
	}

	// prettyLevelColors are used by PrettyFormatter to color log levels.
	prettyLevelColors = map[Level]int{
		LevelTrace:   blue,
		LevelDebug:   0,
		LevelInfo:    green,
		LevelWarning: yellow,
		LevelError:   red,
		LevelFatal:   red,
	}
)

// Formatter declares a formatter to use when preparing to write a log entry line.
type Formatter interface {
	Format(event *Event)
}

// PrettyFormatter is a non-performance focused formatter that is useful during development.
// It supports colored formatting for improved usability when testing out things locally.
type PrettyFormatter struct {
	TimeFormat   string
	AppendSource bool
}

// NewPrettyFormatter  creates a new instance of the PrettyFormatter which output log lines in a pretty format.
func NewPrettyFormatter() *PrettyFormatter {
	return &PrettyFormatter{
		TimeFormat: time.DateTime,
	}
}

func (s *PrettyFormatter) Format(event *Event) {
	s.color(event, grey, event.Time.Format(s.TimeFormat))
	event.buf = append(event.buf, space)
	s.color(event, prettyLevelColors[event.Level], prettyLevelNames[event.Level])
	event.buf = append(event.buf, space)
	s.color(event, grey, bracketLeft)
	if event.Module == "" {
		s.color(event, white, DefaultLoggerName)
	} else {
		s.color(event, white, event.Module)
	}
	s.color(event, grey, bracketRight)
	event.buf = append(event.buf, space)

	if event.Level < LevelInfo {
		// put message in bold
		s.color(event, bold, event.Message)
	} else {
		// normal output of message
		event.buf = append(event.buf, event.Message...)
	}

	// append source information
	if s.AppendSource {
		event.buf = append(event.buf, space)
		s.color(event, cyan, "source=")
		event.buf = append(event.buf, event.Filename...)
		event.buf = append(event.buf, colon)
		event.buf = strconv.AppendInt(event.buf, int64(event.Line), 10)
	}

	event.buf = append(event.buf, newline)
}

func (s *PrettyFormatter) color(e *Event, color int, value string) {
	if color > 0 {
		code := fmt.Sprintf("\x1b[%dm", color)
		e.buf = append(e.buf, code...)
		e.buf = append(e.buf, value...)
		e.buf = append(e.buf, reset...)
	} else {
		// without any color
		e.buf = append(e.buf, value...)
	}
}

// TextFormatter is a performance focused formatter prints out the log lines in a logfmt style.
type TextFormatter struct {
	// field names
	TimestampField string
	LevelField     string
	MessageField   string
	NameField      string
}

// NewTextFormatter creates a new instance of the TextFormatter which outputs log lines in logfmt style.
func NewTextFormatter() *TextFormatter {
	return &TextFormatter{
		TimestampField: "ts",
		LevelField:     "lvl",
		MessageField:   "msg",
		NameField:      "logger",
	}
}

func (t *TextFormatter) Format(event *Event) {
	t.encode(event, t.TimestampField, event.Time.UnixMilli(), false)
	if event.Module != "" {
		// only write module when not empty
		t.encode(event, t.NameField, event.Module, false)
	}
	t.encode(event, t.LevelField, event.Level.String(), false)
	t.encode(event, t.MessageField, event.Message, true)
}

func (t *TextFormatter) encode(e *Event, key string, value interface{}, eol bool) {
	if key == "" {
		return // skip encoding -> key is empty.
	}

	// key=valueString
	t.key(e, key)
	t.equal(e)

	// write value
	switch v := value.(type) {
	case int64:
		t.valueInt64(e, v)
	case string:

		t.valueString(e, v)
	}

	if eol {
		// write newline
		e.buf = append(e.buf, byte(newline))
	} else {
		// write space
		e.buf = append(e.buf, byte(space))
	}
}

func (t *TextFormatter) key(e *Event, key string) {
	e.buf = append(e.buf, key...)
}

func (t *TextFormatter) equal(e *Event) {
	e.buf = append(e.buf, byte(equal))
}

func (t *TextFormatter) valueString(e *Event, value string) {
	if strings.IndexFunc(value, t.needsQuotedValueRune) != -1 {
		e.buf = append(e.buf, byte(quote))
		e.buf = append(e.buf, value...)
		e.buf = append(e.buf, byte(quote))
	} else {
		e.buf = append(e.buf, value...)
	}
}

func (t *TextFormatter) valueInt64(e *Event, value int64) {
	e.buf = strconv.AppendInt(e.buf, value, 10)
}

func (t *TextFormatter) needsQuotedValueRune(r rune) bool {
	return r <= ' ' || r == '=' || r == '"' || r == utf8.RuneError
}

// JournalFormatter is a formatter which prints log lines in the following output:
//
// [module] level - message
type JournalFormatter struct {
}

// NewJournalFormatter creates a new formatter which outputs log lines in a journalctl friendly format.
func NewJournalFormatter() *JournalFormatter {
	return &JournalFormatter{}
}

func (j *JournalFormatter) Format(e *Event) {
	if e.Module != "" {
		e.buf = append(e.buf, bracketLeft...)
		e.buf = append(e.buf, e.Module...)
		e.buf = append(e.buf, bracketRight...)
		e.buf = append(e.buf, space)
	}

	e.buf = append(e.buf, e.Level.String()...)
	e.buf = append(e.buf, space)
	e.buf = append(e.buf, hyphen...)
	e.buf = append(e.buf, space)
	e.buf = append(e.buf, e.Message...)
	e.buf = append(e.buf, newline)
}
