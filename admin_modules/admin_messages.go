package adminmodules

import (
	adminmodels "cointrade/admin_models"
	"cointrade/http/common"

	"github.com/gin-gonic/gin"
)

func (m *AdminUserModule) messageRoutes() common.MODULEHANDLELIST {
	return common.MODULEHANDLELIST{
		&common.ModuleHandles{Method: "post", Path: "/admin/chatlist", Handles: common.HandleArray{m.CheckLogin, m.Chatlist}},
		&common.ModuleHandles{Method: "post", Path: "/admin/custom_msg", Handles: common.HandleArray{m.CheckLogin, m.UserByMessage}},
		&common.ModuleHandles{Method: "post", Path: "/admin/send_msg", Handles: common.HandleArray{m.CheckLogin, m.SendMsg}},
		&common.ModuleHandles{Method: "post", Path: "admin/del_msg", Handles: common.HandleArray{m.CheckLogin, m.Delmsg}},
		&common.ModuleHandles{Method: "post", Path: "admin/msg_list", Handles: common.HandleArray{m.CheckLogin, m.MsgList}},
		&common.ModuleHandles{Method: "post", Path: "admin/send_user_notice", Handles: common.HandleArray{m.CheckLogin, m.SendUserNotice}},
		&common.ModuleHandles{Method: "post", Path: "admin/del_usernotice", Handles: common.HandleArray{m.CheckLogin, m.DelUserNotice}},
	}
}

func (m *AdminUserModule) DelUserNotice(r *gin.Context) {
	id := m.GetInt(r, "id")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.DelUserNotice(id))
}

func (m *AdminUserModule) SendUserNotice(r *gin.Context) {
	rq := new(adminmodels.UserNoticeMsg)
	m.ConvertObject(r, rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.SendUserNotice(rq))
}

func (m *AdminUserModule) MsgList(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.MsgList(rq))
}

func (m *AdminUserModule) Delmsg(r *gin.Context) {
	id := m.GetValue(r, "sn_id")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.SYSTEM_MODEL.Delmsg(id))
}

func (m *AdminUserModule) SendMsg(r *gin.Context) {
	rq := new(adminmodels.CustomMsg)
	m.ConvertObject(r, rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.SendMsg(rq))
}

func (m *AdminUserModule) UserByMessage(r *gin.Context) {
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.UserByMessage(m.GetInt(r, "uid")))
}

func (m *AdminUserModule) Chatlist(r *gin.Context) {
	rq := make(adminmodels.P, 0)
	m.ConvertObject(r, &rq)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, adminmodels.MODEL_USER.CustomServiceList(rq))
}
