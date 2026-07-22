package modules

import "cointrade/http/common"

type MessageModule struct {
	common.ModuleBase
}

func (m *MessageModule) ModuleList() common.MODULEHANDLELIST {
	var routes common.MODULEHANDLELIST
	routes = append(routes, m.noticeRoutes()...)
	routes = append(routes, m.serviceRoutes()...)
	return routes
}
