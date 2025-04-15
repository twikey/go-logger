package logger

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestJournalFormatter_noModule(t *testing.T) {
	formatter := NewJournalFormatter()
	e := &Event{
		Level:   LevelInfo,
		Message: "Hello World",
	}

	want := "info - Hello World\n"
	formatter.Format(e)
	if want != string(e.buf) {
		t.Errorf("\nWant: %s\nGot: %s", want, string(e.buf))
	}
}

func TestJournalFormatter_simple(t *testing.T) {
	formatter := NewJournalFormatter()
	e := &Event{
		Level:   LevelInfo,
		Module:  "main",
		Message: "Hello World",
	}

	want := "[main] info - Hello World\n"
	formatter.Format(e)
	if want != string(e.buf) {
		t.Errorf("\nWant: %s\nGot: %s", want, string(e.buf))
	}
}

func TestTextFormatter_simple(t *testing.T) {
	formatter := NewTextFormatter()
	e := &Event{
		Time:     time.Unix(1, 0),
		Level:    LevelInfo,
		Line:     100,
		Filename: "example.go",
		Message:  "hello",
	}

	want := "ts=1000 lvl=info msg=hello\n"
	formatter.Format(e)
	if want != string(e.buf) {
		t.Errorf("\nWant: %sHave: %s", want, string(e.buf))
	}
}

func TestTextFormatter_quotable(t *testing.T) {
	formatter := NewTextFormatter()
	e := &Event{
		Time:     time.Unix(1, 0),
		Level:    LevelInfo,
		Line:     100,
		Filename: "example.go",
		Message:  "hello world!",
	}

	want := "ts=1000 lvl=info msg=\"hello world!\"\n"
	formatter.Format(e)
	if want != string(e.buf) {
		t.Errorf("\nWant: %sHave: %s", want, string(e.buf))
	}
}

func TestTextFormatter_escapeCharacters(t *testing.T) {
	// TODO: not yet implemented
	t.SkipNow()

	formatter := NewTextFormatter()
	e := &Event{
		buf:      make([]byte, 0, 500),
		Time:     time.Unix(1, 0),
		Level:    LevelInfo,
		Line:     100,
		Filename: "example.go",
		Message:  "hello escape=\"me\"",
	}

	want := "ts=1000 lvl=info msg=\"hello escape=\\\"me\\\"\"\n"
	formatter.Format(e)
	if want != string(e.buf) {
		t.Errorf("\nWant: %sHave: %s", want, string(e.buf))
	}
}

func BenchmarkTextFormatter(b *testing.B) {
	formatter := NewTextFormatter()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		e := &Event{
			buf:      make([]byte, 0, 500),
			Time:     time.Unix(1, 0),
			Module:   "DEFAULT",
			Level:    LevelInfo,
			Line:     100,
			Filename: "example.go",
			Message:  "Hello world!",
		}

		for pb.Next() {
			formatter.Format(e)
		}
	})
}

func BenchmarkJournalFormatter(b *testing.B) {
	formatter := NewJournalFormatter()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		e := &Event{
			buf:     make([]byte, 0, 500),
			Time:    time.Unix(1, 0),
			Module:  "main",
			Level:   LevelInfo,
			Message: "Hello world!",
		}

		for pb.Next() {
			formatter.Format(e)
		}
	})
}

func TestPrettyFormatter(t *testing.T) {
	var buf bytes.Buffer

	formatter := NewPrettyFormatter()
	formatter.AppendSource = false
	log := NewWithOptions(Options{Formatter: formatter, Writer: &buf})

	log.Info("hello world!")

	want := "\u001B[32mINF\u001B[0m \u001B[90m[\u001B[0m\u001B[37mdefault\u001B[0m\u001B[90m]\u001B[0m hello world!\n"
	got := buf.String()
	if !strings.HasSuffix(got, want) {
		t.Errorf("Incorrect log suffix from output -> \nactual: %s", got)
	}
}
