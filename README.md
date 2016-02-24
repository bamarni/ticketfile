# printer

CLI program for ESC/POS printers.

It takes as input a so-called "Ticketfile", and converts it to bytes to be sent to the device.

## Ticketfile reference

    # Clears the print buffer / resets modes to their default values
    INIT

    ALIGN ( "LEFT" | "CENTER" | "RIGHT" )

    FONT ( "A" | "B" | "C" )

    # Black is the default color, some models support an additional color (usually red)
    COLOR ( "BLACK" | "RED" )

    # Default code table is https://en.wikipedia.org/wiki/Code_page_437
    PRINT ascii_char { ascii_char }

    # Line feed(s)
    LF { decimal_digit }

    # Cut paper
    CUT [ "PARTIAL" | "FULL" ]
