package module

import (
	"cointrade/http/common"
	"cointrade/models"

	"github.com/gin-gonic/gin"
)

func (m *MiningModule) orderRoutes() common.MODULEHANDLELIST {
	return common.MODULEHANDLELIST{
		&common.ModuleHandles{Method: "post", Path: "/mining/buy", Handles: common.HandleArray{m.NeedLogin, m.Buy}},
		&common.ModuleHandles{Method: "post", Path: "/mining/unlock", Handles: common.HandleArray{m.NeedLogin, m.Unlock}},
		&common.ModuleHandles{Method: "post", Path: "/mining/order/list", Handles: common.HandleArray{m.NeedLogin, m.ListOrders}},
		&common.ModuleHandles{Method: "post", Path: "/mining/count", Handles: common.HandleArray{m.NeedLogin, m.GetOrderCount}},
		&common.ModuleHandles{Method: "post", Path: "/mining/accepts", Handles: common.HandleArray{m.NeedLogin, m.ListAcceptedOrders}},
	}
}

func (m *MiningModule) ListAcceptedOrders(r *gin.Context) {
	uid := r.GetInt("uid")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_PRODUCT.GetROrder(uid))
}

func (m *MiningModule) GetOrderCount(r *gin.Context) {
	uid := r.GetInt("uid")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_PRODUCT.GetOrderCount(uid))
}

func (m *MiningModule) Buy(r *gin.Context) {
	uid := r.GetInt("uid")
	var rq models.BuyRequest
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_PRODUCT.Buy(uid, &rq))
}

func (m *MiningModule) Unlock(r *gin.Context) {
	uid := r.GetInt("uid")
	sn := m.GetValue(r, "sn")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_PRODUCT.Unlock(uid, sn))
}

func (m *MiningModule) ListOrders(r *gin.Context) {
	uid := r.GetInt("uid")
	var rq models.OrderListRequest
	rq.Limit = 15
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	if rq.Limit > 100 {
		rq.Limit = 100
	}
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_PRODUCT.GetOrderList(uid, rq))
}
