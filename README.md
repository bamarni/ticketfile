# printer [![Build Status](https://travis-ci.org/bamarni/printer.svg?branch=master)](https://travis-ci.org/bamarni/printer)

CLI program for ESC/POS printers.

It takes as input a so-called "Ticketfile", and converts it to bytes to be sent to the device.

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

    # Cut paper
    CUT [ "PARTIAL" | "FULL" ]

A Ticketfile is UTF-8 encoded.
