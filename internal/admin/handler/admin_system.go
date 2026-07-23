package handler

import (
	"cointrade/config"
	"cointrade/http/common"
	adminservice "cointrade/internal/admin/service"
	"cointrade/lib/db"
	"cointrade/models"

	"github.com/gin-gonic/gin"
)

func (m *AdminUserModule) systemRoutes() common.MODULEHANDLELIST {
	return common.MODULEHANDLELIST{
		&common.ModuleHandles{Method: "post", Path: "/admin/coin_list", Handles: common.HandleArray{m.CheckLogin, m.CoinList}},
		&common.ModuleHandles{Method: "post", Path: "/admin/coindesc_list", Handles: common.HandleArray{m.CheckLogin, m.CoinDescList}},
		&common.ModuleHandles{Method: "post", Path: "/admin/op_coindesc", Handles: common.HandleArray{m.CheckLogin, m.SaveCoinDesc}},
		&common.ModuleHandles{Method: "post", Path: "/admin/del_coindesc", Handles: common.HandleArray{m.CheckLogin, m.DeleteCoinDesc}},
		&common.ModuleHandles{Method: "post", Path: "/admin/save_coin", Handles: common.HandleArray{m.CheckLogin, m.SaveCoin}},
		&common.ModuleHandles{Method: "post", Path: "/admin/del_coin", Handles: common.HandleArray{m.CheckLogin, m.DeleteCoin}},
		&common.ModuleHandles{Method: "post", Path: "/admin/currency_list", Handles: common.HandleArray{m.CheckLogin, m.CurrencyList}},
		&common.ModuleHandles{Method: "post", Path: "/admin/op_currency", Handles: common.HandleArray{m.CheckLogin, m.SaveCurrency}},
		&common.ModuleHandles{Method: "post", Path: "/admin/del_currency", Handles: common.HandleArray{m.CheckLogin, m.DeleteCurrency}},
		&common.ModuleHandles{Method: "post", Path: "/admin/minner_open", Handles: common.HandleArray{m.CheckLogin, m.MinnerOpen}},
		&common.ModuleHandles{Method: "post", Path: "/admin/minner_list", Handles: common.HandleArray{m.CheckLogin, m.ListMiningProducts}},
		&common.ModuleHandles{Method: "post", Path: "/admin/minner_op", Handles: common.HandleArray{m.CheckLogin, m.SaveMiningProduct}},
		&common.ModuleHandles{Method: "post", Path: "/admin/minner_del", Handles: common.HandleArray{m.CheckLogin, m.DeleteMiningProduct}},
		&common.ModuleHandles{Method: "post", Path: "/admin/notice_list", Handles: common.HandleArray{m.CheckLogin, m.NoticeList}},
		&common.ModuleHandles{Method: "post", Path: "/admin/notice_op", Handles: common.HandleArray{m.CheckLogin, m.SaveNotice}},
		&common.ModuleHandles{Method: "post", Path: "/admin/notice_del", Handles: common.HandleArray{m.CheckLogin, m.DeleteNotice}},
		&common.ModuleHandles{Method: "post", Path: "/admin/setting", Handles: common.HandleArray{m.CheckLogin, m.Setting}},
		&common.ModuleHandles{Method: "post", Path: "/admin/sitecount", Handles: common.HandleArray{m.CheckLogin, m.SiteCount}},
		&common.ModuleHandles{Method: "post", Path: "/admin/notify_list", Handles: common.HandleArray{m.CheckLogin, m.NotifyList}},
		&common.ModuleHandles{Method: "post", Path: "/admin/clearNotify", Handles: common.HandleArray{m.ClearNotify}},
		&common.ModuleHandles{Method: "post", Path: "/admin/statictis_count", Handles: common.HandleArray{m.CheckLogin, m.StatictisCount}},
		&common.ModuleHandles{Method: "post", Path: "admin/rulelist", Handles: common.HandleArray{m.CheckLogin, m.Rulelist}},
		&common.ModuleHandles{Method: "post", Path: "admin/rulehandler", Handles: common.HandleArray{m.RuleHandler}},
		&common.ModuleHandles{Method: "post", Path: "admin/del_rule", Handles: common.HandleArray{m.CheckLogin, m.DeleteRule}},
		&common.ModuleHandles{Method: "post", Path: "admin/controller_list", Handles: common.HandleArray{m.CheckLogin, m.ControllerList}},
		&common.ModuleHandles{Method: "post", Path: "admin/del_controller", Handles: common.HandleArray{m.CheckLogin, m.DeleteController}},
		&common.ModuleHandles{Method: "post", Path: "admin/kline_controller", Handles: common.HandleArray{m.CheckLogin, m.KlineController}},
		&common.ModuleHandles{Method: "post", Path: "admin/explode_controller", Handles: common.HandleArray{m.CheckLogin, m.ExplodeController}},
		&common.ModuleHandles{Method: "post", Path: "admin/accept_list", Handles: common.HandleArray{m.CheckLogin, m.AcceptList}},
		&common.ModuleHandles{Method: "post", Path: "admin/op_acceptminner", Handles: common.HandleArray{m.CheckLogin, m.SaveMiningAcceptance}},
		&common.ModuleHandles{Method: "post", Path: "admin/del_acceptminner", Handles: common.HandleArray{m.CheckLogin, m.DeleteMiningAcceptance}},
		&common.ModuleHandles{Method: "post", Path: "admin/op_loan", Handles: common.HandleArray{m.CheckLogin, m.LoanSetting}},
		&common.ModuleHandles{Method: "post", Path: "admin/loan", Handles: common.HandleArray{m.CheckLogin, m.LoanList}},
		&common.ModuleHandles{Method: "post", Path: "admin/audit_accept", Handles: common.HandleArray{m.CheckLogin, m.AuditAccept}},
		&common.ModuleHandles{Method: "post", Path: "admin/kline_config", Handles: common.HandleArray{m.CheckLogin, m.kline_config}},
	}
}

func (m *AdminUserModule) AuditAccept(r *gin.Context) {
	rq := make(adminservice.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.SYSTEM_MODEL.AuditAccept(rq))
}

func (m *AdminUserModule) LoanList(r *gin.Context) {
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.SYSTEM_MODEL.LoanList())
}

func (m *AdminUserModule) LoanSetting(r *gin.Context) {
	rq := make(adminservice.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.SYSTEM_MODEL.LoanSetting(rq))
}

func (m *AdminUserModule) DeleteMiningAcceptance(r *gin.Context) {
	id := m.GetInt(r, "id")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.SYSTEM_MODEL.DeleteMiningAcceptance(id))
}

func (m *AdminUserModule) SaveMiningAcceptance(r *gin.Context) {
	rq := make(adminservice.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.SYSTEM_MODEL.SaveMiningAcceptance(rq))
}

func (m *AdminUserModule) AcceptList(r *gin.Context) {
	rq := make(adminservice.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.SYSTEM_MODEL.AcceptList(rq))
}

func (m *AdminUserModule) DeleteController(r *gin.Context) {
	rq := make(adminservice.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.SYSTEM_MODEL.DeleteController(rq))
}

func (m *AdminUserModule) ControllerList(r *gin.Context) {
	rq := make(adminservice.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.SYSTEM_MODEL.ControllerTradeList())
}

func (m *AdminUserModule) ExplodeController(r *gin.Context) {
	rq := make(adminservice.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.SYSTEM_MODEL.ExplodeController(rq))
}

func (m *AdminUserModule) KlineController(r *gin.Context) {
	rq := make(adminservice.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.SYSTEM_MODEL.KlineController(rq))
}

func (m *AdminUserModule) DeleteRule(r *gin.Context) {
	rq := make(adminservice.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.SYSTEM_MODEL.DeleteRule(rq))
}

func (m *AdminUserModule) Rulelist(r *gin.Context) {
	rq := make(adminservice.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.SYSTEM_MODEL.Rulelist(rq))
}

func (m *AdminUserModule) RuleHandler(r *gin.Context) {
	rq := make(adminservice.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.SYSTEM_MODEL.RuleHandler(rq))
}

func (m *AdminUserModule) StatictisCount(r *gin.Context) {
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.SYSTEM_MODEL.StatictisCount())
}

func (m *AdminUserModule) ClearNotify(r *gin.Context) {
	tp := m.GetValue(r, "type")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.SYSTEM_MODEL.ClearNotify(tp))
}

func (m *AdminUserModule) NotifyList(r *gin.Context) {
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.SYSTEM_MODEL.NotifyList())
}

func (m *AdminUserModule) SiteCount(r *gin.Context) {
	rq := make(adminservice.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.AdminResponse{State: 2000, Data: adminservice.SYSTEM_MODEL.SiteCount(rq)})
}

func (m *AdminUserModule) Setting(r *gin.Context) {
	rq := make(adminservice.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.SYSTEM_MODEL.Setting(rq))
}

func (m *AdminUserModule) NoticeList(r *gin.Context) {
	rq := make(adminservice.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.SYSTEM_MODEL.NoticeList(rq))
}

func (m *AdminUserModule) SaveNotice(r *gin.Context) {
	rq := make(adminservice.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.SYSTEM_MODEL.SaveNotice(rq))
}

func (m *AdminUserModule) DeleteNotice(r *gin.Context) {
	id := m.GetInt(r, "id")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.SYSTEM_MODEL.DeleteNotice(id))
}

func (m *AdminUserModule) DeleteMiningProduct(r *gin.Context) {
	id := m.GetInt(r, "id")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.SYSTEM_MODEL.DeleteMiningProduct(id))
}

func (m *AdminUserModule) SaveMiningProduct(r *gin.Context) {
	rq := make(adminservice.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.SYSTEM_MODEL.SaveMiningProduct(rq))
}

func (m *AdminUserModule) MinnerOpen(r *gin.Context) {
	id := m.GetInt(r, "id")
	key := m.GetValue(r, "key")
	isopen := m.GetInt(r, "state")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.SYSTEM_MODEL.MinnerSet(id, key, isopen))
}

func (m *AdminUserModule) ListMiningProducts(r *gin.Context) {
	rq := make(adminservice.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.SYSTEM_MODEL.ListMiningProducts(rq))
}

func (m *AdminUserModule) DeleteCurrency(r *gin.Context) {
	id := m.GetValue(r, "id")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.SYSTEM_MODEL.DeleteCurrency(id))
}

func (m *AdminUserModule) SaveCurrency(r *gin.Context) {
	rq := make(adminservice.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.SYSTEM_MODEL.SaveCurrency(rq))
}

func (m *AdminUserModule) CurrencyList(r *gin.Context) {
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.SYSTEM_MODEL.CurrencyList())
}

func (m *AdminUserModule) DeleteCoin(r *gin.Context) {
	rq := make(adminservice.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.SYSTEM_MODEL.DeleteCoin(rq))
}

func (m *AdminUserModule) SaveCoin(r *gin.Context) {
	rq := make(adminservice.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.SYSTEM_MODEL.SaveCoin(rq))
}

func (m *AdminUserModule) CoinDescList(r *gin.Context) {
	rq := make(adminservice.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.SYSTEM_MODEL.CoinDescList(rq))
}

func (m *AdminUserModule) SaveCoinDesc(r *gin.Context) {
	rq := make(adminservice.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.SYSTEM_MODEL.SaveCoinDesc(rq))
}

func (m *AdminUserModule) DeleteCoinDesc(r *gin.Context) {
	id := m.GetInt(r, "id")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.SYSTEM_MODEL.DeleteCoinDesc(id))
}

func (m *AdminUserModule) CoinList(r *gin.Context) {
	rq := make(adminservice.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.SYSTEM_MODEL.CoinList(rq))
}

func (m *AdminUserModule) kline_config(r *gin.Context) {
	var klineConfig models.CoinKlineConfig
	id := m.GetInt(r, "id")
	klineConfig.BaseAmount = m.GetInt(r, "base_amount")
	klineConfig.Heart = m.GetFloat(r, "heart")
	klineConfig.HighRate = m.GetFloat(r, "high_rate")
	klineConfig.LowRate = m.GetFloat(r, "low_rate")
	klineConfig.MaxPrice = m.GetFloat(r, "max_price")
	klineConfig.MinPrice = m.GetFloat(r, "min_price")
	klineConfig.UpRate = m.GetInt(r, "up_rate")
	config.GlobalDB.UpdateData(models.DB_TABLE_COINS, db.DB_PARAMS{"kline_config": klineConfig}, db.DB_PARAMS{"id": id})
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, &adminservice.AdminResponse{
		State: adminservice.SUCCESS,
		Data:  "已提交控制!",
	})
}
