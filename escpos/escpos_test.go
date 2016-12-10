package escpos

import (
	"fmt"
	"testing"

	"github.com/bamarni/ticketfile"
)

var expectedEscpos = []struct {
	cmdType ticketfile.CommandType
	cmdArg  string
	escpos  string
}{
	// INIT
	{ticketfile.Init, "", "\x1B@"},

	// PRINT, LF, PRINTLF, PRINTRAW
	{ticketfile.Print, "test", "test"},
	{ticketfile.Lf, "", "\n"},
	{ticketfile.Lf, "12", fmt.Sprintf("\x1Bd%c", 12)},
	{ticketfile.Printlf, "test", "test\n"},
	{ticketfile.Printraw, "test\ntest2\n", "test\ntest2\n"},

	// UNITS
	{ticketfile.Units, "5 10", fmt.Sprintf("\x1DP%c%c", 5, 10)},

	// MARGINLEFT
	{ticketfile.Marginleft, "500", fmt.Sprintf("\x1DL%c%c", 244, 1)},
	{ticketfile.Marginleft, "1024", fmt.Sprintf("\x1DL%c%c", 0, 4)},

	// ALIGN
	{ticketfile.Align, "LEFT", fmt.Sprintf("\x1Ba%c", 0)},
	{ticketfile.Align, "CENTER", fmt.Sprintf("\x1Ba%c", 1)},
	{ticketfile.Align, "RIGHT", fmt.Sprintf("\x1Ba%c", 2)},

	// CUT
	{ticketfile.Cut, "", fmt.Sprintf("\x1DV%c", 0)},
	{ticketfile.Cut, "PARTIAL", fmt.Sprintf("\x1DV%c", 1)},
	{ticketfile.Cut, "FULL", fmt.Sprintf("\x1DV%c", 0)},

	// FONT
	{ticketfile.Font, "A", "\x1BM0"},
	{ticketfile.Font, "B", "\x1BM1"},
	{ticketfile.Font, "C", "\x1BM2"},

	// COLOR
	{ticketfile.Color, "BLACK", "\x1Br0"},
	{ticketfile.Color, "RED", "\x1Br1"},

	// BARCODE
	{ticketfile.Barcode, "PRINT CODE39 AZERTY123", fmt.Sprintf("\x1Dk%cAZERTY123\x00", 4)},
	{ticketfile.Barcode, "WIDTH 50", fmt.Sprintf("\x1Dw%c", 50)},
	{ticketfile.Barcode, "HEIGHT 10", fmt.Sprintf("\x1Dh%c", 10)},
	{ticketfile.Barcode, "HRI FONT A", fmt.Sprintf("\x1Df%c", 0)},
	{ticketfile.Barcode, "HRI FONT B", fmt.Sprintf("\x1Df%c", 1)},
	{ticketfile.Barcode, "HRI DISPLAY TOP", fmt.Sprintf("\x1DH%c", 1)},
}

func TestConvert(t *testing.T) {
	conv := NewConverter()
	for _, exp := range expectedEscpos {
		cmd, _ := ticketfile.NewCommand(exp.cmdType, exp.cmdArg)
		escpos, err := conv.Convert(cmd)
		if err != nil {
			t.Error("unexpected error")
		}
		if string(escpos) != exp.escpos {
			t.Error("unexpected ESC/POS command %s, expected %s", escpos, []byte(exp.escpos))
		}
	}
}
