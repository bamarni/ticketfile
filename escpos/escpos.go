package escpos

import (
	"fmt"
	"github.com/bamarni/printer/command"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"strconv"
)

type Escpos struct {
	enc *encoding.Encoder
}

func NewEscpos() *Escpos {
	escpos := Escpos{
		enc: charmap.CodePage437.NewEncoder(),
	}

	return &escpos
}

var dispatchTable map[string]func(*Escpos, command.Command) (string, error)

func init() {
	dispatchTable = map[string]func(*Escpos, command.Command) (string, error){
		command.Align:   handleAlign,
		command.Charset: handleCharset,
		command.Color:   handleColor,
		command.Cut:     handleCut,
		command.Font:    handleFont,
		command.Init:    handleInit,
		command.Lf:      handleLf,
		command.Print:   handlePrint,
	}
}

func (e *Escpos) ToBytes(cmd command.Command) ([]byte, error) {
	if f, ok := dispatchTable[cmd.Name]; ok {
		rawCmd, error := f(e, cmd)

		return []byte(rawCmd), error
	}
	return nil, fmt.Errorf("Command %s is not supported", cmd.Name)
}

func handleAlign(e *Escpos, cmd command.Command) (string, error) {
	switch cmd.Arg {
	case "LEFT":
		return "\x1Ba0", nil
	case "CENTER":
		return "\x1Ba1", nil
	case "RIGHT":
		return "\x1Ba2", nil
	}
	return "", fmt.Errorf("Unsupported alignment %s", cmd.Arg)
}

func handleCut(e *Escpos, cmd command.Command) (string, error) {
	switch cmd.Arg {
	case "FULL":
		return "\x1DVA0", nil
	case "PARTIAL", "":
		return "\x1DVA1", nil
	}
	return "", fmt.Errorf("Unsupported cut %s", cmd.Arg)
}

func handleFont(e *Escpos, cmd command.Command) (string, error) {
	switch cmd.Arg {
	case "A":
		return "\x1BM0", nil
	case "B":
		return "\x1BM1", nil
	case "C":
		return "\x1BM2", nil
	}
	return "", fmt.Errorf("Unsupported font %s", cmd.Arg)
}

func handleColor(e *Escpos, cmd command.Command) (string, error) {
	if cmd.Arg == "RED" {
		return "\x1Br1", nil
	}
	return "\x1Br0", nil
}

func handleCharset(e *Escpos, cmd command.Command) (string, error) {
	var n byte
	switch cmd.Arg {
	case "PC437": // USA: Standard Europe
		e.enc = charmap.CodePage437.NewEncoder()
		n = 0
	case "PC850": // Western Europe
		e.enc = charmap.CodePage850.NewEncoder()
		n = 2
	default:
		return "", fmt.Errorf("Charset %s not supported", cmd.Arg)
	}
	return fmt.Sprintf("\x1Bt%c", n), nil
}

func handleInit(e *Escpos, cmd command.Command) (string, error) {
	return "\x1B@", nil
}

func handleLf(e *Escpos, cmd command.Command) (string, error) {
	if cmd.Arg != "" {
		nb, err := strconv.Atoi(cmd.Arg)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("\x1Bd%c", nb), nil
	}
	return "\n", nil
}

func handlePrint(e *Escpos, cmd command.Command) (string, error) {
	rawCmd, err := e.enc.String(cmd.Arg)
	if err != nil {
		return "", fmt.Errorf("Couldn't encode to charset (%s)", err)
	}
	return rawCmd, nil
}
