package handler

import (
	"cointrade/http/common"
	userdomain "cointrade/internal/domain/user"
	usermarketingrepo "cointrade/internal/usermarketing/repo"
	usermarketingservice "cointrade/internal/usermarketing/service"
	userportalrepo "cointrade/internal/userportal/repo"
	userportalservice "cointrade/internal/userportal/service"
	"cointrade/lib/db"
	"cointrade/models"

	"github.com/gin-gonic/gin"
)

var apiUserMarketingSvc = usermarketingservice.NewService(
	usermarketingrepo.NewDBRepository(),
	apiUserMarketingGateway{},
)

type apiUserMarketingGateway struct{}

func (apiUserMarketingGateway) GetBaseInfo(uid int) *userdomain.UserBaseInfo {
	return models.MODEL_USER.GetBaseInfo(uid)
}

func (apiUserMarketingGateway) AddCredit(uid int, value *userdomain.CreditValue) bool {
	return models.MODEL_USER.AddCredit(uid, value)
}

func (apiUserMarketingGateway) ClearCache(uid int) {
	models.MODEL_USER.ClearCache(uid)
}

func (apiUserMarketingGateway) EncodePassword(password string) string {
	return models.MODEL_USER.EncodePassword(password)
}

var apiUserPortalSvc = userportalservice.NewService(
	userportalrepo.NewDBRepository(),
)

func (m *UserModule) marketingRoutes() common.MODULEHANDLELIST {
	return common.MODULEHANDLELIST{
		&common.ModuleHandles{Method: "post", Path: "/user/crossTrade", Handles: common.HandleArray{m.CrossTrade}},
		&common.ModuleHandles{Method: "post", Path: "/user/welcome", Handles: common.HandleArray{m.Welcome}},
		&common.ModuleHandles{Method: "post", Path: "/user/isNewUser", Handles: common.HandleArray{m.NeedLogin, m.IsNewUser}},
		&common.ModuleHandles{Method: "post", Path: "/user/claim", Handles: common.HandleArray{m.NeedLogin, m.Claim}},
	}
}

func (m *UserModule) CrossTrade(r *gin.Context) {
	uid := r.GetInt("uid")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, apiUserPortalSvc.CrossTrade(uid, db.DB_PARAMS{
		"amount": m.GetFloat(r, "amount"),
	}))
}

func (m *UserModule) Welcome(r *gin.Context) {
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, apiUserPortalSvc.Welcome())
}

func (m *UserModule) Claim(r *gin.Context) {
	uid := r.GetInt("uid")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, apiUserMarketingSvc.Claim(uid))
}

func (m *UserModule) IsNewUser(r *gin.Context) {
	uid := r.GetInt("uid")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, apiUserMarketingSvc.IsNewUser(uid))
}
