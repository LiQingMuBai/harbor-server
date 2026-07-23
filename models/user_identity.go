package models

import userdomain "cointrade/internal/domain/user"

func (m *UserModel) AuthLv1(uid int, authinfo *AuthLv1Request) *BaseResponse {
	return userIdentitySvc.AuthLv1(uid, (*userdomain.AuthLv1Request)(authinfo))
}

func (m *UserModel) AuthLv2(uid int, rq *AuthLv2Request) *BaseResponse {
	return userIdentitySvc.AuthLv2(uid, (*userdomain.AuthLv2Request)(rq))
}

func (m *UserModel) GetAuthInfo(uid int) map[int]interface{} {
	return userIdentitySvc.GetAuthInfo(uid)
}

func (m *UserModel) ChangeCashPassword(uid int, rq *SetCashPasswordRequest) *BaseResponse {
	return userSecuritySvc.ChangeCashPassword(uid, (*userdomain.SetCashPasswordRequest)(rq))
}

func (m *UserModel) UpdateCashPassword(uid int, rq *UpdateCashPasswordRequest) *BaseResponse {
	return userSecuritySvc.UpdateCashPassword(uid, (*userdomain.UpdateCashPasswordRequest)(rq))
}

func (m *UserModel) BindPhone(uid int, phone string, code string) *BaseResponse {
	return userSecuritySvc.BindPhone(uid, phone, code)
}

func (m *UserModel) BindEmail(uid int, email string, code string) *BaseResponse {
	return userSecuritySvc.BindEmail(uid, email, code)
}
