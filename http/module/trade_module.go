package module

import "cointrade/http/common"

// 交易处理MODULE
type TradeModule struct {
	common.ModuleBase
}

func (m *TradeModule) ModuleList() common.MODULEHANDLELIST {
	var routes common.MODULEHANDLELIST
	routes = append(routes, m.delegateRoutes()...)
	routes = append(routes, m.positionRoutes()...)
	return routes
}
