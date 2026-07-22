package models

import (
	"cointrade/config"
	"cointrade/lib/db"
	"cointrade/utils"
	"fmt"
	"math"
	"math/rand"
	"time"
)

//交易模块
type TradeModel struct {
	ModelBase
}

const (
	OPEN_TYPE_KEEP               = 1 //永续
	OPEN_TYPE_EXPLODE            = 2 //交割
	OPEN_TYPE_BB                 = 3 //币币交易
	PRICE_TYPE_LIMIT             = 1 //价格类型 限价
	PRICE_TYPE_MARKET            = 2 //价格类型 市场价格
	DIRECT_TYPE_BIG              = 1 //买涨
	DIRECT_TYPE_SMALL            = 2 //买空
	DELEGATE_TYPE_BUY            = 1 //开仓
	DELEGATE_TYPE_SELL           = 2 //平仓
	TRADE_BUY_PREFIX             = "B"
	TRADE_SELL_PREFIX            = "S"
	DELEGATE_STATE_NOCOIN        = 40001 //持仓不足
	DELEGATE_STATE_CREDIT        = 40002 //余额不足
	DELEGATE_STATE_CLOSETIME     = 40003 //交割合约下平仓时间间隔不合法
	DELEGATE_STATE_NOASSET       = 40004 //没有这个资产 币币交易出售时
	DELEGATE_STATE_MIN           = 40005 //低于交易最小值
	DELEGATE_STATE_TRADE_CLOSED  = 40006 //此类型的交易已关闭
	DELEGATE_STATE_GANGGAN_ERROR = 40007 //杠杆倍率错误
)

type TradeDelegateRequest struct { //委托请求
	OpenType              int     `json:"opentype"`           //开单类型
	DelegateType          int     `json:"closeoropen"`        //委托类型 开仓 还是平仓
	Pair                  string  `json:"pair"`               //交易对
	Coin                  string  `json:"coin"`               //币种
	PriceType             int     `json:"pricetype"`          //价格类型 当前市场价格
	DirectType            int     `json:"directtype"`         //买卖方向
	GangGan               int     `json:"ganggan"`            //杠杆倍率
	Amount                float64 `json:"amount"`             //购买量 永续合约按张数算 1张=1000U 交割合约按 USDT数量算 币币交易按购买的COIN当前汇率或者限价来算
	Price                 float64 `json:"price"`              //限价模式下用户提交的价格
	CloseTime             int     `json:"closetime"`          //平仓时间间隔 交割合约下存在
	StopUpPrice           float64 `json:"stop_up_price"`      //止盈价格 为0为不指定
	StopDownPrice         float64 `json:"stop_down_price"`    //止损价格 为0为不指定
	StopUpDelegatePrice   float64 `json:"stop_up_delegate"`   //止盈委托价格 为0为跟随市场价
	StopDownDelegatePrice float64 `json:"stop_down_delegate"` //止损委托价格 为0为跟随市场价
	Sn                    string  `json:"sn"`                 //订单号 杠杆手动平仓时指定 其他不需要 如果指定了SN其他类型的平仓会返回错误
}

type OpenedInfo struct { //持仓信息结构
	Id            int     `json:"id"`  //系统Id
	Uid           int     `json:"uid"` //用户ID
	Sn            string  `json:"sn"`  //交易单号
	UserType      int     `json:"user_type"`
	TradeType     int     `json:"tradetype"`      //交易类型 交易类型 1 永续合约 2 交割合约 3币币交易
	Flag          int     `json:"flag"`           //方向 多/空
	OpenPrice     float64 `json:"openprice"`      //开仓价格 永续合约为成本价
	ClosePrice    float64 `json:"closeprice"`     //平仓价格
	CoinId        int     `json:"coinid"`         //币系统ID
	CoinPair      string  `json:"pair"`           //币的交易对
	CoinSymbol    string  `json:"symbol"`         //币唯一标识
	CloseTime     int     `json:"closetime"`      //交割合约的下平仓时间间隔
	CloseRealTime int     `json:"close_realtime"` //实际的平仓时间 真实的时间线
	ClearTime     int     `json:"cleartime"`      //结算时间
	CreateTime    int     `json:"createtime"`     //开仓时间
	Ganggan       int     `json:"gangan"`         //杠杆倍率
	Credit        float64 `json:"credit"`         //总投入额
	Profit        float64 `json:"profit"`         //产生的利润
	WinRate       float64 `json:"winrate"`        //交割合约下赢的比例
	LoseRate      float64 `json:"loserate"`       //交割合约下输的比列
	Num           float64 `json:"num"`            //币的总量
	LockNum       float64 `json:"lock_num"`
	UserName      string  `json:"username"`
	Mode          int     `json:"mode"` //交易模式
}

type CloseTrade struct {
	Id         int     `json:"id"`
	Uid        int     `json:"uid"`
	Sn         string  `json:"sn"`
	CoinSymbol string  `json:"coin_symbol"`
	TradeType  string  `json:"trade_type"`
	Flag       int     `json:"flag"`
	Amount     float64 `json:"amount"`
	ClosePrice float64 `json:"close_price"`
	CreateTime int     `json:"createtime"`
	Num        float64 `json:"num"`
	Mode       int     `json:"mode"`
	AllPrice   float64 `json:"allprice"`
	Oprice     float64 `json:"o_price"`

	UserType int     `json:"user_type"`
	Profit   float64 `json:"profit"`
	UserName string  `json:"username"`
}

type DelegateInfo struct {
	Id  int `json:"id"`
	Uid int `json:"uid"`

	UserType     int     `json:"user_type"`
	Sn           string  `json:"sn"`
	DelegameType int     `json:"delegate_type"`
	TradeType    int     `json:"trade_type"`
	Flag         int     `json:"flag"`
	Fee          float64 `json:"fee"`
	Price        float64 `json:"price"`
	CoinId       int     `json:"coinid"`
	CoinPair     string  `json:"coinpair"`
	CoinSymbol   string  `json:"coin_symbol"`
	CloseTime    int     `json:"close_time"`
	Createtime   int     `json:"createtime"`
	Credit       float64 `json:"credit"`
	Num          float64 `json:"num"`
	State        int     `json:"state"`
	Mode         int     `json:"mode"`
	ChangeTime   int     `json:"changetime"`
	UserName     string  `json:"username"`
}
type TradeListRequest struct {
	//委托列表请求
	PageBaseRequest
	TradeType    int    `json:"tradetype"`
	Flag         int    `json:"flag"`
	State        int    `json:"state"`
	Coin         string `json:"coin"`
	DelegateType int    `json:"delegate_type"`
	Ganggan      int    `json:"ganggan"` //是否为杠杆订单 1 为只获取杠杆订单 0 为不指定
}

func (m *TradeModel) MakeSn(uid int, t int) string {
	uidstr := utils.Sup(int64(uid), 10)
	timestr := time.Now().Format("20060102150405")
	switch t {
	case DELEGATE_TYPE_BUY:
		return fmt.Sprintf("%s%s%s%d", TRADE_BUY_PREFIX, timestr, uidstr, 10+rand.Intn(89))
	case DELEGATE_TYPE_SELL:
		return fmt.Sprintf("%s%s%s%d", TRADE_SELL_PREFIX, timestr, uidstr, 10+rand.Intn(89))
	}

	return fmt.Sprintf("%s%s%s%d", TRADE_SELL_PREFIX, timestr, uidstr, 10+rand.Intn(89))
}
func (m *TradeModel) DelegateTrade(uid int, rq *TradeDelegateRequest) *BaseResponse {
	//提交委托单
	rs := new(BaseResponse)
	uinfo := MODEL_USER.GetBaseInfo(uid)
	ntime := utils.GetNow()
	u_credit := uinfo.Credit //用户余额
	//fmt.Println("rq:", rq)
	if rq.GangGan < 1 {
		rq.GangGan = 1
	}
	if rq.GangGan > 100 {
		rs.State = DELEGATE_STATE_GANGGAN_ERROR
		rs.Msg = "error ganggan"
		return rs
	}
	team_log_type := 0
	var team_log_info QueueTeamLog
	if uinfo.Mode == USER_MODE_V {
		u_credit = uinfo.VCredit //虚拟模式时用虚拟余额
	}
	if uinfo == nil {
		rs.State = STATE_SYSTEM_ERROR
		rs.Msg = "system error"
		return rs
	}
	if rq.OpenType != OPEN_TYPE_BB && rq.OpenType != OPEN_TYPE_EXPLODE && rq.OpenType != OPEN_TYPE_KEEP { //开单类型不对
		rs.State = STATE_SYSTEM_ERROR
		rs.Msg = "system error"
		return rs
	}

	coinfo := MODEL_SYSTEM.GetCoinInfo(rq.Coin, rq.Pair)

	if coinfo == nil {
		rs.State = DELEGATE_STATE_NOCOIN
		rs.Msg = "no this coin"
		return rs
	}
	if rq.DirectType != DIRECT_TYPE_BIG { //非多即空
		rq.DirectType = DIRECT_TYPE_SMALL
	}
	if rq.GangGan < 1 {
		rq.GangGan = 1
	}
	if rq.PriceType != PRICE_TYPE_MARKET { //不是市场价就是限价
		rq.PriceType = PRICE_TYPE_LIMIT
	}
	if rq.DelegateType != DELEGATE_TYPE_BUY {
		rq.DelegateType = DELEGATE_TYPE_SELL
	}
	sn := m.MakeSn(uid, rq.DelegateType)
	//coinPriceInfo := config.GlobalMongo.GetOne("lastdata", bson.M{"pair": rq.Pair}, nil) //从MONGO中获取到最新的价格
	coinPriceInfo := MODEL_SYSTEM.GetLastCoinInfo(rq.Pair)
	//allprice := 0.0
	coinChangeType := COIN_LOG_USER_DELEGATE
	//fmt.Println("dsadsa11", coinPriceInfo)
	realprice := 0.0
	insertData := db.DB_PARAMS{}
	insertData["uid"] = uid                       //插入数据
	insertData["delegate_type"] = rq.DelegateType //委托类型 开/平
	insertData["trade_type"] = rq.OpenType        //开仓类型 1 永续 2 交割 3 币币
	insertData["flag"] = rq.DirectType            //方向 1 多 2 空
	insertData["coinid"] = coinfo["id"]           //币系统ID
	insertData["coinpair"] = coinfo["pair"]       //币交易对
	insertData["coin_symbol"] = coinfo["symbol"]  //币的唯一标识符
	insertData["ganggan"] = rq.GangGan            //杠杆
	insertData["state"] = 0                       //状态 0 是委托中
	insertData["mode"] = uinfo.Mode               //用户当前的模式 是虚拟的还是真实的
	insertData["createtime"] = ntime              //委托时间
	insertData["sn"] = sn
	insertData["stop_up_price"] = rq.StopUpPrice                //止盈价格
	insertData["stop_down_price"] = rq.StopDownPrice            //止亏价格
	insertData["stop_up_delegate"] = rq.StopUpDelegatePrice     //止盈委托价格
	insertData["stop_down_delegate"] = rq.StopDownDelegatePrice //止亏委托价格
	coinprice := coinPriceInfo["close"].(float64)

	if rq.PriceType != PRICE_TYPE_MARKET && rq.OpenType != OPEN_TYPE_EXPLODE { //交割合约只按市价
		coinprice = utils.FormatFloatA(rq.Price, utils.GetInt(coinfo["dnum"]))
		if coinprice <= 0 {
			rs.State = DELEGATE_STATE_MIN
			rs.Msg = "trade too small 1"
			return rs
		}
	}

	switch rq.OpenType {
	//针对委托合约的类型不同进行不同的操作
	case OPEN_TYPE_BB: //这里暂时缺少虚拟交易的逻辑流程 后面贯通后记得添加
		//币币合约委托开始
		//if rq.PriceType == PRICE_TYPE_MARKET { //按市价委托
		if coinfo["open_coin2coin"] == "0" {
			rs.State = DELEGATE_STATE_TRADE_CLOSED
			rs.Msg = "trade closed"
			return rs
		}

		coinChangeType = COIN_LOG_BB_TRADE
		rq.Amount = utils.FormatFloatA(rq.Amount, utils.GetInt(coinfo["cnum"])) //数量小数位的处理
		if rq.Amount <= 0 {
			rs.State = DELEGATE_STATE_MIN
			rs.Msg = "trade too small 2"
			return rs
		}
		allprice := rq.Amount * coinprice
		realprice = allprice
		if rq.DelegateType == DELEGATE_TYPE_BUY { //开仓委托时
			if coinfo["isnative"] == "0" && coinfo["f_price"] != "0" { //当自发币 发行价格还为大于0的时候要进行审批
				insertData["is_f"] = 1
			} //如果是自发币则不自动开仓
			if u_credit < allprice { //余额不足
				rs.State = DELEGATE_STATE_CREDIT
				rs.Msg = "not enough credit"
				return rs
			}
		} else {
			//平仓委托时
			userAssets := MODEL_ASSETS.GetAllAssets(uid, uinfo.Mode)
			if _, ok := userAssets[rq.Coin]; !ok { //没有这个资产
				rs.State = DELEGATE_STATE_NOASSET
				rs.Msg = "you dont have this assets"
				return rs
			}
			if userAssets[rq.Coin].Count < rq.Amount { //这个币种的数量不足
				rs.State = DELEGATE_STATE_CREDIT
				rs.Msg = "not enough credit"
				return rs
			}
			if coinfo["isnative"] == "0" {
				if userAssets[rq.Coin].TransOpenTime > ntime {
					rs.State = DELEGATE_STATE_TRADE_CLOSED
					rs.Msg = "locktime"
					return rs
				}
			}
			MODEL_ASSETS.AddAssets(uid, &Assets{ //平仓委托时改变用户的资产 需要先将用户的资产冻结
				Coin:    rq.Coin,
				Pair:    userAssets[rq.Coin].Pair,
				Num:     -1 * rq.Amount,
				LockNum: rq.Amount,
				Price:   coinprice,
				Mode:    uinfo.Mode,
			})

		}

		insertData["price"] = coinprice //单价
		insertData["credit"] = allprice //总额
		insertData["num"] = rq.Amount   //数量

	case OPEN_TYPE_KEEP:
		//永续合约开始
		if coinfo["open_trade"] == "0" {
			rs.State = DELEGATE_STATE_TRADE_CLOSED
			rs.Msg = "trade closed"
			return rs
		}
		allprice := rq.Amount * 1000
		coinChangeType = COIN_LOG_KEEP_TRADE
		if rq.DelegateType == DELEGATE_TYPE_BUY { //开仓时
			fee := allprice * (config.GlobalConfig.GetValue("trade_fee").ToFloat() / float64(100))
			insertData["fee"] = fee      //开仓手续费
			real_price := allprice + fee //实际付出的价格
			realprice = real_price
			if u_credit < real_price {
				rs.State = DELEGATE_STATE_CREDIT
				rs.Msg = "not enough credit"
				return rs
			}
			coincount := allprice / coinprice
			coincount = utils.FormatFloatA(coincount, utils.GetInt(coinfo["cnum"]))
			if coincount <= 0 {
				rs.State = DELEGATE_STATE_MIN
				rs.Msg = "trade too small 3"
				return rs
			}

		} else {
			//平仓时 要先判断是否持仓如果持仓不足的话不允许平仓的
			var one *OpenedInfo
			if rq.Sn != "" {

				one = m.GetOpendBySn(uid, rq.Sn)
				if one.Ganggan <= 1 {
					rs.State = STATE_FAILD
					rs.Msg = "not a ganggan trade"
					return rs
				}
				insertData["ganggan_sn"] = rq.Sn
			} else {
				one = m.GetOpendOne(uid, rq.Coin, rq.OpenType, rq.DirectType, uinfo.Mode, rq.GangGan)
			}

			if one == nil || rq.Amount > one.Num {
				rs.State = DELEGATE_STATE_NOCOIN
				rs.Msg = "no this assets"
				return rs
			}
			config.GlobalDB.AddValue(DB_TABLE_OPENED_TRADE, map[string]float64{"num": -1 * rq.Amount, "lock_num": rq.Amount}, db.DB_PARAMS{"id": one.Id, "mode": uinfo.Mode}) //锁定卖单的张数

		}
		insertData["price"] = coinprice
		insertData["credit"] = allprice
		insertData["num"] = rq.Amount

	case OPEN_TYPE_EXPLODE:
		//交割合约开始 交割合约只有委托的买单 不可能出现主动的卖单
		if coinfo["open_trade"] == "0" {
			rs.State = DELEGATE_STATE_TRADE_CLOSED
			rs.Msg = "trade closed"
			return rs
		}
		coinChangeType = COIN_LOG_EXPLODE_TRADE
		econfig, ok := EXPLODE_CONFIG[rq.CloseTime]
		if !ok {
			rs.State = DELEGATE_STATE_CLOSETIME
			rs.Msg = "incorrect close time"
			return rs
		}
		if rq.Amount < econfig.Minprice {
			rs.State = DELEGATE_STATE_MIN
			rs.Msg = "too small"
			return rs
		}
		if u_credit < rq.Amount {
			rs.State = DELEGATE_STATE_CREDIT
			rs.Msg = "not enough credit"
			return rs
		}
		realprice = rq.Amount
		insertData["price"] = coinprice
		insertData["credit"] = rq.Amount
		insertData["num"] = rq.Amount
		insertData["close_time"] = rq.CloseTime
	}
	var err error
	lock_price := realprice
	if rq.OpenType == OPEN_TYPE_EXPLODE {
		//交割合约直接进场
		lock_price = 0
		insertData := db.DB_PARAMS{}
		insertData["uid"] = uid
		insertData["sn"] = sn
		insertData["trade_type"] = OPEN_TYPE_EXPLODE
		insertData["flag"] = rq.DirectType
		insertData["openprice"] = coinprice
		insertData["closeprice"] = 0
		insertData["coinid"] = coinfo["id"]
		insertData["coinpair"] = coinfo["pair"]
		insertData["coin_symbol"] = coinfo["symbol"]
		insertData["close_time"] = rq.CloseTime
		insertData["close_real_time"] = ntime + rq.CloseTime
		insertData["clear_time"] = 0
		insertData["createtime"] = ntime
		insertData["ganggan"] = 1
		insertData["credit"] = rq.Amount
		insertData["profit"] = 0
		if explodeConfig, ok := EXPLODE_CONFIG[rq.CloseTime]; ok {
			insertData["win_rate"] = explodeConfig.Winrate
			insertData["lose_rate"] = explodeConfig.Loserate
		} else {
			insertData["win_rate"] = 100
			insertData["lose_rate"] = 100
		}

		insertData["num"] = rq.Amount
		insertData["mode"] = uinfo.Mode
		team_log_type = TEAM_LOG_TRADE
		team_log_info.TradeExplode = rq.Amount
		team_log_info.CreateTime = ntime
		_, err = config.GlobalDB.InsertData(DB_TABLE_OPENED_TRADE, insertData) //增加交割合约的持仓

	} else {

		_, err = config.GlobalDB.InsertData(DB_TABLE_DELEGATE_TRADE, insertData) //插入委托表
	}
	//_, err := config.GlobalDB.InsertData(DB_TABLE_DELEGATE_TRADE, insertData) //插入委托表
	if err != nil {
		rs.State = STATE_SYSTEM_ERROR
		rs.Msg = "system error"
		return rs
	}

	if rq.DelegateType == DELEGATE_TYPE_BUY {
		//买单的话需要对用户余额进行处理 卖单委托不需要处理余额
		if uinfo.Mode == USER_MODE_REAL { //后面记得处理账变
			MODEL_USER.AddCredit(uid, &CreditValue{
				Credit:          -1 * realprice,
				LockCredit:      lock_price,
				VCrdit:          0,
				LockVCredit:     0,
				UserCoinLogType: coinChangeType, //委托账变
				UserCoinLogInfo: QueueCreditLog{
					Credit:     -1 * realprice,
					LockCredit: lock_price,
					Sn:         sn,
					CreateTime: ntime,
				},
				TeamCoinLogType: team_log_type,
				TeamCoinLogInfo: team_log_info,
			})
		} else { //虚拟交易 不管账变
			MODEL_USER.AddCredit(uid, &CreditValue{
				Credit:          0,
				LockCredit:      0,
				VCrdit:          -1 * realprice,
				LockVCredit:     lock_price,
				UserCoinLogType: 0,
				UserCoinLogInfo: nil,
				TeamCoinLogType: 0,
				TeamCoinLogInfo: nil,
			})
		}
	}

	rs.State = STATE_SUCCESS
	rs.Msg = utils.GetJsonValue(insertData)
	return rs
}
func (m *TradeModel) GetCloseBySn(uid int, sn string) db.DB_ROW_RESULT {
	one, _ := config.GlobalDB.FetchRow(DB_TABLE_CLOSE_TRADE, db.DB_PARAMS{"sn": sn, "uid": uid}, db.DB_FIELDS{})
	return one
}
func (m *TradeModel) GetOpendBySn(uid int, sn string) *OpenedInfo {
	one, _ := config.GlobalDB.FetchOne(DB_TABLE_OPENED_TRADE, db.DB_PARAMS{"sn": sn, "uid": uid}, db.DB_FIELDS{})
	if one == nil {
		return nil
	}

	rs := new(OpenedInfo)
	rs.Id = one["id"].ToInt()
	rs.Uid = uid
	rs.TradeType = one["trade_type"].ToInt()
	rs.ClearTime = one["clear_time"].ToInt()
	rs.ClosePrice = one["closeprice"].ToFloat()
	rs.CloseRealTime = one["close_real_time"].ToInt()
	rs.CloseTime = one["close_time"].ToInt()
	rs.CoinId = one["coinid"].ToInt()
	rs.CoinPair = one["coinpair"].ToString()
	rs.CoinSymbol = one["coin_symbol"].ToString()
	rs.CreateTime = one["createtime"].ToInt()
	rs.Ganggan = one["ganggan"].ToInt()
	rs.WinRate = one["win_rate"].ToFloat()
	rs.LoseRate = one["lose_rate"].ToFloat()
	rs.Credit = one["credit"].ToFloat()
	rs.Profit = one["profit"].ToFloat()
	rs.Num = one["num"].ToFloat()
	rs.Mode = one["mode"].ToInt()
	rs.Sn = one["sn"].ToString()
	rs.OpenPrice = one["openprice"].ToFloat()
	return rs
}
func (m *TradeModel) GetOpendOne(uid int, coin string, trade_type int, flag int, mode int, ganggan int) *OpenedInfo {
	//取得用户持仓信息 UID 用户ID coin 币种唯一标识 trade_type 交易类型 state 状态 flag 多/空
	one, _ := config.GlobalDB.FetchOne(DB_TABLE_OPENED_TRADE, db.DB_PARAMS{"trade_type": trade_type, "uid": uid, "flag": flag, "coin_symbol": coin, "mode": mode, "ganggan": ganggan}, db.DB_FIELDS{})
	if one == nil {
		return nil
	}

	rs := new(OpenedInfo)
	rs.Id = one["id"].ToInt()
	rs.Uid = uid
	rs.TradeType = one["trade_type"].ToInt()
	rs.ClearTime = one["clear_time"].ToInt()
	rs.ClosePrice = one["closeprice"].ToFloat()
	rs.CloseRealTime = one["close_real_time"].ToInt()
	rs.CloseTime = one["close_time"].ToInt()
	rs.CoinId = one["coinid"].ToInt()
	rs.CoinPair = one["coinpair"].ToString()
	rs.CoinSymbol = one["coin_symbol"].ToString()
	rs.CreateTime = one["createtime"].ToInt()
	rs.Ganggan = one["ganggan"].ToInt()
	rs.WinRate = one["win_rate"].ToFloat()
	rs.LoseRate = one["lose_rate"].ToFloat()
	rs.Credit = one["credit"].ToFloat()
	rs.Profit = one["profit"].ToFloat()
	rs.Num = one["num"].ToFloat()
	rs.Mode = one["mode"].ToInt()
	rs.Sn = one["sn"].ToString()
	rs.OpenPrice = one["openprice"].ToFloat()
	return rs
}
func (m *TradeModel) AddKeepOpend(delegateInfo db.DBValues) {
	//从委托中增加用永续合约持仓
	//coin := delegateInfo["coin_symbol"].ToString()
	ntime := utils.GetNow()
	var openinfo *OpenedInfo
	if delegateInfo["ganggan"].ToInt() > 1 {
		openinfo = nil
	} else {
		openinfo = m.GetOpendOne(delegateInfo["uid"].ToInt(), delegateInfo["coin_symbol"].ToString(), delegateInfo["trade_type"].ToInt(), delegateInfo["flag"].ToInt(), delegateInfo["mode"].ToInt(), delegateInfo["ganggan"].ToInt())
	}

	if openinfo == nil {
		insertData := db.DB_PARAMS{}
		insertData["uid"] = delegateInfo["uid"].ToInt()
		insertData["trade_type"] = delegateInfo["trade_type"].ToInt()
		insertData["closeprice"] = 0
		insertData["flag"] = delegateInfo["flag"].ToInt()
		insertData["openprice"] = delegateInfo["price"].ToFloat()
		insertData["coinid"] = delegateInfo["coinid"].ToInt()
		insertData["coinpair"] = delegateInfo["coinpair"].ToString()
		insertData["coin_symbol"] = delegateInfo["coin_symbol"].ToString()
		insertData["close_time"] = 0
		insertData["close_real_time"] = 0
		insertData["clear_time"] = 0
		insertData["createtime"] = ntime
		insertData["ganggan"] = delegateInfo["ganggan"].ToInt()
		insertData["credit"] = delegateInfo["credit"].ToFloat()
		insertData["profit"] = 0
		insertData["win_rate"] = 0
		insertData["lose_rate"] = 0
		insertData["num"] = delegateInfo["num"].ToFloat()
		insertData["mode"] = delegateInfo["mode"].ToInt()
		insertData["sn"] = delegateInfo["sn"].ToString()
		insertData["stop_up_price"] = delegateInfo["stop_up_price"].ToFloat()
		insertData["stop_down_price"] = delegateInfo["stop_down_price"].ToFloat()
		insertData["stop_up_delegate"] = delegateInfo["stop_up_delegate"].ToFloat()
		insertData["stop_down_delegate"] = delegateInfo["stop_down_delegate"].ToFloat()
		_, err := config.GlobalDB.InsertData(DB_TABLE_OPENED_TRADE, insertData)
		if err == nil {
			config.GlobalDB.UpdateData(DB_TABLE_DELEGATE_TRADE, db.DB_PARAMS{"state": 1, "changetime": ntime}, db.DB_PARAMS{"id": delegateInfo["id"].Value}) //将委托单状态改为已完成
		}
	} else {
		//当已有此币种的持仓的时候 需要开始计算成本价
		o_allprice := openinfo.OpenPrice * openinfo.Num
		n_allprice := delegateInfo["price"].ToFloat() * delegateInfo["num"].ToFloat()
		o_price := (o_allprice + n_allprice) / (delegateInfo["num"].ToFloat() + openinfo.Num)
		config.GlobalDB.AddValue(DB_TABLE_OPENED_TRADE, map[string]float64{"num": delegateInfo["num"].ToFloat(), "credit": delegateInfo["credit"].ToFloat()}, db.DB_PARAMS{"id": openinfo.Id})
		config.GlobalDB.UpdateData(DB_TABLE_OPENED_TRADE, db.DB_PARAMS{"openprice": o_price}, db.DB_PARAMS{"id": openinfo.Id})
		config.GlobalDB.UpdateData(DB_TABLE_DELEGATE_TRADE, db.DB_PARAMS{"state": 1, "changetime": ntime}, db.DB_PARAMS{"id": delegateInfo["id"].Value}) //将委托单状态改为已完成
	}

}
func (m *TradeModel) GetDelegateList(uid int, rq *TradeListRequest) *PageBaseResponse {
	//得到委托列表
	uinfo := MODEL_USER.GetBaseInfo(uid)
	if uinfo == nil {
		return nil
	}
	condition := db.DB_PARAMS{"uid": uid, "mode": uinfo.Mode}
	if rq.Coin != "" {
		condition["coin_symbol"] = rq.Coin
	}
	if rq.Flag != 0 {
		condition["flag"] = rq.Flag
	}
	if rq.State != -1 {
		condition["state"] = rq.State
	}
	if rq.TradeType != 0 {
		condition["trade_type"] = rq.TradeType
	}
	if rq.DelegateType != 0 {
		condition["delegate_type"] = rq.DelegateType
	}
	if rq.Ganggan > 0 {
		condition["_"] = "num>0 and ganggan>1"
	} else {
		condition["_"] = "num>0 and ganggan<=1"
	}
	count := config.GlobalDB.GetCount(DB_TABLE_DELEGATE_TRADE, condition)
	pagesize := int(math.Ceil(float64(count) / float64(rq.Limit)))
	if rq.Page > pagesize {
		rq.Page = pagesize
	}
	if rq.Page <= 0 {
		rq.Page = 1
	}
	limitstr := fmt.Sprintf("limit %d,%d", (rq.Page-1)*rq.Limit, rq.Limit)
	list, _ := config.GlobalDB.FetchRows(DB_TABLE_DELEGATE_TRADE, condition, db.DB_FIELDS{}, "order by id desc", limitstr)
	rs := new(PageBaseResponse)
	rs.Limit = rq.Limit
	rs.Page = rq.Page
	rs.List = list
	rs.PageTotal = pagesize
	rs.State = STATE_SUCCESS
	rs.Msg = ""
	return rs
}

func (m *TradeModel) GetOpendList(uid int, rq *TradeListRequest) *PageBaseResponse {
	//取得持仓列表
	uinfo := MODEL_USER.GetBaseInfo(uid)
	if uinfo == nil {
		return nil
	}

	condition := db.DB_PARAMS{"uid": uid, "mode": uinfo.Mode, "clear_time": 0}
	if rq.Coin != "" {
		condition["coin_symbol"] = rq.Coin
	}
	if rq.Flag != 0 {
		condition["flag"] = rq.Flag
	}
	if rq.TradeType != 0 {
		condition["trade_type"] = rq.TradeType
	}
	if rq.Ganggan > 0 {
		condition["_"] = "num>0 and ganggan>1"
	} else {
		condition["_"] = "num>0 and ganggan<=1"
	}
	count := config.GlobalDB.GetCount(DB_TABLE_OPENED_TRADE, condition)
	pagesize := int(math.Ceil(float64(count) / float64(rq.Limit)))
	if rq.Page > pagesize {
		rq.Page = pagesize
	}
	if rq.Page <= 0 {
		rq.Page = 1
	}
	limitstr := fmt.Sprintf("limit %d,%d", (rq.Page-1)*rq.Limit, rq.Limit)
	list, _ := config.GlobalDB.FetchRows(DB_TABLE_OPENED_TRADE, condition, db.DB_FIELDS{}, "order by id desc", limitstr)

	rs := new(PageBaseResponse)
	rs.Limit = rq.Limit
	rs.Page = rq.Page
	rs.List = list
	rs.PageTotal = pagesize
	rs.State = STATE_SUCCESS
	rs.Msg = fmt.Sprintf("%d", utils.GetNow())
	return rs
}

func (m *TradeModel) GetCloseList(uid int, rq *TradeListRequest) *PageBaseResponse { //取得平仓列表
	uinfo := MODEL_USER.GetBaseInfo(uid)
	if uinfo == nil {
		return nil
	}

	condition := db.DB_PARAMS{"uid": uid, "mode": uinfo.Mode, "_": "num>0"}
	if rq.Coin != "" {
		condition["coin_symbol"] = rq.Coin
	}
	if rq.Flag != 0 {
		condition["flag"] = rq.Flag
	}
	if rq.TradeType != 0 {
		condition["trade_type"] = rq.TradeType
	}
	count := config.GlobalDB.GetCount(DB_TABLE_CLOSE_TRADE, condition)
	pagesize := int(math.Ceil(float64(count) / float64(rq.Limit)))
	if rq.Page > pagesize {
		rq.Page = pagesize
	}
	if rq.Page <= 0 {
		rq.Page = 1
	}
	limitstr := fmt.Sprintf("limit %d,%d", (rq.Page-1)*rq.Limit, rq.Limit)
	list, _ := config.GlobalDB.FetchRows(DB_TABLE_CLOSE_TRADE, condition, db.DB_FIELDS{}, "order by id desc", limitstr)
	rs := new(PageBaseResponse)
	rs.Limit = rq.Limit
	rs.Page = rq.Page
	rs.List = list
	rs.PageTotal = pagesize
	rs.State = STATE_SUCCESS
	rs.Msg = ""
	return rs
}
func (m *TradeModel) CancleDelegate(uid int, sn string) *BaseResponse { //撤单
	rs := new(BaseResponse)
	one, _ := config.GlobalDB.FetchOne(DB_TABLE_DELEGATE_TRADE, db.DB_PARAMS{"uid": uid, "sn": sn, "state": 0}, db.DB_FIELDS{})
	if one == nil {
		rs.State = STATE_FAILD
		rs.Msg = "no this delegate order"
		return rs
	}

	user_credit := one["credit"].ToFloat() + one["fee"].ToFloat()
	user_v_credit := one["credit"].ToFloat() + one["fee"].ToFloat()
	if one["mode"].ToInt() == USER_MODE_REAL {
		user_v_credit = 0
	} else {
		user_credit = 0
	}
	if one["delegate_type"].ToInt() == DELEGATE_TYPE_BUY { //买单 撤单
		if MODEL_USER.AddCredit(uid, &CreditValue{
			Credit:          user_credit,
			LockCredit:      -1 * user_credit,
			VCrdit:          user_v_credit,
			LockVCredit:     -1 * user_v_credit,
			UserCoinLogType: COIN_LOG_USER_CANCLE,
			UserCoinLogInfo: QueueCreditLog{
				Credit:     user_credit,
				LockCredit: -1 * user_credit,
				Sn:         one["sn"].ToString(),
				CreateTime: utils.GetNow(),
			},
			TeamCoinLogType: 0,
			TeamCoinLogInfo: nil,
		}) {
			_, err := config.GlobalDB.UpdateData(DB_TABLE_DELEGATE_TRADE, db.DB_PARAMS{"state": 2}, db.DB_PARAMS{"id": one["id"].Value})
			if err != nil {
				rs.State = STATE_SYSTEM_ERROR
				rs.Msg = err.Error()
				return rs
			}
		}
	} else {
		//卖单撤单
		switch one["trade_type"].ToInt() {
		case OPEN_TYPE_BB: //币币撤单
			if MODEL_ASSETS.AddAssets(uid, &Assets{
				Coin:    one["coin_symbol"].ToString(),
				Pair:    one["coinpair"].ToString(),
				Num:     one["num"].ToFloat(),
				LockNum: -1 * one["num"].ToFloat(),
				Price:   one["price"].ToFloat(),
				Mode:    one["mode"].ToInt(),
			}) {
				_, err := config.GlobalDB.UpdateData(DB_TABLE_DELEGATE_TRADE, db.DB_PARAMS{"state": 2}, db.DB_PARAMS{"id": one["id"].Value})
				if err != nil {
					rs.State = STATE_SYSTEM_ERROR
					rs.Msg = err.Error()
					return rs
				}
			}
		case OPEN_TYPE_KEEP:
			flag := one["flag"].ToInt()
			coin := one["coin_symbol"].ToString()
			uinfo := MODEL_USER.GetBaseInfo(uid)
			num := one["num"].ToFloat()
			opendinfo := m.GetOpendOne(uid, coin, OPEN_TYPE_KEEP, flag, uinfo.Mode, one["ganggan"].ToInt())
			if opendinfo != nil {
				config.GlobalDB.AddValue(DB_TABLE_OPENED_TRADE, map[string]float64{"num": num, "lock_num": -1 * num}, db.DB_PARAMS{"id": opendinfo.Id})
				_, err := config.GlobalDB.UpdateData(DB_TABLE_DELEGATE_TRADE, db.DB_PARAMS{"state": 2}, db.DB_PARAMS{"id": one["id"].Value})
				if err != nil {
					rs.State = STATE_SYSTEM_ERROR
					rs.Msg = err.Error()
					return rs
				}
			}

		}
	}
	/*if one["trade_type"].ToInt() == OPEN_TYPE_BB {
		if MODEL_ASSETS.AddAssets(uid, &Assets{
			Coin:    one["coin_symbol"].ToString(),
			Pair:    one["coinpair"].ToString(),
			Num:     one["num"].ToFloat(),
			LockNum: -1 * one["num"].ToFloat(),
			Price:   one["price"].ToFloat(),
			Mode:    one["mode"].ToInt(),
		}) {
			_, err := config.GlobalDB.UpdateData(DB_TABLE_DELEGATE_TRADE, db.DB_PARAMS{"state": 2}, db.DB_PARAMS{"id": one["id"].Value})
			if err != nil {
				rs.State = STATE_SYSTEM_ERROR
				rs.Msg = err.Error()
				return rs
			}
		}
	} else {
		if MODEL_USER.AddCredit(uid, &CreditValue{
			Credit:          user_credit,
			LockCredit:      -1 * user_credit,
			VCrdit:          user_v_credit,
			LockVCredit:     -1 * user_v_credit,
			UserCoinLogType: COIN_LOG_USER_CANCLE,
			UserCoinLogInfo: QueueCreditLog{
				Credit:     user_credit,
				LockCredit: -1 * user_credit,
				Sn:         one["sn"].ToString(),
				CreateTime: utils.GetNow(),
			},
			TeamCoinLogType: 0,
			TeamCoinLogInfo: nil,
		}) {
			_, err := config.GlobalDB.UpdateData(DB_TABLE_DELEGATE_TRADE, db.DB_PARAMS{"state": 2}, db.DB_PARAMS{"id": one["id"].Value})
			if err != nil {
				rs.State = STATE_SYSTEM_ERROR
				rs.Msg = err.Error()
				return rs
			}
		}
	}*/

	rs.State = STATE_SUCCESS
	rs.Msg = "success"
	return rs
}
