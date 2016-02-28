package main

import (
	"bufio"
	"github.com/bamarni/printer/escpos"
	"github.com/bamarni/printer/command"
	"log"
	"os"
)

func main() {
	w := bufio.NewWriter(os.Stdout)
	driver := escpos.NewEscpos()
	cmds := make(chan command.Command)

	go command.Parse(os.Stdin, cmds)

	for cmd := range cmds {
		log.Printf("Received command : [%s]\n", cmd.Raw)

		rawBytes, err := driver.FromCommand(cmd)
		if err != nil {
			log.Fatalf("Driver error : %s\n", err)
		}

		w.Write(rawBytes)

		if cmd.Name == command.Cut {
			w.Flush()
		}
	}

	w.Flush()
}
