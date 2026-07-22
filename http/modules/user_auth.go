package modules

import (
	"cointrade/http/common"
	"cointrade/models"

	"github.com/gin-gonic/gin"
)

func (m *UserModule) authRoutes() common.MODULEHANDLELIST {
	return common.MODULEHANDLELIST{
		&common.ModuleHandles{Method: "post", Path: "/user/login", Handles: common.HandleArray{m.Login}},
		&common.ModuleHandles{Method: "post", Path: "/user/register", Handles: common.HandleArray{m.Register}},
		&common.ModuleHandles{Method: "post", Path: "/user/register/code", Handles: common.HandleArray{m.SendRegisterCode}},
		&common.ModuleHandles{Method: "post", Path: "/user/auth1", Handles: common.HandleArray{m.NeedLogin, m.AuthLv1}},
		&common.ModuleHandles{Method: "post", Path: "/user/auth2", Handles: common.HandleArray{m.NeedLogin, m.AuthLv2}},
		&common.ModuleHandles{Method: "post", Path: "/user/changemode", Handles: common.HandleArray{m.NeedLogin, m.ChangeMode}},
		&common.ModuleHandles{Method: "post", Path: "/user/login/wallet", Handles: common.HandleArray{m.LoginByWallet}},
	}
}

func (m *UserModule) SendRegisterCode(r *gin.Context) {
	email := m.GetValue(r, "email")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_CODE.SendEmailCodeRegister(email))
}

func (m *UserModule) Login(r *gin.Context) {
	var rq models.LoginRequest
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	rq.ClientIp = r.ClientIP()

	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.Login(&rq))
}

func (m *UserModule) Register(r *gin.Context) {
	var rq models.RegisterRequest
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	rq.ClientIp = r.ClientIP()
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.Register(&rq))
}

func (m *UserModule) AuthLv1(r *gin.Context) {
	var rq models.AuthLv1Request
	uid := r.GetInt("uid")
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.AuthLv1(uid, &rq))
}

func (m *UserModule) AuthLv2(r *gin.Context) {
	var rq models.AuthLv2Request
	uid := r.GetInt("uid")
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.AuthLv2(uid, &rq))
}

func (m *UserModule) ChangeMode(r *gin.Context) {
	uid := r.GetInt("uid")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.ChangeMode(uid))
}

func (m *UserModule) LoginByWallet(r *gin.Context) {
	address := m.GetValue(r, "address")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.LoginByAddress(address, r.ClientIP()))
}
