package models

import (
	"cointrade/config"
	"cointrade/lib/db"
)

func (m *UserModel) Update(uid int, data db.DB_PARAMS) {
	config.GlobalDB.UpdateData(DB_TABLE_USER, data, db.DB_PARAMS{"id": uid})
	m.ClearCache(uid)
}

func (m *UserModel) AddCredit(uid int, credit *CreditValue) bool { //给用户添加各种金额 余额 虚拟余额 冻结金额等等
	err := config.GlobalDB.AddValue(DB_TABLE_USER, map[string]float64{
		"credit":        credit.Credit,
		"v_credit":      credit.VCrdit,
		"lock_credit":   credit.LockCredit,
		"lock_v_credit": credit.LockVCredit,
	}, db.DB_PARAMS{"id": uid})
	MODEL_MESSAGE.PushMessage(uid, MessageCredit{
		Credit:      credit.Credit,
		LockCredit:  credit.LockCredit,
		VCredit:     credit.VCrdit,
		LockVCredit: credit.LockVCredit,
		Text:        nil,
	}, MESSAGE_TYPE_CREDIT)
	m.ClearCache(uid)
	if credit.UserCoinLogInfo != nil {
		MODEL_QUEUE.InputUserQueue(uid, credit.UserCoinLogType, credit.UserCoinLogInfo)
	}
	if credit.TeamCoinLogInfo != nil {
		MODEL_QUEUE.InputTeamQueue(uid, credit.TeamCoinLogType, credit.TeamCoinLogInfo)
	}
	return err == nil
}

func (m *UserModel) ClearCache(uid int) {
	cacheid := m.MakeCacheId(uid)
	config.GlobalRedis.Del(HASH_USER, cacheid)
}
