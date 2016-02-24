package escpos

import (
	"fmt"
	"github.com/bamarni/printer/command"
	"strconv"
)

func ToBytes(cmd command.Command) ([]byte, error) {
	var rawCmd string

	// http://content.epson.de/fileadmin/content/files/RSD/downloads/escpos.pdf
	switch cmd.Name {
	case "INIT":
		rawCmd = "\x1B@"
	case "ALIGN":
		switch cmd.Arg {
		case "LEFT":
			rawCmd = "\x1Ba0"
		case "CENTER":
			rawCmd = "\x1Ba1"
		case "RIGHT":
			rawCmd = "\x1Ba2"
		}
	case "FONT":
		switch cmd.Arg {
		case "A":
			rawCmd = "\x1BM0"
		case "B":
			rawCmd = "\x1BM1"
		case "C":
			rawCmd = "\x1BM2"
		}
	case "COLOR":
		if cmd.Arg == "RED" {
			rawCmd = "\x1Br1"
		} else {
			rawCmd = "\x1Br0"
		}
	case "PRINT":
		rawCmd = cmd.Arg
	case "LF":
		if cmd.Arg != "" {
			nb, _ := strconv.Atoi(cmd.Arg)
			rawCmd = fmt.Sprintf("\x1Bd%c", nb)
		} else {
			rawCmd = "\n"
		}
	case "CUT":
		rawCmd = "\x1DVA0"
	default:
		return nil, fmt.Errorf("Command %s is not supported", cmd)
	}

	return []byte(rawCmd), nil
}
