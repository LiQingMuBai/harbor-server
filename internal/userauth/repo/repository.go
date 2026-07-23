package repo

import (
	"cointrade/config"
	"cointrade/lib/db"
)

const userTable = "users"

type Repository interface {
	FetchInviteUserIDByCode(code string) (db.DBValues, error)
	FetchInvitePoolCode() (db.DBValues, error)
	UpdateInvitePoolByCode(code string, data db.DB_PARAMS) error
	CountUsers(condition db.DB_PARAMS) int
	FetchUser(condition db.DB_PARAMS, fields db.DB_FIELDS) (db.DBValues, error)
	InsertUser(data db.DB_PARAMS) (int64, error)
	UpdateUserByID(id int, data db.DB_PARAMS) error
}

type DBRepository struct{}

func NewDBRepository() Repository {
	return &DBRepository{}
}

func (r *DBRepository) FetchInviteUserIDByCode(code string) (db.DBValues, error) {
	return config.GlobalDB.FetchOne(userTable, db.DB_PARAMS{"invite_code": code}, db.DB_FIELDS{"id"})
}

func (r *DBRepository) FetchInvitePoolCode() (db.DBValues, error) {
	return config.GlobalDB.FetchOne("invitecode_pool", db.DB_PARAMS{"status": 0}, db.DB_FIELDS{"code"}, "limit 0,1")
}

func (r *DBRepository) UpdateInvitePoolByCode(code string, data db.DB_PARAMS) error {
	_, err := config.GlobalDB.UpdateData("invitecode_pool", data, db.DB_PARAMS{"code": code})
	return err
}

func (r *DBRepository) CountUsers(condition db.DB_PARAMS) int {
	return config.GlobalDB.GetCount(userTable, condition)
}

func (r *DBRepository) FetchUser(condition db.DB_PARAMS, fields db.DB_FIELDS) (db.DBValues, error) {
	return config.GlobalDB.FetchOne(userTable, condition, fields, "limit 0,1")
}

func (r *DBRepository) InsertUser(data db.DB_PARAMS) (int64, error) {
	return config.GlobalDB.InsertData(userTable, data)
}

func (r *DBRepository) UpdateUserByID(id int, data db.DB_PARAMS) error {
	_, err := config.GlobalDB.UpdateData(userTable, data, db.DB_PARAMS{"id": id})
	return err
}
