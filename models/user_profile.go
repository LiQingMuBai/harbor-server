package models

import (
	"cointrade/lib/db"
)

func (m *UserModel) UpdateProfile(uid int, rq *UpdateProfileRequest) *BaseResponse {
	return userProfileSvc.UpdateProfile(uid, rq, m.MakeCacheId)
}

func (m *UserModel) GetUserCount(uid int, t int) *UserCount {
	return userProfileSvc.GetUserCount(uid, t)
}

func (m *UserModel) GetBaseInfo(uid int) *UserBaseInfo { //获得单个用户的基础信息
	return userProfileGetBaseInfo(uid)
}

func (m *UserModel) IsNewUser(uid int) *BaseResponse {
	return userMarketingSvc.IsNewUser(uid)
}

func (m *UserModel) Claim(uid int) *BaseResponse {
	return userMarketingSvc.Claim(uid)
}

func (m *UserModel) ClearIncome(uid int, cashpassword string) *BaseResponse {
	return userMarketingSvc.ClearIncome(uid, cashpassword)
}

func (m *UserModel) Welcome2() *WelcomeResponse {
	return userPortalSvc.Welcome2()
}

func (m *UserModel) CrossTrade(uid int, data db.DB_PARAMS) *CrossPlatformTradeResponse {
	return userPortalSvc.CrossTrade(uid, data)
}

func (m *UserModel) Welcome() *WelcomeResponse {
	return userPortalSvc.Welcome()
}

func (m *UserModel) ChangeMode(uid int) *BaseResponse {
	return userProfileSvc.ChangeMode(uid, userProfileGetBaseInfo)
}

func (m *UserModel) GetExplodeState(uid int) *BaseResponse {
	return userProfileSvc.GetExplodeState(uid)
}

func (m *UserModel) ConvertMoney(uid int) db.DB_PARAMS {
	return userProfileSvc.ConvertMoney(uid)
}
