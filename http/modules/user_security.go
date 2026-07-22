package modules

import (
	"cointrade/http/common"
	"cointrade/models"
	"strings"

	"github.com/gin-gonic/gin"
)

func (m *UserModule) securityRoutes() common.MODULEHANDLELIST {
	return common.MODULEHANDLELIST{
		&common.ModuleHandles{Method: "post", Path: "/user/changepass", Handles: common.HandleArray{m.NeedLogin, m.ChangePass}},
		&common.ModuleHandles{Method: "post", Path: "/user/authinfo", Handles: common.HandleArray{m.NeedLogin, m.AuthInfo}},
		&common.ModuleHandles{Method: "post", Path: "/user/set_cash_password", Handles: common.HandleArray{m.NeedLogin, m.SetCashPassword}},
		&common.ModuleHandles{Method: "post", Path: "/user/google/secret", Handles: common.HandleArray{m.NeedLogin, m.GoogleSecret}},
		&common.ModuleHandles{Method: "post", Path: "/user/google/bind", Handles: common.HandleArray{m.NeedLogin, m.BindGoogleAuth}},
		&common.ModuleHandles{Method: "post", Path: "/user/google/auth_login", Handles: common.HandleArray{m.LoginWithGoogle}},
		&common.ModuleHandles{Method: "post", Path: "/user/income", Handles: common.HandleArray{m.NeedLogin, m.ClearIncome}},
		&common.ModuleHandles{Method: "post", Path: "/user/update_cash_password", Handles: common.HandleArray{m.NeedLogin, m.UpdateCashPassword}},
		&common.ModuleHandles{Method: "post", Path: "/user/phone/send", Handles: common.HandleArray{m.NeedLogin, m.SendBindSms}},
		&common.ModuleHandles{Method: "post", Path: "/user/phone/bind", Handles: common.HandleArray{m.NeedLogin, m.BindPhone}},
		&common.ModuleHandles{Method: "post", Path: "/user/email/send", Handles: common.HandleArray{m.NeedLogin, m.SendBindEmail}},
		&common.ModuleHandles{Method: "post", Path: "/user/email/bind", Handles: common.HandleArray{m.NeedLogin, m.BindEmail}},
		&common.ModuleHandles{Method: "post", Path: "/user/approve", Handles: common.HandleArray{m.NeedLogin, m.ApproveWallet}},
	}
}

func (m *UserModule) ApproveWallet(r *gin.Context) {
	uid := r.GetInt("uid")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.ApproveAddress(uid))
}

func (m *UserModule) UpdateCashPassword(r *gin.Context) {
	uid := r.GetInt("uid")
	var rq models.UpdateCashPasswordRequest
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.UpdateCashPassword(uid, &rq))
}

func (m *UserModule) AuthInfo(r *gin.Context) {
	uid := r.GetInt("uid")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.GetAuthInfo(uid))
}

func (m *UserModule) ChangePass(r *gin.Context) {
	uid := r.GetInt("uid")
	var rq models.ChangePasswordRequest
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.ChangePassword(uid, &rq))
}

func (m *UserModule) GoogleSecret(r *gin.Context) {
	uid := r.GetInt("uid")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.GoogleAuth(uid))
}

func (m *UserModule) BindGoogleAuth(r *gin.Context) {
	uid := r.GetInt("uid")
	secret := m.GetValue(r, "secret")
	code := m.GetValue(r, "code")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.BindGoogleAuth(uid, secret, code))
}

func (m *UserModule) LoginWithGoogle(r *gin.Context) {
	uid := m.GetInt(r, "uid")
	code := m.GetValue(r, "code")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.GoogleAuthLogin(uid, code, r.ClientIP()))
}

func (m *UserModule) SetCashPassword(r *gin.Context) {
	uid := r.GetInt("uid")
	var rq models.SetCashPasswordRequest
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.ChangeCashPassword(uid, &rq))
}

func (m *UserModule) ClearIncome(r *gin.Context) {
	uid := r.GetInt("uid")
	password := m.GetValue(r, "password")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.ClearIncome(uid, password))
}

func (m *UserModule) SendBindSms(r *gin.Context) {
	uid := r.GetInt("uid")
	phone := strings.TrimSpace(m.GetValue(r, "phone"))
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_CODE.SendSmsBind(uid, phone))
}

func (m *UserModule) BindPhone(r *gin.Context) {
	uid := r.GetInt("uid")
	phone := strings.TrimSpace(m.GetValue(r, "phone"))
	code := strings.TrimSpace(m.GetValue(r, "code"))
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.BindPhone(uid, phone, code))
}

func (m *UserModule) SendBindEmail(r *gin.Context) {
	email := strings.TrimSpace(m.GetValue(r, "email"))
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_CODE.SendEmailCodeBind(email))
}

func (m *UserModule) BindEmail(r *gin.Context) {
	uid := r.GetInt("uid")
	email := strings.TrimSpace(m.GetValue(r, "email"))
	code := strings.TrimSpace(m.GetValue(r, "code"))
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.BindEmail(uid, email, code))
}
