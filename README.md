# ticketfile [![Build Status](https://travis-ci.org/bamarni/ticketfile.svg?branch=master)](https://travis-ci.org/bamarni/ticketfile)

- [What is a Ticketfile?](#what-is-a-ticketfile)
- [Why would I use Ticketfiles?](#why-would-i-use-ticketfiles)
- [Ticketfile specification](#ticketfile-specification)
- [Golang library usage](#golang-library-usage)
- [Contributing](#contributing)

## What is a Ticketfile?

A Ticketfile is a textual representation of a thermal printer receipt.

It allows you to write simple readable text for your receipts rather than having to deal for example with
[ESC/POS](http://content.epson.de/fileadmin/content/files/RSD/downloads/escpos.pdf)
specification directly.

## Why would I use Ticketfiles?

Maybe you've found an ESC/POS library for your favorite language out there
and you're wondering why you should use Ticketfiles instead?

Even though they're inspired by ESC/POS specification, Ticketfiles are readable 
([here is a Ticketfile](tests/functional/fixtures/Ticketfile) and its [ESC/POS equivalent](tests/functional/fixtures/Ticketfile.expected)) and **manufacturer agnostic**.
In the future if a new standard emerges we'll do our best to support it without changing the spec in a breaking way. 

If you use a library you'll also be tied to a language, Ticketfiles are **language agnostic**, it's just text. 
We provide an official Golang library but you're free to write your own, the format is easily parsable.

Finally, even though they mainly target receipt printers, Ticketfiles are **context agnostic**.
For instance, in the future our Golang library will provide an HTML converter so that a Ticketfile could be displayed in a browser or sent as e-mail attachment.

## Ticketfile specification

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

Full specification are available [here](spec/spec.md).

## Golang library usage

The Golang library contains a parser and various converters.

The most simple way to use it is through the `ticket` command :

    go get -u github.com/bamarni/ticketfile/cmd/ticket

The following command would then convert a Ticketfile into ESC/POS commands and send them to a printer :

    ticket < /path/to/ticketfile > /path/to/device

*In case of a syntax error in your ticketfile, a message would be displayed to stderr while nothing would be sent to the printer device.*

## Contributing

Contributions and new ideas are always welcome.

The Golang library is a first-class citizen so any change in the specification needs to be reflected in the Go engine. 

Released under the MIT license.
