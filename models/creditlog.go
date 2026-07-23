package models

import creditlogdomain "cointrade/internal/domain/creditlog"

// 账变信息包
const (
	COIN_LOG_USER_RECHARGE         = creditlogdomain.COIN_LOG_USER_RECHARGE
	COIN_LOG_USER_WITHDRAW         = creditlogdomain.COIN_LOG_USER_WITHDRAW
	COIN_LOG_USER_PROFIT           = creditlogdomain.COIN_LOG_USER_PROFIT
	COIN_LOG_USER_CLOSE            = creditlogdomain.COIN_LOG_USER_CLOSE
	COIN_LOG_USER_DELEGATE         = creditlogdomain.COIN_LOG_USER_DELEGATE
	COIN_LOG_USER_DELEGATE_SUCCESS = creditlogdomain.COIN_LOG_USER_DELEGATE_SUCCESS
	COIN_LOG_USER_CANCLE           = creditlogdomain.COIN_LOG_USER_CANCLE
	COIN_LOG_USER_BUY_MINING       = creditlogdomain.COIN_LOG_USER_BUY_MINING
	COIN_LOG_USER_MINING_PROFIT    = creditlogdomain.COIN_LOG_USER_MINING_PROFIT

	COIN_LOG_USER_CLEAR_INCOME    = creditlogdomain.COIN_LOG_USER_CLEAR_INCOME
	COIN_LOG_USER_WITHDRAW_FAILD  = creditlogdomain.COIN_LOG_USER_WITHDRAW_FAILD
	COIN_LOG_USER_EXCHANGE        = creditlogdomain.COIN_LOG_USER_EXCHANGE
	COIN_LOG_BB_TRADE             = creditlogdomain.COIN_LOG_BB_TRADE
	COIN_LOG_EXPLODE_TRADE        = creditlogdomain.COIN_LOG_EXPLODE_TRADE
	COIN_LOG_KEEP_TRADE           = creditlogdomain.COIN_LOG_KEEP_TRADE
	COIN_LOG_BACKEND              = creditlogdomain.COIN_LOG_BACKEND
	COIN_LOG_USER_MINING_BACK     = creditlogdomain.COIN_LOG_USER_MINING_BACK
	COIN_LOG_ASSETS_EXCHANGE      = creditlogdomain.COIN_LOG_ASSETS_EXCHANGE
	COIN_LOG_MINING_UNLOCK        = creditlogdomain.COIN_LOG_MINING_UNLOCK
	COIN_LOG_BUY_COIN             = creditlogdomain.COIN_LOG_BUY_COIN
	COIN_LOG_EXCHANGE_ACCOUNT_IN  = creditlogdomain.COIN_LOG_EXCHANGE_ACCOUNT_IN
	COIN_LOG_EXCHANGE_ACCOUNT_OUT = creditlogdomain.COIN_LOG_EXCHANGE_ACCOUNT_OUT
	COIN_LOG_LOAN_BACK            = creditlogdomain.COIN_LOG_LOAN_BACK
	COIN_LOG_LORA_IN              = creditlogdomain.COIN_LOG_LORA_IN
	COIN_LOG_KEEP_BREAK           = creditlogdomain.COIN_LOG_KEEP_BREAK
	COIN_LOG_USER_REVERVATION     = creditlogdomain.COIN_LOG_USER_REVERVATION
	TEAM_LOG_RECHARGE             = creditlogdomain.TEAM_LOG_RECHARGE
	TEAM_LOG_WITHDRAW             = creditlogdomain.TEAM_LOG_WITHDRAW
	TEAM_LOG_MINING               = creditlogdomain.TEAM_LOG_MINING
	TEAM_LOG_MINING_PROFIT        = creditlogdomain.TEAM_LOG_MINING_PROFIT
	TEAM_LOG_TRADE                = creditlogdomain.TEAM_LOG_TRADE
	TEAM_LOG_TRADE_PROFIT         = creditlogdomain.TEAM_LOG_TRADE_PROFIT
	COIN_LOG_SPOT_BACK            = creditlogdomain.COIN_LOG_SPOT_BACK

	LOG_TIMETYPE_ALL   = creditlogdomain.LOG_TIMETYPE_ALL
	LOG_TIMETYPE_DAY   = creditlogdomain.LOG_TIMETYPE_DAY
	LOG_TIMETYPE_MONTH = creditlogdomain.LOG_TIMETYPE_MONTH

	INCOME_TYPE_RECHARGE      = creditlogdomain.INCOME_TYPE_RECHARGE
	INCOME_TYPE_MINING_BUY    = creditlogdomain.INCOME_TYPE_MINING_BUY
	INCOME_TYPE_MINING_PROFIT = creditlogdomain.INCOME_TYPE_MINING_PROFIT
)

type CreditLogModel struct {
	ModelBase
}
type CreditLogInfo = creditlogdomain.CreditLogInfo

type QueueCreditLog = creditlogdomain.QueueCreditLog

type QueueTeamLog = creditlogdomain.QueueTeamLog

type CoinLogRequest = creditlogdomain.CoinLogRequest
