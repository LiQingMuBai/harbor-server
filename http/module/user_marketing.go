package module

import (
	"cointrade/http/common"
	"cointrade/models"

	"github.com/gin-gonic/gin"
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
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.CrossTrade(uid, nil))
}

func (m *UserModule) Welcome(r *gin.Context) {
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.Welcome())
}

func (m *UserModule) Claim(r *gin.Context) {
	uid := r.GetInt("uid")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.Claim(uid))
}

func (m *UserModule) IsNewUser(r *gin.Context) {
	uid := r.GetInt("uid")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.IsNewUser(uid))
}
