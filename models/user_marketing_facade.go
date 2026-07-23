package models

import (
	userdomain "cointrade/internal/domain/user"
	usermarketingrepo "cointrade/internal/usermarketing/repo"
	usermarketingservice "cointrade/internal/usermarketing/service"
)

type userMarketingGateway struct{}

func (userMarketingGateway) GetBaseInfo(uid int) *userdomain.UserBaseInfo {
	return MODEL_USER.GetBaseInfo(uid)
}

func (userMarketingGateway) AddCredit(uid int, value *userdomain.CreditValue) bool {
	return MODEL_USER.AddCredit(uid, value)
}

func (userMarketingGateway) ClearCache(uid int) {
	MODEL_USER.ClearCache(uid)
}

func (userMarketingGateway) EncodePassword(password string) string {
	return MODEL_USER.EncodePassword(password)
}

var userMarketingSvc = usermarketingservice.NewService(
	usermarketingrepo.NewDBRepository(),
	userMarketingGateway{},
)
