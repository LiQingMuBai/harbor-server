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
	log.Printf("==================================================")
	log.Printf("START harbor-server api")
	log.Printf("service=api port=%d rpc_port=%d local_ip=%s", options.ServerPort, options.RPCPort, options.LocalIP)
	log.Printf("==================================================")
	bootstrapapi.Run(options)
}
