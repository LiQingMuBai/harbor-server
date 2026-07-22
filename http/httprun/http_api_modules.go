package httprun

import (
	"cointrade/http/common"
	"cointrade/http/modules"
)

var moduleUser modules.UserModule
var moduleTrade modules.TradeModule
var moduleMining modules.MingingModule
var moduleAssets modules.AssetModule
var moduleMessage modules.MessageModule
var moduleSystem modules.SystemModule
var moduleCredit modules.CreditModule

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
