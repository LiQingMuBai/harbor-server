package module

import (
	"cointrade/http/common"
	"cointrade/models"

	"github.com/gin-gonic/gin"
)

func (m *MessageModule) noticeRoutes() common.MODULEHANDLELIST {
	return common.MODULEHANDLELIST{
		&common.ModuleHandles{Method: "post", Path: "/message/list", Handles: common.HandleArray{m.NeedLogin, m.List}},
		&common.ModuleHandles{Method: "post", Path: "/message/unread", Handles: common.HandleArray{m.NeedLogin, m.UnreadNum}},
		&common.ModuleHandles{Method: "post", Path: "/message/read", Handles: common.HandleArray{m.NeedLogin, m.Read}},
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
