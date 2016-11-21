package escpos

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/bamarni/ticketfile"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
)

type Converter struct {
	enc *encoding.Encoder
}

func NewConverter() *Converter {
	return &Converter{
		enc: charmap.CodePage437.NewEncoder(),
	}
}

var dispatchTable = map[ticketfile.CommandType]func(*Converter, ticketfile.Command) (string, error){
	ticketfile.Align:      handleAlign,
	ticketfile.Charset:    handleCharset,
	ticketfile.Color:      handleColor,
	ticketfile.Cut:        handleCut,
	ticketfile.Font:       handleFont,
	ticketfile.Init:       handleInit,
	ticketfile.Lf:         handleLf,
	ticketfile.Marginleft: handleMarginleft,
	ticketfile.Print:      handlePrint,
	ticketfile.Printlf:    handlePrintlf,
	ticketfile.Printraw:   handlePrintraw,
	ticketfile.Units:      handleUnits,
	ticketfile.Barcode:    handleBarcode,
}

func (c *Converter) Convert(cmd ticketfile.Command) ([]byte, error) {
	if f, ok := dispatchTable[cmd.Type]; ok {
		rawCmd, error := f(c, cmd)

		return []byte(rawCmd), error
	}
	return nil, nil
}

func handleAlign(c *Converter, cmd ticketfile.Command) (string, error) {
	switch cmd.Arg {
	case "LEFT":
		return "\x1Ba0", nil
	case "CENTER":
		return "\x1Ba1", nil
	case "RIGHT":
		return "\x1Ba2", nil
	}
	return "", fmt.Errorf("unsupported alignment %s", cmd.Arg)
}

func handleCut(c *Converter, cmd ticketfile.Command) (string, error) {
	switch cmd.Arg {
	case "FULL":
		return "\x1DVA0", nil
	case "PARTIAL", "":
		return "\x1DVA1", nil
	}
	return "", fmt.Errorf("unsupported cut %s", cmd.Arg)
}

func handleFont(c *Converter, cmd ticketfile.Command) (string, error) {
	switch cmd.Arg {
	case "A":
		return "\x1BM0", nil
	case "B":
		return "\x1BM1", nil
	case "C":
		return "\x1BM2", nil
	}
	return "", fmt.Errorf("unsupported font %s", cmd.Arg)
}

// [Name]	Set left margin
// [Format]
// 	ASCII		GS		L		nL		nH
//	Hex		1D		4C		nL		nH
//	Decimal		29		76		nL		nH
// [Range]	(nL + nH × 256) = 0 – 65535
// [Default]	(nL + nH × 256) = 0
func handleMarginleft(c *Converter, cmd ticketfile.Command) (string, error) {
	margin, err := strconv.Atoi(cmd.Arg)
	if err != nil {
		return "", err
	}
	if margin > 65535 {
		return "", errors.New("invalid left margin")
	}

	return fmt.Sprintf("\x1DL%c%c", margin%256, margin/256), nil
}

// [Name]	Set horizontal and vertical motion units
// [Format]
// 	ASCII	   	GS	  	P	  	x	  	y
// 	Hex		1D		50		x		y
// 	Decimal		29		80		x		y
// [Range]
// 	x = 0 – 255
// 	y = 0 – 255
// [Default]	x, y: different depending on the printers
func handleUnits(c *Converter, cmd ticketfile.Command) (string, error) {
	units := strings.Fields(cmd.Arg)
	x, err := strconv.Atoi(units[0])
	if err != nil {
		return "", err
	}
	if x > 255 {
		return "", errors.New("invalid horizontal unit")
	}
	y, err := strconv.Atoi(units[1])
	if err != nil {
		return "", err
	}
	if y > 255 {
		return "", errors.New("invalid vertical unit")
	}

	return fmt.Sprintf("\x1DP%c%c", x, y), nil
}

// [Name]	Select print color
// [Format]
// 	ASCII	 	ESC	  	r	  	n
// 	Hex		1B		72		n
// 	Decimal		27		114		n
// [Range]	n = 0, 1, 48, 49
// [Default]	n = 0
func handleColor(c *Converter, cmd ticketfile.Command) (string, error) {
	if cmd.Arg == "RED" {
		return "\x1Br1", nil
	}
	return "\x1Br0", nil
}

// [Name]	Select character code table
// [Format]
// 	ASCII	   	ESC	  	t	  	n
// 	Hex		1B		74		n
// 	Decimal		27		116		n
// [Range]	different depending on the printers
// [Default]
// 	n = 20	   	[Thai models]
// 	n = 0	   	[Other models]
func handleCharset(c *Converter, cmd ticketfile.Command) (string, error) {
	var n byte
	switch cmd.Arg {
	case "PC437": // USA: Standard Europe
		c.enc = charmap.CodePage437.NewEncoder()
		n = 0
	case "PC850": // Western Europe
		c.enc = charmap.CodePage850.NewEncoder()
		n = 2
	default:
		return "", fmt.Errorf("charset %s not supported", cmd.Arg)
	}
	return fmt.Sprintf("\x1Bt%c", n), nil
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
	fonts := map[string]int{"A": 0, "B": 1, "C": 2, "D": 3, "E": 4}
	return fmt.Sprintf("\x1Df%c", fonts[value]), nil
}

func handleInit(c *Converter, cmd ticketfile.Command) (string, error) {
	return "\x1B@", nil
}

func handleLf(c *Converter, cmd ticketfile.Command) (string, error) {
	if cmd.Arg != "" {
		nb, err := strconv.Atoi(cmd.Arg)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("\x1Bd%c", nb), nil
	}
	return "\n", nil
}

func handlePrint(c *Converter, cmd ticketfile.Command) (string, error) {
	return c.encode(cmd.Arg)
}

func handlePrintlf(c *Converter, cmd ticketfile.Command) (string, error) {
	s := cmd.Arg + "\n"
	return c.encode(s)
}

func handlePrintraw(c *Converter, cmd ticketfile.Command) (string, error) {
	return c.encode(cmd.Arg)
}

func (c *Converter) encode(s string) (string, error) {
	s, err := c.enc.String(s)
	if err != nil {
		return "", fmt.Errorf("couldn't encode to charset (%s)", err)
	}
	return s, nil
}
