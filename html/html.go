package html

import (
	"strings"

	"github.com/bamarni/ticketfile"
)

type Converter struct {
	enc   *strings.Replacer
	align string
	color string
}

func NewConverter() *Converter {
	c := &Converter{
		enc: strings.NewReplacer(" ", "&nbsp;", "\n", "<br />"),
	}
	c.init()
	return c
}

func (c *Converter) Convert(cmd ticketfile.Command) ([]byte, error) {
	switch cmd.Type {
	case ticketfile.Align:
		c.align = strings.ToLower(cmd.Arg)
	case ticketfile.Color:
		c.color = strings.ToLower(cmd.Arg)
	case ticketfile.Cut:
		return []byte("<hr />"), nil
	case ticketfile.Init:
		c.init()
	case ticketfile.Print, ticketfile.Printraw:
		return c.wrapDiv(cmd.Arg), nil
	case ticketfile.Printlf:
		return c.wrapDiv(cmd.Arg + "\n"), nil
	case ticketfile.Lf:
		return []byte("<br />"), nil
	}
	return nil, nil
}

func (c *Converter) init() {
	c.align = "left"
	c.color = "black"
}

func (c *Converter) wrapDiv(s string) []byte {
	b := []byte("<div style=\"text-align:" + c.align + ";color:" + c.color + "\">")
	return append(append(b, c.enc.Replace(s)...), "</div>"...)
}
