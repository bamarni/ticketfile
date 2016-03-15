package printer

import (
	"bufio"
	"github.com/bamarni/printer/escpos"
	"github.com/bamarni/printer/command"
	"io"
	"fmt"
)

type Printer struct {
	driver *escpos.Escpos
	device *bufio.Writer
}

func NewPrinter(w io.Writer) *Printer {
	printer := Printer{
		driver: escpos.NewEscpos(),
		device: bufio.NewWriter(w),
	}

	return &printer
}

func (p *Printer) Print(r io.Reader) error {
	cmds := make(chan command.Command)

	go command.Scan(r, cmds)

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
