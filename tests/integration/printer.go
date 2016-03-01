package main

import (
	"github.com/bamarni/printer"
	"os"
)

func main() {
	printer := printer.NewPrinter(os.Stdin, os.Stdout)
	printer.Print()
}
