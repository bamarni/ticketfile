package parser

import (
	"bufio"
	"github.com/bamarni/printer/command"
	"io"
	"regexp"
	"strings"
	"unicode"
)

var (
	tokenWhitespace = regexp.MustCompile(`[\t\v\f\r ]+`)
	tokenComment    = regexp.MustCompile(`^#.*$`)
)

func Parse(r io.Reader) []command.Command {
	var cmds []command.Command

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimLeftFunc(scanner.Text(), unicode.IsSpace)

		if tokenComment.MatchString(line) {
			line = tokenComment.ReplaceAllString(line, "")
		}

		if line == "" {
			continue
		}

		cmdSplits := tokenWhitespace.Split(line, 2)
		cmd := command.Command{
			Raw:  line,
			Name: cmdSplits[0],
		}
		if len(cmdSplits) == 2 {
			cmd.Arg = cmdSplits[1]
		}

		cmds = append(cmds, cmd)
	}

	return cmds
}
