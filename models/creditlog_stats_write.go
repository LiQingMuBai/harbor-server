package models

import (
	"cointrade/config"
	"cointrade/lib/db"
	"cointrade/utils"
	"strings"
	"time"
)

func startOfDayUnix(ts int) int64 {
	t := time.Unix(int64(ts), 0)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).Unix()
}

func (m *CreditLogModel) AddUserCount(uid int, createtime int, data map[string]float64) {
	daytime := startOfDayUnix(createtime)
	one, _ := config.GlobalDB.FetchOne(DB_TABLE_USER_COUNT, db.DB_PARAMS{"daytime": daytime, "uid": uid}, db.DB_FIELDS{"id"}, "limit 0,1")
	if one != nil {
		config.GlobalDB.AddValue(DB_TABLE_USER_COUNT, data, db.DB_PARAMS{"id": one["id"].Value})
	} else {
		insertData := db.DB_PARAMS{"uid": uid, "daytime": daytime}
		for k, v := range data {
			insertData[k] = v
		}
		config.GlobalDB.InsertData(DB_TABLE_USER_COUNT, insertData)
	}

	one, _ = config.GlobalDB.FetchOne(DB_TABLE_USER_COUNT_SUM, db.DB_PARAMS{"uid": uid}, db.DB_FIELDS{"id"})
	if one != nil {
		config.GlobalDB.AddValue(DB_TABLE_USER_COUNT_SUM, data, db.DB_PARAMS{"uid": uid})
		return
	}
	insertData := db.DB_PARAMS{"uid": uid}
	for k, v := range data {
		insertData[k] = v
	}
	config.GlobalDB.InsertData(DB_TABLE_USER_COUNT_SUM, insertData)
}

func (m *CreditLogModel) AddLevelCount(uid int, level int, createtime int, data map[string]float64) {
	daytime := startOfDayUnix(createtime)
	condition := db.DB_PARAMS{"daytime": daytime, "uid": uid, "level": level}
	one, _ := config.GlobalDB.FetchOne(DB_TABLE_USER_LEVEL_COUNT, condition, db.DB_FIELDS{"id"}, "limit 0,1")
	uinfo := MODEL_USER.GetBaseInfo(uid)
	if uinfo == nil {
		return
	}
	if one != nil {
		config.GlobalDB.AddValue(DB_TABLE_USER_LEVEL_COUNT, data, db.DB_PARAMS{"id": one["id"].Value})
	} else {
		for k, v := range data {
			condition[k] = v
		}
		config.GlobalDB.InsertData(DB_TABLE_USER_LEVEL_COUNT, condition)
	}

	one, _ = config.GlobalDB.FetchOne(DB_TABLE_USER_LEVEL_COUNT_SUM, db.DB_PARAMS{"uid": uid, "level": level}, db.DB_FIELDS{"id"}, "limit 0,1")
	if one != nil {
		config.GlobalDB.AddValue(DB_TABLE_USER_LEVEL_COUNT_SUM, data, db.DB_PARAMS{"id": one["id"].Value})
		return
	}
	insertData := db.DB_PARAMS{"uid": uid, "level": level, "email": uinfo.Email}
	for k, v := range data {
		insertData[k] = v
	}
	config.GlobalDB.InsertData(DB_TABLE_USER_LEVEL_COUNT_SUM, insertData)
}

func (m *CreditLogModel) AddSiteCount(createtime int, data map[string]float64) {
	daytime := startOfDayUnix(createtime)
	condition := db.DB_PARAMS{"daytime": daytime}
	one, _ := config.GlobalDB.FetchOne(DB_TABLE_SITE_COUNT, condition, db.DB_FIELDS{"id"})
	if one != nil {
		config.GlobalDB.AddValue(DB_TABLE_SITE_COUNT, data, condition)
		return
	}
	for k, v := range data {
		condition[k] = v
	}
	config.GlobalDB.InsertData(DB_TABLE_SITE_COUNT, condition)
}

func (m *CreditLogModel) AddUserCountLog(uid int, logInfo *QueueTeamLog) {
	uinfo := MODEL_USER.GetBaseInfo(uid)
	if uinfo == nil {
		return
	}

	addvalues := map[string]float64{}
	siteAddvalues := map[string]float64{}
	if logInfo.MiningCount > 0 {
		addvalues["mining_count"] = logInfo.MiningCount
		siteAddvalues["minning_count"] = logInfo.MiningCount
	}
	if logInfo.MiningProfit > 0 {
		addvalues["mining_profit"] = logInfo.MiningProfit
		siteAddvalues["minning_profit"] = logInfo.MiningProfit
	}
	if logInfo.Recharge > 0 {
		addvalues["recharge"] = logInfo.Recharge
		siteAddvalues["recharge"] = logInfo.Recharge
		siteAddvalues["register_num"] = 1
		count := config.GlobalDB.GetCount(DB_TABLE_RECHARGE, db.DB_PARAMS{"uid": uid, "state": 1})
		if count == 1 {
			siteAddvalues["pro_num"] = 1
			siteAddvalues["first_recharge"] = logInfo.Recharge
			siteAddvalues["first_recharge_num"] = 1
		}
	}
	if logInfo.TradeBB > 0 {
		addvalues["trade"] = logInfo.TradeBB
		addvalues["trade_bb"] = logInfo.TradeBB
		siteAddvalues["trade"] = logInfo.TradeBB
		siteAddvalues["trade_bb"] = logInfo.TradeBB
	}
	if logInfo.TradeExplode > 0 {
		addvalues["trade"] = logInfo.TradeExplode
		addvalues["trade_explode"] = logInfo.TradeExplode
		siteAddvalues["trade"] = logInfo.TradeExplode
		siteAddvalues["trade_explode"] = logInfo.TradeExplode
	}
	if logInfo.TradeKeep > 0 {
		addvalues["trade"] = logInfo.TradeKeep
		addvalues["trade_keep"] = logInfo.TradeKeep
		siteAddvalues["trade"] = logInfo.TradeKeep
		siteAddvalues["trade_keep"] = logInfo.TradeKeep
	}
	if logInfo.TradeBB_Profit > 0 {
		addvalues["trade_profit"] = logInfo.TradeBB_Profit
		addvalues["trade_bb_profit"] = logInfo.TradeBB_Profit
		siteAddvalues["trade_profit"] = logInfo.TradeBB_Profit
		siteAddvalues["trade_bb_profit"] = logInfo.TradeBB_Profit
	}
	if logInfo.TradeExplode_Profit > 0 {
		addvalues["trade_profit"] = logInfo.TradeExplode_Profit
		addvalues["trade_explode_profit"] = logInfo.TradeExplode_Profit
		siteAddvalues["trade_profit"] = logInfo.TradeExplode_Profit
		siteAddvalues["trade_explode_profit"] = logInfo.TradeExplode_Profit
	}
	if logInfo.TradeKeep_Profit > 0 {
		addvalues["trade_profit"] = logInfo.TradeKeep_Profit
		addvalues["trade_keep_profit"] = logInfo.TradeKeep_Profit
		siteAddvalues["trade_profit"] = logInfo.TradeKeep_Profit
		siteAddvalues["trade_keep_profit"] = logInfo.TradeKeep_Profit
	}
	if logInfo.WithDraw > 0 {
		addvalues["withdraw"] = logInfo.WithDraw
		siteAddvalues["withdraw"] = logInfo.WithDraw
		siteAddvalues["withdraw_num"] = 1
	}
	if uinfo.IsAgent != 1 {
		m.AddSiteCount(logInfo.CreateTime, siteAddvalues)
	}
	if uinfo.ChaneelId != "0" && uinfo.ChaneelId != "" {
		channelID := utils.GetInt(uinfo.ChaneelId)
		m.AddAgentLog(channelID, logInfo.CreateTime, siteAddvalues)
		m.AddAgentLevelLog(channelID, uinfo.ChannelLevel, logInfo.CreateTime, siteAddvalues)
	}
	m.AddUserCount(uid, logInfo.CreateTime, addvalues)
	if uinfo.ParentOrder != "" {
		parentIDs := strings.Split(uinfo.ParentOrder, ",")
		n := len(parentIDs)
		for _, v := range parentIDs {
			m.AddTeamCountLog(utils.GetInt(v), logInfo, n, uid)
			n--
		}
	}
}

func (m *CreditLogModel) AddTeamCountLog(uid int, logInfo *QueueTeamLog, level int, childUID int) {
	childUinfo := MODEL_USER.GetBaseInfo(childUID)
	if childUinfo == nil {
		return
	}
	addvalues := map[string]float64{}
	userCountValues := map[string]float64{}
	if logInfo.MiningCount > 0 {
		count := config.GlobalDB.GetCount(DB_TABLE_MINING_ORDER, db.DB_PARAMS{"uid": childUID, "type": 1})
		if count == 1 {
			if rates, ok := MINING_INCOME_RATES[level]; ok {
				miningIncome := logInfo.MiningCount * rates[1]
				config.GlobalDB.AddValue(DB_TABLE_USER, map[string]float64{"mining_income": miningIncome}, db.DB_PARAMS{"id": uid})
				config.GlobalDB.InsertData(DB_TABLE_INCOME_LOG, db.DB_PARAMS{
					"uid":         uid,
					"income":      miningIncome,
					"type":        INCOME_TYPE_MINING_BUY,
					"child_uid":   childUID,
					"child_email": childUinfo.Email,
					"createtime":  logInfo.CreateTime,
					"level":       level,
				})
			}
		}
		addvalues["mining_count"] = logInfo.MiningCount
		userCountValues["team_mining"] = logInfo.MiningCount
	}
	if logInfo.MiningProfit > 0 {
		addvalues["mining_profit"] = logInfo.MiningProfit
		userCountValues["team_mining_profit"] = logInfo.MiningProfit
	}
	if logInfo.TradeBB > 0 {
		addvalues["trade"] = logInfo.TradeBB
		addvalues["trade_bb"] = logInfo.TradeBB
		userCountValues["team_trade"] = logInfo.TradeBB
		userCountValues["team_trade_bb"] = logInfo.TradeBB
	}
	if logInfo.TradeExplode > 0 {
		addvalues["trade"] = logInfo.TradeExplode
		addvalues["trade_explode"] = logInfo.TradeExplode
		userCountValues["team_trade"] = logInfo.TradeExplode
		userCountValues["team_trade_explode"] = logInfo.TradeExplode
	}
	if logInfo.TradeKeep > 0 {
		addvalues["trade"] = logInfo.TradeKeep
		addvalues["trade_keep"] = logInfo.TradeKeep
		userCountValues["team_trade"] = logInfo.TradeKeep
		userCountValues["team_trade_keep"] = logInfo.TradeKeep
	}
	if logInfo.TradeBB_Profit != 0 {
		addvalues["trade_profit"] = logInfo.TradeBB_Profit
		addvalues["trade_bb_profit"] = logInfo.TradeBB_Profit
		userCountValues["team_trade_profit"] = logInfo.TradeBB_Profit
		userCountValues["team_trade_bb_profit"] = logInfo.TradeBB_Profit
	}
	if logInfo.TradeExplode_Profit != 0 {
		addvalues["trade_profit"] = logInfo.TradeExplode_Profit
		addvalues["trade_explode_profit"] = logInfo.TradeExplode_Profit
		userCountValues["team_trade_profit"] = logInfo.TradeExplode_Profit
		userCountValues["team_trade_explode_profit"] = logInfo.TradeExplode_Profit
	}
	if logInfo.TradeKeep_Profit != 0 {
		addvalues["trade_profit"] = logInfo.TradeKeep_Profit
		addvalues["trade_keep_profit"] = logInfo.TradeKeep_Profit
		userCountValues["team_trade_profit"] = logInfo.TradeKeep_Profit
		userCountValues["team_trade_keep_profit"] = logInfo.TradeKeep_Profit
	}
	if logInfo.Recharge > 0 {
		count := config.GlobalDB.GetCount(DB_TABLE_RECHARGE, db.DB_PARAMS{"uid": childUID, "state": 1})
		if count == 1 {
			addvalues["pro_num"] = 1
			userCountValues["pro_register_num"] = 1
			if level == 1 {
				userCountValues["direct_pro_register_num"] = 1
			}
		}
		if logInfo.Recharge > 1000 {
			if rate, ok := RECHARGE_INCOME_RATES[level]; ok {
				incomeLogInfo, _ := config.GlobalDB.FetchOne(DB_TABLE_INCOME_LOG, db.DB_PARAMS{"uid": uid, "child_uid": childUID}, db.DB_FIELDS{"createtime"}, "order by createtime desc")
				if incomeLogInfo == nil || time.Unix(int64(incomeLogInfo["createtime"].ToInt()), 0).Day() != time.Now().Day() {
					rechargeIncome := logInfo.Recharge * rate
					config.GlobalDB.AddValue(DB_TABLE_USER, map[string]float64{"recharge_income": rechargeIncome}, db.DB_PARAMS{"id": uid})
					config.GlobalDB.InsertData(DB_TABLE_INCOME_LOG, db.DB_PARAMS{
						"uid":         uid,
						"income":      rechargeIncome,
						"type":        INCOME_TYPE_RECHARGE,
						"child_uid":   childUID,
						"child_email": childUinfo.Email,
						"createtime":  logInfo.CreateTime,
						"level":       level,
					})
					MODEL_USER.Update(childUID, db.DB_PARAMS{"income_time": utils.GetNow()})
				}
			}
		}
		addvalues["recharge"] = logInfo.Recharge
		userCountValues["team_recharge"] = logInfo.Recharge
	}
	if logInfo.WithDraw > 0 {
		addvalues["withdraw"] = logInfo.WithDraw
		userCountValues["team_withdraw"] = logInfo.WithDraw
	}
	m.AddLevelCount(uid, level, logInfo.CreateTime, addvalues)
	m.AddUserCount(uid, logInfo.CreateTime, userCountValues)
}
