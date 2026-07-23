package models

import systemdomain "cointrade/internal/domain/system"

const (
	COIN_BUY_STATE_NOTENGOUGH = systemdomain.COIN_BUY_STATE_NOTENGOUGH
	COIN_BUY_STATE_NOMONEY    = systemdomain.COIN_BUY_STATE_NOMONEY
)

type SystemModel struct {
	ModelBase
}

type ExplodeConfig = systemdomain.ExplodeConfig

type RechargeContractConfig = systemdomain.RechargeContractConfig

type RechargeConfig = systemdomain.RechargeConfig

type KlineControlConfig = systemdomain.KlineControlConfig

type CoinKlineConfig = systemdomain.CoinKlineConfig

var PERIOD_LIST = systemdomain.PERIOD_LIST
