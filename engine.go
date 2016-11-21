package ticketfile

import (
	"bufio"
	"fmt"
	"io"
	"sync"
)

type Converter interface {
	Convert(cmd Command) ([]byte, error)
}

type Engine struct {
	conv Converter
	w    *bufio.Writer
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

	cmds, err := parse(r)
	if err != nil {
		return fmt.Errorf("parsing error : %s", err)
	}

	for _, cmd := range cmds {
		rawBytes, err := e.conv.Convert(cmd)
		if err != nil {
			return fmt.Errorf("converter error : %s", err)
		}

		_, err = e.w.Write(rawBytes)
		if err != nil {
			return fmt.Errorf("write error : %s", err)
		}
	}

	return e.w.Flush()
}
