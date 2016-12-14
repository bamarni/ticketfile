package ticketfile

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/bamarni/escpos"
)

const (
	undefined CommandType = iota
	Align
	Charset
	Color
	Cut
	Font
	Init
	Lf
	Marginleft
	Units
	Barcode
	BarcodeFont
	BarcodeHRI
	BarcodeHeight
	BarcodeWidth
	Print
	Printlf
	multiline
	Printraw
)

var commands = map[string]CommandType{
	"ALIGN":          Align,
	"CHARSET":        Charset,
	"COLOR":          Color,
	"CUT":            Cut,
	"FONT":           Font,
	"INIT":           Init,
	"LF":             Lf,
	"MARGINLEFT":     Marginleft,
	"PRINT":          Print,
	"PRINTLF":        Printlf,
	"PRINTRAW":       Printraw,
	"UNITS":          Units,
	"BARCODE":        Barcode,
	"BARCODE_HRI":    BarcodeHRI,
	"BARCODE_HEIGHT": BarcodeHeight,
	"BARCODE_WIDTH":  BarcodeWidth,
	"BARCODE_FONT":   BarcodeFont,
}

type (
	CommandType int
	Command     struct {
		Type   CommandType
		Arg    string
		Opcode []byte
	}
)

func (cmdType CommandType) isMultiline() bool {
	return multiline < cmdType
}

func NewCommand(cmdType CommandType, arg string) (Command, error) {
	var err error
	var opcode []byte

	switch cmdType {
	case Align:
		if arg == "LEFT" {
			opcode = []byte{escpos.AlignLeft}
		} else if arg == "CENTER" {
			opcode = []byte{escpos.AlignCenter}
		} else if arg == "RIGHT" {
			opcode = []byte{escpos.AlignRight}
		} else {
			err = fmt.Errorf("unsupported alignment %s", arg)
		}
	case Color:
		if arg == "BLACK" {
			opcode = []byte{escpos.ColorBlack}
		} else if arg == "RED" {
			opcode = []byte{escpos.ColorRed}
		} else {
			err = fmt.Errorf("unsupported color %s", arg)
		}
	case Cut:
		if arg == "" || arg == "FULL" {
			opcode = []byte{0}
		} else if arg == "PARTIAL" {
			opcode = []byte{1}
		} else {
			err = fmt.Errorf("unsupported cut %s", arg)
		}
	case Font, BarcodeFont:
		if arg == "A" {
			opcode = []byte{escpos.FontA}
		} else if arg == "B" {
			opcode = []byte{escpos.FontB}
		} else if arg == "C" {
			opcode = []byte{escpos.FontC}
		} else {
			err = fmt.Errorf("unsupported font %s", arg)
		}
	case Marginleft:
		margin, err := strconv.ParseUint(arg, 10, 16)
		if err != nil {
			err = errors.New("invalid left margin")
		} else {
			opcode = []byte{byte(margin & 0xff), byte(margin >> 8)}
		}
	case Units:
		units := strings.Fields(arg)
		if len(units) != 2 {
			err = errors.New("expected 2 units")
			break
		}
		x, err := strconv.ParseUint(units[0], 10, 8)
		if err != nil {
			err = errors.New("invalid horizontal unit")
			break
		}
		y, err := strconv.ParseUint(units[1], 10, 8)
		if err != nil {
			err = errors.New("invalid vertical unit")
			break
		}
		opcode = []byte{byte(x), byte(y)}
	case Charset:
		if arg == "PC437" {
			opcode = []byte{escpos.CharsetPC437}
		} else if arg == "PC850" {
			opcode = []byte{escpos.CharsetPC850}
		} else if arg == "PC860" {
			opcode = []byte{escpos.CharsetPC860}
		} else if arg == "PC863" {
			opcode = []byte{escpos.CharsetPC863}
		} else if arg == "PC865" {
			opcode = []byte{escpos.CharsetPC865}
		} else {
			err = fmt.Errorf("charset %s not supported", arg)
		}
	case Barcode:
		args := tokenWhitespace.Split(arg, 2)
		if len(args) != 2 {
			err = errors.New("expected 2 args")
			break
		}
		if args[0] == "CODE39" {
			opcode = []byte{escpos.BarcodeCODE39}
		} else {
			err = fmt.Errorf("barcode system %s not supported", args[0])
		}
		opcode = append(opcode, args[1]...)
	case BarcodeHRI:
		if arg == "TOP" {
			opcode = []byte{escpos.BarcodeHRITop}
		} else if arg == "BOTTOM" {
			opcode = []byte{escpos.BarcodeHRIBottom}
		} else if arg == "BOTH" {
			opcode = []byte{escpos.BarcodeHRIBoth}
		} else if arg == "NONE" {
			opcode = []byte{escpos.BarcodeHRINone}
		} else {
			err = fmt.Errorf("barcode hri position %s not supported", arg)
		}
	case BarcodeWidth, BarcodeHeight:
		n, err := strconv.ParseUint(arg, 10, 8)
		if err != nil {
			err = errors.New("invalid barcode width / height")
			break
		}
		opcode = []byte{byte(n)}
	case Lf:
		if arg != "" {
			n, err := strconv.ParseUint(arg, 10, 8)
			if err != nil {
				err = errors.New("invalid line feed")
			}
			opcode = []byte{byte(n)}
		} else {
			opcode = []byte{1}
		}
	}

	return Command{Type: cmdType, Arg: arg, Opcode: opcode}, err
}
