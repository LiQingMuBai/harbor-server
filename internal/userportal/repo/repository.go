package repo

import (
	"cointrade/config"
	"cointrade/lib/db"
)

type Repository interface {
	FetchWelcome() (db.DBValues, error)
	FetchSystemConfigList() ([]db.DBValues, error)
	FetchUserByID(uid int) (db.DBValues, error)
	UpdateUserByID(uid int, data db.DB_PARAMS) error
	InsertCrossTrade(data db.DB_PARAMS) error
}

type DBRepository struct{}

func NewDBRepository() Repository {
	return &DBRepository{}
}

func (r *DBRepository) FetchWelcome() (db.DBValues, error) {
	return config.GlobalDB.FetchOne("welcome", nil, db.DB_FIELDS{"id", "platform_name", "welcome_page"})
}

func (r *DBRepository) FetchSystemConfigList() ([]db.DBValues, error) {
	return config.GlobalDB.FetchAll("systemconfig", db.DB_PARAMS{}, db.DB_FIELDS{})
}

func (r *DBRepository) FetchUserByID(uid int) (db.DBValues, error) {
	return config.GlobalDB.FetchOne("users", db.DB_PARAMS{"id": uid}, db.DB_FIELDS{})
}

func (r *DBRepository) UpdateUserByID(uid int, data db.DB_PARAMS) error {
	_, err := config.GlobalDB.UpdateData("users", data, db.DB_PARAMS{"id": uid})
	return err
}

func (r *DBRepository) InsertCrossTrade(data db.DB_PARAMS) error {
	_, err := config.GlobalDB.InsertData("cross_exchange_order", data)
	return err
}
