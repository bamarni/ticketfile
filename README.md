# ticketfile [![Build Status](https://travis-ci.org/bamarni/ticketfile.svg?branch=master)](https://travis-ci.org/bamarni/ticketfile)

- [What is a Ticketfile?](#what-is-a-ticketfile)
- [Why should I use Ticketfiles?](#why-should-i-use-ticketfiles)
- [Ticketfile specification](#ticketfile-specification)
- [Golang library usage](#golang-library-usage)
- [Contributing](#contributing)

## What is a Ticketfile?

A Ticketfile is a textual representation of a thermal printer receipt.

It allows you to write simple readable text for your receipts rather than having to deal for example with
[ESC/POS](http://content.epson.de/fileadmin/content/files/RSD/downloads/escpos.pdf)
specification directly.

## Why should I use Ticketfiles?

Maybe you've found an ESC/POS library for your favorite language out there
and you're wondering why you should use Ticketfiles instead?

Even though they're inspired by ESC/POS specification, Ticketfiles are readable 
([here is a Ticketfile](tests/functional/fixtures/Ticketfile) and its [ESC/POS equivalent](tests/functional/fixtures/Ticketfile.expected)) and **manufacturer agnostic**.
In the future if a new standard emerges we'll do our best to support it without changing the spec in a breaking way. 

If you use a library you'll also be tied to a language, Ticketfiles are **language agnostic**. 
We provide an official Golang library but you're free to write your own, the format is easily parsable.

Finally, Ticketfiles are **context agnostic** and not limited to receipt printers.
You could for instance convert a Ticketfile to HTML so it can be displayed in a browser or sent as e-mail attachment.

## Ticketfile specification

The following specification use the [Extended Backus-Naur Form](https://en.wikipedia.org/wiki/Extended_Backus%E2%80%93Naur_form).

    (* Below is a ticketfile comment *)
    "#"{ unicode_char }

    (* Clears the print buffer / resets modes to their default values *)
    "INIT"

    "ALIGN" ( "LEFT" | "CENTER" | "RIGHT" )

    "FONT" ( "A" | "B" | "C" )

    (* Black is the default color, some models support an additional color (usually red) *)
    "COLOR" ( "BLACK" | "RED" )

    (* Default is "PC437" *)
    "CHARSET" ( "PC437" | "Katakana" | "PC850" | "PC860" | ... )

    "PRINT" unicode_char { unicode_char }

    (* Line feed(s) *)
    "LF" { decimal_digit }

    "PRINTLF" unicode_char { unicode_char }

    (* Prints a multiline raw block *)
    "PRINTRAW"
    { unicode_char | "\n" }
    "```"

    (* Cuts paper, default mode is PARTIAL which lefts one point uncut, some models support a FULL cut. *)
    "CUT" [ "PARTIAL" | "FULL" ]

A Ticketfile is UTF-8 encoded. For reference, here is a [Ticketfile example](tests/functional/fixtures/Ticketfile).

## Golang library usage

The Golang library contains a parser and various converters.

The most simple way to use it is to create a CLI program :

```go
package main

import (
	"github.com/bamarni/ticketfile"
	"github.com/bamarni/ticketfile/escpos"
	"os"
)

func main() {
	engine := ticketfile.NewEngine(os.Stdout, escpos.NewConverter())
	if err := engine.Render(os.Stdin); err != nil {
		log.Fatal(err)
	}
}
```
The following command would convert a Ticketfile into ESC/POS commands and send them to a printer :

    program < /path/to/ticketfile > /path/to/device

*In case of a syntax error in your ticketfile, a message would be displayed to stderr while nothing would be sent to the printer device.*

## Contributing

Contributions and new ideas are always welcome.

The Golang library is a first-class citizen so any change in the specification needs to be reflected in the Go engine. 

Released under the MIT license.
