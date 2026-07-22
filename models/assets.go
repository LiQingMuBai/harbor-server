package models

import (
	"cointrade/config"
	"cointrade/lib/db"
	"cointrade/utils"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// 用户资产相关包
const (
	EXCHANGE_STATE_NOCOIN      = 600001 //没有这个币种
	EXCHANGE_STATE_NOTENNOUGH  = 600002 //没有足够的余额支持兑换
	EXCHANGE_STATE_TOOMIN      = 600003 //交易量过小
	EXCHANGE_STATE_NOT_TRANS   = 600004 //未允许交易
	ASSETS_TRANS_TYPE_IN       = 1      //转入
	ASSETS_TRANS_TYPE_OUT      = 2      //转出
	ASSETS_TRANS_TYPE_CONTRACT = 3      //转入交易账户
)

type AssetModel struct {
	ModelBase
}
type AssetInfo struct {
	Symbol        string  `json:"symbol"`          //币种唯一标识
	CoinId        int     `json:"coinid"`          //币的系统ID
	Pair          string  `json:"pair"`            //交易对
	Count         float64 `json:"count"`           //拥有的数量
	O_Price       float64 `json:"o_price"`         //成本单价
	LockCount     float64 `json:"lockcount"`       //锁定数量
	Address       string  `json:"address"`         //用户私有地址
	IsTrans       int     `json:"istrans"`         //是否允许兑换
	TransOpenTime int     `json:"trans_open_time"` //交易开启时间
}
type Assets struct {
	Coin          string  //币种
	Pair          string  //交易对
	Num           float64 //数量
	LockNum       float64 //锁定数量
	Price         float64 //开仓价格
	Mode          int     //模式
	IsTrans       int     //是否可以交易划转
	OpenTransTime int     //交易划转权限开启时间
}
type ExchangeRequest struct { //兑换请求
	From   string  `json:"from"`   //来源币种
	To     string  `json:"to"`     //兑换币种
	Amount float64 `json:"amount"` //兑换数量
}
type AssetsTransRequest struct { //划转请求
	Coin      string  `json:"coin"`       //币种
	Type      int     `json:"type"`       //划转类型
	Amount    float64 `json:"amount"`     //金额
	ToAddress string  `json:"to_address"` //到达地址
}
type QuickExchangeRequest struct {
	//闪兑请求
	Coin   string  `json:"coin"`   //币种
	Amount float64 `json:"amount"` //金额
}

func (m *AssetModel) MakeSn(uid int) string { //创建订单号
	uidstr := utils.Sup(int64(uid), 10)
	timestr := time.Now().Format("20060102150405")

	return fmt.Sprintf("%s%s%s%d", "E", timestr, uidstr, 10+rand.Intn(89))
}
func (m *AssetModel) QuickExchange(uid int, rq *QuickExchangeRequest) *BaseResponse { //闪兑
	ntime := utils.GetNow()
	assetsInfo := m.GetAllAssets(uid, USER_MODE_REAL)
	//fmt.Println(rq)
	rq.Coin = strings.ToLower(rq.Coin)
	coininfo := MODEL_SYSTEM.GetCoinInfo(rq.Coin, rq.Coin+"usdt")
	if coininfo != nil && coininfo["is_market"] == "0" && rq.Coin != "usdt" && rq.Coin != "usdc" {
		return &BaseResponse{
			State: STATE_FAILD,
			Msg:   "faild",
		}
	}
	if rq.Coin == "usdt" {
		return &BaseResponse{
			State: STATE_FAILD,
			Msg:   "faild",
		}
	}
	for _, v := range assetsInfo {
		if v.Symbol == rq.Coin {
			//user_asset_info := v
			if v.IsTrans == 0 {
				return &BaseResponse{
					State: EXCHANGE_STATE_NOT_TRANS,
					Msg:   "not allowed trans",
				}
			}
			if v.Count < rq.Amount {
				return &BaseResponse{
					State: EXCHANGE_STATE_NOTENNOUGH,
					Msg:   "not enough assets",
				}
			}
			pair := rq.Coin + "usdt"
			credit := 0.0
			price := 1.0
			if rq.Coin == "usdc" {
				credit = rq.Amount
			} else {
				coininfo := MODEL_SYSTEM.GetLastCoinInfo(pair)
				credit = rq.Amount * coininfo["close"].(float64)
				price = coininfo["close"].(float64)
			}
			if v.TransOpenTime > 0 && v.TransOpenTime > ntime {
				return &BaseResponse{
					State: EXCHANGE_STATE_NOT_TRANS,
					Msg:   "not allowed trans",
				}
			}
			if m.AddAssets(uid, &Assets{
				Coin:    rq.Coin,
				Pair:    pair,
				Num:     -1 * rq.Amount,
				LockNum: 0,
				Price:   price,
				Mode:    USER_MODE_REAL,
			}) {
				MODEL_USER.AddCredit(uid, &CreditValue{
					Credit:          credit,
					UserCoinLogType: COIN_LOG_USER_EXCHANGE,
					UserCoinLogInfo: QueueCreditLog{
						Credit:     credit,
						CoinType:   "usdt",
						CreateTime: ntime,
						Sn:         m.MakeSn(uid),
					},
				})
			}
			return &BaseResponse{
				State: STATE_SUCCESS,
				Msg:   "ok",
			}

		}
	}
	return &BaseResponse{
		State: STATE_FAILD,
		Msg:   "faild",
	}
}
func (m *AssetModel) GetOneAsset(uid int, coin string) *AssetInfo {
	user_assets := m.GetAllAssets(uid, USER_MODE_REAL)
	if v, ok := user_assets[coin]; ok {
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
func (m *AssetModel) GetAllAssets(uid int, mode int) map[string]AssetInfo { //获取用户所有资产余额
	rs := make(map[string]AssetInfo)
	//cacheid := m.MakeCacheId(uid, mode)
	//err := config.GlobalRedis.GetObject(HASH_USER_ASSETS, cacheid, &rs)
	uinfo := MODEL_USER.GetBaseInfo(uid)
	rs["usdt"] = AssetInfo{
		Symbol:    "usdt",
		Count:     uinfo.Credit,
		LockCount: uinfo.LockCredit,
	}
	//if err == nil && len(rs) > 0 {
	//	return rs
	//}
	list, _ := config.GlobalDB.FetchAll(DB_TABLE_USERASSETS, db.DB_PARAMS{"uid": uid, "mode": mode}, db.DB_FIELDS{})
	for _, v := range list {
		tmp := new(AssetInfo)
		tmp.CoinId = v["coin_id"].ToInt()

		tmp.O_Price = v["o_price"].ToFloat()
		tmp.Pair = v["coin_pair"].ToString()
		tmp.Symbol = v["coin_symbol"].ToString()
		if strings.ToLower(tmp.Symbol) == "usdt" {
			uinfo := MODEL_USER.GetBaseInfo(uid)
			tmp.Count = uinfo.Credit
			tmp.LockCount = uinfo.LockCredit
		} else {
			tmp.Count = v["num"].ToFloat()
			tmp.LockCount = v["lock_num"].ToFloat()
		}

		tmp.Address = v["wallet_address"].ToString()
		tmp.IsTrans = v["is_trans"].ToInt()
		tmp.TransOpenTime = v["trans_open_time"].ToInt()
		rs[tmp.Symbol] = *tmp
	}

	//config.GlobalRedis.SetValue(HASH_USER_ASSETS, cacheid, rs)
	rs["usdt"] = AssetInfo{
		Symbol:    "usdt",
		Count:     uinfo.Credit,
		LockCount: uinfo.LockCredit,
	}
	return rs
}

func (m *AssetModel) AddAssets(uid int, a *Assets) bool { //给用户添加资产 清除用户缓存 这块预估不会出现资产添加冲突 所以使用UPDATE方式来更新
	cacheid := m.MakeCacheId(uid, USER_MODE_REAL)
	//coinInfo := MODEL_SYSTEM.GetCoinInfo(a.Coin, a.Pair)
	if a.Coin == "usdc" {
		a.Price = 1
	}
	if a.Coin == "usdt" {
		fmt.Println("usdt", a.Num)
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
		o_allprice := one["o_price"].ToFloat() * one["num"].ToFloat()
		o_price := ((a.Price * a.Num) + o_allprice) / (one["num"].ToFloat() + a.Num) //计算出当前成本价
		//config.GlobalDB.AddValue(DB_TABLE_USERASSETS, map[string]float64{""})
		//_, err := config.GlobalDB.UpdateData(DB_TABLE_USERASSETS, db.DB_PARAMS{"num": one["num"].ToFloat() + a.Num, "o_price": o_price}, db.DB_PARAMS{"id": one["id"].Value})
		err := config.GlobalDB.AddValue(DB_TABLE_USERASSETS, map[string]float64{"num": a.Num, "lock_num": a.LockNum}, db.DB_PARAMS{"id": one["id"].Value})
		config.GlobalDB.UpdateData(DB_TABLE_USERASSETS, db.DB_PARAMS{"o_price": o_price, "trans_open_time": a.OpenTransTime}, db.DB_PARAMS{"id": one["id"].Value}) //更新成本价
		if err != nil {
			return false
		}
	} else {
		insertData := db.DB_PARAMS{}
		insertData["uid"] = uid
		insertData["coin_symbol"] = a.Coin
		insertData["coin_id"] = coinInfo["id"]
		insertData["coin_pair"] = a.Pair
		insertData["num"] = a.Num
		insertData["lock_num"] = 0
		insertData["o_price"] = a.Price
		insertData["mode"] = a.Mode
		insertData["trans_open_time"] = a.OpenTransTime
		_, err := config.GlobalDB.InsertData(DB_TABLE_USERASSETS, insertData)
		if err != nil {
			return false
		}
	}
	config.GlobalRedis.Del(HASH_USER_ASSETS, cacheid)
	return true
}
func (m *AssetModel) Exchange(uid int, from string, to string, to_amount float64) *BaseResponse {
	//币种兑换
	ntime := utils.GetNow()
	rs := new(BaseResponse)
	uinfo := MODEL_USER.GetBaseInfo(uid)
	if uinfo == nil {
		rs.State = STATE_FAILD
		rs.Msg = "no this user"
		return rs
	}
	if from == to {
		rs.State = STATE_FAILD
		rs.Msg = "same"
		return rs
	}
	from = strings.ToLower(from)
	to = strings.ToLower(to)

	from_coininfo := MODEL_SYSTEM.GetCoinInfo(from, from+"usdt")
	if from == "usdt" {
		from_coininfo = db.DB_ROW_RESULT{"cnum": "8"}
	}
	to_coininfo := MODEL_SYSTEM.GetCoinInfo(to, to+"usdt")
	if to == "usdt" {
		to_coininfo = db.DB_ROW_RESULT{"cnum": "8"}
	}
	if from_coininfo == nil || to_coininfo == nil {
		rs.State = STATE_SYSTEM_ERROR
		rs.Msg = "system error"
		return rs
	}

	to_amount = utils.FormatFloatA(to_amount, utils.GetInt(to_coininfo["cnum"]))
	if to_amount <= 0 {
		rs.State = EXCHANGE_STATE_TOOMIN
		rs.Msg = "too smalll"
		return rs
	}
	//from_priceInfo := config.GlobalMongo.GetOne("lastdata", bson.M{"pair": from_coininfo["pair"]}, nil)
	from_price := float64(1)
	to_price := float64(1)
	if from != "usdt" {
		from_priceInfo := MODEL_SYSTEM.GetLastCoinInfo(from_coininfo["pair"])
		from_price = from_priceInfo["close"].(float64)
	}
	if to != "usdt" {
		to_priceInfo := MODEL_SYSTEM.GetLastCoinInfo(to_coininfo["pair"])

		to_price = to_priceInfo["close"].(float64)
	}

	user_assets := m.GetAllAssets(uid, USER_MODE_REAL)
	from_amount, ok := user_assets[from]

	if !ok && from != "usdt" {
		rs.State = EXCHANGE_STATE_NOTENNOUGH
		rs.Msg = "assets not engough"
		return rs
	}
	to_allprice := to_price * to_amount
	from_allprice := 0.0
	if from == "usdt" {
		from_allprice = uinfo.Credit
	} else {
		if from_amount.IsTrans == 0 {
			rs.State = STATE_SYSTEM_ERROR
			rs.Msg = "from not allowed trans"
			return rs
		}
		from_allprice = from_price * from_amount.Count
	}

	if to_allprice > from_allprice {
		rs.State = EXCHANGE_STATE_NOTENNOUGH
		rs.Msg = "not enough"
		return rs
	}
	from_diff_amount := utils.FormatFloatA(to_allprice/from_price, utils.GetInt(from_coininfo["cnum"]))
	if from_diff_amount <= 0 {
		rs.State = EXCHANGE_STATE_TOOMIN
		rs.Msg = "too smalll"
		return rs
	}
	if from != "usdt" && to != "usdt" {
		if m.AddAssets(uid, &Assets{ //扣除原有币种的持有量
			Coin:    from,
			Pair:    from_coininfo["pair"],
			Num:     -1 * from_diff_amount,
			LockNum: 0,
			Price:   from_price,
			Mode:    USER_MODE_REAL,
		}) {
			MODEL_USER.AddCredit(uid, &CreditValue{
				Credit:          0,
				UserCoinLogType: COIN_LOG_ASSETS_EXCHANGE,
				UserCoinLogInfo: QueueCreditLog{
					Credit:     -1 * from_diff_amount,
					CoinType:   from,
					CreateTime: ntime,
				},
			})
			m.AddAssets(uid, &Assets{ //增加目标币种持有量
				Coin:    to,
				Pair:    to_coininfo["pair"],
				Num:     to_amount,
				LockNum: 0,
				Price:   to_price,
				Mode:    USER_MODE_REAL,
			})
			MODEL_USER.AddCredit(uid, &CreditValue{
				Credit:          0,
				UserCoinLogType: COIN_LOG_ASSETS_EXCHANGE,
				UserCoinLogInfo: QueueCreditLog{
					Credit:     to_amount,
					CoinType:   to,
					CreateTime: ntime,
				},
			})
		}
	} else {
		if from == "usdt" {
			if MODEL_USER.AddCredit(uid, &CreditValue{
				Credit:          -1 * to_allprice,
				LockCredit:      0,
				VCrdit:          0,
				LockVCredit:     0,
				UserCoinLogType: COIN_LOG_ASSETS_EXCHANGE,
				UserCoinLogInfo: QueueCreditLog{
					Credit:     -1 * to_allprice,
					CoinType:   from,
					CreateTime: ntime,
				},
			}) {
				m.AddAssets(uid, &Assets{
					Coin:    to,
					Pair:    to_coininfo["pair"],
					Num:     to_amount,
					LockNum: 0,
					Price:   to_price,
					Mode:    USER_MODE_REAL,
				})
				MODEL_USER.AddCredit(uid, &CreditValue{
					Credit:          0,
					LockCredit:      0,
					VCrdit:          0,
					LockVCredit:     0,
					UserCoinLogType: COIN_LOG_ASSETS_EXCHANGE,
					UserCoinLogInfo: QueueCreditLog{
						Credit:     to_amount,
						CoinType:   to,
						CreateTime: ntime,
					}})

			}
		} else {
			if m.AddAssets(uid, &Assets{ //扣除原有币种的持有量
				Coin:    from,
				Pair:    from_coininfo["pair"],
				Num:     -1 * from_diff_amount,
				LockNum: 0,
				Price:   from_price,
				Mode:    USER_MODE_REAL,
			}) {
				MODEL_USER.AddCredit(uid, &CreditValue{
					UserCoinLogType: COIN_LOG_ASSETS_EXCHANGE,
					UserCoinLogInfo: QueueCreditLog{
						Credit:     -1 * from_diff_amount,
						CoinType:   from,
						CreateTime: ntime,
					}})
				MODEL_USER.AddCredit(uid, &CreditValue{
					Credit:          1 * to_allprice,
					LockCredit:      0,
					VCrdit:          0,
					LockVCredit:     0,
					UserCoinLogType: COIN_LOG_ASSETS_EXCHANGE,
					UserCoinLogInfo: QueueCreditLog{
						Credit:     to_allprice,
						CoinType:   to,
						CreateTime: ntime,
					},
				})
			}
		}

	}

	rs.State = STATE_SUCCESS
	rs.Msg = "success!"
	return rs
}
