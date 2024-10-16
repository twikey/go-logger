package logger

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	space   = ' '
	equal   = '='
	quote   = '"'
	newline = '\n'
)

// Formatter declares a formatter to use when preparing to write a log entry line.
type Formatter interface {
	Format(event *Event)
}

// PrettyFormatter is a non-performance focused formatter that is useful during development.
// It supports colored formatting for improved usability when testing out things locally.
type PrettyFormatter struct {
}

func (s *PrettyFormatter) Format(event *Event) {
	//TODO implement me
	panic("implement me")
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
