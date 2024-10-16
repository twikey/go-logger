package logger

import (
	"sync"
	"time"
)

// Event contains all necessary information when a logging event is emitted from the logger.
// The Formatter of the logger will further handle this event to write this to the io.Writer.
type Event struct {
	buf      []byte
	Time     time.Time
	Module   string
	Level    Level
	Line     int
	Filename string
	Message  string
}

// eventPool is used to efficiently make use of our internal buffer.
var eventPool = &sync.Pool{
	New: func() interface{} {
		return &Event{
			buf: make([]byte, 0, 500),
		}
	},
}

// putEvent will place the Event back into the ring buffer for efficient memory usage.
// Inspiration was token from zerolog in order to reduce the required allocations per log message.
func putEvent(e *Event) {
	// Proper usage of a sync.Pool requires each entry to have approximately
	// the same memory cost. To obtain this property when the stored type
	// contains a variably-sized buffer, we add a hard limit on the maximum buffer
	// to place back in the pool.
	//
	// See https://golang.org/issue/23199
	const maxSize = 1 << 16 // 64KiB
	if cap(e.buf) > maxSize {
		return
	}
	eventPool.Put(e)
}

// getEvent will retrieve an Event from the ring buffer.
func getEvent() *Event {
	e := eventPool.Get().(*Event)
	e.buf = e.buf[:0] // truncate buffer
	return e
}
