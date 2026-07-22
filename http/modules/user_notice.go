package modules

import (
	"cointrade/http/common"
	"cointrade/models"

	"github.com/gin-gonic/gin"
)

func (m *UserModule) noticeRoutes() common.MODULEHANDLELIST {
	return common.MODULEHANDLELIST{
		&common.ModuleHandles{Method: "post", Path: "/user/notice/unread", Handles: common.HandleArray{m.NeedLogin, m.GetNoticeUnRead}},
		&common.ModuleHandles{Method: "post", Path: "/user/notice/list", Handles: common.HandleArray{m.NeedLogin, m.GetNoticeList}},
		&common.ModuleHandles{Method: "post", Path: "/user/notice/detail", Handles: common.HandleArray{m.NeedLogin, m.GetNoticeDetail}},
		&common.ModuleHandles{Method: "post", Path: "/user/notice/read", Handles: common.HandleArray{m.NeedLogin, m.NoticeRead}},
		&common.ModuleHandles{Method: "post", Path: "/user/notice/clear", Handles: common.HandleArray{m.NeedLogin, m.ClearUnreadNotice}},
	}
}

func (m *UserModule) NoticeRead(r *gin.Context) {
	uid := r.GetInt("uid")
	nid := m.GetInt(r, "id")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.ReadNotice(uid, nid))
}

func (m *UserModule) ClearUnreadNotice(r *gin.Context) {
	uid := r.GetInt("uid")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.ClearUnreadNotice(uid))
}

func (m *UserModule) GetNoticeUnRead(r *gin.Context) {
	uid := r.GetInt("uid")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.GetNoticeUnRead(uid))
}

func (m *UserModule) GetNoticeList(r *gin.Context) {
	uid := r.GetInt("uid")
	var rq models.PageBaseRequest
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.GetNoticeList(uid, &rq))
}

func (m *UserModule) GetNoticeDetail(r *gin.Context) {
	uid := r.GetInt("uid")
	nid := m.GetInt(r, "id")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.GetNoticeDetail(uid, nid))
}
