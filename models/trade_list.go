package models

import (
	"cointrade/config"
	"cointrade/lib/db"
	"cointrade/utils"
	"fmt"
	"math"
)

func (m *TradeModel) GetDelegateList(uid int, rq *TradeListRequest) *PageBaseResponse {
	uinfo := MODEL_USER.GetBaseInfo(uid)
	if uinfo == nil {
		return nil
	}
	condition := db.DB_PARAMS{"uid": uid, "mode": uinfo.Mode}
	if rq.Coin != "" {
		condition["coin_symbol"] = rq.Coin
	}
	if rq.Flag != 0 {
		condition["flag"] = rq.Flag
	}
	if rq.State != -1 {
		condition["state"] = rq.State
	}
	if rq.TradeType != 0 {
		condition["trade_type"] = rq.TradeType
	}
	if rq.DelegateType != 0 {
		condition["delegate_type"] = rq.DelegateType
	}
	if rq.Ganggan > 0 {
		condition["_"] = "num>0 and ganggan>1"
	} else {
		condition["_"] = "num>0 and ganggan<=1"
	}
	count := config.GlobalDB.GetCount(DB_TABLE_DELEGATE_TRADE, condition)
	pagesize := int(math.Ceil(float64(count) / float64(rq.Limit)))
	if rq.Page > pagesize {
		rq.Page = pagesize
	}
	if rq.Page <= 0 {
		rq.Page = 1
	}
	limitstr := fmt.Sprintf("limit %d,%d", (rq.Page-1)*rq.Limit, rq.Limit)
	list, _ := config.GlobalDB.FetchRows(DB_TABLE_DELEGATE_TRADE, condition, db.DB_FIELDS{}, "order by id desc", limitstr)
	rs := new(PageBaseResponse)
	rs.Limit = rq.Limit
	rs.Page = rq.Page
	rs.List = list
	rs.PageTotal = pagesize
	rs.State = STATE_SUCCESS
	rs.Msg = ""
	return rs
}

func (m *TradeModel) GetOpendList(uid int, rq *TradeListRequest) *PageBaseResponse {
	uinfo := MODEL_USER.GetBaseInfo(uid)
	if uinfo == nil {
		return nil
	}

	condition := db.DB_PARAMS{"uid": uid, "mode": uinfo.Mode, "clear_time": 0}
	if rq.Coin != "" {
		condition["coin_symbol"] = rq.Coin
	}
	if rq.Flag != 0 {
		condition["flag"] = rq.Flag
	}
	if rq.TradeType != 0 {
		condition["trade_type"] = rq.TradeType
	}
	if rq.Ganggan > 0 {
		condition["_"] = "num>0 and ganggan>1"
	} else {
		condition["_"] = "num>0 and ganggan<=1"
	}
	count := config.GlobalDB.GetCount(DB_TABLE_OPENED_TRADE, condition)
	pagesize := int(math.Ceil(float64(count) / float64(rq.Limit)))
	if rq.Page > pagesize {
		rq.Page = pagesize
	}
	if rq.Page <= 0 {
		rq.Page = 1
	}
	limitstr := fmt.Sprintf("limit %d,%d", (rq.Page-1)*rq.Limit, rq.Limit)
	list, _ := config.GlobalDB.FetchRows(DB_TABLE_OPENED_TRADE, condition, db.DB_FIELDS{}, "order by id desc", limitstr)

	rs := new(PageBaseResponse)
	rs.Limit = rq.Limit
	rs.Page = rq.Page
	rs.List = list
	rs.PageTotal = pagesize
	rs.State = STATE_SUCCESS
	rs.Msg = fmt.Sprintf("%d", utils.GetNow())
	return rs
}

func (m *TradeModel) GetCloseList(uid int, rq *TradeListRequest) *PageBaseResponse {
	uinfo := MODEL_USER.GetBaseInfo(uid)
	if uinfo == nil {
		return nil
	}

	condition := db.DB_PARAMS{"uid": uid, "mode": uinfo.Mode, "_": "num>0"}
	if rq.Coin != "" {
		condition["coin_symbol"] = rq.Coin
	}
	if rq.Flag != 0 {
		condition["flag"] = rq.Flag
	}
	if rq.TradeType != 0 {
		condition["trade_type"] = rq.TradeType
	}
	count := config.GlobalDB.GetCount(DB_TABLE_CLOSE_TRADE, condition)
	pagesize := int(math.Ceil(float64(count) / float64(rq.Limit)))
	if rq.Page > pagesize {
		rq.Page = pagesize
	}
	if rq.Page <= 0 {
		rq.Page = 1
	}
	limitstr := fmt.Sprintf("limit %d,%d", (rq.Page-1)*rq.Limit, rq.Limit)
	list, _ := config.GlobalDB.FetchRows(DB_TABLE_CLOSE_TRADE, condition, db.DB_FIELDS{}, "order by id desc", limitstr)
	rs := new(PageBaseResponse)
	rs.Limit = rq.Limit
	rs.Page = rq.Page
	rs.List = list
	rs.PageTotal = pagesize
	rs.State = STATE_SUCCESS
	rs.Msg = ""
	return rs
}
