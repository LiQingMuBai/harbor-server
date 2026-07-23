package repo

import (
	"cointrade/config"
	"cointrade/lib/db"
	"fmt"
)

const (
	userWithdrawWalletTable = "user_withdraw_wallets"
	transferTable           = "transfer_detail"
)

type TransferRepository interface {
	CountWallet(cointype string, contract string) int
	InsertWallet(data db.DB_PARAMS) error
	DeleteWallet(uid int, id int) error
	FetchWalletList(uid int) (db.DB_LIST_RESULT, error)
	CountTodayOutTransfer(uid int, today int64) int
	InsertTransfer(data db.DB_PARAMS) error
	CountTransfer(condition db.DB_PARAMS) int
	FetchTransferRows(condition db.DB_PARAMS, limitClause string) (db.DB_LIST_RESULT, error)
	FetchTransferDetail(uid int, sn string) (db.DB_ROW_RESULT, error)
}

type DBTransferRepository struct{}

func NewDBTransferRepository() TransferRepository {
	return &DBTransferRepository{}
}

func (r *DBTransferRepository) CountWallet(cointype string, contract string) int {
	return config.GlobalDB.GetCount(userWithdrawWalletTable, db.DB_PARAMS{"cointype": cointype, "contract": contract})
}

func (r *DBTransferRepository) InsertWallet(data db.DB_PARAMS) error {
	_, err := config.GlobalDB.InsertData(userWithdrawWalletTable, data)
	return err
}

func (r *DBTransferRepository) DeleteWallet(uid int, id int) error {
	_, err := config.GlobalDB.Delete(userWithdrawWalletTable, db.DB_PARAMS{"uid": uid, "id": id})
	return err
}

func (r *DBTransferRepository) FetchWalletList(uid int) (db.DB_LIST_RESULT, error) {
	return config.GlobalDB.FetchRows(userWithdrawWalletTable, db.DB_PARAMS{"uid": uid}, db.DB_FIELDS{})
}

func (r *DBTransferRepository) CountTodayOutTransfer(uid int, today int64) int {
	return config.GlobalDB.GetCount(transferTable, db.DB_PARAMS{
		"direction": 2,
		"_":         fmt.Sprintf("createtime >= %d and state != 2 and uid = %d", today, uid),
	})
}

func (r *DBTransferRepository) InsertTransfer(data db.DB_PARAMS) error {
	_, err := config.GlobalDB.InsertData(transferTable, data)
	return err
}

func (r *DBTransferRepository) CountTransfer(condition db.DB_PARAMS) int {
	return config.GlobalDB.GetCount(transferTable, condition)
}

func (r *DBTransferRepository) FetchTransferRows(condition db.DB_PARAMS, limitClause string) (db.DB_LIST_RESULT, error) {
	return config.GlobalDB.FetchRows(transferTable, condition, db.DB_FIELDS{}, "order by id desc", limitClause)
}

func (r *DBTransferRepository) FetchTransferDetail(uid int, sn string) (db.DB_ROW_RESULT, error) {
	return config.GlobalDB.FetchRow(transferTable, db.DB_PARAMS{"uid": uid, "sn": sn}, db.DB_FIELDS{})
}
