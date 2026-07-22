package modules

import "cointrade/http/common"

type UserModule struct {
	common.ModuleBase
}

func (m *UserModule) ModuleList() common.MODULEHANDLELIST {
	var routes common.MODULEHANDLELIST
	routes = append(routes, m.marketingRoutes()...)
	routes = append(routes, m.authRoutes()...)
	routes = append(routes, m.profileRoutes()...)
	routes = append(routes, m.securityRoutes()...)
	routes = append(routes, m.bankRoutes()...)
	routes = append(routes, m.noticeRoutes()...)
	return routes
}
