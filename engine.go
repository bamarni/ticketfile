package ticketfile

import (
	"fmt"
	"io"
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
	cmds chan Command
}

func NewEngine(w io.Writer, c Converter) *Engine {
	return &Engine{
		conv: c,
		w:    w,
		cmds: make(chan Command),
	}
}

func (e *Engine) Render(r io.Reader) error {
	go e.parse(r)

	for cmd := range e.cmds {
		rawBytes, err := e.conv.Convert(cmd)
		if err != nil {
			return fmt.Errorf("encoding error : %s\n", err)
		}

		_, err = e.w.Write(rawBytes)
		if err != nil {
			return fmt.Errorf("write error : %s\n", err)
		}
	}

	return nil
}
