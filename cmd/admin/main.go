package main

import (
	bootstrapadmin "cointrade/internal/bootstrap/admin"
	"cointrade/utils"
	"log"
)

func main() {
	if err := utils.SetupServiceLogger("admin"); err != nil {
		log.Fatal(err)
	}
	options := bootstrapadmin.OptionsFromEnv()
	utils.ServiceStartupBanner("harbor-server admin", "service", "admin", "port", options.Port, "rpc_clients", len(options.RPCClients))
	if err := bootstrapadmin.Run(options); err != nil {
		log.Fatal(err)
	}
}
