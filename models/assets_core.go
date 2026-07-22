package models

import (
	"cointrade/config"
	"cointrade/lib/db"
	"strings"
)

func (m *AssetModel) GetOneAsset(uid int, coin string) *AssetInfo {
	userAssets := m.GetAllAssets(uid, USER_MODE_REAL)
	if v, ok := userAssets[coin]; ok {
		return &v
	}
	return nil
}

func (m *AssetModel) InitUserAssets(uid int) {
	for _, v := range COIN_LIST {
		insertData := db.DB_PARAMS{"uid": uid}
		insertData["coin_symbol"] = v["symbol"]
		insertData["coin_id"] = v["id"]
		insertData["coin_pair"] = v["pair"]
		insertData["wallet_address"] = v["address"]
		insertData["mode"] = USER_MODE_REAL
		config.GlobalDB.InsertData(DB_TABLE_USERASSETS, insertData)
	}
}

func (m *AssetModel) GetAllAssets(uid int, mode int) map[string]AssetInfo {
	rs := make(map[string]AssetInfo)
	uinfo := MODEL_USER.GetBaseInfo(uid)
	rs["usdt"] = AssetInfo{
		Symbol:    "usdt",
		Count:     uinfo.Credit,
		LockCount: uinfo.LockCredit,
	}

	list, _ := config.GlobalDB.FetchAll(DB_TABLE_USERASSETS, db.DB_PARAMS{"uid": uid, "mode": mode}, db.DB_FIELDS{})
	for _, v := range list {
		tmp := AssetInfo{
			CoinId:        v["coin_id"].ToInt(),
			O_Price:       v["o_price"].ToFloat(),
			Pair:          v["coin_pair"].ToString(),
			Symbol:        v["coin_symbol"].ToString(),
			Address:       v["wallet_address"].ToString(),
			IsTrans:       v["is_trans"].ToInt(),
			TransOpenTime: v["trans_open_time"].ToInt(),
		}
		if strings.ToLower(tmp.Symbol) == "usdt" {
			tmp.Count = uinfo.Credit
			tmp.LockCount = uinfo.LockCredit
		} else {
			tmp.Count = v["num"].ToFloat()
			tmp.LockCount = v["lock_num"].ToFloat()
		}
		rs[tmp.Symbol] = tmp
	}

	rs["usdt"] = AssetInfo{
		Symbol:    "usdt",
		Count:     uinfo.Credit,
		LockCount: uinfo.LockCredit,
	}
	return rs
}
