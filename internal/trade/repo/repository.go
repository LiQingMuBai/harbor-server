package repo

import (
	"cointrade/config"
	"cointrade/lib/db"
	"fmt"
)

const (
	tableOpenedTrade   = "open_trade"
	tableDelegateTrade = "delegate_trade"
	tableCloseTrade    = "close_trade"
)

type Repository interface {
	FetchCloseBySN(uid int, sn string) (db.DB_ROW_RESULT, error)
	FetchOpenedBySN(uid int, sn string) (db.DBValues, error)
	FetchOpened(uid int, coin string, tradeType int, flag int, mode int, ganggan int) (db.DBValues, error)
	FetchPendingDelegate(uid int, sn string) (db.DBValues, error)
	FetchPendingDelegates(limit int) ([]db.DBValues, error)
	CountDelegate(condition db.DB_PARAMS) int
	FetchDelegateRows(condition db.DB_PARAMS, order string, limit string) (db.DB_LIST_RESULT, error)
	FetchOpenedByCondition(condition db.DB_PARAMS, limit string) ([]db.DBValues, error)
	CountOpened(condition db.DB_PARAMS) int
	FetchOpenedRows(condition db.DB_PARAMS, order string, limit string) (db.DB_LIST_RESULT, error)
	CountClose(condition db.DB_PARAMS) int
	FetchCloseRows(condition db.DB_PARAMS, order string, limit string) (db.DB_LIST_RESULT, error)
	InsertOpened(data db.DB_PARAMS) error
	InsertClose(data db.DB_PARAMS) error
	InsertDelegate(data db.DB_PARAMS) error
	UpdateDelegateState(id interface{}, state int, changeTime int) error
	AddOpenedPositionValue(id int, data map[string]float64) error
	UpdateOpenedPrice(id int, openPrice float64) error
	AddOpenedLockValue(id int, num float64, lockNum float64, mode int) error
	UpdateOpenedFields(id int, data db.DB_PARAMS) error
}

type DBRepository struct{}

func NewDBRepository() Repository {
	return DBRepository{}
}

func (DBRepository) FetchCloseBySN(uid int, sn string) (db.DB_ROW_RESULT, error) {
	return config.GlobalDB.FetchRow(tableCloseTrade, db.DB_PARAMS{"sn": sn, "uid": uid}, db.DB_FIELDS{})
}

func (DBRepository) FetchOpenedBySN(uid int, sn string) (db.DBValues, error) {
	return config.GlobalDB.FetchOne(tableOpenedTrade, db.DB_PARAMS{"sn": sn, "uid": uid}, db.DB_FIELDS{})
}

func (DBRepository) FetchOpened(uid int, coin string, tradeType int, flag int, mode int, ganggan int) (db.DBValues, error) {
	return config.GlobalDB.FetchOne(
		tableOpenedTrade,
		db.DB_PARAMS{
			"trade_type":  tradeType,
			"uid":         uid,
			"flag":        flag,
			"coin_symbol": coin,
			"mode":        mode,
			"ganggan":     ganggan,
		},
		db.DB_FIELDS{},
	)
}

func (DBRepository) FetchPendingDelegate(uid int, sn string) (db.DBValues, error) {
	return config.GlobalDB.FetchOne(tableDelegateTrade, db.DB_PARAMS{"uid": uid, "sn": sn, "state": 0}, db.DB_FIELDS{})
}

func (DBRepository) FetchPendingDelegates(limit int) ([]db.DBValues, error) {
	return config.GlobalDB.FetchAll(tableDelegateTrade, db.DB_PARAMS{"state": 0, "is_f": 0}, db.DB_FIELDS{}, fmt.Sprintf("limit 0,%d", limit))
}

func (DBRepository) CountDelegate(condition db.DB_PARAMS) int {
	return config.GlobalDB.GetCount(tableDelegateTrade, condition)
}

func (DBRepository) FetchDelegateRows(condition db.DB_PARAMS, order string, limit string) (db.DB_LIST_RESULT, error) {
	return config.GlobalDB.FetchRows(tableDelegateTrade, condition, db.DB_FIELDS{}, order, limit)
}

func (DBRepository) CountOpened(condition db.DB_PARAMS) int {
	return config.GlobalDB.GetCount(tableOpenedTrade, condition)
}

func (DBRepository) FetchOpenedByCondition(condition db.DB_PARAMS, limit string) ([]db.DBValues, error) {
	return config.GlobalDB.FetchAll(tableOpenedTrade, condition, db.DB_FIELDS{}, limit)
}

func (DBRepository) FetchOpenedRows(condition db.DB_PARAMS, order string, limit string) (db.DB_LIST_RESULT, error) {
	return config.GlobalDB.FetchRows(tableOpenedTrade, condition, db.DB_FIELDS{}, order, limit)
}

func (DBRepository) CountClose(condition db.DB_PARAMS) int {
	return config.GlobalDB.GetCount(tableCloseTrade, condition)
}

func (DBRepository) FetchCloseRows(condition db.DB_PARAMS, order string, limit string) (db.DB_LIST_RESULT, error) {
	return config.GlobalDB.FetchRows(tableCloseTrade, condition, db.DB_FIELDS{}, order, limit)
}

func (DBRepository) InsertOpened(data db.DB_PARAMS) error {
	_, err := config.GlobalDB.InsertData(tableOpenedTrade, data)
	return err
}

func (DBRepository) InsertClose(data db.DB_PARAMS) error {
	_, err := config.GlobalDB.InsertData(tableCloseTrade, data)
	return err
}

func (DBRepository) InsertDelegate(data db.DB_PARAMS) error {
	_, err := config.GlobalDB.InsertData(tableDelegateTrade, data)
	return err
}

func (DBRepository) UpdateDelegateState(id interface{}, state int, changeTime int) error {
	updateData := db.DB_PARAMS{"state": state}
	if changeTime > 0 {
		updateData["changetime"] = changeTime
	}
	_, err := config.GlobalDB.UpdateData(tableDelegateTrade, updateData, db.DB_PARAMS{"id": id})
	return err
}

func (DBRepository) AddOpenedPositionValue(id int, data map[string]float64) error {
	return config.GlobalDB.AddValue(tableOpenedTrade, data, db.DB_PARAMS{"id": id})
}

func (DBRepository) UpdateOpenedPrice(id int, openPrice float64) error {
	_, err := config.GlobalDB.UpdateData(tableOpenedTrade, db.DB_PARAMS{"openprice": openPrice}, db.DB_PARAMS{"id": id})
	return err
}

func (DBRepository) AddOpenedLockValue(id int, num float64, lockNum float64, mode int) error {
	return config.GlobalDB.AddValue(tableOpenedTrade, map[string]float64{"num": num, "lock_num": lockNum}, db.DB_PARAMS{"id": id, "mode": mode})
}

func (DBRepository) UpdateOpenedFields(id int, data db.DB_PARAMS) error {
	_, err := config.GlobalDB.UpdateData(tableOpenedTrade, data, db.DB_PARAMS{"id": id})
	return err
}
