package models

import (
	"cointrade/config"
	"cointrade/lib/db"
	"cointrade/lib/notify"
	"cointrade/utils"
	"fmt"
	"math"
	"strings"
)

func (m *CreditModel) CreateWithDraw(uid int, rq *WithDrawRequest) *RechargeResponse {
	rs := new(RechargeResponse)
	sn := m.MakeOrderSn(uid, CREDIT_TYPE_WITHDRAW)
	ntime := utils.GetNow()
	uinfo := MODEL_USER.GetBaseInfo(uid)
	if uinfo == nil {
		rs.State = WIDTHDRAW_STATE_ERROR_USER
		rs.Msg = "error user"
		return rs
	}
	if uinfo.IsWithDraw != 1 {
		rs.State = WIDTHDRAW_STATE_ERROR_LOCKED
		rs.Msg = "user not allowed withdraw"
		return rs
	}
	cointype := strings.ToLower(rq.CoinType)
	if cointype != "bank" {
		rechargeConfig := MODEL_SYSTEM.GetOneRechargeConfig(rq.CoinType, rq.Contract)
		if rechargeConfig == nil {
			rs.State = STATE_SYSTEM_ERROR
			rs.Msg = "system error"
			return rs
		}
		if rq.Amount < rechargeConfig.Min {
			rs.State = WIDTHDRAW_STATE_MIN
			rs.Msg = "too min"
			return rs
		}
	} else if rq.Amount < config.GlobalConfig.GetValue("min_withdraw").ToFloat() {
		rs.State = WIDTHDRAW_STATE_MIN
		rs.Msg = "too min"
		return rs
	}

	rate := 1.0
	bankinfo := m.GetBankInfo(uid)
	if cointype != "usdt" {
		if cointype != "bank" {
			pair := fmt.Sprintf("%susdt", cointype)
			coinPriceInfo := MODEL_SYSTEM.GetLastCoinInfo(pair)
			if coinPriceInfo != nil {
				rate = coinPriceInfo["close"].(float64)
			}
		} else {
			if bankinfo == nil {
				rs.State = WIDTHDRAW_STATE_ERROR_NOTBINDBANK
				rs.Msg = "not bind bank info"
				return rs
			}
			rate = 1
		}
	}

	factCredit := rq.Amount * rate
	if uinfo.Credit < factCredit {
		rs.State = WIDTHDRAW_STATE_NOTENOUGH
		rs.Msg = "no more credit"
		return rs
	}

	insertData := db.DB_PARAMS{
		"uid":         uid,
		"credit":      rq.Amount,
		"rate":        rate,
		"fact_credit": factCredit,
		"cointype":    rq.CoinType,
		"contract":    rq.Contract,
		"address":     rq.Address,
		"fee":         factCredit * config.GlobalConfig.GetValue("withdraw_fee").ToFloat() / float64(100),
		"info":        "",
		"createtime":  ntime,
		"sn":          sn,
		"state":       0,
		"finishtime":  0,
		"memo":        "",
	}
	if cointype == "bank" {
		insertData["type"] = 1
		insertData["bankinfo"] = bankinfo
	}
	_, err := config.GlobalDB.InsertData(DB_TABLE_WITHDRAW, insertData)
	if err != nil {
		rs.State = STATE_SYSTEM_ERROR
		rs.Msg = "system error"
		return rs
	}
	MODEL_USER.AddCredit(uid, &CreditValue{
		Credit:          -1 * factCredit,
		LockCredit:      factCredit,
		LockVCredit:     0,
		VCrdit:          0,
		UserCoinLogType: COIN_LOG_USER_WITHDRAW,
		UserCoinLogInfo: QueueCreditLog{
			Credit:     -1 * factCredit,
			LockCredit: factCredit,
			Sn:         sn,
			CreateTime: ntime,
		},
	})
	notify.NOTIFY.AddNotify(&notify.NotifyItem{Type: 1, Num: 1})
	rs.State = STATE_SUCCESS
	rs.Msg = "ok"
	rs.Info = insertData
	rs.Sn = sn
	return rs
}

func (m *CreditModel) GetWithDrawList(uid int, rq *PageBaseRequest) *PageBaseResponse {
	condition := db.DB_PARAMS{"uid": uid}
	count := config.GlobalDB.GetCount(DB_TABLE_WITHDRAW, condition)
	pagesize := int(math.Ceil(float64(count) / float64(rq.Limit)))
	if rq.Page > pagesize {
		rq.Page = pagesize
	}
	if rq.Page == 0 {
		rq.Page = 1
	}
	limitstr := fmt.Sprintf("limit %d,%d", (rq.Page-1)*rq.Limit, rq.Limit)
	list, _ := config.GlobalDB.FetchRows(DB_TABLE_WITHDRAW, condition, db.DB_FIELDS{}, "order by id desc", limitstr)
	rs := new(PageBaseResponse)
	rs.Limit = rq.Limit
	rs.Page = rq.Page
	rs.PageTotal = pagesize
	rs.Total = count
	rs.List = list
	return rs
}

func (m *CreditModel) WithdrawInfo(uid int, sn string) db.DB_ROW_RESULT {
	one, _ := config.GlobalDB.FetchRow(DB_TABLE_WITHDRAW, db.DB_PARAMS{"uid": uid, "sn": sn}, db.DB_FIELDS{})
	return one
}
