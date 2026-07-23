package repo

import (
	"cointrade/config"
	"cointrade/lib/db"
)

const bankInfoTable = "user_bankinfo"

type BankRepository interface {
	FetchBankByUID(uid int) (db.DBValues, error)
	InsertBank(data db.DB_PARAMS) error
	UpdateBankByID(id interface{}, data db.DB_PARAMS) error
}

type DBBankRepository struct{}

func NewDBBankRepository() BankRepository {
	return &DBBankRepository{}
}

func (r *DBBankRepository) FetchBankByUID(uid int) (db.DBValues, error) {
	return config.GlobalDB.FetchOne(bankInfoTable, db.DB_PARAMS{"uid": uid}, db.DB_FIELDS{})
}

func (r *DBBankRepository) InsertBank(data db.DB_PARAMS) error {
	_, err := config.GlobalDB.InsertData(bankInfoTable, data)
	return err
}

func (r *DBBankRepository) UpdateBankByID(id interface{}, data db.DB_PARAMS) error {
	_, err := config.GlobalDB.UpdateData(bankInfoTable, data, db.DB_PARAMS{"id": id})
	return err
}
