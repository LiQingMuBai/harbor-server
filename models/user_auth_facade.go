package models

import (
	userauthrepo "cointrade/internal/userauth/repo"
	userauthservice "cointrade/internal/userauth/service"
	"cointrade/lib/db"
)

type userAuthUserGateway struct{}

func (userAuthUserGateway) GetBaseInfo(uid int) *UserBaseInfo {
	return MODEL_USER.GetBaseInfo(uid)
}

func (userAuthUserGateway) Update(uid int, data db.DB_PARAMS) {
	MODEL_USER.Update(uid, data)
}

func (userAuthUserGateway) ClearCache(uid int) {
	MODEL_USER.ClearCache(uid)
}

type userAuthWalletGateway struct{}

func (userAuthWalletGateway) RegisterByAddress(address string, ip string) int {
	return MODEL_USER.RegisterByAddress(address, ip)
}

var userAuthSvc = userauthservice.NewService(
	userauthrepo.NewDBRepository(),
	userAuthUserGateway{},
	userAuthWalletGateway{},
)
