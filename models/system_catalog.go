package models

import (
	"cointrade/config"
	"cointrade/lib/db"
)

func (m *SystemModel) GetLoanProductList() map[int]float64 {
	rs := make(map[int]float64)
	list, _ := config.GlobalDB.FetchAll(DB_TABLE_LOAN_PRODUCT, db.DB_PARAMS{"state": 1}, db.DB_FIELDS{})
	for _, v := range list {
		rs[v["circle"].ToInt()] = v["rate"].ToFloat()
	}
	return rs
}

func (m *SystemModel) GetRechargeConfig() map[string]*RechargeConfig {
	list, _ := config.GlobalDB.FetchAll(DB_TABLE_RECHARGE_ADDRESS, db.DB_PARAMS{"state": 1}, db.DB_FIELDS{})
	rs := make(map[string]*RechargeConfig)
	for _, v := range list {
		tmp := &RechargeContractConfig{
			Address:     v["address"].ToString(),
			Contract:    v["contract"].ToString(),
			Min:         v["min"].ToFloat(),
			WithDrawMin: v["withdraw_min"].ToFloat(),
		}
		if _, ok := rs[v["cointype"].ToString()]; ok {
			rs[v["cointype"].ToString()].Contracts = append(rs[v["cointype"].ToString()].Contracts, tmp)
			continue
		}
		rs[v["cointype"].ToString()] = &RechargeConfig{
			CoinType:  v["cointype"].ToString(),
			Logo:      v["logo"].ToString(),
			Contracts: []*RechargeContractConfig{tmp},
		}
	}
	return rs
}

func (m *SystemModel) GetOneRechargeConfig(cointype string, contract string) *RechargeContractConfig {
	coinconfig, ok := RECHARGE_ADDRESS_LIST[cointype]
	if !ok {
		return nil
	}
	for _, v := range coinconfig.Contracts {
		if v.Contract == contract {
			return v
		}
	}
	return nil
}

func (m *SystemModel) GetPairMap() map[string]string {
	rs := make(map[string]string)
	list, _ := config.GlobalDB.FetchAll(DB_TABLE_COINS, db.DB_PARAMS{"isnative": 0}, db.DB_FIELDS{})
	for _, v := range list {
		rs[v["vpair"].ToString()] = v["pair"].ToString()
	}
	return rs
}

func (m *SystemModel) GetAllCoins() db.DB_LIST_RESULT {
	list, _ := config.GlobalDB.FetchRows(DB_TABLE_COINS, db.DB_PARAMS{}, db.DB_FIELDS{})
	return list
}

func (m *SystemModel) GetAllCoinsView() db.DB_LIST_RESULT {
	list, _ := config.GlobalDB.FetchRows(DB_TABLE_COINS, db.DB_PARAMS{}, db.DB_FIELDS{})
	viewList := make([]map[string]string, 0, len(list))
	for _, v := range list {
		v["kline_config"] = ""
		viewList = append(viewList, v)
	}
	return viewList
}

func (m *SystemModel) GetNewCoins() db.DB_LIST_RESULT {
	list, _ := config.GlobalDB.FetchRows(DB_TABLE_COINS, db.DB_PARAMS{"is_new": 1, "is_market": 1}, db.DB_FIELDS{})
	return list
}

func (m *SystemModel) GetBuyCoinList() []map[string]string {
	list, _ := config.GlobalDB.FetchRows(DB_TABLE_COINS, db.DB_PARAMS{"is_f": 1}, db.DB_FIELDS{})
	return list
}

func (m *SystemModel) GetCoinInfo(coin string, pair string) db.DB_ROW_RESULT {
	for _, v := range COIN_LIST {
		if v["symbol"] == coin && v["pair"] == pair {
			return v
		}
	}
	return nil
}

func (m *SystemModel) GetExplodeConfig() map[int]*ExplodeConfig {
	rs := make(map[int]*ExplodeConfig)
	list, _ := config.GlobalDB.FetchAll(DB_TABLE_EXPLODE_CONFIG, db.DB_PARAMS{}, db.DB_FIELDS{})
	for _, v := range list {
		tmp := &ExplodeConfig{
			Loserate: v["lose_rate"].ToFloat(),
			Winrate:  v["win_rate"].ToFloat(),
			Time:     v["time"].ToInt(),
			Minprice: v["minprice"].ToFloat(),
		}
		rs[tmp.Time] = tmp
	}
	return rs
}

func (m *SystemModel) LoadCurrency() map[string]float64 {
	list, _ := config.GlobalDB.FetchAll(DB_TABLE_CURRENCY, db.DB_PARAMS{}, db.DB_FIELDS{})
	rs := make(map[string]float64)
	for _, v := range list {
		rs[v["symbol"].ToString()] = v["rate"].ToFloat()
	}
	return rs
}
