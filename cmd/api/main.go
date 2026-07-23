package main

import (
	bootstrapapi "cointrade/internal/bootstrap/api"
	"cointrade/utils"
	"log"
)

func main() {
	if err := utils.SetupServiceLogger("api"); err != nil {
		log.Fatal(err)
	}
	options := bootstrapapi.OptionsFromEnv()
	utils.ServiceStartupBanner("harbor-server api", "service", "api", "port", options.ServerPort, "rpc_port", options.RPCPort, "local_ip", options.LocalIP)
	if err := bootstrapapi.Run(options); err != nil {
		log.Fatal(err)
	}
}
