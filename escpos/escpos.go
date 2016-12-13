package escpos

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

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

var dispatchTable = map[ticketfile.CommandType]func(*Converter, ticketfile.Command) (string, error){
	ticketfile.Barcode: handleBarcode,
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
	}

	if f, ok := dispatchTable[cmd.Type]; ok {
		rawCmd, error := f(c, cmd)

		return []byte(rawCmd), error
	}
	return nil, nil
}

func handleBarcode(c *Converter, cmd ticketfile.Command) (string, error) {
	args := strings.Fields(cmd.Arg)
	subCmd := args[0]
	switch subCmd {
	case "PRINT":
		return barcodePrint(args[1], args[2])
	case "WIDTH":
		return barcodeWidth(args[1])
	case "HEIGHT":
		return barcodeHeight(args[1])
	case "HRI":
		if args[1] == "FONT" {
			return barcodeHriFont(args[2])
		} else if args[1] == "DISPLAY" {
			return barcodeHriDisplay(args[2])
		} else {
			return "", fmt.Errorf("%s barcode hri subcommand not supported", args[1])
		}
	}

	return "", fmt.Errorf("%s barcode subcommand not supported", subCmd)
}

// [Name] 	Print barcode
// [Format]
// 	(A)	ASCII	    	GS	   	k	   	m	   	d1 ... dk	   	NUL
//		Hex		1D		6B		m		d1 ... dk		NUL
//		Decimal		29		107		m		d1 ... dk		NUL
//
//	(B)	ASCII	    	GS	   	k	   	m	   	n	   	d1 ... dn
//		Hex		1D		6B		m		n		d1 ... dn
//		Decimal		29		107		m		n		d1 ... dn
// [Range] 	m: different depending on the printers d, k of (A), and d, n of (B): different depending on the barcode format. Refer to the tables in the ESC/POS specification.
// [Default]	None
func barcodePrint(format, value string) (string, error) {
	// TODO : validate value according to the format
	var m int
	if format == "CODE39" {
		m = 4
	} else {
		return "", fmt.Errorf("%s barcode format not supported", format)
	}

	// function (A)
	return fmt.Sprintf("\x1Dk%c%s\x00", m, value), nil
}

// [Name]		Set barcode width
// [Format]
//   ASCII		   	GS	  	w	  	n
//   Hex			1D		77		n
//   Decimal		29		119		n
// [Range]		n: different depending on the printers
// [Default]	n: different depending on the printers
func barcodeWidth(value string) (string, error) {
	width, err := strconv.Atoi(value)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("\x1Dw%c", width), nil
}

// [Name]	Set barcode height
// [Format]
//   ASCII		   	GS	  	h	  	n
//   Hex			1D		68		n
//   Decimal		29		104		n
// [Range]		n = 1 – 255
// [Default]	n: different depending on the printers
func barcodeHeight(value string) (string, error) {
	height, err := strconv.Atoi(value)
	if err != nil {
		return "", err
	}
	if height < 1 || height > 255 {
		return "", errors.New("invalid height")
	}
	return fmt.Sprintf("\x1Dh%c", height), nil
}

// [Name]		Select print position of HRI characters
// [Format]
//   ASCII	   	GS	  	H	  	n
//   Hex		1D		48		n
//   Decimal	29		72		n
// [Range]		n = 0 – 3
// 				n = 48 – 51
// [Default]	n = 0
func barcodeHriDisplay(value string) (string, error) {
	var display int
	switch value {
	case "TOP":
		display = 1
	case "BOTTOM":
		display = 2
	case "BOTH":
		display = 3
	}

	return fmt.Sprintf("\x1DH%c", display), nil
}

// [Name]	Select font for HRI characters
// [Format]
//   ASCII		   	GS	  	f	  	n
//   Hex			1D		66		n
//   Decimal		29		102		n
// [Range]		n: different depending on the printers
// [Default]	n = 0
func barcodeHriFont(value string) (string, error) {
	return fmt.Sprintf("\x1Df%c", value[0]-65), nil
}
