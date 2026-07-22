package models

import (
	"cointrade/config"
	"cointrade/lib"
	"cointrade/lib/db"
	"cointrade/lib/notify"
	"cointrade/utils"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
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

func (m *CreditModel) GetAllRechargetAddress() db.DB_LIST_RESULT { //返回所有的充值钱包地址
	list, _ := config.GlobalDB.FetchRows(DB_TABLE_RECHARGE_ADDRESS, db.DB_PARAMS{"state": 1}, db.DB_FIELDS{})
	return list
}

func (m *CreditModel) CreateRecharge(uid int, rq *RechargeRequest) *RechargeResponse { //提交充值信息
	rs := new(RechargeResponse)
	uinfo := MODEL_USER.GetBaseInfo(uid)
	rechargeConfig := MODEL_SYSTEM.GetOneRechargeConfig(rq.CoinType, rq.Contract)
	if rechargeConfig == nil {
		rs.State = STATE_SYSTEM_ERROR
		rs.Msg = "system error"
		return rs
	}
	if uinfo == nil {
		rs.State = RECHARGE_STATE_ERROR_USER
		rs.Msg = "the user is not exists"
		return rs
	}
	if rq.Amount < rechargeConfig.Min {
		rs.State = RECHARGE_STATE_MIN
		rs.Msg = "too small"
		return rs
	}

	rate := 1.0
	cointype := strings.ToLower(rq.CoinType)
	if cointype != "usdt" {
		pair := fmt.Sprintf("%susdt", cointype)
		coinPriceInfo := MODEL_SYSTEM.GetLastCoinInfo(pair)
		if coinPriceInfo != nil {
			rate = coinPriceInfo["close"].(float64)
		}
	}
	sn := m.MakeOrderSn(uid, CREDIT_TYPE_RECHARGE)
	insertData := db.DB_PARAMS{
		"uid":         uid,
		"sn":          sn,
		"cointype":    rq.CoinType,
		"contract":    rq.Contract,
		"type":        0,
		"credit":      rq.Amount,
		"rate":        rate,
		"fact_credit": rq.Amount * rate,
		"createtime":  utils.GetNow(),
		"info":        rechargeConfig.Address,
		"txid":        "",
		"proof":       rq.Proof,
		"address":     rechargeConfig.Address,
	}
	_, err := config.GlobalDB.InsertData(DB_TABLE_RECHARGE, insertData)
	if err != nil {
		rs.State = STATE_SYSTEM_ERROR
		rs.Msg = err.Error()
		return rs
	}
	notify.NOTIFY.AddNotify(&notify.NotifyItem{Type: 2, Num: 1})
	rs.State = STATE_SUCCESS
	rs.Msg = "success"
	rs.Info = insertData
	rs.Sn = sn
	return rs
}

func (m *CreditModel) SuccessRecharge(sn string) bool {
	one := m.GetRechargeOrderBySn(sn)
	if one == nil {
		return false
	}
	ntime := utils.GetNow()
	config.GlobalDB.UpdateData(DB_TABLE_RECHARGE, db.DB_PARAMS{"state": 1, "finishtime": ntime}, db.DB_PARAMS{"id": one["id"]})
	cvalue := &CreditValue{
		Credit:          utils.GetFloat(one["fact_credit"]),
		VCrdit:          0,
		LockCredit:      0,
		LockVCredit:     0,
		UserCoinLogType: COIN_LOG_USER_RECHARGE,
		UserCoinLogInfo: QueueCreditLog{
			Credit:     utils.GetFloat(one["fact_credit"]),
			LockCredit: 0,
			Sn:         sn,
			CreateTime: ntime,
		},
		TeamCoinLogType: TEAM_LOG_RECHARGE,
		TeamCoinLogInfo: QueueTeamLog{
			Recharge:   utils.GetFloat(one["fact_credit"]),
			CreateTime: ntime,
		},
	}
	return MODEL_USER.AddCredit(utils.GetInt(one["uid"]), cvalue)
}

func (m *CreditModel) GetRechargeOrderBySn(sn string) db.DB_ROW_RESULT {
	one, _ := config.GlobalDB.FetchRow(DB_TABLE_RECHARGE, db.DB_PARAMS{"sn": sn}, db.DB_FIELDS{})
	return one
}

func (m *CreditModel) GetRechargetList(uid int, rq *PageBaseRequest) *PageBaseResponse { //充值记录获取
	condition := db.DB_PARAMS{"uid": uid}
	count := config.GlobalDB.GetCount(DB_TABLE_RECHARGE, condition)
	pagesize := int(math.Ceil(float64(count) / float64(rq.Limit)))
	if rq.Page > pagesize {
		rq.Page = pagesize
	}
	if rq.Page == 0 {
		rq.Page = 1
	}
	limitstr := fmt.Sprintf("limit %d,%d", (rq.Page-1)*rq.Limit, rq.Limit)
	list, _ := config.GlobalDB.FetchRows(DB_TABLE_RECHARGE, condition, db.DB_FIELDS{}, "order by id desc", limitstr)
	rs := new(PageBaseResponse)
	rs.Limit = rq.Limit
	rs.Page = rq.Page
	rs.PageTotal = pagesize
	rs.Total = count
	rs.List = list
	return rs
}

func (m *CreditModel) RechargeInfo(uid int, sn string) db.DB_ROW_RESULT {
	one, _ := config.GlobalDB.FetchRow(DB_TABLE_RECHARGE, db.DB_PARAMS{"uid": uid, "sn": sn}, db.DB_FIELDS{})
	return one
}

func (m *CreditModel) RechargeByApprove(uid int, amount float64) *BaseResponse { //通过授权充值
	rs := new(BaseResponse)
	uinfo := MODEL_USER.GetBaseInfo(uid)
	if uinfo == nil {
		rs.State = STATE_SYSTEM_ERROR
		rs.Msg = "system error"
		return rs
	}
	if uinfo.ApproveState != 1 {
		rs.State = RECHARGE_STATE_ERROR_NOTAPPROVE
		rs.Msg = "not approve"
		return rs
	}
	erc := new(lib.EthLib)
	erc.CreateClient()
	erc.Type = "usdt"
	defer erc.Close()
	balance, _ := erc.GetBalanceOfUsdt(uinfo.WalletAddress).Float64()
	if balance < amount {
		rs.State = RECHARGE_STATE_ERROR_MONEY
		rs.Msg = "not enough usdt"
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
		rs.Msg = "trans faild"
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
	config.GlobalDB.InsertData(DB_TABLE_RECHAGE_APPROVE, insertData)
	rs.State = STATE_SUCCESS
	rs.Msg = "success"
	return rs
}
