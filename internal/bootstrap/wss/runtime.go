package wss

import (
	"cointrade/http/common"
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
	httpServer := common.CreateHttp()
	httpServer.Handle.GET("/wss", httprun.CreateWss)
	httprun.DataUpdateFunc()
	go httprun.MessageService()
	go httprun.DBServiceMessageReciveFunc()
	go httprun.DataService()
	httpServer.Run(options.Port)
}
