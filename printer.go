package main

import (
	"bufio"
	"github.com/bamarni/printer/escpos"
	"github.com/bamarni/printer/parser"
	"log"
	"os"
)

func main() {
	w := bufio.NewWriter(os.Stdout)

	for _, cmd := range parser.Parse(os.Stdin) {
		log.Printf("Received command : [%s]\n", cmd.Raw)

		rawBytes, err := escpos.ToBytes(cmd)
		if err != nil {
			log.Fatal(err)
		}

		w.Write(rawBytes)
	}

	w.Flush()
}
