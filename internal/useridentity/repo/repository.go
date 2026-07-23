package repo

import (
	"cointrade/config"
	"cointrade/lib/db"
)

const (
	userAuthLv1Table = "user_auth"
	userAuthLv2Table = "user_auth_2"
)

type Repository interface {
	FetchLv1ByUID(uid int) (db.DBValues, error)
	FetchLv2ByUID(uid int) (db.DBValues, error)
	FetchLv1RowByUID(uid int) (db.DB_ROW_RESULT, error)
	FetchLv2RowByUID(uid int) (db.DB_ROW_RESULT, error)
	FetchLv1ByID(id int) (db.DBValues, error)
	FetchLv2ByID(id int) (db.DBValues, error)
	InsertLv1(data db.DB_PARAMS) (int64, error)
	InsertLv2(data db.DB_PARAMS) (int64, error)
	UpdateLv1ByUID(uid int, data db.DB_PARAMS) error
	UpdateLv2ByID(id interface{}, data db.DB_PARAMS) error
	UpdateLv1ByID(id interface{}, data db.DB_PARAMS) error
	DeleteLv1ByID(id int) error
	DeleteLv2ByID(id int) error
	CountLv1(condition db.DB_PARAMS) int
}

type DBRepository struct{}

func NewDBRepository() Repository {
	return &DBRepository{}
}

func (r *DBRepository) FetchLv1ByUID(uid int) (db.DBValues, error) {
	return config.GlobalDB.FetchOne(userAuthLv1Table, db.DB_PARAMS{"uid": uid}, db.DB_FIELDS{})
}

func (r *DBRepository) FetchLv2ByUID(uid int) (db.DBValues, error) {
	return config.GlobalDB.FetchOne(userAuthLv2Table, db.DB_PARAMS{"uid": uid}, db.DB_FIELDS{})
}

func (r *DBRepository) FetchLv1RowByUID(uid int) (db.DB_ROW_RESULT, error) {
	return config.GlobalDB.FetchRow(userAuthLv1Table, db.DB_PARAMS{"uid": uid}, db.DB_FIELDS{})
}

func (r *DBRepository) FetchLv2RowByUID(uid int) (db.DB_ROW_RESULT, error) {
	return config.GlobalDB.FetchRow(userAuthLv2Table, db.DB_PARAMS{"uid": uid}, db.DB_FIELDS{})
}

func (r *DBRepository) FetchLv1ByID(id int) (db.DBValues, error) {
	return config.GlobalDB.FetchOne(userAuthLv1Table, db.DB_PARAMS{"id": id}, db.DB_FIELDS{})
}

func (r *DBRepository) FetchLv2ByID(id int) (db.DBValues, error) {
	return config.GlobalDB.FetchOne(userAuthLv2Table, db.DB_PARAMS{"id": id}, db.DB_FIELDS{})
}

func (r *DBRepository) InsertLv1(data db.DB_PARAMS) (int64, error) {
	return config.GlobalDB.InsertData(userAuthLv1Table, data)
}

func (r *DBRepository) InsertLv2(data db.DB_PARAMS) (int64, error) {
	return config.GlobalDB.InsertData(userAuthLv2Table, data)
}

func (r *DBRepository) UpdateLv1ByUID(uid int, data db.DB_PARAMS) error {
	_, err := config.GlobalDB.UpdateData(userAuthLv1Table, data, db.DB_PARAMS{"uid": uid})
	return err
}

func (r *DBRepository) UpdateLv1ByID(id interface{}, data db.DB_PARAMS) error {
	_, err := config.GlobalDB.UpdateData(userAuthLv1Table, data, db.DB_PARAMS{"id": id})
	return err
}

func (r *DBRepository) UpdateLv2ByID(id interface{}, data db.DB_PARAMS) error {
	_, err := config.GlobalDB.UpdateData(userAuthLv2Table, data, db.DB_PARAMS{"id": id})
	return err
}

func (r *DBRepository) DeleteLv1ByID(id int) error {
	_, err := config.GlobalDB.Delete(userAuthLv1Table, db.DB_PARAMS{"id": id})
	return err
}

func (r *DBRepository) DeleteLv2ByID(id int) error {
	_, err := config.GlobalDB.Delete(userAuthLv2Table, db.DB_PARAMS{"id": id})
	return err
}

func (r *DBRepository) CountLv1(condition db.DB_PARAMS) int {
	return config.GlobalDB.GetCount(userAuthLv1Table, condition)
}
