package escpos

import (
	"fmt"
	"github.com/bamarni/ticketfile"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"strconv"
)

type Converter struct {
	enc *encoding.Encoder
}

func NewConverter() *Converter {
	return &Converter{
		enc: charmap.CodePage437.NewEncoder(),
	}
}

var dispatchTable map[string]func(*Converter, ticketfile.Command) (string, error)

func init() {
	dispatchTable = map[string]func(*Converter, ticketfile.Command) (string, error){
		ticketfile.Align:    handleAlign,
		ticketfile.Charset:  handleCharset,
		ticketfile.Color:    handleColor,
		ticketfile.Cut:      handleCut,
		ticketfile.Font:     handleFont,
		ticketfile.Init:     handleInit,
		ticketfile.Lf:       handleLf,
		ticketfile.Print:    handlePrint,
		ticketfile.Printlf:  handlePrintlf,
		ticketfile.Printraw: handlePrintraw,
	}
}

func (c *Converter) Convert(cmd ticketfile.Command) ([]byte, error) {
	if f, ok := dispatchTable[cmd.Name]; ok {
		rawCmd, error := f(c, cmd)

		return []byte(rawCmd), error
	}
	return nil, fmt.Errorf("command %s is not supported", cmd.Name)
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

func handleColor(c *Converter, cmd ticketfile.Command) (string, error) {
	if cmd.Arg == "RED" {
		return "\x1Br1", nil
	}
	return "\x1Br0", nil
}

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
