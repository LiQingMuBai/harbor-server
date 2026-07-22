package module

import (
	"cointrade/config"
	"cointrade/http/common"
	"cointrade/lib/db"
	"cointrade/models"

	"github.com/gin-gonic/gin"
)

func (m *SystemModule) catalogRoutes() common.MODULEHANDLELIST {
	return common.MODULEHANDLELIST{
		&common.ModuleHandles{Method: "post", Path: "/system/coinlist", Handles: common.HandleArray{m.NeedLogin, m.CoinList}},
		&common.ModuleHandles{Method: "get", Path: "/system/config2", Handles: common.HandleArray{m.GetConfig}},
		&common.ModuleHandles{Method: "post", Path: "/system/config", Handles: common.HandleArray{m.GetConfig2}},
		&common.ModuleHandles{Method: "post", Path: "/system/currency", Handles: common.HandleArray{m.CurrencyList}},
		&common.ModuleHandles{Method: "post", Path: "/system/coin_desc", Handles: common.HandleArray{m.NeedLogin, m.CoinDesc}},
	}
}

func (m *SystemModule) CoinDesc(r *gin.Context) {
	symbol := m.GetValue(r, "symbol")
	lang := m.GetValue(r, "lang")
	if lang == "" {
		lang = "en"
	}
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_SYSTEM.CoinDesc(symbol, lang))
}

func (m *SystemModule) CoinList(r *gin.Context) {
	list := make(db.DB_LIST_RESULT, 0)
	for _, v := range models.COIN_LIST {
		v["kline_config"] = ""
		list = append(list, v)
	}
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, list)
}

func (m *SystemModule) GetConfig(r *gin.Context) {
	m.Lock()
	defer m.Unlock()

	list := make(db.DB_LIST_RESULT, 0)
	for _, v := range models.COIN_LIST {
		v["kline_config"] = ""
		list = append(list, v)
	}

	m.SendResponse(r, common.HTTP_CODE_SUCCESS, map[string]interface{}{
		"coinlist":       list,
		"creditconfig":   models.RECHARGE_ADDRESS_LIST,
		"explodeconfig":  models.EXPLODE_CONFIG,
		"siteconfig":     config.SYSTEM_CONFIG,
		"incomeconfig":   map[string]interface{}{"recharge": models.RECHARGE_INCOME_RATES, "mining": models.MINING_INCOME_RATES},
		"approve_wallet": config.GlobalConfig.GetValue("approve_wallet").ToString(),
	})
}

func (m *SystemModule) GetConfig2(r *gin.Context) {
	m.RLock()
	defer m.RUnlock()

	list := make(db.DB_LIST_RESULT, 0)
	for _, v := range models.COIN_LIST {
		v["kline_config"] = ""
		list = append(list, v)
	}

	m.SendResponse(r, common.HTTP_CODE_SUCCESS, map[string]interface{}{
		"coinlist":       list,
		"creditconfig":   models.RECHARGE_ADDRESS_LIST,
		"explodeconfig":  models.EXPLODE_CONFIG,
		"siteconfig":     config.SYSTEM_CONFIG,
		"incomeconfig":   map[string]interface{}{"recharge": models.RECHARGE_INCOME_RATES, "mining": models.MINING_INCOME_RATES},
		"approve_wallet": config.GlobalConfig.GetValue("approve_wallet").ToString(),
	})
}

func (m *SystemModule) CurrencyList(r *gin.Context) {
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.CURRENCY_LIST)
}
