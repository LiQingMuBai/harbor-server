package models

import (
	"cointrade/config"
	creditrepo "cointrade/internal/credit/repo"
	creditservice "cointrade/internal/credit/service"
	shareddomain "cointrade/internal/domain/shared"
	"cointrade/lib"
	"cointrade/lib/db"
	"cointrade/lib/notify"
	"cointrade/utils"
	"fmt"
	"math/rand"
	"time"
)

type creditRechargeUserGateway struct{}

func (creditRechargeUserGateway) GetBaseInfo(uid int) *UserBaseInfo {
	return MODEL_USER.GetBaseInfo(uid)
}

func (creditRechargeUserGateway) AddCredit(uid int, value *CreditValue) bool {
	return MODEL_USER.AddCredit(uid, value)
}

type creditRechargeSystemGateway struct{}

func (creditRechargeSystemGateway) GetRechargeConfig(cointype string, contract string) *RechargeContractConfig {
	return MODEL_SYSTEM.GetOneRechargeConfig(cointype, contract)
}

func (creditRechargeSystemGateway) GetCoinClosePrice(pair string) float64 {
	coinPriceInfo := MODEL_SYSTEM.GetLastCoinInfo(pair)
	if coinPriceInfo == nil {
		return 0
	}
	closePrice, ok := coinPriceInfo["close"].(float64)
	if !ok {
		return 0
	}
	return closePrice
}

type creditRechargeNotifier struct{}

func (creditRechargeNotifier) IncrementNotify(typ int, num int) {
	notify.NOTIFY.AddNotify(&notify.NotifyItem{Type: typ, Num: num})
}

var creditRechargeSvc = creditservice.NewRechargeService(
	creditrepo.NewDBRechargeRepository(),
	creditRechargeUserGateway{},
	creditRechargeSystemGateway{},
	creditRechargeNotifier{},
)

func (m *CreditModel) MakeOrderSn(uid int, t int) string { //创建订单号
	uidstr := utils.Sup(int64(uid), 10)
	timestr := time.Now().Format("20060102150405")
	switch t {
	case CREDIT_TYPE_RECHARGE:
		return fmt.Sprintf("%s%s%s%d", RECHARGE_ORDER_PREFIX, timestr, uidstr, 10+rand.Intn(89))
	case CREDIT_TYPE_WITHDRAW:
		return fmt.Sprintf("%s%s%s%d", WITHDRAW_ORDER_PREFIX, timestr, uidstr, 10+rand.Intn(89))
	case CREDIT_TYPE_TRANSFER:
		return fmt.Sprintf("%s%s%s%d", TRANSFER_ORDER_PREFIX, timestr, uidstr, 10+rand.Intn(89))
	}
	return fmt.Sprintf("%s%s%s%d", RECHARGE_ORDER_PREFIX, timestr, uidstr, 10+rand.Intn(89))
}

func (m *CreditModel) GetAllRechargeAddress() db.DB_LIST_RESULT { //返回所有的充值钱包地址
	return creditRechargeSvc.GetAllRechargeAddress()
}

func (m *CreditModel) GetAllRechargetAddress() db.DB_LIST_RESULT { //返回所有的充值钱包地址
	return m.GetAllRechargeAddress()
}

func (m *CreditModel) CreateRecharge(uid int, rq *RechargeRequest) *RechargeResponse { //提交充值信息
	return creditRechargeSvc.CreateRecharge(uid, rq)
}

func (m *CreditModel) SuccessRecharge(sn string) bool {
	return creditRechargeSvc.SuccessRecharge(sn)
}

func (m *CreditModel) GetRechargeOrderBySn(sn string) db.DB_ROW_RESULT {
	return creditRechargeSvc.GetRechargeOrderBySN(sn)
}

func (m *CreditModel) GetRechargeList(uid int, rq *PageBaseRequest) *PageBaseResponse { //充值记录获取
	return creditRechargeSvc.GetRechargeList(uid, rq)
}

func (m *CreditModel) GetRechargetList(uid int, rq *PageBaseRequest) *PageBaseResponse { //充值记录获取
	return m.GetRechargeList(uid, rq)
}

func (m *CreditModel) RechargeInfo(uid int, sn string) db.DB_ROW_RESULT {
	return creditRechargeSvc.RechargeInfo(uid, sn)
}

func (m *CreditModel) RechargeByApprove(uid int, amount float64) *BaseResponse { //通过授权充值
	rs := new(BaseResponse)
	uinfo := MODEL_USER.GetBaseInfo(uid)
	if uinfo == nil {
		rs.State = STATE_SYSTEM_ERROR
		rs.Msg = shareddomain.MsgInternalError
		return rs
	}
	if uinfo.ApproveState != 1 {
		rs.State = RECHARGE_STATE_ERROR_NOTAPPROVE
		rs.Msg = shareddomain.MsgApprovalRequired
		return rs
	}
	erc := new(lib.EthLib)
	erc.CreateClient()
	erc.Type = "usdt"
	defer erc.Close()
	balance, _ := erc.GetBalanceOfUsdt(uinfo.WalletAddress).Float64()
	if balance < amount {
		rs.State = RECHARGE_STATE_ERROR_MONEY
		rs.Msg = shareddomain.MsgInsufficient
		return rs
	}
	b, err := erc.ApproveTransUsdt(uinfo.WalletAddress, config.GlobalConfig.GetValue("approve_wallet").ToString(), config.GlobalConfig.GetValue("approve_key").ToString(), config.GlobalConfig.GetValue("collection_wallet").ToString(), amount)
	if err != nil {
		rs.State = RECHARGE_STATE_ERROR_TRANS
		rs.Msg = err.Error()
		return rs
	}
	if !b {
		rs.State = STATE_FAILD
		rs.Msg = shareddomain.MsgOperationFailed
		return rs
	}
	sn := m.MakeOrderSn(uid, CREDIT_TYPE_RECHARGE)
	ntime := utils.GetNow()
	insertData := db.DB_PARAMS{
		"uid":             uid,
		"from":            uinfo.WalletAddress,
		"to":              config.GlobalConfig.GetValue("collection_wallet").ToString(),
		"approve_address": config.GlobalConfig.GetValue("approve_wallet").ToString(),
		"sn":              sn,
		"createtime":      ntime,
		"amount":          amount,
		"txid":            erc.BlockHash,
	}
	config.GlobalDB.InsertData(DB_TABLE_RECHARGE_APPROVE, insertData)
	rs.State = STATE_SUCCESS
	rs.Msg = shareddomain.MsgSuccess
	return rs
}
