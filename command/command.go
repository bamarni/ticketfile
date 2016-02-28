package command

import (
	"bufio"
	"io"
	"regexp"
	"strings"
	"unicode"
)

const (
	Align   = "ALIGN"
	Charset = "CHARSET"
	Color   = "COLOR"
	Cut     = "CUT"
	Font    = "FONT"
	Init    = "INIT"
	Lf      = "LF"
	Print   = "PRINT"
)

type Command struct {
	Raw  string
	Name string
	Arg  string
}

var (
	tokenWhitespace = regexp.MustCompile(`[\t\v\f\r ]+`)
)

func Parse(r io.Reader, cmds chan Command) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimLeftFunc(scanner.Text(), unicode.IsSpace)

		if line != "" && string(line[0]) == "#" {
			line = ""
		}

		if line == "" {
			continue
		}

		cmdSplits := tokenWhitespace.Split(line, 2)
		cmd := Command{
			Raw:  line,
			Name: cmdSplits[0],
		}
		if len(cmdSplits) == 2 {
			cmd.Arg = cmdSplits[1]
		}

		cmds <- cmd
	}

	close(cmds)
}
