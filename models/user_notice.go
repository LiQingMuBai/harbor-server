package models

import (
	"cointrade/config"
	"cointrade/lib/db"
	"fmt"
	"math"
)

func (m *UserModel) GetNoticeUnRead(uid int) int {
	return config.GlobalDB.GetCount(DB_TABLE_USER_NOTICE, db.DB_PARAMS{"is_read": 0, "uid": uid})
}

func (m *UserModel) GetNoticeDetail(uid int, nid int) db.DB_ROW_RESULT {
	one, _ := config.GlobalDB.FetchRow(DB_TABLE_USER_NOTICE, db.DB_PARAMS{"uid": uid, "id": nid}, db.DB_FIELDS{})
	return one
}

func (m *UserModel) GetNoticeList(uid int, rq *PageBaseRequest) *PageBaseResponse {
	count := config.GlobalDB.GetCount(DB_TABLE_USER_NOTICE, db.DB_PARAMS{"uid": uid})
	pagesize := int(math.Ceil(float64(count) / float64(15)))
	rq.Limit = 15
	if rq.Page > pagesize {
		rq.Page = pagesize
	}
	if rq.Page <= 0 {
		rq.Page = 1
	}
	list, _ := config.GlobalDB.FetchRows(DB_TABLE_USER_NOTICE, db.DB_PARAMS{"_": fmt.Sprintf("(uid = %d OR uid = 0)", uid)}, db.DB_FIELDS{}, "order by createtime desc", fmt.Sprintf("limit %d,%d", (rq.Page-1)*rq.Limit, rq.Limit))
	rs := new(PageBaseResponse)
	rs.State = STATE_SUCCESS
	rs.Msg = "ok"
	rs.Limit = rq.Limit
	rs.Page = rq.Page
	rs.List = list
	rs.PageTotal = pagesize
	return rs
}

func (m *UserModel) ReadNotice(uid int, nid int) *BaseResponse {
	config.GlobalDB.UpdateData(DB_TABLE_USER_NOTICE, db.DB_PARAMS{"is_read": 1}, db.DB_PARAMS{"uid": uid, "id": nid})
	return &BaseResponse{State: STATE_SUCCESS, Msg: "ok"}
}

func (m *UserModel) ClearUnreadNotice(uid int) *BaseResponse {
	config.GlobalDB.UpdateData(DB_TABLE_USER_NOTICE, db.DB_PARAMS{"is_read": 1}, db.DB_PARAMS{"uid": uid})
	return &BaseResponse{State: STATE_SUCCESS, Msg: "ok"}
}
