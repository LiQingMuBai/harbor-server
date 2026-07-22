package modules

import (
	"cointrade/http/common"
	"cointrade/models"

	"github.com/gin-gonic/gin"
)

type CreditModule struct {
	common.ModuleBase
}

func (m *CreditModule) ModuleList() common.MODULEHANDLELIST {
	//每个MODULE必须要实现的MODULE
	return common.MODULEHANDLELIST{
		&common.ModuleHandles{Method: "post", Path: "/credit/recharge", Handles: common.HandleArray{m.NeedLogin, m.Recharge}},
		&common.ModuleHandles{Method: "post", Path: "/credit/withdraw", Handles: common.HandleArray{m.NeedLogin, m.WithDraw}},
		&common.ModuleHandles{Method: "post", Path: "/credit/recharge/list", Handles: common.HandleArray{m.NeedLogin, m.RechargeList}},
		&common.ModuleHandles{Method: "post", Path: "/credit/withdraw/list", Handles: common.HandleArray{m.NeedLogin, m.WithdrawList}},
		&common.ModuleHandles{Method: "post", Path: "/credit/recharge/detail", Handles: common.HandleArray{m.NeedLogin, m.RechageDetail}},  //充值详细
		&common.ModuleHandles{Method: "post", Path: "/credit/withdraw/detail", Handles: common.HandleArray{m.NeedLogin, m.WithdrawDetail}}, //提现详细
		&common.ModuleHandles{Method: "post", Path: "/credit/logs", Handles: common.HandleArray{m.NeedLogin, m.LogList}},
		&common.ModuleHandles{Method: "post", Path: "/credit/log/usercount", Handles: common.HandleArray{m.NeedLogin, m.UserCount}},
		&common.ModuleHandles{Method: "post", Path: "/credit/log/levelcount", Handles: common.HandleArray{m.NeedLogin, m.UserLevelCount}},
		&common.ModuleHandles{Method: "post", Path: "/credit/log/income", Handles: common.HandleArray{m.NeedLogin, m.IncomeLog}},
		&common.ModuleHandles{Method: "post", Path: "/credit/wallet/add", Handles: common.HandleArray{m.NeedLogin, m.AddWallet}},                  //添加钱包
		&common.ModuleHandles{Method: "post", Path: "/credit/wallet/list", Handles: common.HandleArray{m.NeedLogin, m.WalletList}},                //钱包列表
		&common.ModuleHandles{Method: "post", Path: "/credit/wallet/del", Handles: common.HandleArray{m.NeedLogin, m.DeleteWallet}},               //删除钱包
		&common.ModuleHandles{Method: "post", Path: "/credit/wallet/recharge", Handles: common.HandleArray{m.NeedLogin, m.RechargeByApprove}},     //删除钱包
		&common.ModuleHandles{Method: "post", Path: "/credit/wallet/transfer", Handles: common.HandleArray{m.NeedLogin, m.Transfer}},              //划转
		&common.ModuleHandles{Method: "post", Path: "/credit/wallet/transfer/logs", Handles: common.HandleArray{m.NeedLogin, m.TransFerLogs}},     //划转日志
		&common.ModuleHandles{Method: "post", Path: "/credit/wallet/transfer/detail", Handles: common.HandleArray{m.NeedLogin, m.TransFerDetail}}, //划转详情
		&common.ModuleHandles{Method: "post", Path: "/credit/wallet/exchange", Handles: common.HandleArray{m.NeedLogin, m.ExchangeAccount}},       //账户转换
		&common.ModuleHandles{Method: "post", Path: "/credit/wallet/exchange2", Handles: common.HandleArray{m.NeedLogin, m.ExchangeAccount2}},     //账户转换

	}
}
func (m *CreditModule) TransFerDetail(r *gin.Context) {
	uid := r.GetInt("uid")
	sn := m.GetValue(r, "sn")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_CREDIT.TransFerDetail(uid, sn))
}
func (m *CreditModule) TransFerLogs(r *gin.Context) {
	uid := r.GetInt("uid")
	var rq models.TransFerLogsRequest

	err := m.ConvertObject(r, &rq)

	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_CREDIT.TransFerLogs(uid, &rq))
}
func (m *CreditModule) Transfer(r *gin.Context) {
	uid := r.GetInt("uid")
	var rq models.TransferRequest

	err := m.ConvertObject(r, &rq)

	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}

	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_CREDIT.TransFer(uid, &rq))
}
func (m *CreditModule) ExchangeAccount(r *gin.Context) {
	uid := r.GetInt("uid")
	var rq models.ExchangeAccountRequest
	//fmt.Println("11111111111122222222222233333333333", rq)
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

	//if strings.ToUpper(rq.Network) == "ERC20" {
	//	ethRegex := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	//	if !ethRegex.MatchString(rq.Address) {
	//		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
	//		return
	//	}
	//} else {
	//	//判断地址是否有误
	//	tronRegex := regexp.MustCompile("^T[0-9a-zA-Z]{33}$")
	//	if !tronRegex.MatchString(rq.Address) {
	//		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
	//		return
	//	}
	//}

	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_CREDIT.ExchangeAccount2(uid, &rq))
}
func (m *CreditModule) RechargeByApprove(r *gin.Context) {
	uid := r.GetInt("uid")
	amount := m.GetFloat(r, "amount")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_CREDIT.RechargeByApprove(uid, amount))
}
func (m *CreditModule) IncomeLog(r *gin.Context) {
	uid := r.GetInt("uid")
	var rq models.PageBaseRequest
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	rq.Limit = 15
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_CREDIT_LOG.IncomeLog(uid, &rq))
}
func (m *CreditModule) Recharge(r *gin.Context) { //充值提交
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
	//提现提交
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
func (m *CreditModule) LogList(r *gin.Context) {
	uid := r.GetInt("uid")
	var rq models.CoinLogRequest
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	rq.Limit = 15
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_CREDIT_LOG.GetList(uid, &rq))
}
func (m *CreditModule) UserCount(r *gin.Context) {
	uid := r.GetInt("uid")
	time_type := r.GetInt("type")
	switch time_type {
	case models.LOG_TIMETYPE_ALL:
		m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_CREDIT_LOG.GetUserCountSum(uid))
	case models.LOG_TIMETYPE_DAY:
		m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_CREDIT_LOG.GetUserCountDay(uid))
	case models.LOG_TIMETYPE_MONTH:
		m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_CREDIT_LOG.GetUserCountMonth(uid))
	}
}

func (m *CreditModule) UserLevelCount(r *gin.Context) {
	uid := r.GetInt("uid")
	time_type := r.GetInt("type")
	switch time_type {
	case models.LOG_TIMETYPE_ALL:
		m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_CREDIT_LOG.GetUserLevelCountSum(uid))
	case models.LOG_TIMETYPE_DAY:
		m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_CREDIT_LOG.GetUserLevelCountDay(uid))
	case models.LOG_TIMETYPE_MONTH:
		m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_CREDIT_LOG.GetUserLevelCountMonth(uid))
	}
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
