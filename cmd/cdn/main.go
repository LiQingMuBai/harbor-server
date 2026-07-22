package main

import (
	bootstrapcdn "cointrade/internal/bootstrap/cdn"
	"log"
)

func main() {
	options := bootstrapcdn.OptionsFromEnv()
	log.Printf("==================================================")
	log.Printf("START harbor-server cdn")
	log.Printf("service=cdn port=%d domain=%s", options.Port, options.Domain)
	log.Printf("==================================================")
	bootstrapcdn.Run(options)
}
