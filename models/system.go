package models

import (
	"cointrade/config"
	"cointrade/lib/db"
	"cointrade/utils"
	"encoding/json"
	"fmt"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	COIN_BUY_STATE_NOTENGOUGH = 100001 //剩余数量不足
	COIN_BUY_STATE_NOMONEY    = 100002 //余额不足
)

type ExplodeConfig struct {
	//交割全局配置
	Time     int     `json:"time"`
	Winrate  float64 `json:"winrate"`
	Loserate float64 `json:"loserate"`
	Minprice float64 `json:"minprice"`
}
type SystemModel struct {
	ModelBase
}
type RechargeContractConfig struct {
	Contract    string  `json:"contract"`     //合约名称
	Min         float64 `json:"min"`          //充值最小金额
	WithDrawMin float64 `json:"withdraw_min"` //提现最小金额
	Address     string  `json:"address"`      //充值地址
}
type RechargeConfig struct {
	CoinType  string                    `json:"cointype"`  //币种
	Logo      string                    `json:"logo"`      //logo
	Contracts []*RechargeContractConfig `json:"contracts"` //合约类型
}

type KlineControlConfig struct {
	StartTime   int     `json:"starttime"`    //开始时间
	EndTime     int     `json:"endtime"`      //结束时间
	TargetPrice float64 `json:"target_price"` //目标价格
}
type CoinKlineConfig struct { //Kline控制参数
	MaxPrice   float64 `json:"max_price"`   //最高价格
	MinPrice   float64 `json:"min_price"`   //最低价格
	Heart      float64 `json:"heart"`       //每一跳的震荡幅度最大幅度 在震荡幅度内随机
	UpRate     int     `json:"up_rate"`     //看涨几率 1-100
	HighRate   float64 `json:"high_rate"`   //高价幅度 1-5 最合适 百分比
	LowRate    float64 `json:"low_rate"`    //低价幅度 1-5 最合适 百分比
	BaseAmount int     `json:"base_amount"` //购买量基础数值 在数值内随机
}

var PERIOD_LIST []string = []string{"1min", "5min", "15min", "30min", "60min", "4hour", "1day", "1mon", "1week", "1year"}

func (m *SystemModel) GetLoanProductList() map[int]float64 {
	rs := make(map[int]float64)
	list, _ := config.GlobalDB.FetchAll(DB_TABLE_LOAN_PRODUCT, db.DB_PARAMS{"state": 1}, db.DB_FIELDS{})
	for _, v := range list {
		rs[v["circle"].ToInt()] = v["rate"].ToFloat()
	}
	return rs
}
func (m *SystemModel) GetLastCoinInfo(pair string) primitive.M {
	d := config.GlobalMongo.GetOne("lastkline", bson.M{"pair": pair, "period": "1min"}, nil)
	return d
}
func (m *SystemModel) GetLastCoinData() map[string]interface{} {
	//取得最新的行情信息
	rs := make(map[string]interface{})
	for _, v := range COIN_LIST {
		d := config.GlobalMongo.GetOne("lastkline", bson.M{"pair": v["pair"], "period": "1day"}, nil)
		rs[v["pair"]] = d
		/*if v["is_f"] == "1" {
			m, b := rs[v["pair"]].(primitive.M)
			if b {
				if _, ok := m["close"]; ok {
					m["close"] = utils.GetFloat(v["f_price"])
					rs[v["pair"]] = m
				}
			}

		}*/
	}
	return rs
}

func (m *SystemModel) GetCoinHistoryData() map[string]interface{} { //获取行情历史所有信息
	rs := make(map[string]interface{})
	for _, v := range COIN_LIST {
		d := config.GlobalMongo.GetList("tradedata", bson.M{"pair": v["pair"]}, bson.M{"createtime": -1}, 500)
		rs[v["pair"]] = d
	}
	return rs
}
func (m *SystemModel) GetCoinMbp() map[string]interface{} { //获取市场深度所有信息
	rs := make(map[string]interface{})
	for _, v := range COIN_LIST {
		selldata := config.GlobalMongo.GetList(v["pair"]+"_mbp_sell", bson.M{}, nil, 0)
		buydata := config.GlobalMongo.GetList(v["pair"]+"_mbp_buy", bson.M{}, nil, 0)
		d := map[string]interface{}{"sell": selldata, "buy": buydata}
		rs[v["pair"]] = d
	}
	return rs
}
func (m *SystemModel) GetCoinTradeDetail() map[string]interface{} { //获取所有市场交易信息
	rs := make(map[string]interface{})
	for _, v := range COIN_LIST {
		selldata := config.GlobalMongo.GetList(v["pair"]+"_tradedetail_sell", bson.M{}, nil, 0)
		buydata := config.GlobalMongo.GetList(v["pair"]+"_tradedetail_buy", bson.M{}, nil, 0)
		d := map[string]interface{}{"sell": selldata, "buy": buydata}
		rs[v["pair"]] = d
	}
	return rs
}
func (m *SystemModel) GetCoinLastKline() map[string]map[string]interface{} { //获取所有K线最新信息
	rs := make(map[string]map[string]interface{})
	for _, v := range COIN_LIST {
		rs[v["pair"]] = make(map[string]interface{})
		for _, period := range PERIOD_LIST {
			rs[v["pair"]][period] = config.GlobalMongo.GetOne("lastkline", bson.M{"pair": v["pair"], "period": period}, nil)

			/*if v["is_f"] == "1" {

				m, b := rs[v["pair"]][period].(primitive.M)
				if b {
					if _, ok := m["close"]; ok {
						m["close"] = utils.GetFloat(v["f_price"])
						rs[v["pair"]][period] = m
					}
				}

			}*/
		}
	}
	return rs
}
func (m *SystemModel) GetCoinHitoryKline() map[string]map[string]interface{} { //获取所有K线历史信息
	rs := make(map[string]map[string]interface{})
	for _, v := range COIN_LIST {
		rs[v["pair"]] = make(map[string]interface{})
		for _, period := range PERIOD_LIST {
			rs[v["pair"]][period] = config.GlobalMongo.GetList(v["pair"]+"_kline_"+period, bson.M{}, nil, 0)
		}
	}
	return rs
}
func (m *SystemModel) GetRechargeConfig() map[string]*RechargeConfig { //获得支付/提现配置
	list, _ := config.GlobalDB.FetchAll(DB_TABLE_RECHARGE_ADDRESS, db.DB_PARAMS{"state": 1}, db.DB_FIELDS{})
	rs := make(map[string]*RechargeConfig)
	for _, v := range list {
		tmp := new(RechargeContractConfig)
		tmp.Address = v["address"].ToString()
		tmp.Contract = v["contract"].ToString()
		tmp.Min = v["min"].ToFloat()
		tmp.WithDrawMin = v["withdraw_min"].ToFloat()
		if _, ok := rs[v["cointype"].ToString()]; ok {
			rs[v["cointype"].ToString()].Contracts = append(rs[v["cointype"].ToString()].Contracts, tmp)
		} else {
			t := new(RechargeConfig)
			t.CoinType = v["cointype"].ToString()
			t.Contracts = make([]*RechargeContractConfig, 0)
			t.Contracts = append(t.Contracts, tmp)
			t.Logo = v["logo"].ToString()
			rs[v["cointype"].ToString()] = t
		}

	}
	return rs
}

func (m *SystemModel) GetControlExplode(sn string) int32 { // 控制交割合约输赢 1 赢 2 输
	one := config.GlobalMongo.GetOne("explode_control", bson.M{"sn": sn}, nil)
	if one != nil && one["result_time"].(int32) == 0 {
		return one["result"].(int32)
	}
	return 0
}
func (m *SystemModel) GetControlKline(pair string) map[int]float64 {
	one := config.GlobalMongo.GetOne("kline_control", bson.M{"pair": pair}, nil)
	if one == nil {
		return nil
	}
	rs := make([]map[string]interface{}, 0)
	s, _ := json.Marshal(one["timemap"])
	err := json.Unmarshal([]byte(s), &rs)
	re := make(map[int]float64, 0)
	if err == nil && len(rs) > 0 {

		for _, item := range rs {
			now, _ := strconv.Atoi(item["Key"].(string))
			price, _ := item["Value"].(float64)
			re[now] = price
		}
		return re
	}
	return map[int]float64{}
}

func (m *SystemModel) GetOneRechargeConfig(cointype string, contract string) *RechargeContractConfig {
	//取得一个合约配置
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
	//取得所有币的信息
	list, _ := config.GlobalDB.FetchRows(DB_TABLE_COINS, db.DB_PARAMS{}, db.DB_FIELDS{})

	return list
}
func (m *SystemModel) GetAllCoinsView() db.DB_LIST_RESULT {
	//取得所有币的信息
	list, _ := config.GlobalDB.FetchRows(DB_TABLE_COINS, db.DB_PARAMS{}, db.DB_FIELDS{})
	viewList := make([]map[string]string, 0)
	for _, v := range list {
		v["kline_config"] = ""
		viewList = append(viewList, v)
	}
	return viewList
}
func (m *SystemModel) GetNewCoins() db.DB_LIST_RESULT {
	//获取所有新币
	list, _ := config.GlobalDB.FetchRows(DB_TABLE_COINS, db.DB_PARAMS{"is_new": 1, "is_market": 1}, db.DB_FIELDS{})
	return list
}
func (m *SystemModel) GetBuyCoinList() []map[string]string {
	list, _ := config.GlobalDB.FetchRows(DB_TABLE_COINS, db.DB_PARAMS{ "is_f": 1}, db.DB_FIELDS{})
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
		tmp := new(ExplodeConfig)
		tmp.Loserate = v["lose_rate"].ToFloat()
		tmp.Winrate = v["win_rate"].ToFloat()
		tmp.Time = v["time"].ToInt()
		tmp.Minprice = v["minprice"].ToFloat()

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
func (m *SystemModel) RuleText(rule_type string, lang string) db.DB_LIST_RESULT { //获取规则文案
	cache_id := m.MakeCacheId(rule_type, lang)
	one := make(db.DB_LIST_RESULT, 0)
	err := config.GlobalRedis.GetObject(HASH_RULE_TEXT, cache_id, &one)
	if err == nil {
		return one
	}
	one, _ = config.GlobalDB.FetchRows(DB_TABLE_RULE_TEXT, db.DB_PARAMS{"rule_type": rule_type, "lang": lang}, db.DB_FIELDS{})
	if one != nil {
		config.GlobalRedis.SetValue(HASH_RULE_TEXT, cache_id, one)
	}
	return one
}

func (m *SystemModel) CoinDesc(symbol string, lang string) map[string]string {
	cache_id := m.MakeCacheId(symbol, lang)
	cache := make(map[string]string, 0)
	config.GlobalRedis.GetObject(HASH_RULE_TEXT, cache_id, &cache)
	if len(cache) > 0 {
		return cache
	}
	one, _ := config.GlobalDB.FetchRow(DB_TABLE_COIN_DESC, db.DB_PARAMS{"lang": lang, "symbol": symbol}, db.DB_FIELDS{})
	if one != nil {
		config.GlobalRedis.SetValue(DB_TABLE_COIN_DESC, cache_id, one)
	}
	return one
}

func (m *SystemModel) RuleOne(data map[string]interface{}) db.DB_ROW_RESULT {
	one := make(db.DB_ROW_RESULT, 0)
	where := ""
	if id, ok := data["id"]; ok && utils.GetInt(fmt.Sprintf("%v", id)) > 0 {
		where = fmt.Sprintf("id = '%v'", id)
	} else {
		lang := "en"
		if _, ok := data["lang"]; ok {
			lang = fmt.Sprintf("%v", data["lang"])
		}
		if _, ok := data["rule"]; !ok {
			return nil
		}
		where = fmt.Sprintf("lang = '%s' AND rule_type = '%s'", lang, data["rule"])

		if lang=="th"{
			one, _ = config.GlobalDB.FetchRow(DB_TABLE_RULE_TEXT, db.DB_PARAMS{"_": where}, db.DB_FIELDS{})
			return one
		}
	}
	cache_id := m.MakeCacheId("rule", utils.Md5(where))
	err := config.GlobalRedis.GetObject(HASH_RULE_TEXT, cache_id, &one)
	if err == nil {
		return one
	}
	one, _ = config.GlobalDB.FetchRow(DB_TABLE_RULE_TEXT, db.DB_PARAMS{"_": where}, db.DB_FIELDS{})
	if one != nil {
		config.GlobalRedis.SetValue(HASH_RULE_TEXT, cache_id, one)
	}
	return one
}

func (m *SystemModel) BuyCoin(uid int, coin_id int, amount float64) *BaseResponse {
	//新币申购
	ntime := utils.GetNow()
	coininfo, _ := config.GlobalDB.FetchOne(DB_TABLE_COINS, db.DB_PARAMS{"id": coin_id}, db.DB_FIELDS{})
	if coininfo == nil {
		return &BaseResponse{
			State: STATE_FAILD,
			Msg:   "no this coin",
		}
	}
	leave_amount := coininfo["all_amount"].ToInt() - coininfo["selled_amount"].ToInt()
	if amount > float64(leave_amount) {
		return &BaseResponse{
			State: COIN_BUY_STATE_NOTENGOUGH,
			Msg:   "not enough",
		}
	}
	allprice := amount * coininfo["f_price"].ToFloat()
	uinfo := MODEL_USER.GetBaseInfo(uid)
	if allprice <= 0 {
		return &BaseResponse{
			State: STATE_FAILD,
			Msg:   "no this coin",
		}
	}
	if uinfo.Credit < allprice {
		return &BaseResponse{
			State: COIN_BUY_STATE_NOMONEY,
			Msg:   "not enough",
		}
	}
	if MODEL_USER.AddCredit(uid, &CreditValue{
		Credit:          -1 * allprice,
		LockCredit:      allprice,
		UserCoinLogType: COIN_LOG_BUY_COIN,
		UserCoinLogInfo: QueueCreditLog{
			Credit:     -1 * allprice,
			LockCredit: allprice,
			CreateTime: ntime,
			CoinType:   "usdt",
		},
	}) {
		insertData := db.DB_PARAMS{"uid": uid, "coin_id": coininfo["id"].ToInt(), "coin_symbol": coininfo["symbol"].ToString(), "coin_pair": coininfo["pair"].ToString()}
		insertData["amount"] = amount
		insertData["price"] = coininfo["f_price"].ToFloat()
		insertData["all_price"] = allprice
		insertData["createtime"] = ntime
		config.GlobalDB.InsertData(DB_TABLE_BUY_COIN_ORDER, insertData)
		return &BaseResponse{
			State: STATE_SUCCESS,
			Msg:   "ok",
		}
	}
	return nil
}
func (m *SystemModel) GetBuyCoinOrders(uid int) db.DB_LIST_RESULT {
	list, _ := config.GlobalDB.FetchRows(DB_TABLE_BUY_COIN_ORDER, db.DB_PARAMS{"uid": uid}, db.DB_FIELDS{})
	return list
}
