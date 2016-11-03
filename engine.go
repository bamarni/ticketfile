package ticketfile

import (
	"bufio"
	"fmt"
	"io"
	"sync"
)

type Command struct {
	Raw  string
	Name string
	Arg  string
}

type Converter interface {
	Convert(cmd Command) ([]byte, error)
}

type Engine struct {
	conv Converter
	w    io.Writer
	mu   sync.Mutex
}

func NewEngine(w io.Writer, c Converter) *Engine {
	return &Engine{
		conv: c,
		w:    bufio.NewWriter(w),
	}
}

func (e *Engine) Render(r io.Reader) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	cmds := parse(r)

	for cmd := range cmds {
		rawBytes, err := e.conv.Convert(cmd)
		if err != nil {
			return fmt.Errorf("encoding error : %s\n", err)
		}

		_, err = e.w.Write(rawBytes)
		if err != nil {
			return fmt.Errorf("write error : %s\n", err)
		}
	}

	return e.w.Flush()
}
