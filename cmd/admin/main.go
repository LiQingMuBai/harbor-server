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
	log.Printf("==================================================")
	log.Printf("START harbor-server admin")
	log.Printf("service=admin port=%d rpc_clients=%d", options.Port, len(options.RPCClients))
	log.Printf("==================================================")
	bootstrapadmin.Run(options)
}
