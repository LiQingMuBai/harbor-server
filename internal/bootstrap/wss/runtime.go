package wss

import (
	"cointrade/http/httprun"
	"cointrade/internal/bootstrap/shared"
	"cointrade/models"
)

type Options struct {
	Port int
}

func OptionsFromEnv() Options {
	return Options{
		Port: shared.GetenvInt("WSS_PORT", 9088),
	}
}

func Run(options Options) {
	models.InitData()
	httpServer := httprun.CreateWSSHTTPServer()
	httprun.StartWSSBackgroundJobs()
	httpServer.Run(options.Port)
}
