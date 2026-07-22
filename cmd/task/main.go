package main

import (
	bootstraptask "cointrade/internal/bootstrap/task"
	"cointrade/utils"
	"log"
	"os"
)

func main() {
	options, err := bootstraptask.OptionsFromArgs(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	if err = utils.SetupServiceLogger("task/" + options.Mode); err != nil {
		log.Fatal(err)
	}
	log.Printf("==================================================")
	log.Printf("START harbor-server task")
	log.Printf("service=task mode=%s", options.Mode)
	log.Printf("==================================================")
	if err = bootstraptask.Run(options); err != nil {
		log.Fatal(err)
	}
}
