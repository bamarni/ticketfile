package escpos

import (
	"fmt"
	"github.com/bamarni/printer/command"
)

func ToBytes(cmd command.Command) ([]byte, error) {
	var rawCmd string

	// http://content.epson.de/fileadmin/content/files/RSD/downloads/escpos.pdf
	switch cmd.Name {
	case "INIT":
		rawCmd = "\x1B@"
	case "WRITE":
		rawCmd = cmd.Arg
	case "LF":
		if cmd.Arg != "" {
			rawCmd = fmt.Sprintf("\x1Bd%s", cmd.Arg)
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
