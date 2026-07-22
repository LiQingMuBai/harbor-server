package main

import (
	bootstrapwss "cointrade/internal/bootstrap/wss"
	"cointrade/utils"
	"log"
)

func main() {
	if err := utils.SetupServiceLogger("wss"); err != nil {
		log.Fatal(err)
	}
	options := bootstrapwss.OptionsFromEnv()
	log.Printf("==================================================")
	log.Printf("START harbor-server wss")
	log.Printf("service=wss port=%d", options.Port)
	log.Printf("==================================================")
	bootstrapwss.Run(options)
}
