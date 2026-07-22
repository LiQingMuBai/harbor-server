package module

import (
	"cointrade/http/common"
	"sync"
)

type SystemModule struct {
	common.ModuleBase
	sync.RWMutex
}

func (m *SystemModule) ModuleList() common.MODULEHANDLELIST {
	var routes common.MODULEHANDLELIST
	routes = append(routes, m.catalogRoutes()...)
	routes = append(routes, m.contentRoutes()...)
	routes = append(routes, m.financeRoutes()...)
	return routes
}
