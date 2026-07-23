package models

import creditdomain "cointrade/internal/domain/credit"

type CreditModel struct {
	ModelBase
}

const (
	CREDIT_TYPE_RECHARGE         = creditdomain.CREDIT_TYPE_RECHARGE
	CREDIT_TYPE_WITHDRAW         = creditdomain.CREDIT_TYPE_WITHDRAW
	CREDIT_TYPE_TRANSFER         = creditdomain.CREDIT_TYPE_TRANSFER
	RECHARGE_ORDER_PREFIX        = creditdomain.RECHARGE_ORDER_PREFIX
	WITHDRAW_ORDER_PREFIX        = creditdomain.WITHDRAW_ORDER_PREFIX
	TRANSFER_ORDER_PREFIX        = creditdomain.TRANSFER_ORDER_PREFIX
	RECHARGE_STATE_MIN           = creditdomain.RECHARGE_STATE_MIN
	RECHARGE_STATE_ERROR_ADDRESS = creditdomain.RECHARGE_STATE_ERROR_ADDRESS
	RECHARGE_STATE_ERROR_USER    = creditdomain.RECHARGE_STATE_ERROR_USER
	RECHARGE_STATE_ERROR_PROOF   = creditdomain.RECHARGE_STATE_ERROR_PROOF

	WITHDRAW_STATE_MIN                = creditdomain.WITHDRAW_STATE_MIN
	WITHDRAW_STATE_ERROR_USER         = creditdomain.WITHDRAW_STATE_ERROR_USER
	WITHDRAW_STATE_ERROR_ADDRESS      = creditdomain.WITHDRAW_STATE_ERROR_ADDRESS
	WITHDRAW_STATE_NOTENOUGH          = creditdomain.WITHDRAW_STATE_NOTENOUGH
	WITHDRAW_STATE_ERROR_CASHPASSWORD = creditdomain.WITHDRAW_STATE_ERROR_CASHPASSWORD
	WITHDRAW_STATE_ERROR_LOCKED       = creditdomain.WITHDRAW_STATE_ERROR_LOCKED
	WITHDRAW_STATE_ERROR_NOTBINDBANK  = creditdomain.WITHDRAW_STATE_ERROR_NOTBINDBANK
	RECHARGE_STATE_ERROR_NOTAPPROVE    = creditdomain.RECHARGE_STATE_ERROR_NOTAPPROVE
	RECHARGE_STATE_ERROR_MONEY         = creditdomain.RECHARGE_STATE_ERROR_MONEY
	RECHARGE_STATE_ERROR_TRANS         = creditdomain.RECHARGE_STATE_ERROR_TRANS
	RECHARGE_STATE_ERROR_MAX_WITHDRAW  = creditdomain.RECHARGE_STATE_ERROR_MAX_WITHDRAW
	TRANSFER_DIRECTION_IN              = creditdomain.TRANSFER_DIRECTION_IN
	TRANSFER_DIRECTION_OUT             = creditdomain.TRANSFER_DIRECTION_OUT

	EXCHANGE_DIRECTION_CONTRACT = creditdomain.EXCHANGE_DIRECTION_CONTRACT
	EXCHANGE_DIRECTION_ACCOUNT  = creditdomain.EXCHANGE_DIRECTION_ACCOUNT
)

type TransferRequest = creditdomain.TransferRequest

type RechargeRequest = creditdomain.RechargeRequest

type TransferLogsRequest = creditdomain.TransferLogsRequest

type WithdrawRequest = creditdomain.WithdrawRequest

type WalletAddressRequest = creditdomain.WalletAddressRequest

type RechargeResponse = creditdomain.RechargeResponse

type BankInfo = creditdomain.BankInfo

type ExchangeAccountRequest = creditdomain.ExchangeAccountRequest

type ExchangeAccountRequest2 = creditdomain.ExchangeAccountRequest2
