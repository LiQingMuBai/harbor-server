package models

import tradedomain "cointrade/internal/domain/trade"

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
