package handler

import (
	"cointrade/http/common"
	userauthrepo "cointrade/internal/userauth/repo"
	userauthservice "cointrade/internal/userauth/service"
	userdomain "cointrade/internal/domain/user"
	useridentityrepo "cointrade/internal/useridentity/repo"
	useridentityservice "cointrade/internal/useridentity/service"
	"cointrade/lib/db"
	"cointrade/lib/notify"
	"cointrade/models"

	"github.com/gin-gonic/gin"
)

var apiUserAuthSvc = userauthservice.NewService(
	userauthrepo.NewDBRepository(),
	apiUserAuthUserGateway{},
	apiUserAuthWalletGateway{},
)

type apiUserAuthUserGateway struct{}

func (apiUserAuthUserGateway) GetBaseInfo(uid int) *userdomain.UserBaseInfo {
	return models.MODEL_USER.GetBaseInfo(uid)
}

func (apiUserAuthUserGateway) Update(uid int, data db.DB_PARAMS) {
	models.MODEL_USER.Update(uid, data)
}

func (apiUserAuthUserGateway) ClearCache(uid int) {
	models.MODEL_USER.ClearCache(uid)
}

type apiUserAuthWalletGateway struct{}

func (apiUserAuthWalletGateway) RegisterByAddress(address string, ip string) int {
	return models.MODEL_USER.RegisterByAddress(address, ip)
}

var apiUserIdentitySvc = useridentityservice.NewService(
	useridentityrepo.NewDBRepository(),
	apiUserIdentityUserGateway{},
	apiUserIdentityNotifier{},
)

type apiUserIdentityUserGateway struct{}

func (apiUserIdentityUserGateway) GetBaseInfo(uid int) *userdomain.UserBaseInfo {
	return models.MODEL_USER.GetBaseInfo(uid)
}

func (apiUserIdentityUserGateway) Update(uid int, data db.DB_PARAMS) {
	models.MODEL_USER.Update(uid, data)
}

type apiUserIdentityNotifier struct{}

func (apiUserIdentityNotifier) IncrementNotify(typ int, num int) {
	notify.NOTIFY.AddNotify(&notify.NotifyItem{Type: typ, Num: num})
}

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

	m.SendResponse(r, common.HTTP_CODE_SUCCESS, apiUserAuthSvc.Login((*userdomain.LoginRequest)(&rq)))
}

func (m *UserModule) Register(r *gin.Context) {
	var rq models.RegisterRequest
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	rq.ClientIp = r.ClientIP()
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, apiUserAuthSvc.Register((*userdomain.RegisterRequest)(&rq)))
}

func (m *UserModule) AuthLv1(r *gin.Context) {
	var rq models.AuthLv1Request
	uid := r.GetInt("uid")
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, apiUserIdentitySvc.AuthLv1(uid, (*userdomain.AuthLv1Request)(&rq)))
}

func (m *UserModule) AuthLv2(r *gin.Context) {
	var rq models.AuthLv2Request
	uid := r.GetInt("uid")
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, apiUserIdentitySvc.AuthLv2(uid, (*userdomain.AuthLv2Request)(&rq)))
}

func (m *UserModule) ChangeMode(r *gin.Context) {
	uid := r.GetInt("uid")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.ChangeMode(uid))
}

func (m *UserModule) LoginByWallet(r *gin.Context) {
	address := m.GetValue(r, "address")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, apiUserAuthSvc.LoginByAddress(address, r.ClientIP()))
}
