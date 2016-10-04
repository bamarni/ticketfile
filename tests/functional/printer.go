package main

import (
	"github.com/bamarni/ticketfile"
	"github.com/bamarni/ticketfile/escpos"
	"log"
	"os"
)

func main() {
	engine := ticketfile.NewEngine(os.Stdout, escpos.NewConverter())
	if err := engine.Render(os.Stdin); err != nil {
		log.Fatal(err)
	}
}
