package escpos

import (
	"errors"

	"github.com/bamarni/escpos"
	"github.com/bamarni/ticketfile"
)

type Converter struct {
	escpos *escpos.Escpos
}

func NewConverter() *Converter {
	return &Converter{
		escpos: escpos.NewEscpos(),
	}
}

func (c *Converter) Convert(cmd ticketfile.Command) ([]byte, error) {
	switch cmd.Type {
	case ticketfile.Align:
		return c.escpos.Align(cmd.Opcode[0]), nil
	case ticketfile.Color:
		return c.escpos.Color(cmd.Opcode[0]), nil
	case ticketfile.Cut:
		if cmd.Opcode[0] == 1 {
			return c.escpos.Cut(true), nil
		}
		return c.escpos.Cut(false), nil
	case ticketfile.Init:
		return c.escpos.Init(), nil
	case ticketfile.Charset:
		return c.escpos.Charset(cmd.Opcode[0]), nil
	case ticketfile.Marginleft:
		return c.escpos.MarginLeft(uint16(cmd.Opcode[0]) + uint16(cmd.Opcode[1])<<8), nil
	case ticketfile.Print, ticketfile.Printraw:
		return c.escpos.Print(cmd.Arg)
	case ticketfile.Printlf:
		return c.escpos.Print(cmd.Arg + "\n")
	case ticketfile.Lf:
		return c.escpos.Lf(cmd.Opcode[0]), nil
	case ticketfile.Units:
		return c.escpos.Units(cmd.Opcode[0], cmd.Opcode[1]), nil
	case ticketfile.Font:
		return c.escpos.Font(cmd.Opcode[0]), nil
	case ticketfile.Barcode:
		return c.escpos.Barcode(cmd.Opcode[0], string(cmd.Opcode[1:]))
	case ticketfile.BarcodeWidth:
		return c.escpos.BarcodeWidth(cmd.Opcode[0]), nil
	case ticketfile.BarcodeHeight:
		return c.escpos.BarcodeHeight(cmd.Opcode[0]), nil
	case ticketfile.BarcodeHRI:
		return c.escpos.BarcodeHRI(cmd.Opcode[0]), nil
	case ticketfile.BarcodeFont:
		return c.escpos.BarcodeFont(cmd.Opcode[0]), nil
	}
	return nil, errors.New("unsupported command")
}
