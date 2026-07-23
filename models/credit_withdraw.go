package models

import (
	"cointrade/config"
	creditrepo "cointrade/internal/credit/repo"
	creditservice "cointrade/internal/credit/service"
	"cointrade/lib/db"
)

type creditWithdrawUserGateway struct{}

func (creditWithdrawUserGateway) GetBaseInfo(uid int) *UserBaseInfo {
	return MODEL_USER.GetBaseInfo(uid)
}

func (creditWithdrawUserGateway) AddCredit(uid int, value *CreditValue) bool {
	return MODEL_USER.AddCredit(uid, value)
}

type creditWithdrawSystemGateway struct{}

func (creditWithdrawSystemGateway) GetRechargeConfig(cointype string, contract string) *RechargeContractConfig {
	return MODEL_SYSTEM.GetOneRechargeConfig(cointype, contract)
}

func (creditWithdrawSystemGateway) GetCoinClosePrice(pair string) float64 {
	coinPriceInfo := MODEL_SYSTEM.GetLastCoinInfo(pair)
	if coinPriceInfo == nil {
		return 0
	}
	closePrice, ok := coinPriceInfo["close"].(float64)
	if !ok {
		return 0
	}
	return closePrice
}

func (creditWithdrawSystemGateway) GetMinWithdraw() float64 {
	return config.GlobalConfig.GetValue("min_withdraw").ToFloat()
}

func (creditWithdrawSystemGateway) GetWithdrawFee() float64 {
	return config.GlobalConfig.GetValue("withdraw_fee").ToFloat()
}

type creditWithdrawBankGateway struct{}

func (creditWithdrawBankGateway) GetBankInfo(uid int) *BankInfo {
	return MODEL_CREDIT.GetBankInfo(uid)
}

type creditWithdrawNotifier struct{}

func (creditWithdrawNotifier) IncrementNotify(typ int, num int) {
	creditRechargeNotifier{}.IncrementNotify(typ, num)
}

var creditWithdrawSvc = creditservice.NewWithdrawService(
	creditrepo.NewDBWithdrawRepository(),
	creditWithdrawUserGateway{},
	creditWithdrawSystemGateway{},
	creditWithdrawBankGateway{},
	creditWithdrawNotifier{},
)

func (m *CreditModel) CreateWithdraw(uid int, rq *WithdrawRequest) *RechargeResponse {
	return creditWithdrawSvc.CreateWithdraw(uid, rq)
}

func (m *CreditModel) CreateWithDraw(uid int, rq *WithDrawRequest) *RechargeResponse {
	return m.CreateWithdraw(uid, rq)
}

func (m *CreditModel) GetWithdrawList(uid int, rq *PageBaseRequest) *PageBaseResponse {
	return creditWithdrawSvc.GetWithdrawList(uid, rq)
}

func (m *CreditModel) GetWithDrawList(uid int, rq *PageBaseRequest) *PageBaseResponse {
	return m.GetWithdrawList(uid, rq)
}

func (m *CreditModel) WithdrawInfo(uid int, sn string) db.DB_ROW_RESULT {
	return creditWithdrawSvc.WithdrawInfo(uid, sn)
}
