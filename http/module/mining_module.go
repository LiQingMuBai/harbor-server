package module

import "cointrade/http/common"

type MiningModule struct {
	common.ModuleBase
}

func (m *MiningModule) ModuleList() common.MODULEHANDLELIST {
	var routes common.MODULEHANDLELIST
	routes = append(routes, m.productRoutes()...)
	routes = append(routes, m.orderRoutes()...)
	return routes
}
