package modules

import (
	"cointrade/config"
	"cointrade/http/common"
	"cointrade/models"
	"cointrade/utils"
	"fmt"
	"math/rand"

	"github.com/gin-gonic/gin"
)

func (m *SystemModule) contentRoutes() common.MODULEHANDLELIST {
	return common.MODULEHANDLELIST{
		&common.ModuleHandles{Method: "post", Path: "/system/notice", Handles: common.HandleArray{m.Notice}},
		&common.ModuleHandles{Method: "post", Path: "/system/notice/detail", Handles: common.HandleArray{m.NoticeDetail}},
		&common.ModuleHandles{Method: "post", Path: "/system/rule", Handles: common.HandleArray{m.NeedLogin, m.RuleText}},
		&common.ModuleHandles{Method: "post", Path: "/system/rule/detail", Handles: common.HandleArray{m.NeedLogin, m.RuleDetail}},
		&common.ModuleHandles{Method: "get", Path: "/system/test", Handles: common.HandleArray{m.Test}},
	}
}

func (m *SystemModule) Test(r *gin.Context) {
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, "dsadsa123")
}

func (m *SystemModule) NoticeDetail(r *gin.Context) {
	id := m.GetInt(r, "id")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_NOTICE.GetOne(id))
}

func (m *SystemModule) RuleDetail(r *gin.Context) {
	rs := make(map[string]interface{}, 0)
	m.ConvertObject(r, &rs)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_SYSTEM.RuleOne(rs))
}

func (m *SystemModule) Notice(r *gin.Context) {
	var rq models.NoticeRequest
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	rq.Limit = 15
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_NOTICE.GetList(&rq))
}

func (m *SystemModule) TestService(r *gin.Context) {
	uid := r.GetInt("uid")
	config.GlobalRedis.PushQueue(models.HASH_USER_SERVICE_MESSAGE, models.DBServiceMessage{
		Uid:        uid,
		Content:    fmt.Sprintf("测试消息%d", 1000+rand.Intn(9000)),
		CreateTime: utils.GetNow(),
		Flag:       2,
	})
}

func (m *SystemModule) RuleText(r *gin.Context) {
	ruletype := m.GetValue(r, "ruletype")
	lang := m.GetValue(r, "lang")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_SYSTEM.RuleText(ruletype, lang))
}
