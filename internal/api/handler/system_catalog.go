package handler

import (
	"cointrade/config"
	"cointrade/http/common"
	"cointrade/lib/db"
	"cointrade/models"
	"sync/atomic"

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

func buildCoinList() db.DB_LIST_RESULT {
	list := make(db.DB_LIST_RESULT, 0, len(models.COIN_LIST))
	for _, v := range models.COIN_LIST {
		item := make(map[string]string, len(v))
		for k, val := range v {
			item[k] = val
		}
		item["kline_config"] = ""
		list = append(list, item)
	}
	return list
}

func (m *SystemModule) getCoinListCached() db.DB_LIST_RESULT {
	ver := atomic.LoadInt64(&models.ConfigVersion)
	m.RLock()
	if m.coinListCache != nil && m.coinListCacheVer == ver {
		defer m.RUnlock()
		return m.coinListCache
	}
	m.RUnlock()

	m.Lock()
	defer m.Unlock()
	if m.coinListCache != nil && m.coinListCacheVer == ver {
		return m.coinListCache
	}
	m.coinListCache = buildCoinList()
	m.coinListCacheVer = ver
	return m.coinListCache
}

func (m *SystemModule) getConfigCached() map[string]interface{} {
	ver := atomic.LoadInt64(&models.ConfigVersion)
	m.RLock()
	if m.configCache != nil && m.configCacheVer == ver {
		defer m.RUnlock()
		return m.configCache
	}
	m.RUnlock()

	m.Lock()
	defer m.Unlock()
	if m.configCache != nil && m.configCacheVer == ver {
		return m.configCache
	}
	m.configCache = map[string]interface{}{
		"coinlist":       buildCoinList(),
		"creditconfig":   models.RECHARGE_ADDRESS_LIST,
		"explodeconfig":  models.EXPLODE_CONFIG,
		"siteconfig":     config.SYSTEM_CONFIG,
		"incomeconfig":   map[string]interface{}{"recharge": models.RECHARGE_INCOME_RATES, "mining": models.MINING_INCOME_RATES},
		"approve_wallet": config.GlobalConfig.GetValue("approve_wallet").ToString(),
	}
	m.configCacheVer = ver
	return m.configCache
}

func (m *SystemModule) CoinList(r *gin.Context) {
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, m.getCoinListCached())
}

func (m *SystemModule) GetConfig(r *gin.Context) {
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, m.getConfigCached())
}

func (m *SystemModule) GetConfig2(r *gin.Context) {
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, m.getConfigCached())
}

func (m *SystemModule) CurrencyList(r *gin.Context) {
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.CURRENCY_LIST)
}
