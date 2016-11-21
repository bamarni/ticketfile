package ticketfile

const (
	undefined CommandType = iota
	Align
	Charset
	Color
	Cut
	Font
	Init
	Lf
	Marginleft
	Units
	Barcode
	Print
	Printlf
	multiline
	Printraw
)

var commands = map[string]CommandType{
	"ALIGN":      Align,
	"CHARSET":    Charset,
	"COLOR":      Color,
	"CUT":        Cut,
	"FONT":       Font,
	"INIT":       Init,
	"LF":         Lf,
	"MARGINLEFT": Marginleft,
	"PRINT":      Print,
	"PRINTLF":    Printlf,
	"PRINTRAW":   Printraw,
	"UNITS":      Units,
	"BARCODE":    Barcode,
}

type (
	CommandType int
	Command     struct {
		Type CommandType
		Arg  string
	}
)

func (cmdType CommandType) isMultiline() bool {
	return multiline < cmdType
}
