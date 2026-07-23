package models

import (
	"cointrade/config"
	shareddomain "cointrade/internal/domain/shared"
	userdomain "cointrade/internal/domain/user"
	"cointrade/lib/db"
)

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
	rs := new(BaseResponse)
	if rq.Password == "" {
		rs.State = 1
		rs.Msg = shareddomain.MsgPasswordRequired
		return rs
	}
	if len(rq.Password) < 6 || len(rq.Password) > 20 {
		rs.State = 1
		rs.Msg = shareddomain.MsgPasswordLength
		return rs
	}
	if rq.Password != rq.RePassword {
		rs.State = 1
		rs.Msg = shareddomain.MsgPasswordInvalid
		return rs
	}
	m.Update(uid, db.DB_PARAMS{"cash_password": rq.Password})
	rs.State = 0
	rs.Msg = shareddomain.MsgSuccess
	return rs
}

func (m *UserModel) UpdateCashPassword(uid int, rq *UpdateCashPasswordRequest) *BaseResponse {
	uinfo := m.GetBaseInfo(uid)
	rs := new(BaseResponse)
	if uinfo == nil {
		rs.State = STATE_FAILD
		rs.Msg = shareddomain.MsgUserNotFound
		return rs
	}
	if rq.O_Password != uinfo.CashPassword {
		rs.State = CHANGE_PASS_STATE_OLDERROR
		rs.Msg = shareddomain.MsgPasswordInvalid
		return rs
	}
	if len(rq.N_Password) < 6 || len(rq.N_Password) > 20 {
		rs.State = REGISTER_STATE_ERRORPASSWORD
		rs.Msg = shareddomain.MsgPasswordLength
		return rs
	}
	if rq.N_Password != rq.R_Password {
		rs.State = REGISTER_STATE_ERRORPASSWORD
		rs.Msg = shareddomain.MsgPasswordInvalid
		return rs
	}
	m.Update(uid, db.DB_PARAMS{"cash_password": rq.N_Password})
	rs.State = STATE_SUCCESS
	rs.Msg = shareddomain.MsgSuccess
	return rs
}

func (m *UserModel) BindPhone(uid int, phone string, code string) *BaseResponse {
	rs := new(BaseResponse)
	uinfo := m.GetBaseInfo(uid)
	if uinfo.Phone != "" {
		rs.State = BIND_PHONE_STATE_BINDED
		rs.Msg = shareddomain.MsgAlreadyBound
		return rs
	}
	one, _ := config.GlobalDB.FetchOne(DB_TABLE_USER, db.DB_PARAMS{"phone": phone}, db.DB_FIELDS{"id"})
	if one != nil {
		rs.State = BIND_PHONE_STATE_BINDED
		rs.Msg = shareddomain.MsgAlreadyBound
		return rs
	}
	vCode := MODEL_CODE.GetBindSmsCode(uid, phone)
	if vCode != code {
		rs.State = BIND_PHONE_STATE_ERRORCODE
		rs.Msg = shareddomain.MsgVerifyCodeInvalid
		return rs
	}
	m.Update(uid, db.DB_PARAMS{"phone": phone})
	rs.State = STATE_SUCCESS
	rs.Msg = shareddomain.MsgSuccess
	return rs
}

func (m *UserModel) BindEmail(uid int, email string, code string) *BaseResponse {
	rs := new(BaseResponse)
	uinfo := m.GetBaseInfo(uid)
	if uinfo.Email != "" && uinfo.Email != "0" {
		rs.State = BIND_PHONE_STATE_BINDED
		rs.Msg = shareddomain.MsgAlreadyBound
		return rs
	}
	one, _ := config.GlobalDB.FetchOne(DB_TABLE_USER, db.DB_PARAMS{"email": email}, db.DB_FIELDS{"id"})
	if one != nil {
		rs.State = BIND_PHONE_STATE_BINDED
		rs.Msg = shareddomain.MsgAlreadyBound
		return rs
	}
	vCode := MODEL_CODE.GetEmailCodeBind(email)
	if vCode != code {
		rs.State = BIND_PHONE_STATE_ERRORCODE
		rs.Msg = shareddomain.MsgVerifyCodeInvalid
		return rs
	}
	m.Update(uid, db.DB_PARAMS{"email": email})
	rs.State = STATE_SUCCESS
	rs.Msg = shareddomain.MsgSuccess
	return rs
}
