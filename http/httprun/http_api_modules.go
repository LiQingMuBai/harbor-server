package httprun

import (
	"cointrade/http/common"
	"cointrade/http/module"
)

var moduleUser module.UserModule
var moduleTrade module.TradeModule
var moduleMining module.MiningModule
var moduleAssets module.AssetModule
var moduleMessage module.MessageModule
var moduleSystem module.SystemModule
var moduleCredit module.CreditModule

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
