package repo

import (
	"cointrade/config"
	"cointrade/lib/db"
	"fmt"
)

const (
	miningAcceptTable  = "mining_accept"
	miningOrderTable   = "mining_order"
	miningProductTable = "mining_product"
)

type Repository interface {
	FetchAcceptStates(uid int) ([]db.DBValues, error)
	FetchAcceptRecord(uid int, pid int) (db.DBValues, error)
	FetchOpenProducts() ([]db.DBValues, error)
	CountUserProductOrders(uid int, pid int) int
	HasUserProductOrder(uid int, pid int) bool
	InsertOrder(data db.DB_PARAMS) error
	FetchReservedOrderProductIDs(uid int) ([]db.DBValues, error)
	CountOrders(condition db.DB_PARAMS) int
	FetchOrders(condition db.DB_PARAMS, offset int, limit int) ([]db.DBValues, error)
	UpdateOrderByID(id int, data db.DB_PARAMS) error
	FetchOrder(uid int, sn string) (db.DBValues, error)
	FetchOpenOrderSummary(uid int) (db.DBValues, error)
	FetchOrderHistorySummary(uid int) (db.DBValues, error)
}

type DBRepository struct{}

func NewDBRepository() Repository {
	return &DBRepository{}
}

func (r *DBRepository) FetchAcceptStates(uid int) ([]db.DBValues, error) {
	return config.GlobalDB.FetchAll(miningAcceptTable, db.DB_PARAMS{"uid": uid}, db.DB_FIELDS{})
}

func (r *DBRepository) FetchAcceptRecord(uid int, pid int) (db.DBValues, error) {
	return config.GlobalDB.FetchOne(miningAcceptTable, db.DB_PARAMS{"uid": uid, "product_id": pid}, db.DB_FIELDS{})
}

func (r *DBRepository) FetchOpenProducts() ([]db.DBValues, error) {
	return config.GlobalDB.FetchAll(miningProductTable, db.DB_PARAMS{"isopen": 1}, db.DB_FIELDS{})
}

func (r *DBRepository) CountUserProductOrders(uid int, pid int) int {
	return config.GlobalDB.GetCount(miningOrderTable, db.DB_PARAMS{"uid": uid, "pid": pid})
}

func (r *DBRepository) HasUserProductOrder(uid int, pid int) bool {
	one, _ := config.GlobalDB.FetchOne(miningOrderTable, db.DB_PARAMS{"pid": pid, "uid": uid}, db.DB_FIELDS{"id"})
	return one != nil
}

func (r *DBRepository) InsertOrder(data db.DB_PARAMS) error {
	_, err := config.GlobalDB.InsertData(miningOrderTable, data)
	return err
}

func (r *DBRepository) FetchReservedOrderProductIDs(uid int) ([]db.DBValues, error) {
	return config.GlobalDB.FetchAll(miningOrderTable, db.DB_PARAMS{"uid": uid, "_": "(state=2 or state=4)"}, db.DB_FIELDS{"pid"})
}

func (r *DBRepository) CountOrders(condition db.DB_PARAMS) int {
	return config.GlobalDB.GetCount(miningOrderTable, condition)
}

func (r *DBRepository) FetchOrders(condition db.DB_PARAMS, offset int, limit int) ([]db.DBValues, error) {
	limitClause := fmt.Sprintf("limit %d,%d", offset, limit)
	return config.GlobalDB.FetchAll(miningOrderTable, condition, db.DB_FIELDS{}, "order by createtime desc", limitClause)
}

func (r *DBRepository) UpdateOrderByID(id int, data db.DB_PARAMS) error {
	_, err := config.GlobalDB.UpdateData(miningOrderTable, data, db.DB_PARAMS{"id": id})
	return err
}

func (r *DBRepository) FetchOrder(uid int, sn string) (db.DBValues, error) {
	return config.GlobalDB.FetchOne(miningOrderTable, db.DB_PARAMS{"uid": uid, "sn": sn}, db.DB_FIELDS{})
}

func (r *DBRepository) FetchOpenOrderSummary(uid int) (db.DBValues, error) {
	return config.GlobalDB.FetchOne(miningOrderTable, db.DB_PARAMS{"uid": uid, "state": 0}, db.DB_FIELDS{"SUM(profit) as day_profit,sum(amount) as all_amount"})
}

func (r *DBRepository) FetchOrderHistorySummary(uid int) (db.DBValues, error) {
	return config.GlobalDB.FetchOne(miningOrderTable, db.DB_PARAMS{"uid": uid}, db.DB_FIELDS{"SUM(allprofit) as history_profit"})
}
