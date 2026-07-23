package repo

import (
	"cointrade/config"
	"cointrade/lib/db"
)

type Repository interface {
	Count(table string, condition db.DB_PARAMS) int
	Insert(table string, data db.DB_PARAMS) error
	AddValue(table string, values map[string]float64, condition db.DB_PARAMS) error
}

type DBRepository struct{}

func NewDBRepository() Repository {
	return &DBRepository{}
}

func (r *DBRepository) Count(table string, condition db.DB_PARAMS) int {
	return config.GlobalDB.GetCount(table, condition)
}

func (r *DBRepository) Insert(table string, data db.DB_PARAMS) error {
	_, err := config.GlobalDB.InsertData(table, data)
	return err
}

func (r *DBRepository) AddValue(table string, values map[string]float64, condition db.DB_PARAMS) error {
	return config.GlobalDB.AddValue(table, values, condition)
}
