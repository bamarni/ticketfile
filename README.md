# ticketfile [![Build Status](https://travis-ci.org/bamarni/ticketfile.svg?branch=master)](https://travis-ci.org/bamarni/ticketfile)

- [What is a Ticketfile?](#what-is-a-ticketfile)
- [Why should I use Ticketfiles?](#why-should-i-use-ticketfiles)
- [Status of the spec](#status-of-the-spec)
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

If you use a library you'll also be tied to a language, Ticketfiles are **language agnostic**, it's just text. 
We provide an official Golang library but you're free to write your own, the format is easily parsable.

Finally, even though they mainly target receipt printers, Ticketfiles are **context agnostic**.
For instance, in the future our Golang library will provide an HTML converter so that a Ticketfile could be displayed in a browser or sent as e-mail attachment.

## Status of the spec

At the moment specification are still subject to change based on usage feedbacks.
They'll get more stable after a few months.

## Ticketfile specification

A Ticketfile is Unicode text encoded in UTF-8, it contains a set of commands.

Those commands allow you to write to the receipt, cut the paper, define styles, etc.

Before going into the specification, here is a first impression of how a Ticketfile looks like :

    INIT
    
    ALIGN CENTER
    PRINTRAW
    My Shop
    Fifth Avenue
    New York, NY 10020
    >>>
    
    ALIGN LEFT
    FONT B
    PRINTLF Invoice n. 456
    PRINTLF John Smith
    FONT A
    
    ALIGN RIGHT
    PRINTRAW
    8.00
    15.90
    ===
    23.90
    >>>
    
    ALIGN CENTER
    PRINTLF Thank you for your visit!
    CUT
    
    

### Notation

The present specification use the Extended Backus-Naur Form.
More precisely, the exact syntax follows the [Golang specification notation](https://golang.org/ref/spec#Notation).

### Comments

Ticketfile comments are single-line :

``` ebnf
comment = "#" { unicode_char } .
```

### INIT

The INIT command clears the print buffer and resets modes to their default values.

It basically sets the printer in the same state as it would be right after powering it up.

It should typically be at the beginning of a Ticketfile, so that states and styles from
previous Ticketfiles are discarded.

``` ebnf
init_command = "INIT" .
```

### PRINT, LF, PRINTLF and PRINTRAW

Those commands print text and control line feeds.

``` ebnf
print_command    = "PRINT" unicode_char { unicode_char } .
lf_command       = "LF" { decimal_lit } .
printlf_command  = "PRINTLF" unicode_char { unicode_char } .
printraw_command = "PRINTRAW" newline { unicode_char | "\n" } newline ">>>" .
```

For example :

    # First row
    PRINT Hello
    LF 2

    # Third row
    PRINTLF world

    PRINTRAW
    This text can contain newlines.
    This is useful if for instance you're using templating on top of Ticketfiles
    and have multi-line variables to display.
    >>>

### ALIGN

Set alignement for the text.

``` ebnf
align_command = "ALIGN" ( "LEFT" | "CENTER" | "RIGHT" ) .
```

Example :

    ALIGN RIGHT
    PRINTLF To the right!

    ALIGN CENTER
    PRINTLF This is centered.

### UNITS (likely to change)

Sets the vertical and horizontal motion units.

``` ebnf
byte_decimal_lit = "0" … "255" .
units_command   = "UNITS" byte_decimal_lit byte_decimal_lit .
```

The first argument corresponds to the horizontal motion unit, the second to the vertical motion unit. Those units are used for print position related commands, such as `MARGINLEFT`.

The resulting motion in inches is the multiplicative inverse of the provided value. For example :

    UNITS 2 0

Here the horizontal motion unit would be `1 / 2 inches` (approximately 12.7 mm).
The zero value for the vertical motion unit indicates that it should use the printer's default.

*This command will most likely change in the future, as it's currently not abstracted from ESC/POS.*

### MARGINLEFT (likely to change)

Sets left margin.

``` ebnf
two_bytes_decimal_lit = "0" … "65535" .
margin_command   = "MARGINLEFT" two_bytes_decimal_lit .
```

The argument is a number from 0 to 65535, the actual resulting margin is :
`margin * horizontal_motion_unit`.

Example :

    UNITS 2 0
    MARGINLEFT 3

Here the horizontal motion unit is 2, hence the margin would be `3 * 1 / 2` inches (approximately 38.1 mm).

*This command will most likely change in the future, as it's currently not abstracted from ESC/POS.*

### FONT

Sets font style.

``` ebnf
font_command = "FONT" ( "A" | "B" | "C" ) .
```

A is the default one, B is usually smaller text.

### COLOR

Sets text color. Black is the default color, some models support an additional color (usually red).

``` ebnf
color_command = "COLOR" ( "BLACK" | "RED" ) .
```

### CHARSET

Sets the charset. The default is PC437 (USA: Standard Europe), PC850 is for Western Europe.

``` ebnf
charset_command = "CHARSET" ( "PC437" | "PC850" ) .
```

### CUT

Cuts paper, default mode is PARTIAL which lefts one point uncut, some models support a FULL cut.

``` ebnf
cut_command = "CUT" [ "PARTIAL" | "FULL" ] .
```

## Golang library usage

The Golang library contains a parser and various converters.

The most simple way to use it is to create a CLI program :

```go
package main

import (
	"os"

	"github.com/bamarni/ticketfile"
	"github.com/bamarni/ticketfile/escpos"
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
