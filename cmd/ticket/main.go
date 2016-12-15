package main

import (
	"flag"
	"log"
	"os"

	"github.com/bamarni/ticketfile"
	"github.com/bamarni/ticketfile/escpos"
	"github.com/bamarni/ticketfile/html"
)

func main() {
	htmlFlag := flag.Bool("html", false, "Convert into HTML instead of ESC/POS")
	flag.Parse()

	var conv ticketfile.Converter
	if *htmlFlag {
		conv = html.NewConverter()
	} else {
		conv = escpos.NewConverter()
	}

	engine := ticketfile.NewEngine(os.Stdout, conv)
	if err := engine.Render(os.Stdin); err != nil {
		log.Fatal(err)
	}
}
