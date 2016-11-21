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
	{ticketfile.Command{Type: ticketfile.Init}, "\x1B@"},

	// PRINT, LF, PRINTLF, PRINTRAW
	{ticketfile.Command{Type: ticketfile.Print, Arg: "test"}, "test"},
	{ticketfile.Command{Type: ticketfile.Lf}, "\n"},
	{ticketfile.Command{Type: ticketfile.Lf, Arg: "12"}, fmt.Sprintf("\x1Bd%c", 12)},
	{ticketfile.Command{Type: ticketfile.Printlf, Arg: "test"}, "test\n"},
	{ticketfile.Command{Type: ticketfile.Printraw, Arg: "test\ntest2\n"}, "test\ntest2\n"},

	// UNITS
	{ticketfile.Command{Type: ticketfile.Units, Arg: "5 10"}, fmt.Sprintf("\x1DP%c%c", 5, 10)},

	// MARGINLEFT
	{ticketfile.Command{Type: ticketfile.Marginleft, Arg: "500"}, fmt.Sprintf("\x1DL%c%c", 244, 1)},
	{ticketfile.Command{Type: ticketfile.Marginleft, Arg: "1024"}, fmt.Sprintf("\x1DL%c%c", 0, 4)},

	// ALIGN
	{ticketfile.Command{Type: ticketfile.Align, Arg: "LEFT"}, "\x1Ba0"},
	{ticketfile.Command{Type: ticketfile.Align, Arg: "CENTER"}, "\x1Ba1"},
	{ticketfile.Command{Type: ticketfile.Align, Arg: "RIGHT"}, "\x1Ba2"},

	// CUT
	{ticketfile.Command{Type: ticketfile.Cut}, "\x1DVA1"},
	{ticketfile.Command{Type: ticketfile.Cut, Arg: "PARTIAL"}, "\x1DVA1"},
	{ticketfile.Command{Type: ticketfile.Cut, Arg: "FULL"}, "\x1DVA0"},

	// FONT
	{ticketfile.Command{Type: ticketfile.Font, Arg: "A"}, "\x1BM0"},
	{ticketfile.Command{Type: ticketfile.Font, Arg: "B"}, "\x1BM1"},
	{ticketfile.Command{Type: ticketfile.Font, Arg: "C"}, "\x1BM2"},

	// COLOR
	{ticketfile.Command{Type: ticketfile.Color, Arg: "BLACK"}, "\x1Br0"},
	{ticketfile.Command{Type: ticketfile.Color, Arg: "RED"}, "\x1Br1"},

	// BARCODE
	{ticketfile.Command{Type: ticketfile.Barcode, Arg: "PRINT CODE39 AZERTY123"}, fmt.Sprintf("\x1Dk%cAZERTY123\x00", 4)},
}

func TestConvert(t *testing.T) {
	conv := NewConverter()
	for _, exp := range expectedEscpos {
		escpos, err := conv.Convert(exp.command)
		if err != nil {
			t.Error("unexpected error")
		}
		if string(escpos) != exp.escpos {
			t.Error("unexpected ESC/POS command %s, expected %s", escpos, []byte(exp.escpos))
		}
	}
}
