package models

import (
	"cointrade/config"
	"cointrade/lib/db"
	"fmt"
	"math"
)

func (m *CreditLogModel) AddAgentLog(uid int, createtime int, data map[string]float64) {
	daytime := startOfDayUnix(createtime)
	one, _ := config.GlobalDB.FetchOne(DB_TABLE_AGENT_COUNT, db.DB_PARAMS{"uid": uid, "daytime": daytime}, db.DB_FIELDS{"id"})
	if one != nil {
		config.GlobalDB.AddValue(DB_TABLE_AGENT_COUNT, data, db.DB_PARAMS{"id": one["id"].Value})
		return
	}
	insertData := db.DB_PARAMS{"uid": uid, "daytime": daytime}
	for k, v := range data {
		insertData[k] = v
	}
	config.GlobalDB.InsertData(DB_TABLE_AGENT_COUNT, insertData)
}

func (m *CreditLogModel) AddAgentLevelLog(uid int, level int, createtime int, data map[string]float64) {
	daytime := startOfDayUnix(createtime)
	one, _ := config.GlobalDB.FetchOne(DB_TABLE_AGENT_LEVEL_COUNT, db.DB_PARAMS{"uid": uid, "daytime": daytime, "level": level}, db.DB_FIELDS{"id"})
	if one != nil {
		config.GlobalDB.AddValue(DB_TABLE_AGENT_LEVEL_COUNT, data, db.DB_PARAMS{"id": one["id"].Value})
		return
	}
	uinfo := MODEL_USER.GetBaseInfo(uid)
	insertData := db.DB_PARAMS{"uid": uid, "email": uinfo.Email, "daytime": daytime, "level": level}
	for k, v := range data {
		insertData[k] = v
	}
	config.GlobalDB.InsertData(DB_TABLE_AGENT_LEVEL_COUNT, insertData)
}

func (m *CreditLogModel) IncomeLog(uid int, rq *PageBaseRequest) *PageBaseResponse {
	rs := new(PageBaseResponse)
	count := config.GlobalDB.GetCount(DB_TABLE_INCOME_LOG, db.DB_PARAMS{"uid": uid})
	pagesize := int(math.Ceil(float64(count) / float64(rq.Limit)))
	if rq.Page > pagesize {
		rq.Page = pagesize
	}
	if rq.Page == 0 {
		rq.Page = 1
	}
	limitstr := fmt.Sprintf("limit %d,%d", (rq.Page-1)*rq.Limit, rq.Limit)
	list, _ := config.GlobalDB.FetchRows(DB_TABLE_INCOME_LOG, db.DB_PARAMS{"uid": uid}, db.DB_FIELDS{}, "order by id desc", limitstr)
	rs.State = STATE_SUCCESS
	rs.Limit = rq.Limit
	rs.Page = rq.Page
	rs.Msg = "success"
	rs.PageTotal = pagesize
	rs.Total = count
	rs.List = list
	return rs
}
