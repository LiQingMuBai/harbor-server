package modules

import "cointrade/http/common"

// 资产相关
type AssetModule struct {
	common.ModuleBase
}

func (m *AssetModule) ModuleList() common.MODULEHANDLELIST {
	return m.exchangeRoutes()
}
