package adminmodules

import (
	adminmodels "cointrade/admin_models"
	"cointrade/config"
	"cointrade/http/common"
	"cointrade/lib/db"
	"cointrade/utils"
	"regexp"

	"github.com/gin-gonic/gin"
)

func (m *AdminUserModule) authRoutes() common.MODULEHANDLELIST {
	return common.MODULEHANDLELIST{
		&common.ModuleHandles{Method: "post", Path: "/admin/login", Handles: common.HandleArray{m.Login}},
		&common.ModuleHandles{Method: "post", Path: "/admin/logout", Handles: common.HandleArray{m.Logout}},
		&common.ModuleHandles{Method: "get", Path: "/admin/token_info", Handles: common.HandleArray{m.TokenSid}},
		&common.ModuleHandles{Method: "post", Path: "/admin/admin_list", Handles: common.HandleArray{m.CheckLogin, m.AdminList}},
		&common.ModuleHandles{Method: "post", Path: "/admin/add_manage", Handles: common.HandleArray{m.CheckLogin, m.AddAdmin}},
		&common.ModuleHandles{Method: "post", Path: "/admin/mean_list", Handles: common.HandleArray{m.CheckLogin, m.MeanRouter}},
		&common.ModuleHandles{Method: "get", Path: "/admin/mean_router", Handles: common.HandleArray{m.CheckLogin, m.MeanRouter}},
		&common.ModuleHandles{Method: "get", Path: "/admin/auth_router", Handles: common.HandleArray{m.AuthRouter}},
		&common.ModuleHandles{Method: "post", Path: "/admin/role_list", Handles: common.HandleArray{m.CheckLogin, m.RoleList}},
		&common.ModuleHandles{Method: "post", Path: "/admin/handler_role", Handles: common.HandleArray{m.CheckLogin, m.HandlerRole}},
		&common.ModuleHandles{Method: "post", Path: "/admin/handler_mean", Handles: common.HandleArray{m.CheckLogin, m.HandlerMean}},
		&common.ModuleHandles{Method: "post", Path: "/admin/del_mean", Handles: common.HandleArray{m.CheckLogin, m.DelMean}},
		&common.ModuleHandles{Method: "post", Path: "admin/del_admin", Handles: common.HandleArray{m.CheckLogin, m.DelAdmin}},
	}
}

func (m *AdminUserModule) Login(r *gin.Context) {
	var request adminmodels.LoginRequest
	m.ConvertObject(r, &request)
	request.ClientIp = r.ClientIP()
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.Login(request))
}

func (m *AdminUserModule) Logout(r *gin.Context) {
	sid := r.GetString("sid")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.AdminResponse{State: adminmodels.SUCCESS, Data: adminmodels.MODEL_USER.Logout(sid)})
}

func (m *AdminUserModule) TokenSid(r *gin.Context) {
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.TokenInfo(r.GetString("sid")))
}

func (m *AdminUserModule) AddAdmin(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.AddAdmin(rq))
}

func (m *AdminUserModule) MeanRouter(r *gin.Context) {
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.MeanRouter())
}

func (m *AdminUserModule) RoleList(r *gin.Context) {
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.RoleList())
}

func (m *AdminUserModule) AuthRouter(r *gin.Context) {
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.AuthRouter())
}

func (m *AdminUserModule) AdminList(r *gin.Context) {
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.AdminList())
}

func (m *AdminUserModule) DelAdmin(r *gin.Context) {
	id := m.GetInt(r, "id")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.DelAdmin(id))
}

func (m *AdminUserModule) DelMean(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.DelMean(rq))
}

func (m *AdminUserModule) HandlerMean(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.HandlerMean(rq))
}

func (m *AdminUserModule) HandlerRole(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.HandlerRole(rq))
}

func (m *AdminUserModule) CheckLogin(r *gin.Context) {
	sid := r.GetString("sid")
	u := adminmodels.MODEL_USER.SidInfo(sid)
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	if u == nil {
		m.SendResponse(r, common.HTTP_CODE_SUCCESS, &adminmodels.AdminResponse{
			State: 50008,
			Data:  "当前用户的sid无法获取到值",
		})
		r.Abort()
	}
	notify, _ := regexp.Match(`admin\/notify_list`, []byte(r.Request.URL.Path))
	if notify {
		return
	}
	config.GlobalDB.InsertData("operation_log", db.DB_PARAMS{
		"username":   u.UserName,
		"path":       r.Request.URL.Path,
		"createtime": utils.GetNow(),
		"data":       rq,
		"ip":         r.ClientIP(),
	})
}
