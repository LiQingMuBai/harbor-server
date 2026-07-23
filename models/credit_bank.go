package models

import (
	"cointrade/config"
	creditrepo "cointrade/internal/credit/repo"
	creditservice "cointrade/internal/credit/service"
)

type creditBankCacheGateway struct{}

func (creditBankCacheGateway) GetBankInfo(uid int) *BankInfo {
	var bankInfo BankInfo
	cacheID := MODEL_CREDIT.MakeCacheId("bank", uid)
	err := config.GlobalRedis.GetObject(HASH_USER_BANK, cacheID, &bankInfo)
	if err == nil && bankInfo.Account != "" {
		return &bankInfo
	}
	return nil
}

func (creditBankCacheGateway) SetBankInfo(uid int, info *BankInfo) {
	cacheID := MODEL_CREDIT.MakeCacheId("bank", uid)
	config.GlobalRedis.SetValue(HASH_USER_BANK, cacheID, info)
}

func (creditBankCacheGateway) DeleteBankInfo(uid int) {
	cacheID := MODEL_CREDIT.MakeCacheId("bank", uid)
	config.GlobalRedis.Del(HASH_USER_BANK, cacheID)
}

var creditBankSvc = creditservice.NewBankService(
	creditrepo.NewDBBankRepository(),
	creditBankCacheGateway{},
)

func (m *CreditModel) BindBank(uid int, rq *BankInfo) *BaseResponse { //用户绑定银行卡
	return creditBankSvc.BindBank(uid, rq)
}

func (m *CreditModel) GetBankInfo(uid int) *BankInfo {
	return creditBankSvc.GetBankInfo(uid)
}
