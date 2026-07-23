package handler

import (
	"cointrade/http/common"
	"cointrade/models"

	"github.com/gin-gonic/gin"
)

func (m *MessageModule) serviceRoutes() common.MODULEHANDLELIST {
	return common.MODULEHANDLELIST{
		&common.ModuleHandles{Method: "post", Path: "/message/service/list", Handles: common.HandleArray{m.NeedLogin, m.ServiceMessageList}},
		&common.ModuleHandles{Method: "post", Path: "/message/service/unread", Handles: common.HandleArray{m.NeedLogin, m.GetServiceUnreadCount}},
		&common.ModuleHandles{Method: "post", Path: "/message/service/unread/clear", Handles: common.HandleArray{m.NeedLogin, m.ClearServiceUnread}},
	}
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
