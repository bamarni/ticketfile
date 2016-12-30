package escpos

import (
	"bytes"

	"github.com/bamarni/escpos"
	"github.com/bamarni/ticketfile"
)

type Converter struct {
	escpos *escpos.Escpos
	buf    *bytes.Buffer
}

func NewConverter() *Converter {
	var buf bytes.Buffer
	return &Converter{
		escpos: escpos.NewEscpos(&buf),
		buf:    &buf,
	}
}

func (c *Converter) Convert(cmd ticketfile.Command) ([]byte, error) {
	var err error
	switch cmd.Type {
	case ticketfile.Align:
		err = c.escpos.Align(cmd.Opcode[0])
	case ticketfile.Color:
		err = c.escpos.Color(cmd.Opcode[0])
	case ticketfile.Cut:
		if cmd.Opcode[0] == 1 {
			err = c.escpos.CutB(true, 3)
		} else {
			err = c.escpos.CutB(false, 3)
		}
	case ticketfile.Init:
		err = c.escpos.Init()
	case ticketfile.Charset:
		err = c.escpos.Charset(cmd.Opcode[0])
	case ticketfile.Marginleft:
		err = c.escpos.MarginLeft(uint16(cmd.Opcode[0]) + uint16(cmd.Opcode[1])<<8)
	case ticketfile.Print, ticketfile.Printraw:
		err = c.escpos.Print(cmd.Arg)
	case ticketfile.Printlf:
		err = c.escpos.Print(cmd.Arg + "\n")
	case ticketfile.Lf:
		err = c.escpos.Lf(cmd.Opcode[0])
	case ticketfile.Units:
		err = c.escpos.Units(cmd.Opcode[0], cmd.Opcode[1])
	case ticketfile.Font:
		err = c.escpos.Font(cmd.Opcode[0])
	case ticketfile.Barcode:
		err = c.escpos.Barcode(cmd.Opcode[0], string(cmd.Opcode[1:]))
	case ticketfile.BarcodeWidth:
		err = c.escpos.BarcodeWidth(cmd.Opcode[0])
	case ticketfile.BarcodeHeight:
		err = c.escpos.BarcodeHeight(cmd.Opcode[0])
	case ticketfile.BarcodeHRI:
		err = c.escpos.BarcodeHRI(cmd.Opcode[0])
	case ticketfile.BarcodeFont:
		err = c.escpos.BarcodeFont(cmd.Opcode[0])
	case ticketfile.Tab:
		err = c.escpos.Tab()
	case ticketfile.Tabs:
		err = c.escpos.TabPositions(cmd.Opcode...)
	case ticketfile.Printmode:
		err = c.escpos.PrintMode(cmd.Opcode[0])
	case ticketfile.Width:
		err = c.escpos.Width(uint16(cmd.Opcode[0]) + uint16(cmd.Opcode[1])<<8)
	}
	b := c.buf.Bytes()
	c.buf.Reset()
	return b, err
}
