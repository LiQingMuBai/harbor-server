package modules

import "cointrade/http/common"

type MingingModule struct {
	common.ModuleBase
}

func (m *MingingModule) ModuleList() common.MODULEHANDLELIST {
	var routes common.MODULEHANDLELIST
	routes = append(routes, m.productRoutes()...)
	routes = append(routes, m.orderRoutes()...)
	return routes
}
