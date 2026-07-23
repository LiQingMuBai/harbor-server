package models

import (
	"cointrade/lib/db"

	tradedomain "cointrade/internal/domain/trade"
	traderepo "cointrade/internal/trade/repo"
	tradeservice "cointrade/internal/trade/service"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 交易模块
type TradeModel struct {
	ModelBase
}

const (
	OPEN_TYPE_KEEP               = tradedomain.OPEN_TYPE_KEEP
	OPEN_TYPE_EXPLODE            = tradedomain.OPEN_TYPE_EXPLODE
	OPEN_TYPE_BB                 = tradedomain.OPEN_TYPE_BB
	PRICE_TYPE_LIMIT             = tradedomain.PRICE_TYPE_LIMIT
	PRICE_TYPE_MARKET            = tradedomain.PRICE_TYPE_MARKET
	DIRECT_TYPE_BIG              = tradedomain.DIRECT_TYPE_BIG
	DIRECT_TYPE_SMALL            = tradedomain.DIRECT_TYPE_SMALL
	DELEGATE_TYPE_BUY            = tradedomain.DELEGATE_TYPE_BUY
	DELEGATE_TYPE_SELL           = tradedomain.DELEGATE_TYPE_SELL
	TRADE_BUY_PREFIX             = tradedomain.TRADE_BUY_PREFIX
	TRADE_SELL_PREFIX            = tradedomain.TRADE_SELL_PREFIX
	DELEGATE_STATE_NOCOIN        = tradedomain.DELEGATE_STATE_NOCOIN
	DELEGATE_STATE_CREDIT        = tradedomain.DELEGATE_STATE_CREDIT
	DELEGATE_STATE_CLOSETIME     = tradedomain.DELEGATE_STATE_CLOSETIME
	DELEGATE_STATE_NOASSET       = tradedomain.DELEGATE_STATE_NOASSET
	DELEGATE_STATE_MIN           = tradedomain.DELEGATE_STATE_MIN
	DELEGATE_STATE_TRADE_CLOSED  = tradedomain.DELEGATE_STATE_TRADE_CLOSED
	DELEGATE_STATE_GANGGAN_ERROR = tradedomain.DELEGATE_STATE_GANGGAN_ERROR
)

type TradeDelegateRequest = tradedomain.TradeDelegateRequest

type OpenedInfo = tradedomain.OpenedInfo

type CloseTrade = tradedomain.CloseTrade

type DelegateInfo = tradedomain.DelegateInfo

type TradeListRequest = tradedomain.TradeListRequest

type tradeUserGateway struct{}

func (tradeUserGateway) GetBaseInfo(uid int) *UserBaseInfo {
	return MODEL_USER.GetBaseInfo(uid)
}

func (tradeUserGateway) AddCredit(uid int, value *CreditValue) bool {
	return MODEL_USER.AddCredit(uid, value)
}

func (tradeUserGateway) GetExplodeState(uid int) int {
	return MODEL_USER.GetExplodeState(uid).State
}

type tradeAssetGateway struct{}

func (tradeAssetGateway) GetAllAssets(uid int, mode int) map[string]AssetInfo {
	return MODEL_ASSETS.GetAllAssets(uid, mode)
}

func (tradeAssetGateway) AddAssets(uid int, asset *Assets) bool {
	return MODEL_ASSETS.AddAssets(uid, asset)
}

type tradeMarketGateway struct{}

func (tradeMarketGateway) GetCoinInfo(coin string, pair string) db.DB_ROW_RESULT {
	return MODEL_SYSTEM.GetCoinInfo(coin, pair)
}

func (tradeMarketGateway) GetLastCoinInfo(pair string) primitive.M {
	return MODEL_SYSTEM.GetLastCoinInfo(pair)
}

func (tradeMarketGateway) GetExplodeConfig() map[int]*ExplodeConfig {
	return MODEL_SYSTEM.GetExplodeConfig()
}

func (tradeMarketGateway) GetControlExplode(sn string) int32 {
	return MODEL_SYSTEM.GetControlExplode(sn)
}

var tradeSvc = tradeservice.New(
	traderepo.NewDBRepository(),
	tradeUserGateway{},
	tradeAssetGateway{},
	tradeMarketGateway{},
)

func (m *TradeModel) OperateDelegateTrade(one db.DBValues) bool {
	return tradeSvc.OperateDelegateTrade(one)
}

func (m *TradeModel) EqualsPrice(delegatePrice float64, coinprice float64, isbig bool) bool {
	return tradeSvc.EqualsPrice(delegatePrice, coinprice, isbig)
}

func (m *TradeModel) CheckStop(v db.DBValues, coinprice float64, isbig bool) bool {
	return tradeSvc.CheckStop(v, coinprice, isbig)
}

func (m *TradeModel) SettleKeepCross(v db.DBValues, coinprice float64, ntime int) bool {
	return tradeSvc.SettleKeepCross(v, coinprice, ntime)
}

func (m *TradeModel) SettleExplodeTrade(v db.DBValues, ntime int) bool {
	return tradeSvc.SettleExplodeTrade(v, ntime, COIN_CONTROLLER)
}

func (m *TradeModel) ListPendingDelegates(limit int) ([]db.DBValues, error) {
	return tradeSvc.ListPendingDelegates(limit)
}

func (m *TradeModel) ListKeepCrossOpened(page int, limit int) ([]db.DBValues, int) {
	return tradeSvc.ListKeepCrossOpened(page, limit)
}

func (m *TradeModel) ListDueExplodeTrades(limit int, ntime int) ([]db.DBValues, error) {
	return tradeSvc.ListDueExplodeTrades(limit, ntime)
}
