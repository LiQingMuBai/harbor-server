package api

import (
	"cointrade/http/httprun"
	"cointrade/internal/bootstrap/shared"
	"cointrade/models"
)

type Options struct {
	ServerPort int
	LocalIP    string
	RPCPort    int
}

func OptionsFromEnv() Options {
	return Options{
		ServerPort: shared.GetenvInt("API_PORT", 9001),
		LocalIP:    shared.Getenv("API_LOCAL_IP", "127.0.0.1"),
		RPCPort:    shared.GetenvInt("API_RPC_PORT", 9010),
	}
}

func Run(options Options) {
	models.InitData()
	httprun.StartAPIBackgroundJobs(options.LocalIP, options.RPCPort)
	httpServer := httprun.CreateAPIHTTPServer()
	httpServer.Run(options.ServerPort)
}
