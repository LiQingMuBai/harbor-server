package models

import (
	"cointrade/config"
	"cointrade/lib/db"
	"cointrade/utils"
)

func (m *AssetModel) AddAssets(uid int, a *Assets) bool {
	cacheid := m.MakeCacheId(uid, USER_MODE_REAL)
	if a.Coin == "usdc" {
		a.Price = 1
	}
	if a.Coin == "usdt" {
		MODEL_USER.AddCredit(uid, &CreditValue{
			Credit:          a.Num,
			UserCoinLogType: COIN_LOG_USER_RECHARGE,
			UserCoinLogInfo: QueueCreditLog{
				Credit:     a.Num,
				CoinType:   "usdt",
				CreateTime: utils.GetNow(),
			},
		})
		config.GlobalRedis.Del(HASH_USER_ASSETS, cacheid)
		return true
	}

	coinInfo, _ := config.GlobalDB.FetchRow(DB_TABLE_COINS, db.DB_PARAMS{"symbol": a.Coin, "pair": a.Pair}, db.DB_FIELDS{})
	if coinInfo == nil {
		return false
	}
	one, _ := config.GlobalDB.FetchOne(DB_TABLE_USERASSETS, db.DB_PARAMS{"uid": uid, "coin_symbol": a.Coin, "mode": a.Mode}, db.DB_FIELDS{})
	if one != nil {
		oAllPrice := one["o_price"].ToFloat() * one["num"].ToFloat()
		oPrice := ((a.Price * a.Num) + oAllPrice) / (one["num"].ToFloat() + a.Num)
		err := config.GlobalDB.AddValue(DB_TABLE_USERASSETS, map[string]float64{"num": a.Num, "lock_num": a.LockNum}, db.DB_PARAMS{"id": one["id"].Value})
		config.GlobalDB.UpdateData(DB_TABLE_USERASSETS, db.DB_PARAMS{"o_price": oPrice, "trans_open_time": a.OpenTransTime}, db.DB_PARAMS{"id": one["id"].Value})
		if err != nil {
			return false
		}
	} else {
		insertData := db.DB_PARAMS{
			"uid":             uid,
			"coin_symbol":     a.Coin,
			"coin_id":         coinInfo["id"],
			"coin_pair":       a.Pair,
			"num":             a.Num,
			"lock_num":        0,
			"o_price":         a.Price,
			"mode":            a.Mode,
			"trans_open_time": a.OpenTransTime,
		}
		_, err := config.GlobalDB.InsertData(DB_TABLE_USERASSETS, insertData)
		if err != nil {
			return false
		}
	}
	config.GlobalRedis.Del(HASH_USER_ASSETS, cacheid)
	return true
}
