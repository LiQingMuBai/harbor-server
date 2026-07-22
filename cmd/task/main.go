package main

import (
	bootstraptask "cointrade/internal/bootstrap/task"
	"log"
	"os"
)

func main() {
	options, err := bootstraptask.OptionsFromArgs(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	if err = bootstraptask.Run(options); err != nil {
		log.Fatal(err)
	}
}
