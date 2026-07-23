package models

import (
	useridentityrepo "cointrade/internal/useridentity/repo"
	useridentityservice "cointrade/internal/useridentity/service"
	userdomain "cointrade/internal/domain/user"
	"cointrade/lib/db"
	"cointrade/lib/notify"
)

type userIdentityUserGateway struct{}

func (userIdentityUserGateway) GetBaseInfo(uid int) *userdomain.UserBaseInfo {
	return MODEL_USER.GetBaseInfo(uid)
}

func (userIdentityUserGateway) Update(uid int, data db.DB_PARAMS) {
	MODEL_USER.Update(uid, data)
}

type userIdentityNotifier struct{}

func (userIdentityNotifier) IncrementNotify(typ int, num int) {
	notify.NOTIFY.AddNotify(&notify.NotifyItem{Type: typ, Num: num})
}

var userIdentitySvc = useridentityservice.NewService(
	useridentityrepo.NewDBRepository(),
	userIdentityUserGateway{},
	userIdentityNotifier{},
)
