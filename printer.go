package printer

import (
	"bufio"
	"github.com/bamarni/printer/escpos"
	"github.com/bamarni/printer/command"
	"log"
	"io"
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

func (p *Printer) Print() {
	cmds := make(chan command.Command)

	go command.Scan(p.input, cmds)

	for cmd := range cmds {
		log.Printf("Received command : [%s]\n", cmd.Raw)

		rawBytes, err := p.driver.ToBytes(cmd)
		if err != nil {
			log.Fatalf("Driver error : %s\n", err)
		}

		_, err = p.device.Write(rawBytes)
		if err != nil {
			log.Fatalf("Write error : %s\n", err)
		}

		if cmd.Name == command.Cut {
			err := p.device.Flush()
			if err != nil {
				log.Fatalf("Flush error : %s\n", err)
			}
		}
	}

	err := p.device.Flush()
	if err != nil {
		log.Fatalf("Flush error : %s\n", err)
	}
}
