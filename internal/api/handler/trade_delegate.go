package handler

import (
	"cointrade/http/common"
	"cointrade/models"

	"github.com/gin-gonic/gin"
)

func (m *TradeModule) delegateRoutes() common.MODULEHANDLELIST {
	return common.MODULEHANDLELIST{
		&common.ModuleHandles{Method: "post", Path: "/trade/delegate", Handles: common.HandleArray{m.NeedLogin, m.Delegate}},
		&common.ModuleHandles{Method: "post", Path: "/trade/delegate/list", Handles: common.HandleArray{m.NeedLogin, m.DelegateList}},
		&common.ModuleHandles{Method: "post", Path: "/trade/cancle", Handles: common.HandleArray{m.NeedLogin, m.CancleDelegate}},
	}
}

func (m *TradeModule) Delegate(r *gin.Context) {
	uid := r.GetInt("uid")
	var rq models.TradeDelegateRequest
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_TRADE.DelegateTrade(uid, &rq))
}

func (m *TradeModule) DelegateList(r *gin.Context) {
	uid := r.GetInt("uid")
	var rq models.TradeListRequest
	rq.Limit = 15
	err := m.ConvertObject(r, &rq)
	if rq.Limit > 100 {
		rq.Limit = 100
	}
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_TRADE.GetDelegateList(uid, &rq))
}

func (m *TradeModule) CancleDelegate(r *gin.Context) {
	uid := r.GetInt("uid")
	sn := m.GetValue(r, "sn")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_TRADE.CancleDelegate(uid, sn))
}
