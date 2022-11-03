package astool

import (
	"bytes"
	"fmt"
)

type Bytes struct {
	*bytes.Buffer
}

func NewBytes() *Bytes {
	return &Bytes{
		Buffer: bytes.NewBuffer(make([]byte, 0, 1024)),
	}
}

func (b *Bytes) P(args ...any) {
	b.WriteString(fmt.Sprintln(args...))
}

func (b *Bytes) Pf(format string, args ...any) {
	b.WriteString(fmt.Sprintf(format, args...))
	b.WriteByte('\n')
}
