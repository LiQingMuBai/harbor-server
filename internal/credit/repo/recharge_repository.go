package repo

import (
	"cointrade/config"
	"cointrade/lib/db"
)

const (
	rechargeAddressTable = "recharge_address"
	rechargeTable        = "recharge"
)

type RechargeRepository interface {
	FetchRechargeAddresses() (db.DB_LIST_RESULT, error)
	InsertRecharge(data db.DB_PARAMS) error
	FetchRechargeBySN(sn string) (db.DBValues, error)
	FetchRechargeRowBySN(sn string) (db.DB_ROW_RESULT, error)
	UpdateRechargeByID(id interface{}, data db.DB_PARAMS) error
	CountRecharge(condition db.DB_PARAMS) int
	FetchRechargeRows(condition db.DB_PARAMS, limitClause string) (db.DB_LIST_RESULT, error)
	FetchRechargeInfo(uid int, sn string) (db.DB_ROW_RESULT, error)
}

type DBRechargeRepository struct{}

func NewDBRechargeRepository() RechargeRepository {
	return &DBRechargeRepository{}
}

func (r *DBRechargeRepository) FetchRechargeAddresses() (db.DB_LIST_RESULT, error) {
	return config.GlobalDB.FetchRows(rechargeAddressTable, db.DB_PARAMS{"state": 1}, db.DB_FIELDS{})
}

func (r *DBRechargeRepository) InsertRecharge(data db.DB_PARAMS) error {
	_, err := config.GlobalDB.InsertData(rechargeTable, data)
	return err
}

func (r *DBRechargeRepository) FetchRechargeBySN(sn string) (db.DBValues, error) {
	return config.GlobalDB.FetchOne(rechargeTable, db.DB_PARAMS{"sn": sn}, db.DB_FIELDS{})
}

func (r *DBRechargeRepository) FetchRechargeRowBySN(sn string) (db.DB_ROW_RESULT, error) {
	return config.GlobalDB.FetchRow(rechargeTable, db.DB_PARAMS{"sn": sn}, db.DB_FIELDS{})
}

func (r *DBRechargeRepository) UpdateRechargeByID(id interface{}, data db.DB_PARAMS) error {
	_, err := config.GlobalDB.UpdateData(rechargeTable, data, db.DB_PARAMS{"id": id})
	return err
}

func (r *DBRechargeRepository) CountRecharge(condition db.DB_PARAMS) int {
	return config.GlobalDB.GetCount(rechargeTable, condition)
}

func (r *DBRechargeRepository) FetchRechargeRows(condition db.DB_PARAMS, limitClause string) (db.DB_LIST_RESULT, error) {
	return config.GlobalDB.FetchRows(rechargeTable, condition, db.DB_FIELDS{}, "order by id desc", limitClause)
}

func (r *DBRechargeRepository) FetchRechargeInfo(uid int, sn string) (db.DB_ROW_RESULT, error) {
	return config.GlobalDB.FetchRow(rechargeTable, db.DB_PARAMS{"uid": uid, "sn": sn}, db.DB_FIELDS{})
}
