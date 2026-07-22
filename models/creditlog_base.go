package models

import (
	"cointrade/config"
	"cointrade/lib/db"
	"fmt"
	"math"
)

func (m *CreditLogModel) Add(info *CreditLogInfo) {
	uinfo := MODEL_USER.GetBaseInfo(info.Uid)
	if uinfo == nil {
		return
	}

	insertData := db.DB_PARAMS{
		"uid":          info.Uid,
		"credit":       info.Credit,
		"lock_credit":  info.LockCredit,
		"mode":         info.Mode,
		"sn":           info.Sn,
		"type":         info.Type,
		"createtime":   info.Createtime,
		"after_credit": uinfo.Credit,
		"cointype":     info.CoinType,
	}
	config.GlobalDB.InsertData(DB_TABLE_CREDIT_LOG, insertData)
}

func (m *CreditLogModel) GetList(uid int, rq *CoinLogRequest) *PageBaseResponse {
	condition := db.DB_PARAMS{"uid": uid, "display": 1}
	if rq.Type != 0 {
		condition["type"] = rq.Type
	}
	count := config.GlobalDB.GetCount(DB_TABLE_CREDIT_LOG, condition)
	pagesize := int(math.Ceil(float64(count) / float64(rq.Limit)))
	if rq.Page > pagesize {
		rq.Page = pagesize
	}
	if rq.Page == 0 {
		rq.Page = 1
	}
	limitstr := fmt.Sprintf("limit %d,%d", (rq.Page-1)*rq.Limit, rq.Limit)
	list, _ := config.GlobalDB.FetchRows(DB_TABLE_CREDIT_LOG, condition, db.DB_FIELDS{}, "order by id desc", limitstr)
	rs := new(PageBaseResponse)
	rs.State = STATE_SUCCESS
	rs.Msg = "success"
	rs.Limit = rq.Limit
	rs.Page = rq.Page
	rs.PageTotal = pagesize
	rs.Total = count
	rs.List = list
	return rs
}
