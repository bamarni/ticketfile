package main

import (
	"github.com/bamarni/printer"
	"os"
	"log"
)

func main() {
	printer := printer.NewPrinter(os.Stdout)
	err := printer.Print(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
}
