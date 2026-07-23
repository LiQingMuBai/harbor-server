package handler

import "cointrade/http/common"

type CreditModule struct {
	common.ModuleBase
}

func (m *CreditModule) ModuleList() common.MODULEHANDLELIST {
	var routes common.MODULEHANDLELIST
	routes = append(routes, m.fundingRoutes()...)
	routes = append(routes, m.logRoutes()...)
	routes = append(routes, m.walletRoutes()...)
	return routes
}
