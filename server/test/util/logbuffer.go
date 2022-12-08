package util

import (
	"fmt"
	"github.com/rs/zerolog"
)

type BufferWriter struct {
	cw *zerolog.ConsoleWriter

	buffer []func(cw *zerolog.ConsoleWriter)
}

func (b *BufferWriter) Write(p []byte) (n int, err error) {
	b.buffer = append(b.buffer, func(cw *zerolog.ConsoleWriter) {
		_, _ = cw.Write(p)
	})
	return len(p), nil
}

func (b *BufferWriter) Dump() {
	for _, v := range b.buffer {
		v(b.cw)
	}
	b.buffer = nil
}

func (b *BufferWriter) Println(a ...any) {
	b.buffer = append(b.buffer, func(cw *zerolog.ConsoleWriter) {
		fmt.Println(a...)
	})
}

func (b *BufferWriter) Printf(format string, a ...any) {
	b.buffer = append(b.buffer, func(cw *zerolog.ConsoleWriter) {
		fmt.Printf(format, a...)
	})
}
