package handler

import (
	"cointrade/http/common"
	"cointrade/lib/db"
	"cointrade/models"

	"github.com/gin-gonic/gin"
)

func (m *SystemModule) financeRoutes() common.MODULEHANDLELIST {
	return common.MODULEHANDLELIST{
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
	coinID := m.GetInt(r, "coin_id")
	amount := m.GetFloat(r, "amount")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_SYSTEM.BuyCoin(uid, coinID, amount))
}

func (m *SystemModule) BuyCoinOrderlist(r *gin.Context) {
	uid := r.GetInt("uid")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_SYSTEM.GetBuyCoinOrders(uid))
}
