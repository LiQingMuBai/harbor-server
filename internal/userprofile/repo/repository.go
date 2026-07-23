package repo

import (
	"cointrade/config"
	"cointrade/lib/db"
)

type Repository interface {
	FetchUserByID(uid int) (db.DBValues, error)
	FetchUser(condition db.DB_PARAMS, fields db.DB_FIELDS) (db.DBValues, error)
	FetchCountSummary(table string, condition db.DB_PARAMS) (db.DBValues, error)
}

type DBRepository struct{}

func NewDBRepository() Repository {
	return &DBRepository{}
}

func (r *DBRepository) FetchUserByID(uid int) (db.DBValues, error) {
	return config.GlobalDB.FetchOne("users", db.DB_PARAMS{"id": uid}, db.DB_FIELDS{})
}

func (r *DBRepository) FetchUser(condition db.DB_PARAMS, fields db.DB_FIELDS) (db.DBValues, error) {
	return config.GlobalDB.FetchOne("users", condition, fields)
}

func (r *DBRepository) FetchCountSummary(table string, condition db.DB_PARAMS) (db.DBValues, error) {
	return config.GlobalDB.FetchOne(table, condition, db.DB_FIELDS{})
}
