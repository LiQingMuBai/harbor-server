package repo

import (
	"cointrade/config"
	"cointrade/lib/db"
)

const withdrawTable = "withdraw"

type WithdrawRepository interface {
	InsertWithdraw(data db.DB_PARAMS) error
	CountWithdraw(condition db.DB_PARAMS) int
	FetchWithdrawRows(condition db.DB_PARAMS, limitClause string) (db.DB_LIST_RESULT, error)
	FetchWithdrawInfo(uid int, sn string) (db.DB_ROW_RESULT, error)
}

type DBWithdrawRepository struct{}

func NewDBWithdrawRepository() WithdrawRepository {
	return &DBWithdrawRepository{}
}

func (r *DBWithdrawRepository) InsertWithdraw(data db.DB_PARAMS) error {
	_, err := config.GlobalDB.InsertData(withdrawTable, data)
	return err
}

func (r *DBWithdrawRepository) CountWithdraw(condition db.DB_PARAMS) int {
	return config.GlobalDB.GetCount(withdrawTable, condition)
}

func (r *DBWithdrawRepository) FetchWithdrawRows(condition db.DB_PARAMS, limitClause string) (db.DB_LIST_RESULT, error) {
	return config.GlobalDB.FetchRows(withdrawTable, condition, db.DB_FIELDS{}, "order by id desc", limitClause)
}

func (r *DBWithdrawRepository) FetchWithdrawInfo(uid int, sn string) (db.DB_ROW_RESULT, error) {
	return config.GlobalDB.FetchRow(withdrawTable, db.DB_PARAMS{"uid": uid, "sn": sn}, db.DB_FIELDS{})
}
