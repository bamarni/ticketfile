package ticketfile

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
	"unicode"
)

var (
	tokenWhitespace   = regexp.MustCompile(`[\t\v\f\r ]+`)
	tokenMultilineEnd = ">>>"
)

func parse(r io.Reader) ([]Command, error) {
	var (
		cmdName string
		cmdType CommandType
		cmdArg  string
		cmds    []Command
	)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		if cmdType.isMultiline() {
			if line != tokenMultilineEnd {
				cmdArg = cmdArg + line + "\n"
				continue
			}
		} else {
			line = strings.TrimLeftFunc(line, unicode.IsSpace)
			if line == "" || string(line[0]) == "#" {
				continue
			}
			cmdSplits := tokenWhitespace.Split(line, 2)
			cmdName = cmdSplits[0]
			cmdType = commands[cmdName]
			if cmdType == undefined {
				return nil, fmt.Errorf("undefined command %s", cmdName)
			}
			if cmdType.isMultiline() {
				continue
			}
			if len(cmdSplits) == 2 {
				cmdArg = cmdSplits[1]
			}
		}
		cmds = append(cmds, Command{Type: cmdType, Arg: cmdArg})
		cmdName = ""
		cmdType = undefined
		cmdArg = ""
	}
	return cmds, nil
}
