package ticketfile

import (
	"reflect"
	"strings"
	"testing"
)

var (
	ticketfile = `
INIT

PRINTLF print this

# This is a comment

ALIGN RIGHT
PRINTRAW 
This is
multiline
>>>

	`
	expectedCommands = []Command{
		{Type: Init},
		{Type: Printlf, Arg: "print this"},
		{Type: Align, Arg: "RIGHT"},
		{Type: Printraw, Arg: "This is\nmultiline\n"},
	}
)

func TestParse(t *testing.T) {
	cmds, err := parse(strings.NewReader(ticketfile))
	if err != nil {
		t.Errorf("unexpected error : %s", err)
	}
	if !reflect.DeepEqual(expectedCommands, cmds) {
		t.Error("unexpected commands")
	}
}
