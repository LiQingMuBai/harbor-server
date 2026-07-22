package main

import (
	bootstrapadmin "cointrade/internal/bootstrap/admin"
	"log"
)

func main() {
	options := bootstrapadmin.OptionsFromEnv()
	log.Printf("==================================================")
	log.Printf("START harbor-server admin")
	log.Printf("service=admin port=%d rpc_clients=%d", options.Port, len(options.RPCClients))
	log.Printf("==================================================")
	bootstrapadmin.Run(options)
}
