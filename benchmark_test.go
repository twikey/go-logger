package logger

import (
	"io"
	"testing"
)

var (
	fakeMessage = "Test logging, but use a somewhat realistic message length."
)

func BenchmarkLogEmpty(b *testing.B) {
	logger := New(io.Discard)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("")
		}
	})
}

func BenchmarkDisabled(b *testing.B) {
	logger := New(io.Discard)
	logger.SetLogLevel(LevelFatal)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Debug(fakeMessage)
		}
	})
}

func BenchmarkInfo(b *testing.B) {
	logger := New(io.Discard)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info(fakeMessage)
		}
	})
}

func BenchmarkFormatted(b *testing.B) {
	logger := New(io.Discard)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// bool:              %t
			// int, int8 etc.:    %d
			// float32:           %g
			// string:            %s
			logger.Infof("bool=%t int=%d float=%g string=%s", true, 100, 22.23, "hello")
		}
	})
}

func BenchmarkPrettyFormatter(b *testing.B) {
	formatter := NewPrettyFormatter()
	logger := NewWithOptions(Options{Writer: io.Discard, Formatter: formatter})
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("Hello world!")
		}
	})
}

func BenchmarkPrettyFormatterWithCallerSource(b *testing.B) {
	formatter := NewPrettyFormatter()
	formatter.AppendSource = true
	logger := NewWithOptions(Options{Writer: io.Discard, Formatter: formatter})
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("Hello world!")
		}
	})
}

func BenchmarkLoggerNew(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger := New(io.Discard)
			if logger == nil {
				b.Fatal("unable to create logger")
			}
		}
	})
}
