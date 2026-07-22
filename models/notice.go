package models

import (
	"cointrade/config"
	"cointrade/lib/db"
	"fmt"
	"math"
)

type NoticeModel struct {
	ModelBase
}
type NoticeRequest struct {
	PageBaseRequest
	Pos  string `json:"pos"`  //公告类型
	Lang string `json:"lang"` //公告语言
}

func (m *NoticeModel) GetList(rq *NoticeRequest) *PageBaseResponse { //获取公告
	condition := db.DB_PARAMS{}
	if rq.Pos != "" {
		condition["pos"] = rq.Pos
	}
	if rq.Lang != "" {
		condition["lang"] = rq.Lang
	} else {
		condition["lang"] = "global"
	}
	count := config.GlobalDB.GetCount(DB_TABLE_NOTICE, condition)
	pagesize := int(math.Ceil(float64(count) / float64(rq.Limit)))
	if rq.Page > pagesize {
		rq.Page = pagesize
	}
	if rq.Page == 0 {
		rq.Page = 1
	}
	limitstr := fmt.Sprintf("limit %d,%d", (rq.Page-1)*rq.Limit, rq.Limit)
	list, _ := config.GlobalDB.FetchRows(DB_TABLE_NOTICE, condition, db.DB_FIELDS{}, "order by pubtime desc", limitstr)
	rs := new(PageBaseResponse)
	rs.Limit = rq.Limit
	rs.Page = rq.Page
	rs.List = list
	rs.PageTotal = pagesize
	rs.Total = count
	rs.State = STATE_SUCCESS
	rs.Msg = "success"
	return rs
}

func (m *NoticeModel) GetOne(id int) db.DB_ROW_RESULT {
	cacheId := m.MakeCacheId("notice", id)
	var rs = make(db.DB_ROW_RESULT)
	err := config.GlobalRedis.GetObject(HASH_NOTICE, cacheId, &rs)
	if err == nil && rs["id"] != "" {
		return rs
	}
	one, _ := config.GlobalDB.FetchRow(DB_TABLE_NOTICE, db.DB_PARAMS{"id": id}, db.DB_FIELDS{})
	if one != nil {
		config.GlobalRedis.SetValue(HASH_NOTICE, cacheId, one)
	}
	return one
}
