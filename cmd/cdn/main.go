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
	options := bootstrapcdn.OptionsFromEnv()
	log.Printf("==================================================")
	log.Printf("START harbor-server cdn")
	log.Printf("service=cdn port=%d domain=%s", options.Port, options.Domain)
	log.Printf("==================================================")
	bootstrapcdn.Run(options)
}
