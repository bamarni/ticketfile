package escpos

import (
	"fmt"
	"testing"

	"github.com/bamarni/ticketfile"
)

var expectedEscpos = []struct {
	command ticketfile.Command
	escpos  string
}{
	// INIT
	{ticketfile.Command{Name: "INIT"}, "\x1B@"},

	// PRINT, LF, PRINTLF, PRINTRAW
	{ticketfile.Command{Name: "PRINT", Arg: "test"}, "test"},
	{ticketfile.Command{Name: "LF"}, "\n"},
	{ticketfile.Command{Name: "LF", Arg: "12"}, fmt.Sprintf("\x1Bd%c", 12)},
	{ticketfile.Command{Name: "PRINTLF", Arg: "test"}, "test\n"},
	{ticketfile.Command{Name: "PRINTRAW", Arg: "test\ntest2\n"}, "test\ntest2\n"},

	// MARGINLEFT
	{ticketfile.Command{Name: "MARGINLEFT", Arg: "500"}, fmt.Sprintf("\x1DL%c%c", 244, 1)},
	{ticketfile.Command{Name: "MARGINLEFT", Arg: "1024"}, fmt.Sprintf("\x1DL%c%c", 0, 4)},

	// ALIGN
	{ticketfile.Command{Name: "ALIGN", Arg: "LEFT"}, "\x1Ba0"},
	{ticketfile.Command{Name: "ALIGN", Arg: "CENTER"}, "\x1Ba1"},
	{ticketfile.Command{Name: "ALIGN", Arg: "RIGHT"}, "\x1Ba2"},

	// CUT
	{ticketfile.Command{Name: "CUT"}, "\x1DVA1"},
	{ticketfile.Command{Name: "CUT", Arg: "PARTIAL"}, "\x1DVA1"},
	{ticketfile.Command{Name: "CUT", Arg: "FULL"}, "\x1DVA0"},

	// FONT
	{ticketfile.Command{Name: "FONT", Arg: "A"}, "\x1BM0"},
	{ticketfile.Command{Name: "FONT", Arg: "B"}, "\x1BM1"},
	{ticketfile.Command{Name: "FONT", Arg: "C"}, "\x1BM2"},

	// COLOR
	{ticketfile.Command{Name: "COLOR", Arg: "BLACK"}, "\x1Br0"},
	{ticketfile.Command{Name: "COLOR", Arg: "RED"}, "\x1Br1"},
}

func TestConvert(t *testing.T) {
	conv := NewConverter()
	for _, exp := range expectedEscpos {
		escpos, err := conv.Convert(exp.command)
		if err != nil {
			t.Error("unexpected error")
		}
		if string(escpos) != exp.escpos {
			t.Error("unexpected ESC/POS command %s", escpos)
		}
	}
}
