package httprun

import (
	"cointrade/http/common"
	apihandler "cointrade/internal/api/handler"
)

var moduleUser apihandler.UserModule
var moduleTrade apihandler.TradeModule
var moduleMining apihandler.MiningModule
var moduleAssets apihandler.AssetModule
var moduleMessage apihandler.MessageModule
var moduleSystem apihandler.SystemModule
var moduleCredit apihandler.CreditModule

func registerAPIModules(httpServer *common.HttpModules) {
	registerCoreModules(httpServer)
	registerMessageModules(httpServer)
}

func registerCoreModules(httpServer *common.HttpModules) {
	httpServer.LoadModule(&moduleUser)
	httpServer.LoadModule(&moduleTrade)
	httpServer.LoadModule(&moduleAssets)
	httpServer.LoadModule(&moduleMining)
	httpServer.LoadModule(&moduleSystem)
	httpServer.LoadModule(&moduleCredit)
}

func registerMessageModules(httpServer *common.HttpModules) {
	httpServer.LoadModule(&moduleMessage)
}
