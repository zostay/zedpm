package metal

import (
	"io"
	"sync"
)

// SyncBuffer is a simple synchronized buffer.
type SyncBuffer struct {
	sync.Mutex
	w io.Writer
}

// NewSyncBuffer creates a new syncrhonized buffer wrapping the given io.Writer.
func NewSyncBuffer(w io.Writer) *SyncBuffer {
	return &SyncBuffer{w: w}
}

// Write calls the Write method of the wrapped io.Writer, but acuires a write
// lock before performing the operation, ensuring that only a single write is
// occurring at a given time.
func (b *SyncBuffer) Write(p []byte) (int, error) {
	b.Lock()
	defer b.Unlock()
	return b.w.Write(p)
}
