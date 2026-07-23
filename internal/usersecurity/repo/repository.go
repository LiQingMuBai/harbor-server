package repo

import (
	"cointrade/config"
	"cointrade/lib/db"
)

const userTable = "users"

type Repository interface {
	FetchUser(condition db.DB_PARAMS, fields db.DB_FIELDS) (db.DBValues, error)
}

type DBRepository struct{}

func NewDBRepository() Repository {
	return &DBRepository{}
}

func (r *DBRepository) FetchUser(condition db.DB_PARAMS, fields db.DB_FIELDS) (db.DBValues, error) {
	return config.GlobalDB.FetchOne(userTable, condition, fields)
}
