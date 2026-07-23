package handler

import (
	"cointrade/http/common"
	adminservice "cointrade/internal/admin/service"
	"cointrade/models"

	"github.com/gin-gonic/gin"
)

func (m *AdminUserModule) userRoutes() common.MODULEHANDLELIST {
	return common.MODULEHANDLELIST{
		&common.ModuleHandles{Method: "post", Path: "/admin/userlist", Handles: common.HandleArray{m.CheckLogin, m.UserList}},
		&common.ModuleHandles{Method: "post", Path: "/admin/op_user", Handles: common.HandleArray{m.CheckLogin, m.SaveUser}},
		&common.ModuleHandles{Method: "post", Path: "/admin/usercoinlog", Handles: common.HandleArray{m.CheckLogin, m.UserCoinLog}},
		&common.ModuleHandles{Method: "post", Path: "/admin/coin/change", Handles: common.HandleArray{m.CheckLogin, m.ChangeCoin}},
		&common.ModuleHandles{Method: "post", Path: "/admin/user_asset", Handles: common.HandleArray{m.CheckLogin, m.UserAsset}},
		&common.ModuleHandles{Method: "post", Path: "/admin/userinfo", Handles: common.HandleArray{m.CheckLogin, m.UserInfo}},
		&common.ModuleHandles{Method: "post", Path: "/admin/user_levelcount", Handles: common.HandleArray{m.CheckLogin, m.UserLevelCount}},
		&common.ModuleHandles{Method: "post", Path: "admin/agent_count", Handles: common.HandleArray{m.CheckLogin, m.AgentCount}},
		&common.ModuleHandles{Method: "post", Path: "admin/agent_list", Handles: common.HandleArray{m.CheckLogin, m.AgentList}},
		&common.ModuleHandles{Method: "post", Path: "admin/employer_list", Handles: common.HandleArray{m.CheckLogin, m.EmployerList}},
		&common.ModuleHandles{Method: "post", Path: "admin/del_agent", Handles: common.HandleArray{m.CheckLogin, m.DeleteAgent}},
		&common.ModuleHandles{Method: "post", Path: "admin/userexplode_controller", Handles: common.HandleArray{m.CheckLogin, m.UserControllerExp}},
		&common.ModuleHandles{Method: "post", Path: "admin/asset", Handles: common.HandleArray{m.CheckLogin, m.Assets}},
		&common.ModuleHandles{Method: "post", Path: "admin/save_parentmemo", Handles: common.HandleArray{m.CheckLogin, m.SaveParentMemo}},
		&common.ModuleHandles{Method: "post", Path: "admin/kick_user", Handles: common.HandleArray{m.CheckLogin, m.KickUser}},
	}
}

func (m *AdminUserModule) UserInfo(r *gin.Context) {
	uid := m.GetInt(r, "uid")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, &adminservice.AdminResponse{
		State: adminservice.SUCCESS,
		Data:  models.MODEL_USER.GetBaseInfo(uid),
	})
}

func (m *AdminUserModule) SaveParentMemo(r *gin.Context) {
	rq := make(adminservice.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.MODEL_USER.SaveParentMemo(rq))
}

func (m *AdminUserModule) Assets(r *gin.Context) {
	uid := m.GetInt(r, "uid")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.MODEL_USER.GetUserAssetOverview(uid))
}

func (m *AdminUserModule) UserControllerExp(r *gin.Context) {
	rq := make(adminservice.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.MODEL_USER.UserControllerExp(rq))
}

func (m *AdminUserModule) DeleteAgent(r *gin.Context) {
	id := m.GetInt(r, "id")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.MODEL_AGENT.DeleteAgent(id))
}

func (m *AdminUserModule) UserCoinLog(r *gin.Context) {
	rq := make(adminservice.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.MODEL_USER.UserCoinLog(rq))
}

func (m *AdminUserModule) EmployerList(r *gin.Context) {
	rq := make(adminservice.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.MODEL_AGENT.EmpoyerList(rq))
}

func (m *AdminUserModule) AgentList(r *gin.Context) {
	rq := make(adminservice.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.MODEL_AGENT.AgentList(rq))
}

func (m *AdminUserModule) AgentCount(r *gin.Context) {
	rq := make(adminservice.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.MODEL_AGENT.AgentCountList(rq))
}

func (m *AdminUserModule) ChangeCoin(r *gin.Context) {
	rq := make(adminservice.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.MODEL_USER.AdjustUserCredit(rq))
}

func (m *AdminUserModule) KickUser(r *gin.Context) {
	uid := m.GetInt(r, "uid")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.SYSTEM_MODEL.KickUser(uid))
}

func (m *AdminUserModule) SaveUser(r *gin.Context) {
	rq := make(adminservice.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.MODEL_USER.SaveUser(rq))
}

func (m *AdminUserModule) UserLevelCount(r *gin.Context) {
	rq := make(adminservice.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.MODEL_USER.UserTeamLevelCount(rq))
}

func (m *AdminUserModule) UserAsset(r *gin.Context) {
	rq := make(adminservice.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.MODEL_USER.UserAssetList(rq))
}

func (m *AdminUserModule) UserList(r *gin.Context) {
	rq := make(adminservice.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminservice.MODEL_USER.UserList(rq))
}
