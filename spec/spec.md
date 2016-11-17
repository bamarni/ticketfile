# Ticketfile specification

A Ticketfile is Unicode text encoded in UTF-8, it contains a set of commands.

Those commands allow you to write to the receipt, cut the paper, define styles, etc.

Before going into the specification, here is a first impression of how a Ticketfile might look like :

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

    LF
 
    ALIGN CENTER
    PRINTLF Thank you for your visit!
    CUT

*At the moment specification are still subject to change based on usage feedbacks.
They'll get more stable after a few months.*

## Notation

The present specification use the Extended Backus-Naur Form.
More precisely, the syntax follows the [Golang specification notation](https://golang.org/ref/spec#Notation).

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
