package models

import (
	"cointrade/config"
	"cointrade/lib/db"
	"cointrade/utils"
	"fmt"
	"math/rand"
	"time"
)

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
	rs := new(BaseResponse)
	uinfo := MODEL_USER.GetBaseInfo(uid)
	ntime := utils.GetNow()
	uCredit := uinfo.Credit
	if rq.GangGan < 1 {
		rq.GangGan = 1
	}
	if rq.GangGan > 100 {
		rs.State = DELEGATE_STATE_GANGGAN_ERROR
		rs.Msg = "error ganggan"
		return rs
	}
	teamLogType := 0
	var teamLogInfo QueueTeamLog
	if uinfo.Mode == USER_MODE_V {
		uCredit = uinfo.VCredit
	}
	if uinfo == nil {
		rs.State = STATE_SYSTEM_ERROR
		rs.Msg = "system error"
		return rs
	}
	if rq.OpenType != OPEN_TYPE_BB && rq.OpenType != OPEN_TYPE_EXPLODE && rq.OpenType != OPEN_TYPE_KEEP {
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
	if rq.DirectType != DIRECT_TYPE_BIG {
		rq.DirectType = DIRECT_TYPE_SMALL
	}
	if rq.GangGan < 1 {
		rq.GangGan = 1
	}
	if rq.PriceType != PRICE_TYPE_MARKET {
		rq.PriceType = PRICE_TYPE_LIMIT
	}
	if rq.DelegateType != DELEGATE_TYPE_BUY {
		rq.DelegateType = DELEGATE_TYPE_SELL
	}

	sn := m.MakeSn(uid, rq.DelegateType)
	coinPriceInfo := MODEL_SYSTEM.GetLastCoinInfo(rq.Pair)
	coinChangeType := COIN_LOG_USER_DELEGATE
	realprice := 0.0
	insertData := db.DB_PARAMS{
		"uid":                uid,
		"delegate_type":      rq.DelegateType,
		"trade_type":         rq.OpenType,
		"flag":               rq.DirectType,
		"coinid":             coinfo["id"],
		"coinpair":           coinfo["pair"],
		"coin_symbol":        coinfo["symbol"],
		"ganggan":            rq.GangGan,
		"state":              0,
		"mode":               uinfo.Mode,
		"createtime":         ntime,
		"sn":                 sn,
		"stop_up_price":      rq.StopUpPrice,
		"stop_down_price":    rq.StopDownPrice,
		"stop_up_delegate":   rq.StopUpDelegatePrice,
		"stop_down_delegate": rq.StopDownDelegatePrice,
	}
	coinprice := coinPriceInfo["close"].(float64)

	if rq.PriceType != PRICE_TYPE_MARKET && rq.OpenType != OPEN_TYPE_EXPLODE {
		coinprice = utils.FormatFloatA(rq.Price, utils.GetInt(coinfo["dnum"]))
		if coinprice <= 0 {
			rs.State = DELEGATE_STATE_MIN
			rs.Msg = "trade too small 1"
			return rs
		}
	}

	switch rq.OpenType {
	case OPEN_TYPE_BB:
		if coinfo["open_coin2coin"] == "0" {
			rs.State = DELEGATE_STATE_TRADE_CLOSED
			rs.Msg = "trade closed"
			return rs
		}

		coinChangeType = COIN_LOG_BB_TRADE
		rq.Amount = utils.FormatFloatA(rq.Amount, utils.GetInt(coinfo["cnum"]))
		if rq.Amount <= 0 {
			rs.State = DELEGATE_STATE_MIN
			rs.Msg = "trade too small 2"
			return rs
		}
		allprice := rq.Amount * coinprice
		realprice = allprice
		if rq.DelegateType == DELEGATE_TYPE_BUY {
			if coinfo["isnative"] == "0" && coinfo["f_price"] != "0" {
				insertData["is_f"] = 1
			}
			if uCredit < allprice {
				rs.State = DELEGATE_STATE_CREDIT
				rs.Msg = "not enough credit"
				return rs
			}
		} else {
			userAssets := MODEL_ASSETS.GetAllAssets(uid, uinfo.Mode)
			if _, ok := userAssets[rq.Coin]; !ok {
				rs.State = DELEGATE_STATE_NOASSET
				rs.Msg = "you dont have this assets"
				return rs
			}
			if userAssets[rq.Coin].Count < rq.Amount {
				rs.State = DELEGATE_STATE_CREDIT
				rs.Msg = "not enough credit"
				return rs
			}
			if coinfo["isnative"] == "0" && userAssets[rq.Coin].TransOpenTime > ntime {
				rs.State = DELEGATE_STATE_TRADE_CLOSED
				rs.Msg = "locktime"
				return rs
			}
			MODEL_ASSETS.AddAssets(uid, &Assets{
				Coin:    rq.Coin,
				Pair:    userAssets[rq.Coin].Pair,
				Num:     -1 * rq.Amount,
				LockNum: rq.Amount,
				Price:   coinprice,
				Mode:    uinfo.Mode,
			})
		}

		insertData["price"] = coinprice
		insertData["credit"] = allprice
		insertData["num"] = rq.Amount

	case OPEN_TYPE_KEEP:
		if coinfo["open_trade"] == "0" {
			rs.State = DELEGATE_STATE_TRADE_CLOSED
			rs.Msg = "trade closed"
			return rs
		}
		allprice := rq.Amount * 1000
		coinChangeType = COIN_LOG_KEEP_TRADE
		if rq.DelegateType == DELEGATE_TYPE_BUY {
			fee := allprice * (config.GlobalConfig.GetValue("trade_fee").ToFloat() / float64(100))
			insertData["fee"] = fee
			realPrice := allprice + fee
			realprice = realPrice
			if uCredit < realPrice {
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
			config.GlobalDB.AddValue(DB_TABLE_OPENED_TRADE, map[string]float64{"num": -1 * rq.Amount, "lock_num": rq.Amount}, db.DB_PARAMS{"id": one.Id, "mode": uinfo.Mode})
		}
		insertData["price"] = coinprice
		insertData["credit"] = allprice
		insertData["num"] = rq.Amount

	case OPEN_TYPE_EXPLODE:
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
		if uCredit < rq.Amount {
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
	lockPrice := realprice
	if rq.OpenType == OPEN_TYPE_EXPLODE {
		lockPrice = 0
		openData := db.DB_PARAMS{
			"uid":             uid,
			"sn":              sn,
			"trade_type":      OPEN_TYPE_EXPLODE,
			"flag":            rq.DirectType,
			"openprice":       coinprice,
			"closeprice":      0,
			"coinid":          coinfo["id"],
			"coinpair":        coinfo["pair"],
			"coin_symbol":     coinfo["symbol"],
			"close_time":      rq.CloseTime,
			"close_real_time": ntime + rq.CloseTime,
			"clear_time":      0,
			"createtime":      ntime,
			"ganggan":         1,
			"credit":          rq.Amount,
			"profit":          0,
			"num":             rq.Amount,
			"mode":            uinfo.Mode,
		}
		if explodeConfig, ok := EXPLODE_CONFIG[rq.CloseTime]; ok {
			openData["win_rate"] = explodeConfig.Winrate
			openData["lose_rate"] = explodeConfig.Loserate
		} else {
			openData["win_rate"] = 100
			openData["lose_rate"] = 100
		}
		teamLogType = TEAM_LOG_TRADE
		teamLogInfo.TradeExplode = rq.Amount
		teamLogInfo.CreateTime = ntime
		_, err = config.GlobalDB.InsertData(DB_TABLE_OPENED_TRADE, openData)
	} else {
		_, err = config.GlobalDB.InsertData(DB_TABLE_DELEGATE_TRADE, insertData)
	}
	if err != nil {
		rs.State = STATE_SYSTEM_ERROR
		rs.Msg = "system error"
		return rs
	}

	if rq.DelegateType == DELEGATE_TYPE_BUY {
		if uinfo.Mode == USER_MODE_REAL {
			MODEL_USER.AddCredit(uid, &CreditValue{
				Credit:          -1 * realprice,
				LockCredit:      lockPrice,
				VCrdit:          0,
				LockVCredit:     0,
				UserCoinLogType: coinChangeType,
				UserCoinLogInfo: QueueCreditLog{
					Credit:     -1 * realprice,
					LockCredit: lockPrice,
					Sn:         sn,
					CreateTime: ntime,
				},
				TeamCoinLogType: teamLogType,
				TeamCoinLogInfo: teamLogInfo,
			})
		} else {
			MODEL_USER.AddCredit(uid, &CreditValue{
				Credit:          0,
				LockCredit:      0,
				VCrdit:          -1 * realprice,
				LockVCredit:     lockPrice,
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
