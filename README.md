# go-logger

[![Go Reference](https://pkg.go.dev/badge/github.com/jorenkoyen/go-logger.svg)](https://pkg.go.dev/github.com/jorenkoyen/go-logger)

A simple logger with support for leveled logging.

# Installation

```bash
go get -u github.com/jorenkoyen/go-logger
```

# Getting Started

## Simple Logging Example

For simple logging you can import the global logger from `github.com/jorenkoyen/go-logger/log`

```go
package main

import (
	"github.com/jorenkoyen/go-logger/log"
)

func main() {
	log.Info("hello world")
}

// Output: ts=1729066152742 lvl=info msg="hello world"
```

## Bring Your Own Writer

You can create a logger instance and specify the writer it should use for outputting the log lines.

```go
package main

import (
	"os"

	"github.com/jorenkoyen/go-logger"
)

func main() {
	log := logger.New(os.Stdout)
	log.Info("hello world")
}

// Output: ts=1729066279358 logger=default lvl=info msg="hello world"
```

## Customize Your Logger

Create a fully customized logger with all your preferred options and a custom Formatter.

```go
package main

import (
	"os"

	"github.com/jorenkoyen/go-logger"
)

func main() {
	formatter := logger.NewTextFormatter()
	formatter.TimestampField = "" // disables timestamp printing

	log := logger.NewWithOptions(logger.Options{
		Name:      "my-logger",
		Writer:    os.Stdout,
		Level:     logger.LevelTrace, // explicit logger level overwrites logger.GlobalLevel
		Formatter: formatter,
	})

	log.Info("hello world")
}

// Output: logger=my-logger lvl=info msg="hello world"
```

## Pretty Formatting

You can make use of pretty formatting in environments where optimisation is not crucial such as during development.

```go
package main

import (
	"os"

	"github.com/jorenkoyen/go-logger"
)

func main() {
	formatter := logger.NewPrettyFormatter()
	log := logger.NewWithOptions(logger.Options{
		Writer:    os.Stdout,
		Formatter: formatter,
		Name:      "my-logger",
	})

	log.Info("hello world")
}

// Output: 2024-10-16 14:18:03 INF [default] hello world!
```

## Overwrite Default Logger

You can overwrite the default logger if you prefer to use other settings

```go
package main

import (
	"os"

	"github.com/jorenkoyen/go-logger"
	"github.com/jorenkoyen/go-logger/log"
)

func main() {
	formatter := logger.NewTextFormatter()
	formatter.LevelField = "l"
	formatter.TimestampField = "t"
	formatter.NameField = "n"

	log.SetDefaultLogger(logger.NewWithOptions(logger.Options{
		Name:      "default-logger",
		Writer:    os.Stdout,
		Level:     logger.LevelTrace, // explicit logger level overwrites logger.GlobalLevel
		Formatter: formatter,
	}))
	
	// assign to default logger
	log.Info("hello world")

	// clone logger with new name
	cloned := log.WithName("clone")
	cloned.Debug("hello from clone")
}

// Output: 
// t=1729067411654 n=default-logger l=info msg="hello world"
// t=1729067411654 n=clone l=debug msg="hello from clone"
```

# Tests

Run:

- `make test` to run all test.
- `make cov` to run coverage tests.
- `make bench` to run benchmark tests.

# Benchmark Results

```bash
BenchmarkLogEmpty-8             53943338                22.97  ns/op          0   B/op          0 allocs/op
BenchmarkDisabled-8             1000000000              0.3272 ns/op          0   B/op          0 allocs/op
BenchmarkInfo-8                 47163553                24.77  ns/op          0   B/op          0 allocs/op
BenchmarkFormatted-8            18861511                68.44  ns/op          48  B/op          1 allocs/op
BenchmarkLoggerNew-8            1000000000              0.3190 ns/op          0   B/op          0 allocs/op
BenchmarkTextFormatter-8        53173220                21.42  ns/op          282 B/op          0 allocs/op
```

## License

The [BSD 3-Clause license](http://opensource.org/licenses/BSD-3-Clause), the same as
the [Go language](http://golang.org/LICENSE).
