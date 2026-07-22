package modules

import (
	"cointrade/http/common"
	"cointrade/models"

	"github.com/gin-gonic/gin"
)

func (m *CreditModule) fundingRoutes() common.MODULEHANDLELIST {
	return common.MODULEHANDLELIST{
		&common.ModuleHandles{Method: "post", Path: "/credit/recharge", Handles: common.HandleArray{m.NeedLogin, m.Recharge}},
		&common.ModuleHandles{Method: "post", Path: "/credit/withdraw", Handles: common.HandleArray{m.NeedLogin, m.WithDraw}},
		&common.ModuleHandles{Method: "post", Path: "/credit/recharge/list", Handles: common.HandleArray{m.NeedLogin, m.RechargeList}},
		&common.ModuleHandles{Method: "post", Path: "/credit/withdraw/list", Handles: common.HandleArray{m.NeedLogin, m.WithdrawList}},
		&common.ModuleHandles{Method: "post", Path: "/credit/recharge/detail", Handles: common.HandleArray{m.NeedLogin, m.RechageDetail}},
		&common.ModuleHandles{Method: "post", Path: "/credit/withdraw/detail", Handles: common.HandleArray{m.NeedLogin, m.WithdrawDetail}},
	}
}

func (m *CreditModule) Recharge(r *gin.Context) {
	uid := r.GetInt("uid")
	var rq models.RechargeRequest
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_CREDIT.CreateRecharge(uid, &rq))
}

func (m *CreditModule) WithDraw(r *gin.Context) {
	uid := r.GetInt("uid")
	var rq models.WithDrawRequest
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_CREDIT.CreateWithDraw(uid, &rq))
}

func (m *CreditModule) RechargeList(r *gin.Context) {
	uid := r.GetInt("uid")
	var rq models.PageBaseRequest
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	rq.Limit = 15
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_CREDIT.GetRechargetList(uid, &rq))
}

func (m *CreditModule) WithdrawList(r *gin.Context) {
	uid := r.GetInt("uid")
	var rq models.PageBaseRequest
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	rq.Limit = 15
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_CREDIT.GetWithDrawList(uid, &rq))
}

func (m *CreditModule) RechageDetail(r *gin.Context) {
	uid := r.GetInt("uid")
	sn := m.GetValue(r, "sn")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_CREDIT.RechargeInfo(uid, sn))
}

func (m *CreditModule) WithdrawDetail(r *gin.Context) {
	uid := r.GetInt("uid")
	sn := m.GetValue(r, "sn")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_CREDIT.WithdrawInfo(uid, sn))
}
