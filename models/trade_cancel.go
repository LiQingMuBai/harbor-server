package models

import (
	"cointrade/config"
	"cointrade/lib/db"
	"cointrade/utils"
)

func (m *TradeModel) CancleDelegate(uid int, sn string) *BaseResponse {
	rs := new(BaseResponse)
	one, _ := config.GlobalDB.FetchOne(DB_TABLE_DELEGATE_TRADE, db.DB_PARAMS{"uid": uid, "sn": sn, "state": 0}, db.DB_FIELDS{})
	if one == nil {
		rs.State = STATE_FAILD
		rs.Msg = "no this delegate order"
		return rs
	}

	userCredit := one["credit"].ToFloat() + one["fee"].ToFloat()
	userVCredit := one["credit"].ToFloat() + one["fee"].ToFloat()
	if one["mode"].ToInt() == USER_MODE_REAL {
		userVCredit = 0
	} else {
		userCredit = 0
	}
	if one["delegate_type"].ToInt() == DELEGATE_TYPE_BUY {
		if MODEL_USER.AddCredit(uid, &CreditValue{
			Credit:          userCredit,
			LockCredit:      -1 * userCredit,
			VCrdit:          userVCredit,
			LockVCredit:     -1 * userVCredit,
			UserCoinLogType: COIN_LOG_USER_CANCLE,
			UserCoinLogInfo: QueueCreditLog{
				Credit:     userCredit,
				LockCredit: -1 * userCredit,
				Sn:         one["sn"].ToString(),
				CreateTime: utils.GetNow(),
			},
			TeamCoinLogType: 0,
			TeamCoinLogInfo: nil,
		}) {
			_, err := config.GlobalDB.UpdateData(DB_TABLE_DELEGATE_TRADE, db.DB_PARAMS{"state": 2}, db.DB_PARAMS{"id": one["id"].Value})
			if err != nil {
				rs.State = STATE_SYSTEM_ERROR
				rs.Msg = err.Error()
				return rs
			}
		}
		rs.State = STATE_SUCCESS
		rs.Msg = "success"
		return rs
	}

	switch one["trade_type"].ToInt() {
	case OPEN_TYPE_BB:
		if MODEL_ASSETS.AddAssets(uid, &Assets{
			Coin:    one["coin_symbol"].ToString(),
			Pair:    one["coinpair"].ToString(),
			Num:     one["num"].ToFloat(),
			LockNum: -1 * one["num"].ToFloat(),
			Price:   one["price"].ToFloat(),
			Mode:    one["mode"].ToInt(),
		}) {
			_, err := config.GlobalDB.UpdateData(DB_TABLE_DELEGATE_TRADE, db.DB_PARAMS{"state": 2}, db.DB_PARAMS{"id": one["id"].Value})
			if err != nil {
				rs.State = STATE_SYSTEM_ERROR
				rs.Msg = err.Error()
				return rs
			}
		}
	case OPEN_TYPE_KEEP:
		flag := one["flag"].ToInt()
		coin := one["coin_symbol"].ToString()
		uinfo := MODEL_USER.GetBaseInfo(uid)
		num := one["num"].ToFloat()
		opendinfo := m.GetOpendOne(uid, coin, OPEN_TYPE_KEEP, flag, uinfo.Mode, one["ganggan"].ToInt())
		if opendinfo != nil {
			config.GlobalDB.AddValue(DB_TABLE_OPENED_TRADE, map[string]float64{"num": num, "lock_num": -1 * num}, db.DB_PARAMS{"id": opendinfo.Id})
			_, err := config.GlobalDB.UpdateData(DB_TABLE_DELEGATE_TRADE, db.DB_PARAMS{"state": 2}, db.DB_PARAMS{"id": one["id"].Value})
			if err != nil {
				rs.State = STATE_SYSTEM_ERROR
				rs.Msg = err.Error()
				return rs
			}
		}
	}

	rs.State = STATE_SUCCESS
	rs.Msg = "success"
	return rs
}
