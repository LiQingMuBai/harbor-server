package adminmodules

import (
	adminmodels "cointrade/admin_models"
	"cointrade/config"
	"cointrade/http/common"
	"cointrade/lib/db"
	"cointrade/models"

	"github.com/gin-gonic/gin"
)

func (m *AdminUserModule) systemRoutes() common.MODULEHANDLELIST {
	return common.MODULEHANDLELIST{
		&common.ModuleHandles{Method: "post", Path: "/admin/coin_list", Handles: common.HandleArray{m.CheckLogin, m.CoinList}},
		&common.ModuleHandles{Method: "post", Path: "/admin/coindesc_list", Handles: common.HandleArray{m.CheckLogin, m.CoinDescList}},
		&common.ModuleHandles{Method: "post", Path: "/admin/op_coindesc", Handles: common.HandleArray{m.CheckLogin, m.OpCoinDesc}},
		&common.ModuleHandles{Method: "post", Path: "/admin/del_coindesc", Handles: common.HandleArray{m.CheckLogin, m.DelCoinDesc}},
		&common.ModuleHandles{Method: "post", Path: "/admin/save_coin", Handles: common.HandleArray{m.CheckLogin, m.SaveCoin}},
		&common.ModuleHandles{Method: "post", Path: "/admin/del_coin", Handles: common.HandleArray{m.CheckLogin, m.DelCoin}},
		&common.ModuleHandles{Method: "post", Path: "/admin/currency_list", Handles: common.HandleArray{m.CheckLogin, m.CurrencyList}},
		&common.ModuleHandles{Method: "post", Path: "/admin/op_currency", Handles: common.HandleArray{m.CheckLogin, m.OpCurrency}},
		&common.ModuleHandles{Method: "post", Path: "/admin/del_currency", Handles: common.HandleArray{m.CheckLogin, m.DelCurrency}},
		&common.ModuleHandles{Method: "post", Path: "/admin/minner_open", Handles: common.HandleArray{m.CheckLogin, m.MinnerOpen}},
		&common.ModuleHandles{Method: "post", Path: "/admin/minner_list", Handles: common.HandleArray{m.CheckLogin, m.MinnerList}},
		&common.ModuleHandles{Method: "post", Path: "/admin/minner_op", Handles: common.HandleArray{m.CheckLogin, m.OpMinner}},
		&common.ModuleHandles{Method: "post", Path: "/admin/minner_del", Handles: common.HandleArray{m.CheckLogin, m.DelMinner}},
		&common.ModuleHandles{Method: "post", Path: "/admin/notice_list", Handles: common.HandleArray{m.CheckLogin, m.NoticeList}},
		&common.ModuleHandles{Method: "post", Path: "/admin/notice_op", Handles: common.HandleArray{m.CheckLogin, m.OPNotice}},
		&common.ModuleHandles{Method: "post", Path: "/admin/notice_del", Handles: common.HandleArray{m.CheckLogin, m.DelNotice}},
		&common.ModuleHandles{Method: "post", Path: "/admin/setting", Handles: common.HandleArray{m.CheckLogin, m.Setting}},
		&common.ModuleHandles{Method: "post", Path: "/admin/sitecount", Handles: common.HandleArray{m.CheckLogin, m.SiteCount}},
		&common.ModuleHandles{Method: "post", Path: "/admin/notify_list", Handles: common.HandleArray{m.CheckLogin, m.NotifyList}},
		&common.ModuleHandles{Method: "post", Path: "/admin/clearNotify", Handles: common.HandleArray{m.ClearNotify}},
		&common.ModuleHandles{Method: "post", Path: "/admin/statictis_count", Handles: common.HandleArray{m.CheckLogin, m.StatictisCount}},
		&common.ModuleHandles{Method: "post", Path: "admin/rulelist", Handles: common.HandleArray{m.CheckLogin, m.Rulelist}},
		&common.ModuleHandles{Method: "post", Path: "admin/rulehandler", Handles: common.HandleArray{m.RuleHandler}},
		&common.ModuleHandles{Method: "post", Path: "admin/del_rule", Handles: common.HandleArray{m.CheckLogin, m.DelRule}},
		&common.ModuleHandles{Method: "post", Path: "admin/controller_list", Handles: common.HandleArray{m.CheckLogin, m.ControllerList}},
		&common.ModuleHandles{Method: "post", Path: "admin/del_controller", Handles: common.HandleArray{m.CheckLogin, m.DelController}},
		&common.ModuleHandles{Method: "post", Path: "admin/kline_controller", Handles: common.HandleArray{m.CheckLogin, m.KlineController}},
		&common.ModuleHandles{Method: "post", Path: "admin/explode_controller", Handles: common.HandleArray{m.CheckLogin, m.ExplodeController}},
		&common.ModuleHandles{Method: "post", Path: "admin/accept_list", Handles: common.HandleArray{m.CheckLogin, m.AcceptList}},
		&common.ModuleHandles{Method: "post", Path: "admin/op_acceptminner", Handles: common.HandleArray{m.CheckLogin, m.OpAcceptMinner}},
		&common.ModuleHandles{Method: "post", Path: "admin/del_acceptminner", Handles: common.HandleArray{m.CheckLogin, m.DelAcceptMinner}},
		&common.ModuleHandles{Method: "post", Path: "admin/op_loan", Handles: common.HandleArray{m.CheckLogin, m.LoanSetting}},
		&common.ModuleHandles{Method: "post", Path: "admin/loan", Handles: common.HandleArray{m.CheckLogin, m.LoanList}},
		&common.ModuleHandles{Method: "post", Path: "admin/audit_accept", Handles: common.HandleArray{m.CheckLogin, m.AuditAccept}},
		&common.ModuleHandles{Method: "post", Path: "admin/kline_config", Handles: common.HandleArray{m.CheckLogin, m.kline_config}},
	}
}

func (m *AdminUserModule) AuditAccept(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.AuditAccept(rq))
}

func (m *AdminUserModule) LoanList(r *gin.Context) {
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.LoanList())
}

func (m *AdminUserModule) LoanSetting(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.LoanSetting(rq))
}

func (m *AdminUserModule) DelAcceptMinner(r *gin.Context) {
	id := m.GetInt(r, "id")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.DelMinneAccept(id))
}

func (m *AdminUserModule) OpAcceptMinner(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.OpMinnerAccept(rq))
}

func (m *AdminUserModule) AcceptList(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.AcceptList(rq))
}

func (m *AdminUserModule) DelController(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.DelController(rq))
}

func (m *AdminUserModule) ControllerList(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.ControllerTradeList())
}

func (m *AdminUserModule) ExplodeController(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.ExplodeController(rq))
}

func (m *AdminUserModule) KlineController(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.KlineController(rq))
}

func (m *AdminUserModule) DelRule(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.DelRule(rq))
}

func (m *AdminUserModule) Rulelist(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.Rulelist(rq))
}

func (m *AdminUserModule) RuleHandler(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.RuleHandler(rq))
}

func (m *AdminUserModule) StatictisCount(r *gin.Context) {
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.StatictisCount())
}

func (m *AdminUserModule) ClearNotify(r *gin.Context) {
	tp := m.GetValue(r, "type")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.ClearNotify(tp))
}

func (m *AdminUserModule) NotifyList(r *gin.Context) {
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.NotifyList())
}

func (m *AdminUserModule) SiteCount(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.AdminResponse{State: 2000, Data: adminmodels.SYSTEM_MODEL.SiteCount(rq)})
}

func (m *AdminUserModule) Setting(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.Setting(rq))
}

func (m *AdminUserModule) NoticeList(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.NoticeList(rq))
}

func (m *AdminUserModule) OPNotice(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.OpNotice(rq))
}

func (m *AdminUserModule) DelNotice(r *gin.Context) {
	id := m.GetInt(r, "id")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.DelNotice(id))
}

func (m *AdminUserModule) DelMinner(r *gin.Context) {
	id := m.GetInt(r, "id")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.DelMinner(id))
}

func (m *AdminUserModule) OpMinner(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.OpMinner(rq))
}

func (m *AdminUserModule) MinnerOpen(r *gin.Context) {
	id := m.GetInt(r, "id")
	key := m.GetValue(r, "key")
	isopen := m.GetInt(r, "state")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.MinnerSet(id, key, isopen))
}

func (m *AdminUserModule) MinnerList(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.MinnerList(rq))
}

func (m *AdminUserModule) DelCurrency(r *gin.Context) {
	id := m.GetValue(r, "id")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.DelCurrency(id))
}

func (m *AdminUserModule) OpCurrency(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.OpCurrency(rq))
}

func (m *AdminUserModule) CurrencyList(r *gin.Context) {
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.CurrencyList())
}

func (m *AdminUserModule) DelCoin(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.DelCoin(rq))
}

func (m *AdminUserModule) SaveCoin(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.OpCoin(rq))
}

func (m *AdminUserModule) CoinDescList(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.CoinDescList(rq))
}

func (m *AdminUserModule) OpCoinDesc(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.OpCoinDesc(rq))
}

func (m *AdminUserModule) DelCoinDesc(r *gin.Context) {
	id := m.GetInt(r, "id")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.DelCoinDesc(id))
}

func (m *AdminUserModule) CoinList(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.CoinList(rq))
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
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, &adminmodels.AdminResponse{
		State: adminmodels.SUCCESS,
		Data:  "已提交控制!",
	})
}
