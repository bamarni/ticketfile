package main

import (
	"github.com/bamarni/printer"
	"os"
	"log"
)

func main() {
	printer := printer.NewPrinter(os.Stdin, os.Stdout)
	err := printer.Print()
	if err != nil {
		log.Fatal(err)
	}
}
