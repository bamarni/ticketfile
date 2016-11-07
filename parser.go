package ticketfile

import (
	"bufio"
	"io"
	"regexp"
	"strings"
	"unicode"
)

const (
	Align    = "ALIGN"
	Charset  = "CHARSET"
	Color    = "COLOR"
	Cut      = "CUT"
	Font     = "FONT"
	Init     = "INIT"
	Lf       = "LF"
	Print    = "PRINT"
	Printlf  = "PRINTLF"
	Printraw = "PRINTRAW"
)

var (
	tokenWhitespace   = regexp.MustCompile(`[\t\v\f\r ]+`)
	tokenMultilineEnd = ">>>"
)

type context struct {
	raw       string
	cmd       string
	arg       string
	multiline bool
}

func parse(r io.Reader) <-chan Command {
	cmds := make(chan Command)

	go func() {
		ctx := &context{}
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			line := scanner.Text()

			if ctx.multiline {
				if line != tokenMultilineEnd {
					ctx.raw = ctx.raw + line + "\n"
					ctx.arg = ctx.arg + line + "\n"
					continue
				}
				ctx.raw = ctx.raw + tokenMultilineEnd + "\n"
			} else {
				line = trimLine(line)
				if line == "" {
					continue
				}
				cmdSplits := tokenWhitespace.Split(line, 2)
				ctx.raw = line

				// TODO : this should throw a parse error
				// in case the command doesn't exist
				ctx.cmd = cmdSplits[0]

				// TODO : handle in a more generic way
				if ctx.cmd == "PRINTRAW" {
					ctx.multiline = true
					continue
				}
				if len(cmdSplits) == 2 {
					ctx.arg = cmdSplits[1]
				}
			}

			cmds <- Command{
				Raw:  ctx.raw,
				Name: ctx.cmd,
				Arg:  ctx.arg,
			}
			ctx = &context{}
		}

		close(cmds)
	}()

	return cmds
}

func trimLine(line string) string {
	line = strings.TrimLeftFunc(line, unicode.IsSpace)

	// comment line
	if line != "" && string(line[0]) == "#" {
		return ""
	}

	return line
}
