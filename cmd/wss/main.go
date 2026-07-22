package main

import (
	bootstrapwss "cointrade/internal/bootstrap/wss"
	"log"
)

func main() {
	options := bootstrapwss.OptionsFromEnv()
	log.Printf("==================================================")
	log.Printf("START harbor-server wss")
	log.Printf("service=wss port=%d", options.Port)
	log.Printf("==================================================")
	bootstrapwss.Run(options)
}
