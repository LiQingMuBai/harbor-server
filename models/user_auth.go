package models

import (
	userdomain "cointrade/internal/domain/user"
	"cointrade/utils"
	"strings"
)

func (m *UserModel) EncodePassword(password string) string {
	return userAuthSvc.EncodePassword(password)
}

func (m *UserModel) GetInviteUser(code string) int {
	return userAuthSvc.GetInviteUser(code)
}

func (m *UserModel) Register(rq *RegisterRequest) *BaseResponse { //注册
	return userAuthSvc.Register((*userdomain.RegisterRequest)(rq))
}

func (m *UserModel) GetChanelLevel(channelID int, parentOrder string, n int) int {
	list := strings.Split(parentOrder, ",")
	j := 0
	l := len(list)
	if l == 0 {
		return 0
	}
	for _, v := range list {
		if utils.GetInt(v) == channelID {
			return n - j
		}
		j++
	}
	n += l
	topuser := m.GetBaseInfo(utils.GetInt(list[0]))
	if topuser == nil {
		return 0
	}
	return m.GetChanelLevel(channelID, topuser.ParentOrder, n)
}

func (m *UserModel) GetUidByInvateCode(code string) int {
	return userAuthSvc.GetUIDByInviteCode(code)
}

func (m *UserModel) GetInvateCode() string {
	return userAuthSvc.GetInviteCode()
}

func (m *UserModel) Login(rq *LoginRequest) *LoginResponse {
	return userAuthSvc.Login((*userdomain.LoginRequest)(rq))
}

func (m *UserModel) GoogleAuthLogin(uid int, verdifyCode string, ip string) *LoginResponse { //GOOGLE验证器登陆
	return userAuthSvc.GoogleAuthLogin(uid, verdifyCode, ip)
}

func (m *UserModel) MakeSessionId(uid int) string {
	return userAuthSvc.MakeSessionID(uid)
}

func (m *UserModel) AfterLogin(uid int, clientip string, sid string) *UserBaseInfo {
	return userAuthSvc.AfterLogin(uid, clientip, sid)
}

func (m *UserModel) CheckSessionId(sid string) int {
	return userAuthSvc.CheckSessionID(sid)
}

func (m *UserModel) ChangePassword(uid int, rq *ChangePasswordRequest) *BaseResponse {
	return userAuthSvc.ChangePassword(uid, (*userdomain.ChangePasswordRequest)(rq))
}

func (m *UserModel) GoogleAuth(uid int) map[string]string { //绑定GOOGLE验证器
	return userAuthSvc.GoogleAuth(uid)
}

func (m *UserModel) BindGoogleAuth(uid int, secret string, verdifycode string) *BaseResponse {
	return userAuthSvc.BindGoogleAuth(uid, secret, verdifycode)
}
