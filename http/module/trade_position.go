package module

import (
	"cointrade/http/common"
	"cointrade/models"

	"github.com/gin-gonic/gin"
)

func (m *TradeModule) positionRoutes() common.MODULEHANDLELIST {
	return common.MODULEHANDLELIST{
		&common.ModuleHandles{Method: "post", Path: "/trade/opend/list", Handles: common.HandleArray{m.NeedLogin, m.OpendList}},
		&common.ModuleHandles{Method: "post", Path: "/trade/close/list", Handles: common.HandleArray{m.NeedLogin, m.CloseList}},
		&common.ModuleHandles{Method: "post", Path: "/trade/close/detail", Handles: common.HandleArray{m.NeedLogin, m.GetClose}},
	}
}

func (m *TradeModule) GetClose(r *gin.Context) {
	uid := r.GetInt("uid")
	sn := m.GetValue(r, "sn")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_TRADE.GetCloseBySn(uid, sn))
}

func (m *TradeModule) OpendList(r *gin.Context) {
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
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_TRADE.GetOpendList(uid, &rq))
}

func (m *TradeModule) CloseList(r *gin.Context) {
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
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_TRADE.GetCloseList(uid, &rq))
}
