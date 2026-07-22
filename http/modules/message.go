package modules

import (
	"cointrade/http/common"
	"cointrade/models"

	"github.com/gin-gonic/gin"
)

type MessageModule struct {
	common.ModuleBase
}

func (m *MessageModule) ModuleList() common.MODULEHANDLELIST {
	//每个MODULE必须要实现的MODULE
	return common.MODULEHANDLELIST{
		&common.ModuleHandles{Method: "post", Path: "/message/list", Handles: common.HandleArray{m.NeedLogin, m.List}},
		&common.ModuleHandles{Method: "post", Path: "/message/unread", Handles: common.HandleArray{m.NeedLogin, m.UnreadNum}},
		&common.ModuleHandles{Method: "post", Path: "/message/read", Handles: common.HandleArray{m.NeedLogin, m.Read}},
		&common.ModuleHandles{Method: "post", Path: "/message/service/list", Handles: common.HandleArray{m.NeedLogin, m.ServiceMessageList}},
		&common.ModuleHandles{Method: "post", Path: "/message/service/unread", Handles: common.HandleArray{m.NeedLogin, m.GetServiceUnreadCount}},
		&common.ModuleHandles{Method: "post", Path: "/message/service/unread/clear", Handles: common.HandleArray{m.NeedLogin, m.ClearServiceUnread}},
	}
}
func (m *MessageModule) List(r *gin.Context) {
	uid := r.GetInt("uid")
	var rq models.PageBaseRequest
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_MESSAGE.GetList(uid, &rq))
}
func (m *MessageModule) UnreadNum(r *gin.Context) {
	uid := r.GetInt("uid")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_MESSAGE.GetUnreadNum(uid))
}
func (m *MessageModule) Read(r *gin.Context) {
	uid := r.GetInt("uid")
	id := m.GetInt(r, "id")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_MESSAGE.ChangeState(uid, id))
}
func (m *MessageModule) ServiceMessageList(r *gin.Context) {
	uid := r.GetInt("uid")
	var rq models.PageBaseRequest
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_MESSAGE.GetServiceList(uid, &rq))
}
func (m *MessageModule) ClearServiceUnread(r *gin.Context) {
	uid := r.GetInt("uid")

	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_MESSAGE.ClearServiceUnread(uid))
}
func (m *MessageModule) GetServiceUnreadCount(r *gin.Context) {
	uid := r.GetInt("uid")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_MESSAGE.GetServiceUnreadCount(uid))
}
