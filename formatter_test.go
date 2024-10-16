package logger

import (
	"testing"
	"time"
)

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
