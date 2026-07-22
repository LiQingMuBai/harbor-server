package models

import (
	"cointrade/config"
	"cointrade/lib/db"
	"cointrade/utils"
	"fmt"
	"math"
	"strings"
	"time"
)

// 账变信息包
const (
	COIN_LOG_USER_RECHARGE         = 5100001 //用户充值
	COIN_LOG_USER_WITHDRAW         = 5100002 //用户提现
	COIN_LOG_USER_PROFIT           = 5100003 //用户矿机收益
	COIN_LOG_USER_CLOSE            = 5100004 //用户平仓
	COIN_LOG_USER_DELEGATE         = 5100005 //用户委托
	COIN_LOG_USER_DELEGATE_SUCCESS = 5100010 //用户委托成功
	COIN_LOG_USER_CANCLE           = 5100007 //用户撤单
	COIN_LOG_USER_BUY_MINING       = 5100008 //用户购买矿机
	COIN_LOG_USER_MINING_PROFIT    = 5100009 //用户挖矿获利

	COIN_LOG_USER_CLEAR_INCOME    = 5100011 //用户提取下级返利收入
	COIN_LOG_USER_WITHDRAW_FAILD  = 5100012 //用户提现失败
	COIN_LOG_USER_EXCHANGE        = 5100013 //用户兑换
	COIN_LOG_BB_TRADE             = 5100014 //币币交易
	COIN_LOG_EXPLODE_TRADE        = 5100015 //交割交易
	COIN_LOG_KEEP_TRADE           = 5100016 //永续交易
	COIN_LOG_BACKEND              = 5100017 //后台划转
	COIN_LOG_USER_MINING_BACK     = 5100018 //用户矿机本金返还
	COIN_LOG_ASSETS_EXCHANGE      = 5100019 //币种兑换
	COIN_LOG_MINING_UNLOCK        = 5100020 //矿机解锁
	COIN_LOG_BUY_COIN             = 5100021 //用户申购新币
	COIN_LOG_EXCHANGE_ACCOUNT_IN  = 5100022 //用户资产转换 资金账户到合约账户
	COIN_LOG_EXCHANGE_ACCOUNT_OUT = 5100023 //用户资产转换 合约账户到资金账户
	COIN_LOG_LOAN_BACK            = 5100025 //归还贷款
	COIN_LOG_LORA_IN              = 5100024 //贷款成功
	COIN_LOG_KEEP_BREAK           = 5100026 //杠杆穿仓
	COIN_LOG_USER_REVERVATION     = 5100027 // 预约扣除金额
	TEAM_LOG_RECHARGE             = 5200001 //用户充值-团队
	TEAM_LOG_WITHDRAW             = 5200002 //提现-团队
	TEAM_LOG_MINING               = 5200003 //挖矿总投入
	TEAM_LOG_MINING_PROFIT        = 5200004 //挖矿总收益
	TEAM_LOG_TRADE                = 5200005 //交易总投入
	TEAM_LOG_TRADE_PROFIT         = 5200006 //交易总收益
	COIN_LOG_SPOT_BACK            = 5200028 //现货申购退还

	LOG_TIMETYPE_ALL   = 0 //全部
	LOG_TIMETYPE_DAY   = 1 //当日
	LOG_TIMETYPE_MONTH = 2 //当月

	INCOME_TYPE_RECHARGE      = 1 //下级充值返利
	INCOME_TYPE_MINING_BUY    = 2 //下级购买矿机返利
	INCOME_TYPE_MINING_PROFIT = 3 //下级矿机收入返利
)

type CreditLogModel struct {
	ModelBase
}
type CreditLogInfo struct { //用户账变添加结构
	Uid        int     `json:"uid"`
	Credit     float64 `json:"credit"`
	LockCredit float64 `json:"lockcredit"`
	Mode       int     `josn:"mode"`
	Sn         string  `josn:"sn"`
	Type       int     `json:"type"`
	CoinType   string  `json:"cointype"`
	Createtime int     `json:"credittime"`
}
type QueueCreditLog struct { //用户账变入队结构
	Credit     float64 `json:"credit"`
	LockCredit float64 `json:"lockcredit"`
	CoinType   string  `json:"cointype"`
	Sn         string  `json:"sn"`
	CreateTime int     `json:"createtime"`
}
type QueueTeamLog struct { //团队统计入队结构
	Recharge            float64 `json:"recharge"`             //个人充值量
	WithDraw            float64 `json:"withdraw"`             //个人提现量
	Trade               float64 `json:"trade"`                //个人交易量
	TradeProfit         float64 `json:"trade_profit"`         //个人交易获利
	MiningCount         float64 `json:"mining_count"`         //个人矿机投入额度
	MiningProfit        float64 `json:"mining_profit"`        //个人矿机获利
	CreateTime          int     `json:"createtime"`           //发生时间
	TradeBB             float64 `json:"trade_bb"`             //币币交易
	TradeExplode        float64 `json:"trade_explode"`        //交割交易
	TradeKeep           float64 `json:"trade_keep"`           //永续合约
	TradeBB_Profit      float64 `json:"trade_bb_profit"`      //币币利润
	TradeExplode_Profit float64 `json:"trade_explode_profit"` //交割利润
	TradeKeep_Profit    float64 `json:"trade_keep_profit"`    //永续利润
	//IsFirstRecharge int     `json:"is_first_recharge"` //是否首充
	/*RegisterNum          int     `json:"register_num"`            //下级注册数量
	ProRegisterNum       int     `json:"pro_register_num"`        //下级有效注册量
	DirectRegisterNum    int     `json:"direct_register_num"`     //直属注册量
	DirectProRegisterNum int     `json:"direct_pro_register_num"` //直属有效注册
	TeamRecharge         float64 `json:"team_recharge"`           //团队充值总额
	TeamWithDraw         float64 `json:"team_withdraw"`           //团队提现总额
	TeamMining           float64 `json:"team_mining"`             //团队挖矿投入总额
	TeamMiningProfit     float64 `json:"team_mining_profit"`      //团队挖矿总利润
	TeamTrade            float64 `json:"team_trade"`              //团队交易总投入
	TeamTradeProfit      float64 `json:"team_trade_profit"`       //团队交易总获利*/
}
type CoinLogRequest struct {
	PageBaseRequest
	Type int `json:"type"`
}

func (m *CreditLogModel) Add(info *CreditLogInfo) {
	insertData := db.DB_PARAMS{}
	uinfo := MODEL_USER.GetBaseInfo(info.Uid)
	if uinfo == nil {
		return
	}

	insertData["uid"] = info.Uid
	insertData["credit"] = info.Credit
	insertData["lock_credit"] = info.LockCredit
	insertData["mode"] = info.Mode
	insertData["sn"] = info.Sn
	insertData["type"] = info.Type
	insertData["createtime"] = info.Createtime
	insertData["after_credit"] = uinfo.Credit
	insertData["cointype"] = info.CoinType
	config.GlobalDB.InsertData(DB_TABLE_CREDIT_LOG, insertData)
}

func (m *CreditLogModel) GetList(uid int, rq *CoinLogRequest) *PageBaseResponse {
	//取得账变列表
	condition := db.DB_PARAMS{"uid": uid}
	if rq.Type != 0 {
		condition["type"] = rq.Type
	}
	//1 代表显示，0代表不显示
	condition["display"] = 1
	count := config.GlobalDB.GetCount(DB_TABLE_CREDIT_LOG, condition)
	pagesize := int(math.Ceil(float64(count) / float64(rq.Limit)))
	if rq.Page > pagesize {
		rq.Page = pagesize
	}
	if rq.Page == 0 {
		rq.Page = 1
	}
	limitstr := fmt.Sprintf("limit %d,%d", (rq.Page-1)*rq.Limit, rq.Limit)
	list, _ := config.GlobalDB.FetchRows(DB_TABLE_CREDIT_LOG, condition, db.DB_FIELDS{}, "order by id desc", limitstr)
	rs := new(PageBaseResponse)
	rs.State = STATE_SUCCESS
	rs.Msg = "success"
	rs.Limit = rq.Limit
	rs.Page = rq.Page
	rs.PageTotal = pagesize
	rs.Total = count
	rs.List = list
	return rs
}
func (m *CreditLogModel) AddUserCount(uid int, createtime int, data map[string]float64) {
	//增加团队统计

	t := time.Unix(int64(createtime), 0)
	daytime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).Unix()
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
	} else {
		insertData := db.DB_PARAMS{"uid": uid}
		for k, v := range data {
			insertData[k] = v
		}
		config.GlobalDB.InsertData(DB_TABLE_USER_COUNT_SUM, insertData)
	}
}
func (m *CreditLogModel) AddLevelCount(uid int, level int, createtime int, data map[string]float64) {
	//增加用户层级统计
	t := time.Unix(int64(createtime), 0)
	daytime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).Unix()
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
		_, err := config.GlobalDB.InsertData(DB_TABLE_USER_LEVEL_COUNT, condition)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	one, _ = config.GlobalDB.FetchOne(DB_TABLE_USER_LEVEL_COUNT_SUM, db.DB_PARAMS{"uid": uid, "level": level}, db.DB_FIELDS{"id"}, "limit 0,1")
	if one != nil {
		config.GlobalDB.AddValue(DB_TABLE_USER_LEVEL_COUNT_SUM, data, db.DB_PARAMS{"id": one["id"].Value})
	} else {
		insertData := db.DB_PARAMS{"uid": uid, "level": level, "email": uinfo.Email}
		for k, v := range data {
			insertData[k] = v
		}
		_, err := config.GlobalDB.InsertData(DB_TABLE_USER_LEVEL_COUNT_SUM, insertData)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

}
func (m *CreditLogModel) AddSiteCount(createtime int, data map[string]float64) {
	//增加全站统计
	t := time.Unix(int64(createtime), 0)
	daytime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).Unix()
	condition := db.DB_PARAMS{"daytime": daytime}
	one, _ := config.GlobalDB.FetchOne(DB_TABLE_SITE_COUNT, condition, db.DB_FIELDS{"id"})
	if one != nil {
		config.GlobalDB.AddValue(DB_TABLE_SITE_COUNT, data, condition)
	} else {
		for k, v := range data {
			condition[k] = v
		}
		config.GlobalDB.InsertData(DB_TABLE_SITE_COUNT, condition)
	}
}
func (m *CreditLogModel) AddUserCountLog(uid int, logInfo *QueueTeamLog) {
	//增加团队统计
	uinfo := MODEL_USER.GetBaseInfo(uid)
	if uinfo == nil {
		return
	}

	addvalues := map[string]float64{}
	site_addvalues := map[string]float64{}
	if logInfo.MiningCount > 0 {
		addvalues["mining_count"] = logInfo.MiningCount
		site_addvalues["minning_count"] = logInfo.MiningCount
	}
	if logInfo.MiningProfit > 0 {
		addvalues["mining_profit"] = logInfo.MiningProfit
		site_addvalues["minning_profit"] = logInfo.MiningProfit
	}
	if logInfo.Recharge > 0 {
		addvalues["recharge"] = logInfo.Recharge
		site_addvalues["recharge"] = logInfo.Recharge
		site_addvalues["register_num"] = 1
		count := config.GlobalDB.GetCount(DB_TABLE_RECHARGE, db.DB_PARAMS{"uid": uid, "state": 1})
		if count == 1 {
			site_addvalues["pro_num"] = 1
			site_addvalues["first_recharge"] = logInfo.Recharge
			site_addvalues["first_recharge_num"] = 1
		}
	}
	if logInfo.TradeBB > 0 {
		addvalues["trade"] = logInfo.TradeBB
		addvalues["trade_bb"] = logInfo.TradeBB
		site_addvalues["trade"] = logInfo.TradeBB
		site_addvalues["trade_bb"] = logInfo.TradeBB
	}
	if logInfo.TradeExplode > 0 {
		addvalues["trade"] = logInfo.TradeExplode
		addvalues["trade_explode"] = logInfo.TradeExplode
		site_addvalues["trade"] = logInfo.TradeExplode
		site_addvalues["trade_explode"] = logInfo.TradeExplode
	}
	if logInfo.TradeKeep > 0 {
		addvalues["trade"] = logInfo.TradeKeep
		addvalues["trade_keep"] = logInfo.TradeKeep
		site_addvalues["trade"] = logInfo.TradeKeep
		site_addvalues["trade_keep"] = logInfo.TradeKeep
	}
	/*if logInfo.TradeProfit > 0 {
		addvalues["trade_profit"] = logInfo.TradeProfit
		site_addvalues["trade_profit"] = logInfo.TradeProfit
	}*/
	if logInfo.TradeBB_Profit > 0 {
		addvalues["trade_profit"] = logInfo.TradeBB_Profit
		addvalues["trade_bb_profit"] = logInfo.TradeBB_Profit
		site_addvalues["trade_profit"] = logInfo.TradeBB_Profit
		site_addvalues["trade_bb_profit"] = logInfo.TradeBB_Profit
	}
	if logInfo.TradeExplode_Profit > 0 {
		addvalues["trade_profit"] = logInfo.TradeExplode_Profit
		addvalues["trade_explode_profit"] = logInfo.TradeExplode_Profit
		site_addvalues["trade_profit"] = logInfo.TradeExplode_Profit
		site_addvalues["trade_explode_profit"] = logInfo.TradeExplode_Profit
	}
	if logInfo.TradeKeep_Profit > 0 {
		addvalues["trade_profit"] = logInfo.TradeKeep_Profit
		addvalues["trade_keep_profit"] = logInfo.TradeKeep_Profit
		site_addvalues["trade_profit"] = logInfo.TradeKeep_Profit
		site_addvalues["trade_keep_profit"] = logInfo.TradeKeep_Profit
	}
	if logInfo.WithDraw > 0 {
		addvalues["withdraw"] = logInfo.WithDraw

		site_addvalues["withdraw"] = logInfo.WithDraw
		site_addvalues["withdraw_num"] = 1
	}
	if uinfo.IsAgent != 1 {
		m.AddSiteCount(logInfo.CreateTime, site_addvalues)
	}
	if uinfo.ChaneelId != "0" && uinfo.ChaneelId != "" {
		//utils.WriteLog("/home/agent.log", site_addvalues)
		m.AddAgentLog(utils.GetInt(uinfo.ChaneelId), logInfo.CreateTime, site_addvalues)
		m.AddAgentLevelLog(utils.GetInt(uinfo.ChaneelId), uinfo.ChannelLevel, logInfo.CreateTime, site_addvalues)
	}
	m.AddUserCount(uid, logInfo.CreateTime, addvalues)
	if uinfo.ParentOrder != "" {
		parentIds := strings.Split(uinfo.ParentOrder, ",")
		n := len(parentIds)
		for _, v := range parentIds {
			m.AddTeamCountLog(utils.GetInt(v), logInfo, n, uid)
			n--
		}
	}
}
func (m *CreditLogModel) AddTeamCountLog(uid int, logInfo *QueueTeamLog, level int, child_uid int) { //用户ID 统计数据 层级 子用户ID
	//增加团队层级统计
	child_uinfo := MODEL_USER.GetBaseInfo(child_uid)
	if child_uinfo == nil {
		return
	}
	addvalues := map[string]float64{}
	user_count_values := map[string]float64{}
	if logInfo.MiningCount > 0 {
		count := config.GlobalDB.GetCount(DB_TABLE_MINING_ORDER, db.DB_PARAMS{"uid": child_uid, "type": 1}) //判断子用户是不是首次认购
		if count == 1 {
			if rates, ok := MINING_INCOME_RATES[level]; ok {

				mining_income := logInfo.MiningCount * rates[1]
				config.GlobalDB.AddValue(DB_TABLE_USER, map[string]float64{"mining_income": mining_income}, db.DB_PARAMS{"id": uid}) //增加充值返利
				config.GlobalDB.InsertData(DB_TABLE_INCOME_LOG, db.DB_PARAMS{
					"uid":         uid,
					"income":      mining_income,
					"type":        INCOME_TYPE_MINING_BUY,
					"child_uid":   child_uid,
					"child_email": child_uinfo.Email,
					"createtime":  logInfo.CreateTime,
					"level":       level,
				})
			}
		}
		addvalues["mining_count"] = logInfo.MiningCount
		user_count_values["team_mining"] = logInfo.MiningCount
	}
	if logInfo.MiningProfit > 0 {
		addvalues["mining_profit"] = logInfo.MiningProfit
		user_count_values["team_mining_profit"] = logInfo.MiningProfit
	}
	/*if logInfo.Trade > 0 {
		addvalues["trade"] = logInfo.Trade
		user_count_values["team_trade"] = logInfo.Trade
	}
	if logInfo.TradeProfit > 0 {
		addvalues["trade_profit"] = logInfo.TradeProfit
		user_count_values["team_trade_profit"] = logInfo.TradeProfit
	}*/
	if logInfo.TradeBB > 0 {
		addvalues["trade"] = logInfo.TradeBB
		addvalues["trade_bb"] = logInfo.TradeBB
		user_count_values["team_trade"] = logInfo.TradeBB
		user_count_values["team_trade_bb"] = logInfo.TradeBB
	}
	if logInfo.TradeExplode > 0 {
		addvalues["trade"] = logInfo.TradeExplode
		addvalues["trade_explode"] = logInfo.TradeExplode
		user_count_values["team_trade"] = logInfo.TradeExplode
		user_count_values["team_trade_explode"] = logInfo.TradeExplode
	}
	if logInfo.TradeKeep > 0 {
		addvalues["trade"] = logInfo.TradeKeep
		addvalues["trade_keep"] = logInfo.TradeKeep
		user_count_values["team_trade"] = logInfo.TradeKeep
		user_count_values["team_trade_keep"] = logInfo.TradeKeep
	}

	if logInfo.TradeBB_Profit != 0 {
		addvalues["trade_profit"] = logInfo.TradeBB_Profit
		addvalues["trade_bb_profit"] = logInfo.TradeBB_Profit
		user_count_values["team_trade_profit"] = logInfo.TradeBB_Profit
		user_count_values["team_trade_bb_profit"] = logInfo.TradeBB_Profit
	}
	if logInfo.TradeExplode_Profit != 0 {
		addvalues["trade_profit"] = logInfo.TradeExplode_Profit
		addvalues["trade_explode_profit"] = logInfo.TradeExplode_Profit
		user_count_values["team_trade_profit"] = logInfo.TradeExplode_Profit
		user_count_values["team_trade_explode_profit"] = logInfo.TradeExplode_Profit
	}
	if logInfo.TradeKeep_Profit != 0 {
		addvalues["trade_profit"] = logInfo.TradeKeep_Profit
		addvalues["trade_keep_profit"] = logInfo.TradeKeep_Profit
		user_count_values["team_trade_profit"] = logInfo.TradeKeep_Profit
		user_count_values["team_trade_keep_profit"] = logInfo.TradeKeep_Profit
	}
	if logInfo.Recharge > 0 {
		count := config.GlobalDB.GetCount(DB_TABLE_RECHARGE, db.DB_PARAMS{"uid": child_uid, "state": 1})

		if count == 1 {
			addvalues["pro_num"] = 1
			user_count_values["pro_register_num"] = 1
			if level == 1 {
				user_count_values["direct_pro_register_num"] = 1
			}
		}
		if logInfo.Recharge > 1000 { //大于1000的充值开始返利
			if rate, ok := RECHARGE_INCOME_RATES[level]; ok {
				income_log_info, _ := config.GlobalDB.FetchOne(DB_TABLE_INCOME_LOG, db.DB_PARAMS{"uid": uid, "child_uid": child_uid}, db.DB_FIELDS{"createtime"}, "order by createtime desc")
				if income_log_info == nil || time.Unix(int64(income_log_info["createtime"].ToInt()), 0).Day() != time.Now().Day() {
					recharge_income := logInfo.Recharge * rate

					config.GlobalDB.AddValue(DB_TABLE_USER, map[string]float64{"recharge_income": recharge_income}, db.DB_PARAMS{"id": uid}) //增加充值返利
					config.GlobalDB.InsertData(DB_TABLE_INCOME_LOG, db.DB_PARAMS{
						"uid":         uid,
						"income":      recharge_income,
						"type":        INCOME_TYPE_RECHARGE,
						"child_uid":   child_uid,
						"child_email": child_uinfo.Email,
						"createtime":  logInfo.CreateTime,
						"level":       level,
					})
					MODEL_USER.Update(child_uid, db.DB_PARAMS{"income_time": utils.GetNow()})
				}

			}
		}
		addvalues["recharge"] = logInfo.Recharge
		user_count_values["team_recharge"] = logInfo.Recharge
	}
	if logInfo.WithDraw > 0 {
		addvalues["withdraw"] = logInfo.WithDraw
		user_count_values["team_withdraw"] = logInfo.WithDraw
	}
	m.AddLevelCount(uid, level, logInfo.CreateTime, addvalues)
	m.AddUserCount(uid, logInfo.CreateTime, user_count_values)
}
func (m *CreditLogModel) GetUserCountDay(uid int) db.DB_ROW_RESULT {
	//获取用户统计列表 天
	t := time.Now()
	daytime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).Unix()
	one, _ := config.GlobalDB.FetchRow(DB_TABLE_USER_COUNT, db.DB_PARAMS{"daytime": daytime}, db.DB_FIELDS{})
	return one
}
func (m *CreditLogModel) GetUserCountSum(uid int) db.DB_ROW_RESULT { //获得用户历史统计
	one, _ := config.GlobalDB.FetchRow(DB_TABLE_USER_COUNT_SUM, db.DB_PARAMS{"uid": uid}, db.DB_FIELDS{})
	return one
}

func (m *CreditLogModel) GetUserCountMonth(uid int) db.DB_LIST_RESULT { //获取用户统计列表 月
	t := time.Now()
	year := t.Year()
	month := t.Month()
	starttime := time.Date(year, month, 1, 0, 0, 0, 0, time.Local).Unix()
	nextmonth := month + 1
	if nextmonth > 12 {
		nextmonth = 1
		year = year + 1
	}
	endtime := time.Date(year, nextmonth, 1, 0, 0, 0, 0, time.Local).Unix()
	list, _ := config.GlobalDB.FetchRows(DB_TABLE_USER_COUNT, db.DB_PARAMS{"_": fmt.Sprintf("daytime>=%d and daytime<%d", starttime, endtime)}, db.DB_FIELDS{})
	return list
}

func (m *CreditLogModel) GetUserLevelCountDay(uid int) map[int]db.DB_ROW_RESULT {
	t := time.Now()
	daytime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).Unix()
	condition := db.DB_PARAMS{"uid": uid, "daytime": daytime}
	rs := make(map[int]db.DB_ROW_RESULT)
	list, _ := config.GlobalDB.FetchRows(DB_TABLE_USER_LEVEL_COUNT, condition, db.DB_FIELDS{})
	for _, v := range list {
		rs[utils.GetInt(v["level"])] = v
	}
	return rs
}

func (m *CreditLogModel) GetUserLevelCountSum(uid int) map[int]db.DB_ROW_RESULT {

	condition := db.DB_PARAMS{"uid": uid}
	rs := make(map[int]db.DB_ROW_RESULT)
	list, _ := config.GlobalDB.FetchRows(DB_TABLE_USER_LEVEL_COUNT_SUM, condition, db.DB_FIELDS{})
	for _, v := range list {
		rs[utils.GetInt(v["level"])] = v
	}
	return rs
}

func (m *CreditLogModel) GetUserLevelCountMonth(uid int) map[int]db.DB_LIST_RESULT {
	t := time.Now()
	year := t.Year()
	month := t.Month()
	starttime := time.Date(year, month, 1, 0, 0, 0, 0, time.Local).Unix()
	nextmonth := month + 1
	if nextmonth > 12 {
		nextmonth = 1
		year = year + 1
	}
	endtime := time.Date(year, nextmonth, 1, 0, 0, 0, 0, time.Local).Unix()
	list, _ := config.GlobalDB.FetchRows(DB_TABLE_USER_LEVEL_COUNT, db.DB_PARAMS{"_": fmt.Sprintf("daytime>=%d and daytime<%d", starttime, endtime)}, db.DB_FIELDS{})
	rs := make(map[int]db.DB_LIST_RESULT)
	for _, v := range list {
		level := utils.GetInt(v["level"])
		_, ok := rs[level]
		if !ok {
			rs[level] = make(db.DB_LIST_RESULT, 0)
		}
		rs[level] = append(rs[level], v)
	}
	return rs
}

func (m *CreditLogModel) AddAgentLog(uid int, createtime int, data map[string]float64) {
	//增加渠道统计
	t := time.Unix(int64(createtime), 0)
	daytime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).Unix()
	one, _ := config.GlobalDB.FetchOne(DB_TABLE_AGENT_COUNT, db.DB_PARAMS{"uid": uid, "daytime": daytime}, db.DB_FIELDS{"id"})
	if one != nil {
		config.GlobalDB.AddValue(DB_TABLE_AGENT_COUNT, data, db.DB_PARAMS{"id": one["id"].Value})
	} else {
		//uinfo := MODEL_USER.GetBaseInfo(uid)
		insertData := db.DB_PARAMS{"uid": uid, "daytime": daytime}
		for k, v := range data {
			insertData[k] = v
		}
		config.GlobalDB.InsertData(DB_TABLE_AGENT_COUNT, insertData)

	}

}
func (m *CreditLogModel) AddAgentLevelLog(uid int, level int, createtime int, data map[string]float64) {
	//增加渠道层级统计
	t := time.Unix(int64(createtime), 0)
	daytime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).Unix()
	one, _ := config.GlobalDB.FetchOne(DB_TABLE_AGENT_LEVEL_COUNT, db.DB_PARAMS{"uid": uid, "daytime": daytime, "level": level}, db.DB_FIELDS{"id"})
	if one != nil {
		config.GlobalDB.AddValue(DB_TABLE_AGENT_LEVEL_COUNT, data, db.DB_PARAMS{"id": one["id"].Value})
	} else {
		uinfo := MODEL_USER.GetBaseInfo(uid)
		insertData := db.DB_PARAMS{"uid": uid, "email": uinfo.Email, "daytime": daytime, "level": level}
		for k, v := range data {
			insertData[k] = v
		}
		config.GlobalDB.InsertData(DB_TABLE_AGENT_LEVEL_COUNT, insertData)

	}
}
func (m *CreditLogModel) IncomeLog(uid int, rq *PageBaseRequest) *PageBaseResponse {
	rs := new(PageBaseResponse)
	count := config.GlobalDB.GetCount(DB_TABLE_INCOME_LOG, db.DB_PARAMS{"uid": uid})
	pagesize := int(math.Ceil(float64(count) / float64(rq.Limit)))
	if rq.Page > pagesize {
		rq.Page = pagesize
	}
	if rq.Page == 0 {
		rq.Page = 1
	}
	limitstr := fmt.Sprintf("limit %d,%d", (rq.Page-1)*rq.Limit, rq.Limit)
	list, _ := config.GlobalDB.FetchRows(DB_TABLE_INCOME_LOG, db.DB_PARAMS{"uid": uid}, db.DB_FIELDS{}, "order by id desc", limitstr)
	rs.State = STATE_SUCCESS
	rs.Limit = rq.Limit
	rs.Page = rq.Page
	rs.Msg = "success"
	rs.PageTotal = pagesize
	rs.Total = count
	rs.List = list
	return rs
}
