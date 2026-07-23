package models

import assetsdomain "cointrade/internal/domain/assets"

// 用户资产相关包
const (
	EXCHANGE_STATE_NOCOIN      = assetsdomain.EXCHANGE_STATE_NOCOIN
	EXCHANGE_STATE_NOTENNOUGH  = assetsdomain.EXCHANGE_STATE_NOTENNOUGH
	EXCHANGE_STATE_TOOMIN      = assetsdomain.EXCHANGE_STATE_TOOMIN
	EXCHANGE_STATE_NOT_TRANS   = assetsdomain.EXCHANGE_STATE_NOT_TRANS
	ASSETS_TRANS_TYPE_IN       = assetsdomain.ASSETS_TRANS_TYPE_IN
	ASSETS_TRANS_TYPE_OUT      = assetsdomain.ASSETS_TRANS_TYPE_OUT
	ASSETS_TRANS_TYPE_CONTRACT = assetsdomain.ASSETS_TRANS_TYPE_CONTRACT
)

type AssetModel struct {
	ModelBase
}
type AssetInfo = assetsdomain.AssetInfo

type Assets = assetsdomain.Assets

type ExchangeRequest = assetsdomain.ExchangeRequest

type AssetsTransRequest = assetsdomain.AssetsTransRequest

type QuickExchangeRequest = assetsdomain.QuickExchangeRequest
