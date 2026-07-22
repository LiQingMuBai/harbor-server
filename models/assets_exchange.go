package models

import (
	"cointrade/lib/db"
	"cointrade/utils"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func (m *AssetModel) MakeSn(uid int) string {
	uidstr := utils.Sup(int64(uid), 10)
	timestr := time.Now().Format("20060102150405")
	return fmt.Sprintf("%s%s%s%d", "E", timestr, uidstr, 10+rand.Intn(89))
}

func (m *AssetModel) QuickExchange(uid int, rq *QuickExchangeRequest) *BaseResponse {
	ntime := utils.GetNow()
	assetsInfo := m.GetAllAssets(uid, USER_MODE_REAL)
	rq.Coin = strings.ToLower(rq.Coin)
	coininfo := MODEL_SYSTEM.GetCoinInfo(rq.Coin, rq.Coin+"usdt")
	if coininfo != nil && coininfo["is_market"] == "0" && rq.Coin != "usdt" && rq.Coin != "usdc" {
		return &BaseResponse{State: STATE_FAILD, Msg: "faild"}
	}
	if rq.Coin == "usdt" {
		return &BaseResponse{State: STATE_FAILD, Msg: "faild"}
	}
	for _, v := range assetsInfo {
		if v.Symbol != rq.Coin {
			continue
		}
		if v.IsTrans == 0 {
			return &BaseResponse{State: EXCHANGE_STATE_NOT_TRANS, Msg: "not allowed trans"}
		}
		if v.Count < rq.Amount {
			return &BaseResponse{State: EXCHANGE_STATE_NOTENNOUGH, Msg: "not enough assets"}
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
			return &BaseResponse{State: EXCHANGE_STATE_NOT_TRANS, Msg: "not allowed trans"}
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
		return &BaseResponse{State: STATE_SUCCESS, Msg: "ok"}
	}
	return &BaseResponse{State: STATE_FAILD, Msg: "faild"}
}

func (m *AssetModel) Exchange(uid int, from string, to string, toAmount float64) *BaseResponse {
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

	fromCoinInfo := MODEL_SYSTEM.GetCoinInfo(from, from+"usdt")
	if from == "usdt" {
		fromCoinInfo = db.DB_ROW_RESULT{"cnum": "8"}
	}
	toCoinInfo := MODEL_SYSTEM.GetCoinInfo(to, to+"usdt")
	if to == "usdt" {
		toCoinInfo = db.DB_ROW_RESULT{"cnum": "8"}
	}
	if fromCoinInfo == nil || toCoinInfo == nil {
		rs.State = STATE_SYSTEM_ERROR
		rs.Msg = "system error"
		return rs
	}

	toAmount = utils.FormatFloatA(toAmount, utils.GetInt(toCoinInfo["cnum"]))
	if toAmount <= 0 {
		rs.State = EXCHANGE_STATE_TOOMIN
		rs.Msg = "too smalll"
		return rs
	}

	fromPrice := 1.0
	toPrice := 1.0
	if from != "usdt" {
		fromPriceInfo := MODEL_SYSTEM.GetLastCoinInfo(fromCoinInfo["pair"])
		fromPrice = fromPriceInfo["close"].(float64)
	}
	if to != "usdt" {
		toPriceInfo := MODEL_SYSTEM.GetLastCoinInfo(toCoinInfo["pair"])
		toPrice = toPriceInfo["close"].(float64)
	}

	userAssets := m.GetAllAssets(uid, USER_MODE_REAL)
	fromAmount, ok := userAssets[from]
	if !ok && from != "usdt" {
		rs.State = EXCHANGE_STATE_NOTENNOUGH
		rs.Msg = "assets not engough"
		return rs
	}

	toAllPrice := toPrice * toAmount
	fromAllPrice := 0.0
	if from == "usdt" {
		fromAllPrice = uinfo.Credit
	} else {
		if fromAmount.IsTrans == 0 {
			rs.State = STATE_SYSTEM_ERROR
			rs.Msg = "from not allowed trans"
			return rs
		}
		fromAllPrice = fromPrice * fromAmount.Count
	}

	if toAllPrice > fromAllPrice {
		rs.State = EXCHANGE_STATE_NOTENNOUGH
		rs.Msg = "not enough"
		return rs
	}
	fromDiffAmount := utils.FormatFloatA(toAllPrice/fromPrice, utils.GetInt(fromCoinInfo["cnum"]))
	if fromDiffAmount <= 0 {
		rs.State = EXCHANGE_STATE_TOOMIN
		rs.Msg = "too smalll"
		return rs
	}

	if from != "usdt" && to != "usdt" {
		if m.AddAssets(uid, &Assets{
			Coin:    from,
			Pair:    fromCoinInfo["pair"],
			Num:     -1 * fromDiffAmount,
			LockNum: 0,
			Price:   fromPrice,
			Mode:    USER_MODE_REAL,
		}) {
			MODEL_USER.AddCredit(uid, &CreditValue{
				Credit:          0,
				UserCoinLogType: COIN_LOG_ASSETS_EXCHANGE,
				UserCoinLogInfo: QueueCreditLog{
					Credit:     -1 * fromDiffAmount,
					CoinType:   from,
					CreateTime: ntime,
				},
			})
			m.AddAssets(uid, &Assets{
				Coin:    to,
				Pair:    toCoinInfo["pair"],
				Num:     toAmount,
				LockNum: 0,
				Price:   toPrice,
				Mode:    USER_MODE_REAL,
			})
			MODEL_USER.AddCredit(uid, &CreditValue{
				Credit:          0,
				UserCoinLogType: COIN_LOG_ASSETS_EXCHANGE,
				UserCoinLogInfo: QueueCreditLog{
					Credit:     toAmount,
					CoinType:   to,
					CreateTime: ntime,
				},
			})
		}
	} else if from == "usdt" {
		if MODEL_USER.AddCredit(uid, &CreditValue{
			Credit:          -1 * toAllPrice,
			LockCredit:      0,
			VCrdit:          0,
			LockVCredit:     0,
			UserCoinLogType: COIN_LOG_ASSETS_EXCHANGE,
			UserCoinLogInfo: QueueCreditLog{
				Credit:     -1 * toAllPrice,
				CoinType:   from,
				CreateTime: ntime,
			},
		}) {
			m.AddAssets(uid, &Assets{
				Coin:    to,
				Pair:    toCoinInfo["pair"],
				Num:     toAmount,
				LockNum: 0,
				Price:   toPrice,
				Mode:    USER_MODE_REAL,
			})
			MODEL_USER.AddCredit(uid, &CreditValue{
				Credit:          0,
				LockCredit:      0,
				VCrdit:          0,
				LockVCredit:     0,
				UserCoinLogType: COIN_LOG_ASSETS_EXCHANGE,
				UserCoinLogInfo: QueueCreditLog{
					Credit:     toAmount,
					CoinType:   to,
					CreateTime: ntime,
				},
			})
		}
	} else {
		if m.AddAssets(uid, &Assets{
			Coin:    from,
			Pair:    fromCoinInfo["pair"],
			Num:     -1 * fromDiffAmount,
			LockNum: 0,
			Price:   fromPrice,
			Mode:    USER_MODE_REAL,
		}) {
			MODEL_USER.AddCredit(uid, &CreditValue{
				UserCoinLogType: COIN_LOG_ASSETS_EXCHANGE,
				UserCoinLogInfo: QueueCreditLog{
					Credit:     -1 * fromDiffAmount,
					CoinType:   from,
					CreateTime: ntime,
				},
			})
			MODEL_USER.AddCredit(uid, &CreditValue{
				Credit:          toAllPrice,
				LockCredit:      0,
				VCrdit:          0,
				LockVCredit:     0,
				UserCoinLogType: COIN_LOG_ASSETS_EXCHANGE,
				UserCoinLogInfo: QueueCreditLog{
					Credit:     toAllPrice,
					CoinType:   to,
					CreateTime: ntime,
				},
			})
		}
	}

	rs.State = STATE_SUCCESS
	rs.Msg = "success!"
	return rs
}
