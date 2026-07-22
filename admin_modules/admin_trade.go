package adminmodules

import (
	adminmodels "cointrade/admin_models"
	"cointrade/http/common"

	"github.com/gin-gonic/gin"
)

func (m *AdminUserModule) tradeRoutes() common.MODULEHANDLELIST {
	return common.MODULEHANDLELIST{
		&common.ModuleHandles{Method: "post", Path: "/admin/explode_list", Handles: common.HandleArray{m.ExplodeList}},
		&common.ModuleHandles{Method: "post", Path: "/admin/op_explode", Handles: common.HandleArray{m.CheckLogin, m.OpExplodeTrade}},
		&common.ModuleHandles{Method: "post", Path: "/admin/del_explode", Handles: common.HandleArray{m.CheckLogin, m.DelExplode}},
		&common.ModuleHandles{Method: "post", Path: "/admin/morder_list", Handles: common.HandleArray{m.CheckLogin, m.MorderList}},
		&common.ModuleHandles{Method: "post", Path: "/admin/stop_minner", Handles: common.HandleArray{m.CheckLogin, m.StopMinner}},
		&common.ModuleHandles{Method: "post", Path: "/admin/close_trade_list", Handles: common.HandleArray{m.CheckLogin, m.CloseTrade}},
		&common.ModuleHandles{Method: "post", Path: "/admin/hold_trade_list", Handles: common.HandleArray{m.CheckLogin, m.HoldTradeList}},
		&common.ModuleHandles{Method: "post", Path: "/admin/delegate_now", Handles: common.HandleArray{m.CheckLogin, m.DelegateNow}},
		&common.ModuleHandles{Method: "post", Path: "/admin/manual_delegate_trade", Handles: common.HandleArray{m.CheckLogin, m.ManualDelegateTrade}},
		&common.ModuleHandles{Method: "post", Path: "/admin/delegate_history", Handles: common.HandleArray{m.CheckLogin, m.DelegateHistory}},
		&common.ModuleHandles{Method: "post", Path: "/admin/delegate_del", Handles: common.HandleArray{m.CheckLogin, m.DelegateHistoryDel}},
		&common.ModuleHandles{Method: "post", Path: "admin/spot_delegate", Handles: common.HandleArray{m.CheckLogin, m.SpotDelegate}},
		&common.ModuleHandles{Method: "post", Path: "admin/opspot", Handles: common.HandleArray{m.CheckLogin, m.OpSpot}},
	}
}

func (m *AdminUserModule) OpSpot(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_TRADE.OpSpot(rq))
}

func (m *AdminUserModule) SpotDelegate(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_TRADE.HistoryDelegateList(rq, true, ""))
}

func (m *AdminUserModule) DelegateHistoryDel(r *gin.Context) {
	id := m.GetInt(r, "id")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_TRADE.DelegateHistoryDel(id))
}

func (m *AdminUserModule) DelegateHistory(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	rq["state"] = 1
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_TRADE.HistoryDelegateList(rq, true, ""))
}

func (m *AdminUserModule) DelegateNow(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	rq["state"] = 0
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_TRADE.HistoryDelegateList(rq, true, ""))
}

func (m *AdminUserModule) ManualDelegateTrade(r *gin.Context) {
	uid := m.GetInt(r, "uid")
	sn := m.GetValue(r, "sn")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_TRADE.ManualOperationTrade(uid, sn))
}

func (m *AdminUserModule) HoldTradeList(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	rq["state"] = 1
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_TRADE.TradeList(rq, true, ""))
}

func (m *AdminUserModule) CloseTrade(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_TRADE.HistoryCloseRradeList(rq, true, ""))
}

func (m *AdminUserModule) StopMinner(r *gin.Context) {
	id := m.GetInt(r, "id")
	pass := m.GetValue(r, "pass")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.StopMinner(id, pass))
}

func (m *AdminUserModule) MorderList(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.MorderList(rq))
}

func (m *AdminUserModule) DelExplode(r *gin.Context) {
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.DelExplodeTrade(m.GetInt(r, "id")))
}

func (m *AdminUserModule) OpExplodeTrade(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.OpExplodeTrade(rq))
}

func (m *AdminUserModule) ExplodeList(r *gin.Context) {
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.ExplodeTradeList())
}
