package models

import (
	"cointrade/config"
	creditrepo "cointrade/internal/credit/repo"
	creditservice "cointrade/internal/credit/service"
	"cointrade/lib/db"
)

type creditTransferUserGateway struct{}

func (creditTransferUserGateway) GetBaseInfo(uid int) *UserBaseInfo {
	return MODEL_USER.GetBaseInfo(uid)
}

func (creditTransferUserGateway) AddCredit(uid int, value *CreditValue) bool {
	return MODEL_USER.AddCredit(uid, value)
}

func (creditTransferUserGateway) RechargeByApprove(uid int, amount float64) *BaseResponse {
	return MODEL_CREDIT.RechargeByApprove(uid, amount)
}

type creditTransferAssetGateway struct{}

func (creditTransferAssetGateway) GetOneAsset(uid int, coin string) *AssetInfo {
	return MODEL_ASSETS.GetOneAsset(uid, coin)
}

func (creditTransferAssetGateway) AddAssets(uid int, asset *Assets) bool {
	return MODEL_ASSETS.AddAssets(uid, asset)
}

type creditTransferConfigGateway struct{}

func (creditTransferConfigGateway) GetMaxWithdrawNum() int {
	return config.GlobalConfig.GetValue("max_withdrawnum").ToInt()
}

var creditTransferSvc = creditservice.NewTransferService(
	creditrepo.NewDBTransferRepository(),
	creditTransferUserGateway{},
	creditTransferAssetGateway{},
	creditTransferConfigGateway{},
)

func (m *CreditModel) AddWallet(uid int, rq *WalletAddressRequest) *BaseResponse {
	return creditTransferSvc.AddWallet(uid, rq)
}

func (m *CreditModel) DeleteWallet(uid, id int) *BaseResponse {
	return creditTransferSvc.DeleteWallet(uid, id)
}

func (m *CreditModel) GetWalletList(uid int) db.DB_LIST_RESULT {
	return creditTransferSvc.GetWalletList(uid)
}

func (m *CreditModel) TransFer(uid int, trans *TransferRequest) *BaseResponse {
	return creditTransferSvc.Transfer(uid, trans)
}

func (m *CreditModel) ExchangeAccount(uid int, rq *ExchangeAccountRequest) *BaseResponse {
	return creditTransferSvc.ExchangeAccount(uid, rq)
}

func (m *CreditModel) ExchangeAccount2(uid int, rq *ExchangeAccountRequest2) *BaseResponse {
	return creditTransferSvc.ExchangeAccount2(uid, rq)
}

func (m *CreditModel) TransFerLogs(uid int, rq *TransFerLogsRequest) *PageBaseResponse {
	return creditTransferSvc.TransferLogs(uid, rq)
}

func (m *CreditModel) TransFerDetail(uid int, sn string) db.DB_ROW_RESULT {
	return creditTransferSvc.TransferDetail(uid, sn)
}
