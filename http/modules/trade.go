package modules

import (
	"cointrade/http/common"
	"cointrade/models"

	"github.com/gin-gonic/gin"
)

//交易处理MODULE
type TradeModule struct {
	common.ModuleBase
}

func (m *TradeModule) ModuleList() common.MODULEHANDLELIST {
	//每个MODULE必须要实现的MODULE
	return common.MODULEHANDLELIST{
		&common.ModuleHandles{Method: "post", Path: "/trade/delegate", Handles: common.HandleArray{m.NeedLogin, m.Delegate}},
		&common.ModuleHandles{Method: "post", Path: "/trade/delegate/list", Handles: common.HandleArray{m.NeedLogin, m.DelegateList}}, //委托列表
		&common.ModuleHandles{Method: "post", Path: "/trade/opend/list", Handles: common.HandleArray{m.NeedLogin, m.OpendList}},       //持仓列表
		&common.ModuleHandles{Method: "post", Path: "/trade/close/list", Handles: common.HandleArray{m.NeedLogin, m.CloseList}},       //平仓列表
		&common.ModuleHandles{Method: "post", Path: "/trade/close/detail", Handles: common.HandleArray{m.NeedLogin, m.GetClose}},      //获取单个平仓详情
		&common.ModuleHandles{Method: "post", Path: "/trade/cancle", Handles: common.HandleArray{m.NeedLogin, m.CancleDelegate}},      //撤单
	}
}
func (m *TradeModule) GetClose(r *gin.Context) {
	uid := r.GetInt("uid")
	sn := m.GetValue(r, "sn")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_TRADE.GetCloseBySn(uid, sn))
}
func (m *TradeModule) Delegate(r *gin.Context) {
	//委托下单
	uid := r.GetInt("uid")
	var rq models.TradeDelegateRequest
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_TRADE.DelegateTrade(uid, &rq))
}
func (m *TradeModule) DelegateList(r *gin.Context) { //委托列表
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
func (m *TradeModule) OpendList(r *gin.Context) { //持仓列表
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
	//平仓列表
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
func (m *TradeModule) CancleDelegate(r *gin.Context) {
	//撤单
	uid := r.GetInt("uid")
	sn := m.GetValue(r, "sn")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_TRADE.CancleDelegate(uid, sn))
}
