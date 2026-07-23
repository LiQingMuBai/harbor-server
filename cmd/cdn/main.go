package main

import (
	bootstrapcdn "cointrade/internal/bootstrap/cdn"
	"cointrade/utils"
	"log"
)

func main() {
	if err := utils.SetupServiceLogger("cdn"); err != nil {
		log.Fatal(err)
	}
	options, err := bootstrapcdn.OptionsFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("==================================================")
	log.Printf("START harbor-server cdn")
	log.Printf("service=cdn port=%d domain=%s", options.Port, options.Domain)
	log.Printf("==================================================")
	if err := bootstrapcdn.Run(options); err != nil {
		log.Fatal(err)
	}
}
