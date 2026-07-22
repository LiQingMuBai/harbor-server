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
		&common.ModuleHandles{Method: "get", Path: "/admin/token_info", Handles: common.HandleArray{m.GetTokenInfo}},
		&common.ModuleHandles{Method: "post", Path: "/admin/admin_list", Handles: common.HandleArray{m.CheckLogin, m.ListAdmins}},
		&common.ModuleHandles{Method: "post", Path: "/admin/add_manage", Handles: common.HandleArray{m.CheckLogin, m.SaveAdmin}},
		&common.ModuleHandles{Method: "post", Path: "/admin/mean_list", Handles: common.HandleArray{m.CheckLogin, m.GetMenuTree}},
		&common.ModuleHandles{Method: "get", Path: "/admin/mean_router", Handles: common.HandleArray{m.CheckLogin, m.GetMenuTree}},
		&common.ModuleHandles{Method: "get", Path: "/admin/auth_router", Handles: common.HandleArray{m.GetRoleMenuTree}},
		&common.ModuleHandles{Method: "post", Path: "/admin/role_list", Handles: common.HandleArray{m.CheckLogin, m.RoleList}},
		&common.ModuleHandles{Method: "post", Path: "/admin/handler_role", Handles: common.HandleArray{m.CheckLogin, m.SaveRole}},
		&common.ModuleHandles{Method: "post", Path: "/admin/handler_mean", Handles: common.HandleArray{m.CheckLogin, m.SaveMenu}},
		&common.ModuleHandles{Method: "post", Path: "/admin/del_mean", Handles: common.HandleArray{m.CheckLogin, m.DeleteMenu}},
		&common.ModuleHandles{Method: "post", Path: "admin/del_admin", Handles: common.HandleArray{m.CheckLogin, m.DeleteAdmin}},
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

func (m *AdminUserModule) GetTokenInfo(r *gin.Context) {
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.GetTokenInfo(r.GetString("sid")))
}

func (m *AdminUserModule) SaveAdmin(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.SaveAdmin(rq))
}

func (m *AdminUserModule) GetMenuTree(r *gin.Context) {
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.GetMenuTree())
}

func (m *AdminUserModule) RoleList(r *gin.Context) {
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.RoleList())
}

func (m *AdminUserModule) GetRoleMenuTree(r *gin.Context) {
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.GetRoleMenuTree())
}

func (m *AdminUserModule) ListAdmins(r *gin.Context) {
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.ListAdmins())
}

func (m *AdminUserModule) DeleteAdmin(r *gin.Context) {
	id := m.GetInt(r, "id")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.DeleteAdmin(id))
}

func (m *AdminUserModule) DeleteMenu(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.DeleteMenu(rq))
}

func (m *AdminUserModule) SaveMenu(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.SaveMenu(rq))
}

func (m *AdminUserModule) SaveRole(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.SaveRole(rq))
}

func (m *AdminUserModule) CheckLogin(r *gin.Context) {
	sid := r.GetString("sid")
	u := adminmodels.MODEL_USER.GetAdminBySession(sid)
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
