package modules

import (
	"cointrade/config"
	"cointrade/http/common"
	"cointrade/lib/db"
	"cointrade/models"
	"cointrade/utils"
	"fmt"
	"math/rand"
	"sync"

	"github.com/gin-gonic/gin"
)

type SystemModule struct {
	common.ModuleBase
	sync.RWMutex
}

func (m *SystemModule) ModuleList() common.MODULEHANDLELIST {
	//每个MODULE必须要实现的MODULE
	return common.MODULEHANDLELIST{
		&common.ModuleHandles{Method: "post", Path: "/system/coinlist", Handles: common.HandleArray{m.NeedLogin, m.CoinList}},
		&common.ModuleHandles{Method: "get", Path: "/system/config2", Handles: common.HandleArray{m.GetConfig}},
		&common.ModuleHandles{Method: "post", Path: "/system/config", Handles: common.HandleArray{m.GetConfig2}},
		&common.ModuleHandles{Method: "post", Path: "/system/notice", Handles: common.HandleArray{m.Notice}},
		&common.ModuleHandles{Method: "post", Path: "/system/notice/detail", Handles: common.HandleArray{m.NoticeDetail}},
		&common.ModuleHandles{Method: "post", Path: "/system/currency", Handles: common.HandleArray{m.CurrencyList}},
		&common.ModuleHandles{Method: "post", Path: "/system/coin_desc", Handles: common.HandleArray{m.NeedLogin, m.CoinDesc}},
		//&common.ModuleHandles{Method: "post", Path: "/system/test/service", Handles: common.HandleArray{m.NeedLogin, m.TestService}},
		&common.ModuleHandles{Method: "post", Path: "/system/rule", Handles: common.HandleArray{m.NeedLogin, m.RuleText}},
		&common.ModuleHandles{Method: "post", Path: "/system/rule/detail", Handles: common.HandleArray{m.NeedLogin, m.RuleDetail}},
		&common.ModuleHandles{Method: "get", Path: "/system/test", Handles: common.HandleArray{m.Test}},
		&common.ModuleHandles{Method: "post", Path: "/loan/product", Handles: common.HandleArray{m.NeedLogin, m.GetLoanProdut}},
		&common.ModuleHandles{Method: "post", Path: "/loan/order", Handles: common.HandleArray{m.NeedLogin, m.Loan}},
		&common.ModuleHandles{Method: "post", Path: "/loan/order/list", Handles: common.HandleArray{m.NeedLogin, m.LoanOrderList}},
		&common.ModuleHandles{Method: "post", Path: "/loan/order/info", Handles: common.HandleArray{m.NeedLogin, m.LoanInfo}},
		&common.ModuleHandles{Method: "post", Path: "/loan/order/detail", Handles: common.HandleArray{m.NeedLogin, m.LoanOrderInfo}},
		&common.ModuleHandles{Method: "post", Path: "/coin/buy/list", Handles: common.HandleArray{m.NeedLogin, m.GetBuyCoinList}},
		&common.ModuleHandles{Method: "post", Path: "/coin/buy/order", Handles: common.HandleArray{m.NeedLogin, m.BuyCoin}},
		&common.ModuleHandles{Method: "post", Path: "/coin/new/list", Handles: common.HandleArray{m.NeedLogin, m.NewCoinList}},
		&common.ModuleHandles{Method: "post", Path: "/coin/buy/order/list", Handles: common.HandleArray{m.NeedLogin, m.BuyCoinOrderlist}},
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
func (m *SystemModule) NewCoinList(r *gin.Context) {
	list := make(db.DB_LIST_RESULT, 0)
	for _, v := range models.NEW_COIN_LIST {
		v["kline_config"] = ""
		list = append(list, v)
	}
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, list)
}
func (m *SystemModule) LoanInfo(r *gin.Context) {
	uid := r.GetInt("uid")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_LOAN.GetAllLoanAmount(uid))
}
func (m *SystemModule) LoanOrderInfo(r *gin.Context) {
	uid := r.GetInt("uid")
	sn := m.GetValue(r, "sn")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_LOAN.GetLoanInfo(uid, sn))
}
func (m *SystemModule) LoanOrderList(r *gin.Context) {
	uid := r.GetInt("uid")
	var rq models.LoanOrderListRequest
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_LOAN.GetOrderList(uid, &rq))
}
func (m *SystemModule) Loan(r *gin.Context) {
	uid := r.GetInt("uid")
	var rq models.LoanOrderRequest
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_LOAN.Loan(uid, &rq))
}
func (m *SystemModule) Test(r *gin.Context) {
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, "dsadsa123")
}
func (m *SystemModule) NoticeDetail(r *gin.Context) {
	id := m.GetInt(r, "id")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_NOTICE.GetOne(id))
}
func (m *SystemModule) RuleDetail(r *gin.Context) {
	rs := make(map[string]interface{}, 0)
	m.ConvertObject(r, &rs)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_SYSTEM.RuleOne(rs))
}
func (m *SystemModule) CoinList(r *gin.Context) {
	list := make(db.DB_LIST_RESULT, 0)
	for _, v := range models.COIN_LIST {
		v["kline_config"] = ""
		list = append(list, v)
	}
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, list)
}

var CoinListMu sync.RWMutex // 新增这一把锁
func (m *SystemModule) GetConfig(r *gin.Context) {
	//CoinListMu.Lock()
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
	//m.Unlock()
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
func (m *SystemModule) Notice(r *gin.Context) {
	var rq models.NoticeRequest
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	rq.Limit = 15
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_NOTICE.GetList(&rq))
}
func (m *SystemModule) TestService(r *gin.Context) {
	uid := r.GetInt("uid")
	config.GlobalRedis.PushQueue(models.HASH_USER_SERVICE_MESSAGE, models.DBServiceMessage{
		Uid:        uid,
		Content:    fmt.Sprintf("测试消息%d", 1000+rand.Intn(9000)),
		CreateTime: utils.GetNow(),
		Flag:       2,
	})
}
func (m *SystemModule) RuleText(r *gin.Context) {
	ruletype := m.GetValue(r, "ruletype")
	lang := m.GetValue(r, "lang")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_SYSTEM.RuleText(ruletype, lang))
}
func (m *SystemModule) GetLoanProdut(r *gin.Context) {
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.LOAN_PRODUCT_LIST)
}
func (m *SystemModule) GetBuyCoinList(r *gin.Context) {
	list := make(db.DB_LIST_RESULT, 0)
	for _, v := range models.BUY_COIN_LIST {
		v["kline_config"] = ""
		list = append(list, v)
	}
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, list)
}
func (m *SystemModule) BuyCoin(r *gin.Context) {
	uid := r.GetInt("uid")
	coin_id := m.GetInt(r, "coin_id")
	amount := m.GetFloat(r, "amount")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_SYSTEM.BuyCoin(uid, coin_id, amount))

}
func (m *SystemModule) BuyCoinOrderlist(r *gin.Context) {
	uid := r.GetInt("uid")

	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_SYSTEM.GetBuyCoinOrders(uid))

}
