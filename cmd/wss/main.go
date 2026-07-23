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
	utils.ServiceStartupBanner("harbor-server wss", "service", "wss", "port", options.Port)
	if err := bootstrapwss.Run(options); err != nil {
		log.Fatal(err)
	}
}
