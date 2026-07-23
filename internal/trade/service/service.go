package service

import (
	"cointrade/config"
	assetsdomain "cointrade/internal/domain/assets"
	creditlogdomain "cointrade/internal/domain/creditlog"
	shareddomain "cointrade/internal/domain/shared"
	systemdomain "cointrade/internal/domain/system"
	tradedomain "cointrade/internal/domain/trade"
	userdomain "cointrade/internal/domain/user"
	traderepo "cointrade/internal/trade/repo"
	"cointrade/lib/db"
	"cointrade/utils"
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	stateSuccess     = 0
	stateFailed      = 1
	stateSystemError = 9999999
)

type UserGateway interface {
	GetBaseInfo(uid int) *userdomain.UserBaseInfo
	AddCredit(uid int, value *userdomain.CreditValue) bool
	GetExplodeState(uid int) int
}

type AssetGateway interface {
	GetAllAssets(uid int, mode int) map[string]assetsdomain.AssetInfo
	AddAssets(uid int, asset *assetsdomain.Assets) bool
}

type MarketGateway interface {
	GetCoinInfo(coin string, pair string) db.DB_ROW_RESULT
	GetLastCoinInfo(pair string) primitive.M
	GetExplodeConfig() map[int]*systemdomain.ExplodeConfig
	GetControlExplode(sn string) int32
}

type Service struct {
	repo   traderepo.Repository
	user   UserGateway
	asset  AssetGateway
	market MarketGateway
}

func New(repo traderepo.Repository, user UserGateway, asset AssetGateway, market MarketGateway) *Service {
	return &Service{
		repo:   repo,
		user:   user,
		asset:  asset,
		market: market,
	}
}

func (s *Service) MakeSn(uid int, delegateType int) string {
	uidstr := utils.Sup(int64(uid), 10)
	timestr := time.Now().Format("20060102150405")
	switch delegateType {
	case tradedomain.DELEGATE_TYPE_BUY:
		return fmt.Sprintf("%s%s%s%d", tradedomain.TRADE_BUY_PREFIX, timestr, uidstr, 10+rand.Intn(89))
	case tradedomain.DELEGATE_TYPE_SELL:
		return fmt.Sprintf("%s%s%s%d", tradedomain.TRADE_SELL_PREFIX, timestr, uidstr, 10+rand.Intn(89))
	}
	return fmt.Sprintf("%s%s%s%d", tradedomain.TRADE_SELL_PREFIX, timestr, uidstr, 10+rand.Intn(89))
}

func (s *Service) GetCloseBySn(uid int, sn string) db.DB_ROW_RESULT {
	one, _ := s.repo.FetchCloseBySN(uid, sn)
	return one
}

func (s *Service) GetOpenedBySn(uid int, sn string) *tradedomain.OpenedInfo {
	one, _ := s.repo.FetchOpenedBySN(uid, sn)
	return buildOpenedInfo(uid, one)
}

func (s *Service) GetOpenedOne(uid int, coin string, tradeType int, flag int, mode int, ganggan int) *tradedomain.OpenedInfo {
	one, _ := s.repo.FetchOpened(uid, coin, tradeType, flag, mode, ganggan)
	return buildOpenedInfo(uid, one)
}

func (s *Service) AddKeepOpened(delegateInfo db.DBValues) {
	ntime := utils.GetNow()
	var openinfo *tradedomain.OpenedInfo
	if delegateInfo["ganggan"].ToInt() > 1 {
		openinfo = nil
	} else {
		openinfo = s.GetOpenedOne(
			delegateInfo["uid"].ToInt(),
			delegateInfo["coin_symbol"].ToString(),
			delegateInfo["trade_type"].ToInt(),
			delegateInfo["flag"].ToInt(),
			delegateInfo["mode"].ToInt(),
			delegateInfo["ganggan"].ToInt(),
		)
	}

	if openinfo == nil {
		insertData := db.DB_PARAMS{
			"uid":                delegateInfo["uid"].ToInt(),
			"trade_type":         delegateInfo["trade_type"].ToInt(),
			"closeprice":         0,
			"flag":               delegateInfo["flag"].ToInt(),
			"openprice":          delegateInfo["price"].ToFloat(),
			"coinid":             delegateInfo["coinid"].ToInt(),
			"coinpair":           delegateInfo["coinpair"].ToString(),
			"coin_symbol":        delegateInfo["coin_symbol"].ToString(),
			"close_time":         0,
			"close_real_time":    0,
			"clear_time":         0,
			"createtime":         ntime,
			"ganggan":            delegateInfo["ganggan"].ToInt(),
			"credit":             delegateInfo["credit"].ToFloat(),
			"profit":             0,
			"win_rate":           0,
			"lose_rate":          0,
			"num":                delegateInfo["num"].ToFloat(),
			"mode":               delegateInfo["mode"].ToInt(),
			"sn":                 delegateInfo["sn"].ToString(),
			"stop_up_price":      delegateInfo["stop_up_price"].ToFloat(),
			"stop_down_price":    delegateInfo["stop_down_price"].ToFloat(),
			"stop_up_delegate":   delegateInfo["stop_up_delegate"].ToFloat(),
			"stop_down_delegate": delegateInfo["stop_down_delegate"].ToFloat(),
		}
		if err := s.repo.InsertOpened(insertData); err == nil {
			_ = s.repo.UpdateDelegateState(delegateInfo["id"].Value, 1, ntime)
		}
		return
	}

	oldAllPrice := openinfo.OpenPrice * openinfo.Num
	newAllPrice := delegateInfo["price"].ToFloat() * delegateInfo["num"].ToFloat()
	openPrice := (oldAllPrice + newAllPrice) / (delegateInfo["num"].ToFloat() + openinfo.Num)
	_ = s.repo.AddOpenedPositionValue(openinfo.Id, map[string]float64{
		"num":    delegateInfo["num"].ToFloat(),
		"credit": delegateInfo["credit"].ToFloat(),
	})
	_ = s.repo.UpdateOpenedPrice(openinfo.Id, openPrice)
	_ = s.repo.UpdateDelegateState(delegateInfo["id"].Value, 1, ntime)
}

func (s *Service) GetDelegateList(uid int, rq *tradedomain.TradeListRequest) *shareddomain.PageBaseResponse {
	uinfo := s.user.GetBaseInfo(uid)
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
	count := s.repo.CountDelegate(condition)
	page, pageTotal := normalizePage(rq.Page, rq.Limit, count)
	list, _ := s.repo.FetchDelegateRows(condition, "order by id desc", fmt.Sprintf("limit %d,%d", (page-1)*rq.Limit, rq.Limit))
	return &shareddomain.PageBaseResponse{
		BaseResponse: shareddomain.BaseResponse{State: stateSuccess, Msg: shareddomain.MsgOK},
		Limit:        rq.Limit,
		Page:         page,
		PageTotal:    pageTotal,
		List:         list,
	}
}

func (s *Service) GetOpenedList(uid int, rq *tradedomain.TradeListRequest) *shareddomain.PageBaseResponse {
	uinfo := s.user.GetBaseInfo(uid)
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
	count := s.repo.CountOpened(condition)
	page, pageTotal := normalizePage(rq.Page, rq.Limit, count)
	list, _ := s.repo.FetchOpenedRows(condition, "order by id desc", fmt.Sprintf("limit %d,%d", (page-1)*rq.Limit, rq.Limit))
	return &shareddomain.PageBaseResponse{
		BaseResponse: shareddomain.BaseResponse{State: stateSuccess, Msg: fmt.Sprintf("%d", utils.GetNow())},
		Limit:        rq.Limit,
		Page:         page,
		PageTotal:    pageTotal,
		List:         list,
	}
}

func (s *Service) GetCloseList(uid int, rq *tradedomain.TradeListRequest) *shareddomain.PageBaseResponse {
	uinfo := s.user.GetBaseInfo(uid)
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
	count := s.repo.CountClose(condition)
	page, pageTotal := normalizePage(rq.Page, rq.Limit, count)
	list, _ := s.repo.FetchCloseRows(condition, "order by id desc", fmt.Sprintf("limit %d,%d", (page-1)*rq.Limit, rq.Limit))
	return &shareddomain.PageBaseResponse{
		BaseResponse: shareddomain.BaseResponse{State: stateSuccess, Msg: shareddomain.MsgOK},
		Limit:        rq.Limit,
		Page:         page,
		PageTotal:    pageTotal,
		List:         list,
	}
}

func (s *Service) DelegateTrade(uid int, rq *tradedomain.TradeDelegateRequest) *shareddomain.BaseResponse {
	rs := new(shareddomain.BaseResponse)
	uinfo := s.user.GetBaseInfo(uid)
	if uinfo == nil {
		rs.State = stateSystemError
		rs.Msg = shareddomain.MsgInternalError
		return rs
	}
	ntime := utils.GetNow()
	userCredit := uinfo.Credit
	if rq.GangGan < 1 {
		rq.GangGan = 1
	}
	if rq.GangGan > 100 {
		rs.State = tradedomain.DELEGATE_STATE_GANGGAN_ERROR
		rs.Msg = shareddomain.MsgLeverageInvalid
		return rs
	}
	teamLogType := 0
	var teamLogInfo creditlogdomain.QueueTeamLog
	if uinfo.Mode == 2 {
		userCredit = uinfo.VCredit
	}
	if rq.OpenType != tradedomain.OPEN_TYPE_BB && rq.OpenType != tradedomain.OPEN_TYPE_EXPLODE && rq.OpenType != tradedomain.OPEN_TYPE_KEEP {
		rs.State = stateSystemError
		rs.Msg = shareddomain.MsgInternalError
		return rs
	}

	coinInfo := s.market.GetCoinInfo(rq.Coin, rq.Pair)
	if coinInfo == nil {
		rs.State = tradedomain.DELEGATE_STATE_NOCOIN
		rs.Msg = shareddomain.MsgCoinNotFound
		return rs
	}
	if rq.DirectType != tradedomain.DIRECT_TYPE_BIG {
		rq.DirectType = tradedomain.DIRECT_TYPE_SMALL
	}
	if rq.PriceType != tradedomain.PRICE_TYPE_MARKET {
		rq.PriceType = tradedomain.PRICE_TYPE_LIMIT
	}
	if rq.DelegateType != tradedomain.DELEGATE_TYPE_BUY {
		rq.DelegateType = tradedomain.DELEGATE_TYPE_SELL
	}

	sn := s.MakeSn(uid, rq.DelegateType)
	coinPriceInfo := s.market.GetLastCoinInfo(rq.Pair)
	coinChangeType := creditlogdomain.COIN_LOG_USER_DELEGATE
	realprice := 0.0
	insertData := db.DB_PARAMS{
		"uid":                uid,
		"delegate_type":      rq.DelegateType,
		"trade_type":         rq.OpenType,
		"flag":               rq.DirectType,
		"coinid":             coinInfo["id"],
		"coinpair":           coinInfo["pair"],
		"coin_symbol":        coinInfo["symbol"],
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
	coinprice := utils.GetFloat(utils.GetJsonValue(coinPriceInfo["close"]))

	if rq.PriceType != tradedomain.PRICE_TYPE_MARKET && rq.OpenType != tradedomain.OPEN_TYPE_EXPLODE {
		coinprice = utils.FormatFloatA(rq.Price, utils.GetInt(coinInfo["dnum"]))
		if coinprice <= 0 {
			rs.State = tradedomain.DELEGATE_STATE_MIN
			rs.Msg = shareddomain.MsgTradeTooSmall
			return rs
		}
	}

	switch rq.OpenType {
	case tradedomain.OPEN_TYPE_BB:
		if coinInfo["open_coin2coin"] == "0" {
			rs.State = tradedomain.DELEGATE_STATE_TRADE_CLOSED
			rs.Msg = shareddomain.MsgTradeClosed
			return rs
		}
		coinChangeType = creditlogdomain.COIN_LOG_BB_TRADE
		rq.Amount = utils.FormatFloatA(rq.Amount, utils.GetInt(coinInfo["cnum"]))
		if rq.Amount <= 0 {
			rs.State = tradedomain.DELEGATE_STATE_MIN
			rs.Msg = shareddomain.MsgTradeTooSmall
			return rs
		}
		allprice := rq.Amount * coinprice
		realprice = allprice
		if rq.DelegateType == tradedomain.DELEGATE_TYPE_BUY {
			if coinInfo["isnative"] == "0" && coinInfo["f_price"] != "0" {
				insertData["is_f"] = 1
			}
			if userCredit < allprice {
				rs.State = tradedomain.DELEGATE_STATE_CREDIT
				rs.Msg = shareddomain.MsgInsufficient
				return rs
			}
		} else {
			userAssets := s.asset.GetAllAssets(uid, uinfo.Mode)
			assetInfo, ok := userAssets[rq.Coin]
			if !ok {
				rs.State = tradedomain.DELEGATE_STATE_NOASSET
				rs.Msg = shareddomain.MsgAssetNotOwned
				return rs
			}
			if assetInfo.Count < rq.Amount {
				rs.State = tradedomain.DELEGATE_STATE_CREDIT
				rs.Msg = shareddomain.MsgInsufficient
				return rs
			}
			if coinInfo["isnative"] == "0" && assetInfo.TransOpenTime > ntime {
				rs.State = tradedomain.DELEGATE_STATE_TRADE_CLOSED
				rs.Msg = shareddomain.MsgAssetLockedByTime
				return rs
			}
			s.asset.AddAssets(uid, &assetsdomain.Assets{
				Coin:    rq.Coin,
				Pair:    assetInfo.Pair,
				Num:     -1 * rq.Amount,
				LockNum: rq.Amount,
				Price:   coinprice,
				Mode:    uinfo.Mode,
			})
		}
		insertData["price"] = coinprice
		insertData["credit"] = allprice
		insertData["num"] = rq.Amount

	case tradedomain.OPEN_TYPE_KEEP:
		if coinInfo["open_trade"] == "0" {
			rs.State = tradedomain.DELEGATE_STATE_TRADE_CLOSED
			rs.Msg = shareddomain.MsgTradeClosed
			return rs
		}
		allprice := rq.Amount * 1000
		coinChangeType = creditlogdomain.COIN_LOG_KEEP_TRADE
		if rq.DelegateType == tradedomain.DELEGATE_TYPE_BUY {
			fee := allprice * (config.GlobalConfig.GetValue("trade_fee").ToFloat() / float64(100))
			insertData["fee"] = fee
			realPrice := allprice + fee
			realprice = realPrice
			if userCredit < realPrice {
				rs.State = tradedomain.DELEGATE_STATE_CREDIT
				rs.Msg = shareddomain.MsgInsufficient
				return rs
			}
			coinCount := allprice / coinprice
			coinCount = utils.FormatFloatA(coinCount, utils.GetInt(coinInfo["cnum"]))
			if coinCount <= 0 {
				rs.State = tradedomain.DELEGATE_STATE_MIN
				rs.Msg = shareddomain.MsgTradeTooSmall
				return rs
			}
		} else {
			var one *tradedomain.OpenedInfo
			if rq.Sn != "" {
				one = s.GetOpenedBySn(uid, rq.Sn)
				if one == nil {
					rs.State = tradedomain.DELEGATE_STATE_NOCOIN
					rs.Msg = shareddomain.MsgAssetNotOwned
					return rs
				}
				if one.Ganggan <= 1 {
					rs.State = stateFailed
					rs.Msg = shareddomain.MsgLeverageOrderOnly
					return rs
				}
				insertData["ganggan_sn"] = rq.Sn
			} else {
				one = s.GetOpenedOne(uid, rq.Coin, rq.OpenType, rq.DirectType, uinfo.Mode, rq.GangGan)
			}
			if one == nil || rq.Amount > one.Num {
				rs.State = tradedomain.DELEGATE_STATE_NOCOIN
				rs.Msg = shareddomain.MsgAssetNotOwned
				return rs
			}
			_ = s.repo.AddOpenedLockValue(one.Id, -1*rq.Amount, rq.Amount, uinfo.Mode)
		}
		insertData["price"] = coinprice
		insertData["credit"] = allprice
		insertData["num"] = rq.Amount

	case tradedomain.OPEN_TYPE_EXPLODE:
		if coinInfo["open_trade"] == "0" {
			rs.State = tradedomain.DELEGATE_STATE_TRADE_CLOSED
			rs.Msg = shareddomain.MsgTradeClosed
			return rs
		}
		coinChangeType = creditlogdomain.COIN_LOG_EXPLODE_TRADE
		explodeConfig, ok := s.market.GetExplodeConfig()[rq.CloseTime]
		if !ok {
			rs.State = tradedomain.DELEGATE_STATE_CLOSETIME
			rs.Msg = shareddomain.MsgCloseTimeInvalid
			return rs
		}
		if rq.Amount < explodeConfig.Minprice {
			rs.State = tradedomain.DELEGATE_STATE_MIN
			rs.Msg = shareddomain.MsgTradeTooSmall
			return rs
		}
		if userCredit < rq.Amount {
			rs.State = tradedomain.DELEGATE_STATE_CREDIT
			rs.Msg = shareddomain.MsgInsufficient
			return rs
		}
		realprice = rq.Amount
		insertData["price"] = coinprice
		insertData["credit"] = rq.Amount
		insertData["num"] = rq.Amount
		insertData["close_time"] = rq.CloseTime
	}

	lockPrice := realprice
	var err error
	if rq.OpenType == tradedomain.OPEN_TYPE_EXPLODE {
		lockPrice = 0
		openData := db.DB_PARAMS{
			"uid":             uid,
			"sn":              sn,
			"trade_type":      tradedomain.OPEN_TYPE_EXPLODE,
			"flag":            rq.DirectType,
			"openprice":       coinprice,
			"closeprice":      0,
			"coinid":          coinInfo["id"],
			"coinpair":        coinInfo["pair"],
			"coin_symbol":     coinInfo["symbol"],
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
		if explodeConfig, ok := s.market.GetExplodeConfig()[rq.CloseTime]; ok {
			openData["win_rate"] = explodeConfig.Winrate
			openData["lose_rate"] = explodeConfig.Loserate
		} else {
			openData["win_rate"] = 100
			openData["lose_rate"] = 100
		}
		teamLogType = creditlogdomain.TEAM_LOG_TRADE
		teamLogInfo.TradeExplode = rq.Amount
		teamLogInfo.CreateTime = ntime
		err = s.repo.InsertOpened(openData)
	} else {
		err = s.repo.InsertDelegate(insertData)
	}
	if err != nil {
		rs.State = stateSystemError
		rs.Msg = shareddomain.MsgInternalError
		return rs
	}

	if rq.DelegateType == tradedomain.DELEGATE_TYPE_BUY {
		if uinfo.Mode == 1 {
			s.user.AddCredit(uid, &userdomain.CreditValue{
				Credit:          -1 * realprice,
				LockCredit:      lockPrice,
				VCrdit:          0,
				LockVCredit:     0,
				UserCoinLogType: coinChangeType,
				UserCoinLogInfo: creditlogdomain.QueueCreditLog{
					Credit:     -1 * realprice,
					LockCredit: lockPrice,
					Sn:         sn,
					CreateTime: ntime,
				},
				TeamCoinLogType: teamLogType,
				TeamCoinLogInfo: teamLogInfo,
			})
		} else {
			s.user.AddCredit(uid, &userdomain.CreditValue{
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

	rs.State = stateSuccess
	rs.Msg = utils.GetJsonValue(insertData)
	return rs
}

func (s *Service) CancleDelegate(uid int, sn string) *shareddomain.BaseResponse {
	rs := new(shareddomain.BaseResponse)
	one, _ := s.repo.FetchPendingDelegate(uid, sn)
	if one == nil {
		rs.State = stateFailed
		rs.Msg = shareddomain.MsgOrderNotFound
		return rs
	}

	userCredit := one["credit"].ToFloat() + one["fee"].ToFloat()
	userVCredit := one["credit"].ToFloat() + one["fee"].ToFloat()
	if one["mode"].ToInt() == 1 {
		userVCredit = 0
	} else {
		userCredit = 0
	}

	if one["delegate_type"].ToInt() == tradedomain.DELEGATE_TYPE_BUY {
		if s.user.AddCredit(uid, &userdomain.CreditValue{
			Credit:          userCredit,
			LockCredit:      -1 * userCredit,
			VCrdit:          userVCredit,
			LockVCredit:     -1 * userVCredit,
			UserCoinLogType: creditlogdomain.COIN_LOG_USER_CANCLE,
			UserCoinLogInfo: creditlogdomain.QueueCreditLog{
				Credit:     userCredit,
				LockCredit: -1 * userCredit,
				Sn:         one["sn"].ToString(),
				CreateTime: utils.GetNow(),
			},
			TeamCoinLogType: 0,
			TeamCoinLogInfo: nil,
		}) {
			if err := s.repo.UpdateDelegateState(one["id"].Value, 2, 0); err != nil {
				rs.State = stateSystemError
				rs.Msg = err.Error()
				return rs
			}
		}
		rs.State = stateSuccess
		rs.Msg = shareddomain.MsgSuccess
		return rs
	}

	switch one["trade_type"].ToInt() {
	case tradedomain.OPEN_TYPE_BB:
		if s.asset.AddAssets(uid, &assetsdomain.Assets{
			Coin:    one["coin_symbol"].ToString(),
			Pair:    one["coinpair"].ToString(),
			Num:     one["num"].ToFloat(),
			LockNum: -1 * one["num"].ToFloat(),
			Price:   one["price"].ToFloat(),
			Mode:    one["mode"].ToInt(),
		}) {
			if err := s.repo.UpdateDelegateState(one["id"].Value, 2, 0); err != nil {
				rs.State = stateSystemError
				rs.Msg = err.Error()
				return rs
			}
		}
	case tradedomain.OPEN_TYPE_KEEP:
		flag := one["flag"].ToInt()
		coin := one["coin_symbol"].ToString()
		uinfo := s.user.GetBaseInfo(uid)
		if uinfo == nil {
			rs.State = stateSystemError
			rs.Msg = shareddomain.MsgInternalError
			return rs
		}
		num := one["num"].ToFloat()
		opendInfo := s.GetOpenedOne(uid, coin, tradedomain.OPEN_TYPE_KEEP, flag, uinfo.Mode, one["ganggan"].ToInt())
		if opendInfo != nil {
			_ = s.repo.AddOpenedPositionValue(opendInfo.Id, map[string]float64{"num": num, "lock_num": -1 * num})
			if err := s.repo.UpdateDelegateState(one["id"].Value, 2, 0); err != nil {
				rs.State = stateSystemError
				rs.Msg = err.Error()
				return rs
			}
		}
	}

	rs.State = stateSuccess
	rs.Msg = shareddomain.MsgSuccess
	return rs
}

func (s *Service) OperateDelegateTrade(one db.DBValues) bool {
	ntime := utils.GetNow()
	uid := one["uid"].ToInt()
	if ntime-one["createtime"].ToInt() >= 7*24*60*60 {
		s.CancleDelegate(uid, one["sn"].ToString())
		return false
	}

	tradeType := one["trade_type"].ToInt()
	mode := one["mode"].ToInt()
	coin := one["coin_symbol"].ToString()
	coinPriceInfo := s.market.GetLastCoinInfo(one["coinpair"].ToString())
	if coinPriceInfo["close"] == nil {
		return false
	}
	coinprice, ok := coinPriceInfo["close"].(float64)
	if !ok {
		return false
	}

	virtualCredit := 0.0
	credit := 0.0
	lockVirtualCredit := 0.0
	lockCredit := 0.0
	if mode == 1 {
		credit = one["credit"].ToFloat()
		lockCredit = one["credit"].ToFloat() + one["fee"].ToFloat()
	} else {
		virtualCredit = one["credit"].ToFloat()
		lockVirtualCredit = one["credit"].ToFloat() + one["fee"].ToFloat()
	}

	switch tradeType {
	case tradedomain.OPEN_TYPE_BB:
		return s.operateBBDelegate(one, coin, coinprice, credit, virtualCredit, lockCredit, lockVirtualCredit, ntime)
	case tradedomain.OPEN_TYPE_EXPLODE:
		return s.operateExplodeDelegate(one, coinprice, lockCredit, lockVirtualCredit, ntime)
	case tradedomain.OPEN_TYPE_KEEP:
		return s.operateKeepDelegate(one, coin, coinprice, credit, virtualCredit, lockCredit, lockVirtualCredit, ntime)
	}
	return false
}

func (s *Service) EqualsPrice(delegatePrice float64, coinprice float64, isbig bool) bool {
	if isbig {
		return delegatePrice <= coinprice
	}
	return delegatePrice >= coinprice
}

func (s *Service) CheckStop(v db.DBValues, coinprice float64, isbig bool) bool {
	if v["stop_up_price"].ToFloat() > 0 && s.EqualsPrice(v["stop_up_price"].ToFloat(), coinprice, isbig) {
		delegatePrice := v["stop_up_delegate"].ToFloat()
		if delegatePrice == 0 {
			delegatePrice = coinprice
		}
		if v["stop_up_delegate"].ToFloat() > 0 {
			rs := s.DelegateTrade(v["uid"].ToInt(), &tradedomain.TradeDelegateRequest{
				OpenType:     tradedomain.OPEN_TYPE_KEEP,
				DelegateType: tradedomain.DELEGATE_TYPE_SELL,
				Pair:         v["coinpair"].ToString(),
				Coin:         v["coin_symbol"].ToString(),
				PriceType:    tradedomain.PRICE_TYPE_LIMIT,
				GangGan:      v["ganggan"].ToInt(),
				Amount:       v["num"].ToFloat(),
				Price:        delegatePrice,
				Sn:           v["sn"].ToString(),
			})
			if rs != nil && rs.State == stateSuccess {
				_ = s.repo.UpdateOpenedFields(v["id"].ToInt(), db.DB_PARAMS{"auto_delegate": 1})
			}
			return true
		}
	}

	if v["stop_down_price"].ToFloat() > 0 && s.EqualsPrice(v["stop_up_price"].ToFloat(), coinprice, !isbig) {
		delegatePrice := v["stop_down_delegate"].ToFloat()
		if delegatePrice == 0 {
			delegatePrice = coinprice
		}
		if v["stop_down_delegate"].ToFloat() > 0 {
			rs := s.DelegateTrade(v["uid"].ToInt(), &tradedomain.TradeDelegateRequest{
				OpenType:     tradedomain.OPEN_TYPE_KEEP,
				DelegateType: tradedomain.DELEGATE_TYPE_SELL,
				Pair:         v["coinpair"].ToString(),
				Coin:         v["coin_symbol"].ToString(),
				PriceType:    tradedomain.PRICE_TYPE_LIMIT,
				GangGan:      v["ganggan"].ToInt(),
				Amount:       v["num"].ToFloat(),
				Price:        delegatePrice,
				Sn:           v["sn"].ToString(),
			})
			if rs != nil && rs.State == stateSuccess {
				_ = s.repo.UpdateOpenedFields(v["id"].ToInt(), db.DB_PARAMS{"auto_delegate": 1})
			}
			return true
		}
	}
	return false
}

func (s *Service) SettleKeepCross(v db.DBValues, coinprice float64, ntime int) bool {
	profit := 0.0
	if v["flag"].ToInt() == tradedomain.DIRECT_TYPE_BIG {
		if s.CheckStop(v, coinprice, true) {
			return false
		}
		profit = ((coinprice - v["openprice"].ToFloat()) / v["openprice"].ToFloat()) * v["num"].ToFloat() * 1000
	} else {
		if s.CheckStop(v, coinprice, false) {
			return false
		}
		profit = ((v["num"].ToFloat()*1000)/coinprice)*v["openprice"].ToFloat() - v["num"].ToFloat()*1000
	}
	profit = profit * v["ganggan"].ToFloat()
	if profit >= 0 || (math.Abs(profit)/(v["num"].ToFloat()*1000)) < 1 {
		return false
	}

	allRealCredit := profit + v["num"].ToFloat()*1000
	insertData := db.DB_PARAMS{
		"uid":         v["uid"].ToInt(),
		"sn":          v["sn"].ToString(),
		"coin_symbol": v["coin_symbol"].ToString(),
		"trade_type":  v["trade_type"].ToInt(),
		"flag":        v["flag"].ToInt(),
		"amount":      v["num"].ToFloat(),
		"close_price": coinprice,
		"createtime":  ntime,
		"num":         v["num"].ToFloat(),
		"mode":        1,
		"allprice":    allRealCredit,
		"profit":      profit,
		"o_price":     v["openprice"].ToFloat(),
	}
	if err := s.repo.InsertClose(insertData); err != nil {
		return false
	}
	_ = s.repo.AddOpenedPositionValue(v["id"].ToInt(), map[string]float64{
		"num":    -1 * v["num"].ToFloat(),
		"credit": -1 * v["num"].ToFloat() * 1000,
	})
	s.user.AddCredit(v["uid"].ToInt(), &userdomain.CreditValue{
		Credit:          allRealCredit,
		LockCredit:      0,
		VCrdit:          0,
		LockVCredit:     0,
		UserCoinLogType: creditlogdomain.COIN_LOG_KEEP_BREAK,
		UserCoinLogInfo: creditlogdomain.QueueCreditLog{
			Credit:     allRealCredit,
			Sn:         v["sn"].ToString(),
			CreateTime: ntime,
		},
		TeamCoinLogType: creditlogdomain.TEAM_LOG_TRADE_PROFIT,
		TeamCoinLogInfo: creditlogdomain.QueueTeamLog{
			TradeKeep_Profit: profit,
			CreateTime:       ntime,
		},
	})
	return true
}

func (s *Service) SettleExplodeTrade(v db.DBValues, ntime int, controllerCollection string) bool {
	coinPriceInfo := s.market.GetLastCoinInfo(v["coinpair"].ToString())
	coinprice, ok := coinPriceInfo["close"].(float64)
	if !ok {
		return false
	}
	coinInfo := s.market.GetCoinInfo(v["coin_symbol"].ToString(), v["coinpair"].ToString())
	if coinInfo == nil {
		return false
	}

	iswin := (v["flag"].ToInt() == tradedomain.DIRECT_TYPE_BIG && coinprice > v["openprice"].ToFloat()) ||
		(v["flag"].ToInt() == tradedomain.DIRECT_TYPE_SMALL && coinprice < v["openprice"].ToFloat())
	control := s.resolveExplodeControl(v)
	diffPrice := s.resolveExplodeDiffPrice(coinInfo)
	if control == 1 {
		iswin = true
		if v["flag"].ToInt() == tradedomain.DIRECT_TYPE_BIG {
			coinprice = v["openprice"].ToFloat() + diffPrice
		} else {
			coinprice = v["openprice"].ToFloat() - diffPrice
		}
	} else if control == 2 {
		iswin = false
		if v["flag"].ToInt() == tradedomain.DIRECT_TYPE_BIG {
			coinprice = v["openprice"].ToFloat() - diffPrice
		} else {
			coinprice = v["openprice"].ToFloat() + diffPrice
		}
	}

	config.GlobalMongo.DBHandle.Collection(controllerCollection).DeleteOne(context.TODO(), bson.M{"sn": v["sn"].ToString()})

	profit := -1 * v["credit"].ToFloat() * (v["lose_rate"].ToFloat() / float64(100))
	if iswin {
		profit = v["credit"].ToFloat() * (v["win_rate"].ToFloat() / float64(100))
	}
	backCredit := v["credit"].ToFloat() + profit

	if v["mode"].ToInt() == 1 {
		s.user.AddCredit(v["uid"].ToInt(), &userdomain.CreditValue{
			Credit:          backCredit,
			LockCredit:      0,
			VCrdit:          0,
			LockVCredit:     0,
			UserCoinLogType: creditlogdomain.COIN_LOG_USER_CLOSE,
			UserCoinLogInfo: creditlogdomain.QueueCreditLog{
				Credit:     backCredit,
				LockCredit: 0,
				Sn:         v["sn"].ToString(),
				CreateTime: ntime,
			},
			TeamCoinLogType: creditlogdomain.TEAM_LOG_TRADE_PROFIT,
			TeamCoinLogInfo: creditlogdomain.QueueTeamLog{
				TradeExplode_Profit: profit,
				CreateTime:          ntime,
			},
		})
	} else {
		s.user.AddCredit(v["uid"].ToInt(), &userdomain.CreditValue{
			Credit:          0,
			LockCredit:      0,
			VCrdit:          backCredit,
			LockVCredit:     0,
			UserCoinLogType: 0,
			UserCoinLogInfo: nil,
			TeamCoinLogType: 0,
			TeamCoinLogInfo: nil,
		})
	}

	_ = s.repo.UpdateOpenedFields(v["id"].ToInt(), db.DB_PARAMS{
		"clear_time": v["close_real_time"].ToInt(),
		"profit":     profit,
		"closeprice": coinprice,
	})
	_ = s.repo.InsertClose(db.DB_PARAMS{
		"uid":         v["uid"].ToInt(),
		"sn":          v["sn"].ToString(),
		"coin_symbol": v["coin_symbol"].ToString(),
		"trade_type":  v["trade_type"].ToInt(),
		"flag":        v["flag"].ToInt(),
		"amount":      v["num"].ToFloat(),
		"close_price": coinprice,
		"createtime":  v["close_real_time"].ToInt(),
		"num":         v["num"].ToFloat(),
		"mode":        v["mode"].ToInt(),
		"allprice":    backCredit,
		"o_price":     v["openprice"].ToFloat(),
		"profit":      profit,
	})
	return true
}

func (s *Service) ListPendingDelegates(limit int) ([]db.DBValues, error) {
	if limit <= 0 {
		limit = 500
	}
	return s.repo.FetchPendingDelegates(limit)
}

func (s *Service) ListKeepCrossOpened(page int, limit int) ([]db.DBValues, int) {
	if limit <= 0 {
		limit = 1000
	}
	condition := db.DB_PARAMS{"trade_type": tradedomain.OPEN_TYPE_KEEP, "_": "num>0 and ganggan>1 and auto_delegate=0"}
	count := s.repo.CountOpened(condition)
	page, pageTotal := normalizePage(page, limit, count)
	list, _ := s.repo.FetchOpenedByCondition(condition, fmt.Sprintf("limit %d,%d", (page-1)*limit, limit))
	return list, pageTotal
}

func (s *Service) ListDueExplodeTrades(limit int, ntime int) ([]db.DBValues, error) {
	if limit <= 0 {
		limit = 500
	}
	return s.repo.FetchOpenedByCondition(
		db.DB_PARAMS{"trade_type": tradedomain.OPEN_TYPE_EXPLODE, "clear_time": 0, "_": fmt.Sprintf("close_real_time<=%d", ntime)},
		fmt.Sprintf("limit 0,%d", limit),
	)
}

func buildOpenedInfo(uid int, one db.DBValues) *tradedomain.OpenedInfo {
	if one == nil {
		return nil
	}
	rs := new(tradedomain.OpenedInfo)
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

func (s *Service) operateBBDelegate(one db.DBValues, coin string, coinprice float64, credit float64, virtualCredit float64, lockCredit float64, lockVirtualCredit float64, ntime int) bool {
	uid := one["uid"].ToInt()
	mode := one["mode"].ToInt()
	if one["delegate_type"].ToInt() == tradedomain.DELEGATE_TYPE_BUY {
		if coinprice <= one["price"].ToFloat() {
			if s.asset.AddAssets(uid, &assetsdomain.Assets{
				Coin:    coin,
				Pair:    one["coinpair"].ToString(),
				Num:     one["num"].ToFloat(),
				LockNum: 0,
				Price:   one["price"].ToFloat(),
				Mode:    mode,
			}) {
				s.user.AddCredit(uid, &userdomain.CreditValue{
					Credit:          0,
					LockCredit:      -1 * lockCredit,
					VCrdit:          0,
					LockVCredit:     -1 * lockVirtualCredit,
					UserCoinLogType: creditlogdomain.COIN_LOG_BB_TRADE,
					UserCoinLogInfo: creditlogdomain.QueueCreditLog{
						Credit:     one["num"].ToFloat(),
						CoinType:   one["coin_symbol"].ToString(),
						CreateTime: ntime,
					},
					TeamCoinLogType: creditlogdomain.TEAM_LOG_TRADE,
					TeamCoinLogInfo: creditlogdomain.QueueTeamLog{
						TradeBB:    lockCredit,
						CreateTime: ntime,
					},
				})
				_ = s.repo.UpdateDelegateState(one["id"].ToInt(), 1, ntime)
				return true
			}
		}
		return false
	}

	if coinprice >= one["price"].ToFloat() {
		assetInfoMap := s.asset.GetAllAssets(uid, mode)
		assetInfo, ok := assetInfoMap[one["coin_symbol"].ToString()]
		if !ok {
			return false
		}
		bbProfit := credit - one["num"].ToFloat()*assetInfo.O_Price
		if s.asset.AddAssets(uid, &assetsdomain.Assets{
			Coin:    coin,
			Pair:    one["coinpair"].ToString(),
			Num:     0,
			LockNum: -1 * one["num"].ToFloat(),
			Price:   -1 * one["price"].ToFloat(),
			Mode:    mode,
		}) {
			s.user.AddCredit(uid, &userdomain.CreditValue{
				Credit:          credit,
				LockCredit:      0,
				VCrdit:          virtualCredit,
				LockVCredit:     0,
				UserCoinLogType: creditlogdomain.COIN_LOG_USER_CLOSE,
				UserCoinLogInfo: creditlogdomain.QueueCreditLog{
					Credit:     credit,
					LockCredit: 0,
					Sn:         one["sn"].ToString(),
					CreateTime: ntime,
				},
				TeamCoinLogType: creditlogdomain.TEAM_LOG_TRADE_PROFIT,
				TeamCoinLogInfo: creditlogdomain.QueueTeamLog{
					CreateTime:     ntime,
					TradeBB_Profit: bbProfit,
				},
			})
			s.user.AddCredit(uid, &userdomain.CreditValue{
				Credit:          0,
				LockCredit:      0,
				VCrdit:          0,
				LockVCredit:     0,
				UserCoinLogType: creditlogdomain.COIN_LOG_USER_CLOSE,
				UserCoinLogInfo: creditlogdomain.QueueCreditLog{
					Credit:     -1 * one["num"].ToFloat(),
					Sn:         one["sn"].ToString(),
					CoinType:   one["coin_symbol"].ToString(),
					CreateTime: ntime,
				},
			})
			_ = s.repo.UpdateDelegateState(one["id"].ToInt(), 1, ntime)
			return true
		}
	}
	return false
}

func (s *Service) operateExplodeDelegate(one db.DBValues, coinprice float64, lockCredit float64, lockVirtualCredit float64, ntime int) bool {
	if (one["flag"].ToInt() == tradedomain.DIRECT_TYPE_BIG && coinprice <= one["price"].ToFloat()) ||
		(one["flag"].ToInt() == tradedomain.DIRECT_TYPE_SMALL && coinprice >= one["price"].ToFloat()) {
		explodeConfig, ok := s.market.GetExplodeConfig()[one["close_time"].ToInt()]
		if !ok {
			return false
		}
		insertData := db.DB_PARAMS{
			"uid":             one["uid"].ToInt(),
			"sn":              one["sn"].ToString(),
			"trade_type":      tradedomain.OPEN_TYPE_EXPLODE,
			"flag":            one["flag"].ToInt(),
			"openprice":       one["price"].ToFloat(),
			"closeprice":      0,
			"coinid":          one["coinid"].ToInt(),
			"coinpair":        one["coinpair"].ToString(),
			"coin_symbol":     one["coin_symbol"].ToString(),
			"close_time":      one["close_time"].ToInt(),
			"close_real_time": ntime + one["close_time"].ToInt(),
			"clear_time":      0,
			"createtime":      ntime,
			"ganggan":         one["ganggan"].ToInt(),
			"credit":          one["credit"].ToFloat(),
			"profit":          0,
			"win_rate":        explodeConfig.Winrate,
			"lose_rate":       explodeConfig.Loserate,
			"num":             one["num"].ToFloat(),
			"mode":            one["mode"].ToInt(),
		}
		if err := s.repo.InsertOpened(insertData); err == nil {
			if s.user.AddCredit(one["uid"].ToInt(), &userdomain.CreditValue{
				Credit:          0,
				LockCredit:      -1 * lockCredit,
				VCrdit:          0,
				LockVCredit:     -1 * lockVirtualCredit,
				UserCoinLogType: 0,
				UserCoinLogInfo: nil,
				TeamCoinLogType: creditlogdomain.TEAM_LOG_TRADE,
				TeamCoinLogInfo: creditlogdomain.QueueTeamLog{
					TradeExplode: lockCredit,
					CreateTime:   ntime,
				},
			}) {
				_ = s.repo.UpdateDelegateState(one["id"].ToInt(), 1, ntime)
				return true
			}
		}
	}
	return false
}

func (s *Service) operateKeepDelegate(one db.DBValues, coin string, coinprice float64, credit float64, virtualCredit float64, lockCredit float64, lockVirtualCredit float64, ntime int) bool {
	uid := one["uid"].ToInt()
	if one["delegate_type"].ToInt() == tradedomain.DELEGATE_TYPE_BUY {
		if (one["flag"].ToInt() == tradedomain.DIRECT_TYPE_BIG && coinprice <= one["price"].ToFloat()) ||
			(one["flag"].ToInt() == tradedomain.DIRECT_TYPE_SMALL && coinprice >= one["price"].ToFloat()) {
			if s.user.AddCredit(uid, &userdomain.CreditValue{
				Credit:          0,
				LockCredit:      -1 * lockCredit,
				VCrdit:          0,
				LockVCredit:     -1 * lockVirtualCredit,
				UserCoinLogType: 0,
				UserCoinLogInfo: nil,
				TeamCoinLogType: creditlogdomain.TEAM_LOG_TRADE,
				TeamCoinLogInfo: creditlogdomain.QueueTeamLog{
					TradeKeep:  lockCredit,
					CreateTime: ntime,
				},
			}) {
				s.AddKeepOpened(one)
				return true
			}
		}
		return false
	}

	if (one["flag"].ToInt() == tradedomain.DIRECT_TYPE_BIG && coinprice >= one["price"].ToFloat()) ||
		(one["flag"].ToInt() == tradedomain.DIRECT_TYPE_SMALL && coinprice <= one["price"].ToFloat()) {
		var appendInfo *tradedomain.OpenedInfo
		if one["ganggan_sn"].ToString() != "0" && one["ganggan"].ToInt() > 1 {
			appendInfo = s.GetOpenedBySn(uid, one["ganggan_sn"].ToString())
		} else {
			appendInfo = s.GetOpenedOne(uid, coin, one["trade_type"].ToInt(), one["flag"].ToInt(), one["mode"].ToInt(), one["ganggan"].ToInt())
		}
		if appendInfo == nil {
			return false
		}

		profit := ((one["price"].ToFloat() - appendInfo.OpenPrice) / appendInfo.OpenPrice) * one["num"].ToFloat() * 1000
		if one["flag"].ToInt() != tradedomain.DIRECT_TYPE_BIG {
			profit = ((one["num"].ToFloat()*1000)/one["price"].ToFloat())*appendInfo.OpenPrice - one["num"].ToFloat()*1000
		}
		profit = profit * float64(appendInfo.Ganggan)
		allRealCredit := profit + one["num"].ToFloat()*1000
		insertData := db.DB_PARAMS{
			"uid":         uid,
			"sn":          one["sn"].ToString(),
			"coin_symbol": one["coin_symbol"].ToString(),
			"trade_type":  one["trade_type"].ToInt(),
			"flag":        one["flag"].ToInt(),
			"amount":      one["num"].ToFloat(),
			"close_price": one["price"].ToFloat(),
			"createtime":  ntime,
			"num":         one["num"].ToFloat(),
			"mode":        one["mode"].ToInt(),
			"allprice":    allRealCredit,
			"profit":      profit,
			"o_price":     appendInfo.OpenPrice,
		}
		if err := s.repo.InsertClose(insertData); err == nil {
			_ = s.repo.AddOpenedPositionValue(appendInfo.Id, map[string]float64{
				"lock_num": -1 * one["num"].ToFloat(),
				"credit":   -1 * one["num"].ToFloat() * 1000,
			})
			if one["mode"].ToInt() == 1 {
				s.user.AddCredit(uid, &userdomain.CreditValue{
					Credit:          allRealCredit,
					LockCredit:      0,
					VCrdit:          0,
					LockVCredit:     0,
					UserCoinLogType: creditlogdomain.COIN_LOG_USER_CLOSE,
					UserCoinLogInfo: creditlogdomain.QueueCreditLog{
						Credit:     allRealCredit,
						Sn:         one["sn"].ToString(),
						CreateTime: ntime,
					},
					TeamCoinLogType: creditlogdomain.TEAM_LOG_TRADE_PROFIT,
					TeamCoinLogInfo: creditlogdomain.QueueTeamLog{
						TradeKeep_Profit: profit,
						CreateTime:       ntime,
					},
				})
			} else {
				s.user.AddCredit(uid, &userdomain.CreditValue{
					Credit:          0,
					LockCredit:      0,
					VCrdit:          allRealCredit,
					LockVCredit:     0,
					UserCoinLogType: 0,
					UserCoinLogInfo: nil,
					TeamCoinLogType: 0,
					TeamCoinLogInfo: nil,
				})
			}
			_ = s.repo.UpdateDelegateState(one["id"].ToInt(), 1, ntime)
			return true
		}
	}
	_ = credit
	_ = virtualCredit
	return false
}

func normalizePage(page int, limit int, count int) (int, int) {
	pageTotal := int(math.Ceil(float64(count) / float64(limit)))
	if page > pageTotal {
		page = pageTotal
	}
	if page <= 0 {
		page = 1
	}
	return page, pageTotal
}

func (s *Service) resolveExplodeControl(v db.DBValues) int32 {
	control := int32(0)
	switch s.user.GetExplodeState(v["uid"].ToInt()) {
	case 1:
		control = 1
	case 2:
		control = 2
	case 3:
		if v["flag"].ToInt() == tradedomain.DIRECT_TYPE_BIG {
			control = 1
		}
	case 4:
		if v["flag"].ToInt() == tradedomain.DIRECT_TYPE_SMALL {
			control = 1
		}
	case 5:
		if v["flag"].ToInt() == tradedomain.DIRECT_TYPE_BIG {
			control = 1
		} else {
			control = 2
		}
	case 6:
		if v["flag"].ToInt() == tradedomain.DIRECT_TYPE_BIG {
			control = 2
		} else {
			control = 1
		}
	}
	snControl := s.market.GetControlExplode(v["sn"].ToString())
	if snControl > 0 {
		control = snControl
	}
	return control
}

func (s *Service) resolveExplodeDiffPrice(coinInfo db.DB_ROW_RESULT) float64 {
	diffPrice := float64(1+rand.Intn(100)) / float64(math.Pow10(utils.GetInt(coinInfo["dnum"])))
	if controlPriceMin, ok := coinInfo["contorl_price_min"]; ok {
		if controlPriceMax, okMax := coinInfo["contorl_price_max"]; okMax {
			minPrice := utils.GetFloat(controlPriceMin)
			maxPrice := utils.GetFloat(controlPriceMax)
			if minPrice > 0 && maxPrice > 0 {
				digitPow := math.Pow10(utils.GetInt(coinInfo["dnum"]))
				startMin := minPrice * digitPow
				endMax := maxPrice * digitPow
				diffPrice = (startMin + float64(rand.Intn(int(endMax-startMin)))) / digitPow
			}
		}
	}
	return diffPrice
}
