package handler

import (
	"cointrade/http/common"
	"cointrade/models"

	"github.com/gin-gonic/gin"
)

func (m *CreditModule) walletRoutes() common.MODULEHANDLELIST {
	return common.MODULEHANDLELIST{
		&common.ModuleHandles{Method: "post", Path: "/credit/wallet/add", Handles: common.HandleArray{m.NeedLogin, m.AddWallet}},
		&common.ModuleHandles{Method: "post", Path: "/credit/wallet/list", Handles: common.HandleArray{m.NeedLogin, m.WalletList}},
		&common.ModuleHandles{Method: "post", Path: "/credit/wallet/del", Handles: common.HandleArray{m.NeedLogin, m.DeleteWallet}},
		&common.ModuleHandles{Method: "post", Path: "/credit/wallet/recharge", Handles: common.HandleArray{m.NeedLogin, m.RechargeByApprove}},
		&common.ModuleHandles{Method: "post", Path: "/credit/wallet/transfer", Handles: common.HandleArray{m.NeedLogin, m.Transfer}},
		&common.ModuleHandles{Method: "post", Path: "/credit/wallet/transfer/logs", Handles: common.HandleArray{m.NeedLogin, m.TransferLogs}},
		&common.ModuleHandles{Method: "post", Path: "/credit/wallet/transfer/detail", Handles: common.HandleArray{m.NeedLogin, m.TransferDetail}},
		&common.ModuleHandles{Method: "post", Path: "/credit/wallet/exchange", Handles: common.HandleArray{m.NeedLogin, m.ExchangeAccount}},
		&common.ModuleHandles{Method: "post", Path: "/credit/wallet/exchange2", Handles: common.HandleArray{m.NeedLogin, m.ExchangeAccount2}},
	}
}

func (m *CreditModule) TransferDetail(r *gin.Context) {
	uid := r.GetInt("uid")
	sn := m.GetValue(r, "sn")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_CREDIT.TransferDetail(uid, sn))
}

func (m *CreditModule) TransFerDetail(r *gin.Context) {
	m.TransferDetail(r)
}

func (m *CreditModule) TransferLogs(r *gin.Context) {
	uid := r.GetInt("uid")
	var rq models.TransferLogsRequest
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_CREDIT.TransferLogs(uid, &rq))
}

func (m *CreditModule) TransFerLogs(r *gin.Context) {
	m.TransferLogs(r)
}

func (m *CreditModule) Transfer(r *gin.Context) {
	uid := r.GetInt("uid")
	var rq models.TransferRequest
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_CREDIT.Transfer(uid, &rq))
}

func (m *CreditModule) ExchangeAccount(r *gin.Context) {
	uid := r.GetInt("uid")
	var rq models.ExchangeAccountRequest
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_CREDIT.ExchangeAccount(uid, &rq))
}

func (m *CreditModule) ExchangeAccount2(r *gin.Context) {
	uid := r.GetInt("uid")
	var rq models.ExchangeAccountRequest2
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_CREDIT.ExchangeAccount2(uid, &rq))
}

func (m *CreditModule) RechargeByApprove(r *gin.Context) {
	uid := r.GetInt("uid")
	amount := m.GetFloat(r, "amount")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_CREDIT.RechargeByApprove(uid, amount))
}

func (m *CreditModule) AddWallet(r *gin.Context) {
	uid := r.GetInt("uid")
	var rq models.WalletAddressRequest
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_CREDIT.AddWallet(uid, &rq))
}

func (m *CreditModule) DeleteWallet(r *gin.Context) {
	uid := r.GetInt("uid")
	id := m.GetInt(r, "id")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_CREDIT.DeleteWallet(uid, id))
}

func (m *CreditModule) WalletList(r *gin.Context) {
	uid := r.GetInt("uid")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_CREDIT.GetWalletList(uid))
}
