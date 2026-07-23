package models

import (
	userdomain "cointrade/internal/domain/user"
	usersecurityrepo "cointrade/internal/usersecurity/repo"
	usersecurityservice "cointrade/internal/usersecurity/service"
	"cointrade/lib/db"
)

type userSecurityUserGateway struct{}

func (userSecurityUserGateway) GetBaseInfo(uid int) *userdomain.UserBaseInfo {
	return MODEL_USER.GetBaseInfo(uid)
}

func (userSecurityUserGateway) Update(uid int, data db.DB_PARAMS) {
	MODEL_USER.Update(uid, data)
}

type userSecurityCodeGateway struct{}

func (userSecurityCodeGateway) GetBindSMSCode(uid int, phone string) string {
	return MODEL_CODE.GetBindSmsCode(uid, phone)
}

func (userSecurityCodeGateway) GetEmailBindCode(email string) string {
	return MODEL_CODE.GetEmailCodeBind(email)
}

var userSecuritySvc = usersecurityservice.NewService(
	usersecurityrepo.NewDBRepository(),
	userSecurityUserGateway{},
	userSecurityCodeGateway{},
)
