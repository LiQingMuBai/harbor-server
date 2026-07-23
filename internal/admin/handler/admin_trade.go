package handler

import (
	"cointrade/http/common"
	adminservice "cointrade/internal/admin/service"

	"github.com/gin-gonic/gin"
)

func (m *AdminUserModule) tradeRoutes() common.MODULEHANDLELIST {
	return common.MODULEHANDLELIST{
		&common.ModuleHandles{Method: "post", Path: "/admin/explode_list", Handles: common.HandleArray{m.ExplodeList}},
		&common.ModuleHandles{Method: "post", Path: "/admin/op_explode", Handles: common.HandleArray{m.CheckLogin, m.SaveExplodeTrade}},
		&common.ModuleHandles{Method: "post", Path: "/admin/del_explode", Handles: common.HandleArray{m.CheckLogin, m.DeleteExplodeTrade}},
		&common.ModuleHandles{Method: "post", Path: "/admin/morder_list", Handles: common.HandleArray{m.CheckLogin, m.ListMiningOrders}},
		&common.ModuleHandles{Method: "post", Path: "/admin/stop_minner", Handles: common.HandleArray{m.CheckLogin, m.StopMiningOrder}},
		&common.ModuleHandles{Method: "post", Path: "/admin/close_trade_list", Handles: common.HandleArray{m.CheckLogin, m.ListClosedTrades}},
		&common.ModuleHandles{Method: "post", Path: "/admin/hold_trade_list", Handles: common.HandleArray{m.CheckLogin, m.ListOpenTrades}},
		&common.ModuleHandles{Method: "post", Path: "/admin/delegate_now", Handles: common.HandleArray{m.CheckLogin, m.DelegateNow}},
		&common.ModuleHandles{Method: "post", Path: "/admin/manual_delegate_trade", Handles: common.HandleArray{m.CheckLogin, m.ExecuteManualDelegateTrade}},
		&common.ModuleHandles{Method: "post", Path: "/admin/delegate_history", Handles: common.HandleArray{m.CheckLogin, m.ListDelegateHistory}},
		&common.ModuleHandles{Method: "post", Path: "/admin/delegate_del", Handles: common.HandleArray{m.CheckLogin, m.DeleteDelegateHistory}},
		&common.ModuleHandles{Method: "post", Path: "admin/spot_delegate", Handles: common.HandleArray{m.CheckLogin, m.SpotDelegate}},
		&common.ModuleHandles{Method: "post", Path: "admin/opspot", Handles: common.HandleArray{m.CheckLogin, m.ReviewSpotDelegate}},
	}
}

func (m *AdminUserModule) ReviewSpotDelegate(r *gin.Context) {
	rq := make(adminservice.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.MODEL_TRADE.ReviewSpotDelegate(rq))
}

func (m *AdminUserModule) SpotDelegate(r *gin.Context) {
	rq := make(adminservice.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.MODEL_TRADE.ListDelegates(rq, true, ""))
}

func (m *AdminUserModule) DeleteDelegateHistory(r *gin.Context) {
	id := m.GetInt(r, "id")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.MODEL_TRADE.DeleteDelegateHistory(id))
}

func (m *AdminUserModule) ListDelegateHistory(r *gin.Context) {
	rq := make(adminservice.P, 0)
	m.ConvertObject(r, &rq)
	rq["state"] = 1
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.MODEL_TRADE.ListDelegates(rq, true, ""))
}

func (m *AdminUserModule) DelegateNow(r *gin.Context) {
	rq := make(adminservice.P, 0)
	m.ConvertObject(r, &rq)
	rq["state"] = 0
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.MODEL_TRADE.ListDelegates(rq, true, ""))
}

func (m *AdminUserModule) ExecuteManualDelegateTrade(r *gin.Context) {
	uid := m.GetInt(r, "uid")
	sn := m.GetValue(r, "sn")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.MODEL_TRADE.ExecuteManualDelegateTrade(uid, sn))
}

func (m *AdminUserModule) ListOpenTrades(r *gin.Context) {
	rq := make(adminservice.P, 0)
	m.ConvertObject(r, &rq)
	rq["state"] = 1
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.MODEL_TRADE.ListOpenTrades(rq, true, ""))
}

func (m *AdminUserModule) ListClosedTrades(r *gin.Context) {
	rq := make(adminservice.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.MODEL_TRADE.ListClosedTrades(rq, true, ""))
}

func (m *AdminUserModule) StopMiningOrder(r *gin.Context) {
	id := m.GetInt(r, "id")
	pass := m.GetValue(r, "pass")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.MODEL_USER.StopMiningOrder(id, pass))
}

func (m *AdminUserModule) ListMiningOrders(r *gin.Context) {
	rq := make(adminservice.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.MODEL_USER.ListMiningOrders(rq))
}

func (m *AdminUserModule) DeleteExplodeTrade(r *gin.Context) {
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.SYSTEM_MODEL.DeleteExplodeTrade(m.GetInt(r, "id")))
}

func (m *AdminUserModule) SaveExplodeTrade(r *gin.Context) {
	rq := make(adminservice.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.SYSTEM_MODEL.SaveExplodeTrade(rq))
}

func (m *AdminUserModule) ExplodeList(r *gin.Context) {
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.SYSTEM_MODEL.ExplodeTradeList())
}
