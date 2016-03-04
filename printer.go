package printer

import (
	"bufio"
	"github.com/bamarni/printer/escpos"
	"github.com/bamarni/printer/command"
	"io"
	"fmt"
)

type Printer struct {
	input io.Reader
	driver *escpos.Escpos
	device *bufio.Writer
}

func NewPrinter(r io.Reader, w io.Writer) *Printer {
	printer := Printer{
		input: r,
		driver: escpos.NewEscpos(),
		device: bufio.NewWriter(w),
	}

	return &printer
}

func (p *Printer) Print() error {
	cmds := make(chan command.Command)

	go command.Scan(p.input, cmds)

	for cmd := range cmds {
		rawBytes, err := p.driver.ToBytes(cmd)
		if err != nil {
			return fmt.Errorf("Driver error : %s\n", err)
		}

		_, err = p.device.Write(rawBytes)
		if err != nil {
			return fmt.Errorf("Write error : %s\n", err)
		}

		if cmd.Name == command.Cut {
			if err := p.device.Flush(); err != nil {
				return fmt.Errorf("Flush error : %s\n", err)
			}
		}
	}

	return p.device.Flush()
}
