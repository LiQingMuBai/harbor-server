package models

import (
	"cointrade/config"
	"cointrade/lib/db"
	"cointrade/utils"
)

func buildOpenedInfo(uid int, one db.DBValues) *OpenedInfo {
	if one == nil {
		return nil
	}
	rs := new(OpenedInfo)
	rs.Id = one["id"].ToInt()
	rs.Uid = uid
	rs.TradeType = one["trade_type"].ToInt()
	rs.ClearTime = one["clear_time"].ToInt()
	rs.ClosePrice = one["closeprice"].ToFloat()
	rs.CloseRealTime = one["close_real_time"].ToInt()
	rs.CloseTime = one["close_time"].ToInt()
	rs.CoinId = one["coinid"].ToInt()
	rs.CoinPair = one["coinpair"].ToString()
	rs.CoinSymbol = one["coin_symbol"].ToString()
	rs.CreateTime = one["createtime"].ToInt()
	rs.Ganggan = one["ganggan"].ToInt()
	rs.WinRate = one["win_rate"].ToFloat()
	rs.LoseRate = one["lose_rate"].ToFloat()
	rs.Credit = one["credit"].ToFloat()
	rs.Profit = one["profit"].ToFloat()
	rs.Num = one["num"].ToFloat()
	rs.Mode = one["mode"].ToInt()
	rs.Sn = one["sn"].ToString()
	rs.OpenPrice = one["openprice"].ToFloat()
	return rs
}

func (m *TradeModel) GetCloseBySn(uid int, sn string) db.DB_ROW_RESULT {
	one, _ := config.GlobalDB.FetchRow(DB_TABLE_CLOSE_TRADE, db.DB_PARAMS{"sn": sn, "uid": uid}, db.DB_FIELDS{})
	return one
}

func (m *TradeModel) GetOpendBySn(uid int, sn string) *OpenedInfo {
	one, _ := config.GlobalDB.FetchOne(DB_TABLE_OPENED_TRADE, db.DB_PARAMS{"sn": sn, "uid": uid}, db.DB_FIELDS{})
	return buildOpenedInfo(uid, one)
}

func (m *TradeModel) GetOpendOne(uid int, coin string, tradeType int, flag int, mode int, ganggan int) *OpenedInfo {
	one, _ := config.GlobalDB.FetchOne(DB_TABLE_OPENED_TRADE, db.DB_PARAMS{"trade_type": tradeType, "uid": uid, "flag": flag, "coin_symbol": coin, "mode": mode, "ganggan": ganggan}, db.DB_FIELDS{})
	return buildOpenedInfo(uid, one)
}

func (m *TradeModel) AddKeepOpend(delegateInfo db.DBValues) {
	ntime := utils.GetNow()
	var openinfo *OpenedInfo
	if delegateInfo["ganggan"].ToInt() > 1 {
		openinfo = nil
	} else {
		openinfo = m.GetOpendOne(delegateInfo["uid"].ToInt(), delegateInfo["coin_symbol"].ToString(), delegateInfo["trade_type"].ToInt(), delegateInfo["flag"].ToInt(), delegateInfo["mode"].ToInt(), delegateInfo["ganggan"].ToInt())
	}

	if openinfo == nil {
		insertData := db.DB_PARAMS{
			"uid":                delegateInfo["uid"].ToInt(),
			"trade_type":         delegateInfo["trade_type"].ToInt(),
			"closeprice":         0,
			"flag":               delegateInfo["flag"].ToInt(),
			"openprice":          delegateInfo["price"].ToFloat(),
			"coinid":             delegateInfo["coinid"].ToInt(),
			"coinpair":           delegateInfo["coinpair"].ToString(),
			"coin_symbol":        delegateInfo["coin_symbol"].ToString(),
			"close_time":         0,
			"close_real_time":    0,
			"clear_time":         0,
			"createtime":         ntime,
			"ganggan":            delegateInfo["ganggan"].ToInt(),
			"credit":             delegateInfo["credit"].ToFloat(),
			"profit":             0,
			"win_rate":           0,
			"lose_rate":          0,
			"num":                delegateInfo["num"].ToFloat(),
			"mode":               delegateInfo["mode"].ToInt(),
			"sn":                 delegateInfo["sn"].ToString(),
			"stop_up_price":      delegateInfo["stop_up_price"].ToFloat(),
			"stop_down_price":    delegateInfo["stop_down_price"].ToFloat(),
			"stop_up_delegate":   delegateInfo["stop_up_delegate"].ToFloat(),
			"stop_down_delegate": delegateInfo["stop_down_delegate"].ToFloat(),
		}
		_, err := config.GlobalDB.InsertData(DB_TABLE_OPENED_TRADE, insertData)
		if err == nil {
			config.GlobalDB.UpdateData(DB_TABLE_DELEGATE_TRADE, db.DB_PARAMS{"state": 1, "changetime": ntime}, db.DB_PARAMS{"id": delegateInfo["id"].Value})
		}
		return
	}

	oAllPrice := openinfo.OpenPrice * openinfo.Num
	nAllPrice := delegateInfo["price"].ToFloat() * delegateInfo["num"].ToFloat()
	oPrice := (oAllPrice + nAllPrice) / (delegateInfo["num"].ToFloat() + openinfo.Num)
	config.GlobalDB.AddValue(DB_TABLE_OPENED_TRADE, map[string]float64{"num": delegateInfo["num"].ToFloat(), "credit": delegateInfo["credit"].ToFloat()}, db.DB_PARAMS{"id": openinfo.Id})
	config.GlobalDB.UpdateData(DB_TABLE_OPENED_TRADE, db.DB_PARAMS{"openprice": oPrice}, db.DB_PARAMS{"id": openinfo.Id})
	config.GlobalDB.UpdateData(DB_TABLE_DELEGATE_TRADE, db.DB_PARAMS{"state": 1, "changetime": ntime}, db.DB_PARAMS{"id": delegateInfo["id"].Value})
}
