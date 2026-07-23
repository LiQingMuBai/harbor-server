package creditlog

import "cointrade/internal/domain/shared"

const (
	COIN_LOG_USER_RECHARGE         = 5100001
	COIN_LOG_USER_WITHDRAW         = 5100002
	COIN_LOG_USER_PROFIT           = 5100003
	COIN_LOG_USER_CLOSE            = 5100004
	COIN_LOG_USER_DELEGATE         = 5100005
	COIN_LOG_USER_DELEGATE_SUCCESS = 5100010
	COIN_LOG_USER_CANCLE           = 5100007
	COIN_LOG_USER_BUY_MINING       = 5100008
	COIN_LOG_USER_MINING_PROFIT    = 5100009

	COIN_LOG_USER_CLEAR_INCOME    = 5100011
	COIN_LOG_USER_WITHDRAW_FAILD  = 5100012
	COIN_LOG_USER_EXCHANGE        = 5100013
	COIN_LOG_BB_TRADE             = 5100014
	COIN_LOG_EXPLODE_TRADE        = 5100015
	COIN_LOG_KEEP_TRADE           = 5100016
	COIN_LOG_BACKEND              = 5100017
	COIN_LOG_USER_MINING_BACK     = 5100018
	COIN_LOG_ASSETS_EXCHANGE      = 5100019
	COIN_LOG_MINING_UNLOCK        = 5100020
	COIN_LOG_BUY_COIN             = 5100021
	COIN_LOG_EXCHANGE_ACCOUNT_IN  = 5100022
	COIN_LOG_EXCHANGE_ACCOUNT_OUT = 5100023
	COIN_LOG_LOAN_BACK            = 5100025
	COIN_LOG_LORA_IN              = 5100024
	COIN_LOG_KEEP_BREAK           = 5100026
	COIN_LOG_USER_REVERVATION     = 5100027
	TEAM_LOG_RECHARGE             = 5200001
	TEAM_LOG_WITHDRAW             = 5200002
	TEAM_LOG_MINING               = 5200003
	TEAM_LOG_MINING_PROFIT        = 5200004
	TEAM_LOG_TRADE                = 5200005
	TEAM_LOG_TRADE_PROFIT         = 5200006
	COIN_LOG_SPOT_BACK            = 5200028

	LOG_TIMETYPE_ALL   = 0
	LOG_TIMETYPE_DAY   = 1
	LOG_TIMETYPE_MONTH = 2

	INCOME_TYPE_RECHARGE      = 1
	INCOME_TYPE_MINING_BUY    = 2
	INCOME_TYPE_MINING_PROFIT = 3
)

type CreditLogInfo struct {
	Uid        int     `json:"uid"`
	Credit     float64 `json:"credit"`
	LockCredit float64 `json:"lockcredit"`
	Mode       int     `josn:"mode"`
	Sn         string  `josn:"sn"`
	Type       int     `json:"type"`
	CoinType   string  `json:"cointype"`
	Createtime int     `json:"credittime"`
}

type QueueCreditLog struct {
	Credit     float64 `json:"credit"`
	LockCredit float64 `json:"lockcredit"`
	CoinType   string  `json:"cointype"`
	Sn         string  `json:"sn"`
	CreateTime int     `json:"createtime"`
}

type QueueTeamLog struct {
	Recharge            float64 `json:"recharge"`
	Withdraw            float64 `json:"withdraw"`
	Trade               float64 `json:"trade"`
	TradeProfit         float64 `json:"trade_profit"`
	MiningCount         float64 `json:"mining_count"`
	MiningProfit        float64 `json:"mining_profit"`
	CreateTime          int     `json:"createtime"`
	TradeBB             float64 `json:"trade_bb"`
	TradeExplode        float64 `json:"trade_explode"`
	TradeKeep           float64 `json:"trade_keep"`
	TradeBB_Profit      float64 `json:"trade_bb_profit"`
	TradeExplode_Profit float64 `json:"trade_explode_profit"`
	TradeKeep_Profit    float64 `json:"trade_keep_profit"`
}

type CoinLogRequest struct {
	shared.PageBaseRequest
	Type int `json:"type"`
}
