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
	utils.ServiceStartupBanner("harbor-server cdn", "service", "cdn", "port", options.Port, "domain", options.Domain)
	if err := bootstrapcdn.Run(options); err != nil {
		log.Fatal(err)
	}
}
