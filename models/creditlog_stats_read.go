package models

import (
	"cointrade/config"
	"cointrade/lib/db"
	"cointrade/utils"
	"fmt"
	"time"
)

func monthRange() (int64, int64) {
	t := time.Now()
	year := t.Year()
	month := t.Month()
	starttime := time.Date(year, month, 1, 0, 0, 0, 0, time.Local).Unix()
	nextmonth := month + 1
	if nextmonth > 12 {
		nextmonth = 1
		year++
	}
	endtime := time.Date(year, nextmonth, 1, 0, 0, 0, 0, time.Local).Unix()
	return starttime, endtime
}

func (m *CreditLogModel) GetUserCountDay(uid int) db.DB_ROW_RESULT {
	daytime := startOfDayUnix(utils.GetNow())
	one, _ := config.GlobalDB.FetchRow(DB_TABLE_USER_COUNT, db.DB_PARAMS{"daytime": daytime}, db.DB_FIELDS{})
	return one
}

func (m *CreditLogModel) GetUserCountSum(uid int) db.DB_ROW_RESULT {
	one, _ := config.GlobalDB.FetchRow(DB_TABLE_USER_COUNT_SUM, db.DB_PARAMS{"uid": uid}, db.DB_FIELDS{})
	return one
}

func (m *CreditLogModel) GetUserCountMonth(uid int) db.DB_LIST_RESULT {
	starttime, endtime := monthRange()
	list, _ := config.GlobalDB.FetchRows(DB_TABLE_USER_COUNT, db.DB_PARAMS{"_": fmt.Sprintf("daytime>=%d and daytime<%d", starttime, endtime)}, db.DB_FIELDS{})
	return list
}

func (m *CreditLogModel) GetUserLevelCountDay(uid int) map[int]db.DB_ROW_RESULT {
	daytime := startOfDayUnix(utils.GetNow())
	condition := db.DB_PARAMS{"uid": uid, "daytime": daytime}
	rs := make(map[int]db.DB_ROW_RESULT)
	list, _ := config.GlobalDB.FetchRows(DB_TABLE_USER_LEVEL_COUNT, condition, db.DB_FIELDS{})
	for _, v := range list {
		rs[utils.GetInt(v["level"])] = v
	}
	return rs
}

func (m *CreditLogModel) GetUserLevelCountSum(uid int) map[int]db.DB_ROW_RESULT {
	condition := db.DB_PARAMS{"uid": uid}
	rs := make(map[int]db.DB_ROW_RESULT)
	list, _ := config.GlobalDB.FetchRows(DB_TABLE_USER_LEVEL_COUNT_SUM, condition, db.DB_FIELDS{})
	for _, v := range list {
		rs[utils.GetInt(v["level"])] = v
	}
	return rs
}

func (m *CreditLogModel) GetUserLevelCountMonth(uid int) map[int]db.DB_LIST_RESULT {
	starttime, endtime := monthRange()
	list, _ := config.GlobalDB.FetchRows(DB_TABLE_USER_LEVEL_COUNT, db.DB_PARAMS{"_": fmt.Sprintf("daytime>=%d and daytime<%d", starttime, endtime)}, db.DB_FIELDS{})
	rs := make(map[int]db.DB_LIST_RESULT)
	for _, v := range list {
		level := utils.GetInt(v["level"])
		if _, ok := rs[level]; !ok {
			rs[level] = make(db.DB_LIST_RESULT, 0)
		}
		rs[level] = append(rs[level], v)
	}
	return rs
}
