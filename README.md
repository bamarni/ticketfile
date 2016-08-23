# printer [![Build Status](https://travis-ci.org/bamarni/printer.svg?branch=master)](https://travis-ci.org/bamarni/printer)

Library to control receipt printers supporting ESC/POS commands.

It takes as input a so-called "Ticketfile", and converts it to bytes to be sent to the device.

## Usage

The most simple way to use it is to create a CLI program :

```go
package main

import (
	"github.com/bamarni/printer"
	"os"
)

func main() {
	printer := printer.NewPrinter(os.Stdout)
	printer.Print(os.Stdin)
}
```

You can then run : `program < /path/to/ticketfile > /path/to/device`

## Ticketfile reference

    # Clears the print buffer / resets modes to their default values
    INIT

    ALIGN ( "LEFT" | "CENTER" | "RIGHT" )

    FONT ( "A" | "B" | "C" )

    # Black is the default color, some models support an additional color (usually red)
    COLOR ( "BLACK" | "RED" )

    # Default is "PC437"
    CHARSET ( "PC437" | "Katakana" | "PC850" | "PC860" | ... )

    PRINT unicode_char { unicode_char }

    # Line feed(s)
    LF { decimal_digit }

    # Cuts paper, default mode is PARTIAL which lefts one point uncut, some models support a FULL cut.
    CUT [ "PARTIAL" | "FULL" ]

A Ticketfile is UTF-8 encoded.
