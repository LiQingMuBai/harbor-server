package models

import userdomain "cointrade/internal/domain/user"

// 用户相关的模块
const (
	REGISTER_STATE_ERROREMAIL    = 2001 //错误的邮件
	REGISTER_STATE_ERRORPASSWORD = 2002 //错误的密码
	REGISTER_STATE_ERRORVERDIFY  = 2003 //错误的验证码
	REGISTER_STATE_NOREPASSWORD  = 2004 //两次输入密码不一致
	REGISTER_STATE_INVITEERROR   = 2005 //错误的邀请码
	REGISTER_STATE_EMAILEXISTS   = 2006 //邮件地址已存在
	LOGIN_STATE_GOOGLE_AUTH      = 3001 //二次验证
	LOGIN_STATE_LOCKED           = 3002 //被封禁

	BIND_PHONE_STATE_BINDED    = 4001 //已经绑定过了
	BIND_PHONE_STATE_ERRORCODE = 4002 //错误的验证码

	CHANGE_PASS_STATE_REERROR  = 1001 //两次输入的密码不一样
	CHANGE_PASS_STATE_OLDERROR = 1002 //原密码错误
	USER_MODE_REAL             = 1    //用户模式 真实
	USER_MODE_V                = 2    //用户模式 虚拟
)

type UserModel struct {
	ModelBase
}

type SetCashPasswordRequest = userdomain.SetCashPasswordRequest

type CreditValue = userdomain.CreditValue

type UpdateProfileRequest = userdomain.UpdateProfileRequest

type UserBaseInfo = userdomain.UserBaseInfo

type WelcomeInfo = userdomain.WelcomeInfo

type RegisterRequest = userdomain.RegisterRequest

type UpdateCashPasswordRequest = userdomain.UpdateCashPasswordRequest

type LoginRequest = userdomain.LoginRequest

type LoginResponse = userdomain.LoginResponse

type WelcomeResponse = userdomain.WelcomeResponse

type CrossPlatformTradeResponse = userdomain.CrossPlatformTradeResponse

type AuthLv1Request = userdomain.AuthLv1Request

type UserCount = userdomain.UserCount

type AuthLv2Request = userdomain.AuthLv2Request

type ChangePasswordRequest = userdomain.ChangePasswordRequest
