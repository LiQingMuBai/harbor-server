package handler

import (
	"cointrade/http/common"
	"cointrade/lib/db"
	"sync"
)

type SystemModule struct {
	common.ModuleBase
	sync.RWMutex

	coinListCache    db.DB_LIST_RESULT
	coinListCacheVer int64

	configCache    map[string]interface{}
	configCacheVer int64
}

func (m *SystemModule) ModuleList() common.MODULEHANDLELIST {
	var routes common.MODULEHANDLELIST
	routes = append(routes, m.catalogRoutes()...)
	routes = append(routes, m.contentRoutes()...)
	routes = append(routes, m.financeRoutes()...)
	return routes
}
