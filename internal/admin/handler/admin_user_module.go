package handler

import "cointrade/http/common"

type AdminUserModule struct {
	common.ModuleBase
}

func (m *AdminUserModule) ModuleList() common.MODULEHANDLELIST {
	var routes common.MODULEHANDLELIST
	routes = append(routes, m.authRoutes()...)
	routes = append(routes, m.userRoutes()...)
	routes = append(routes, m.financeRoutes()...)
	routes = append(routes, m.tradeRoutes()...)
	routes = append(routes, m.systemRoutes()...)
	routes = append(routes, m.messageRoutes()...)
	return routes
}
