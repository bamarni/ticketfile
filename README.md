# ticketfile [![Build Status](https://travis-ci.org/bamarni/ticketfile.svg?branch=master)](https://travis-ci.org/bamarni/ticketfile)

- [What is a TicketFile?](#what-is-a-ticketfile)
- [Ticketfile specification](#ticketfile-specification)
- [Golang library usage](#golang-library-usage)
- [Contributing](#contributing)

## What is a Ticketfile?

A Ticketfile is a textual representation of a thermal printer receipt.

This allows you to write simple readable text for your receipts rather than having to deal with ESC/POS commands directly, which are binary specification.

It can also be converted into different representations. For instance, you could convert the same Ticketfile into ESC/POS commands to be sent to a retail shop's thermal printer, or convert it to an HTML view to be displayed in a browser.

## Ticketfile specification

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

## Golang library usage

The Golang library contains a parser and various converters.

The most simple way to use it is to create a CLI program :

```go
package main

import (
	"github.com/bamarni/ticketfile"
	"os"
)

func main() {
	engine := ticketfile.NewEngine(os.Stdout, escpos.NewConverter())
	if err := engine.Render(os.Stdin); err != nil {
		log.Fatal(err)
	}
}
```
The following command would convert a Ticketfile into ESC/POS commands :

    program < /path/to/ticketfile > /path/to/device

## Contributing

This package is under development and not yet stable, the main features to come are :
- specification for bar codes and images
- an html converter

Contributions and other ideas are also welcome.

The Golang library is a first-class citizen so any change in the specification needs to be reflected in the Go engine. 

Released under the MIT license.
