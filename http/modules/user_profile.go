package modules

import (
	"cointrade/http/common"
	"cointrade/models"
	"fmt"

	"github.com/gin-gonic/gin"
)

func (m *UserModule) profileRoutes() common.MODULEHANDLELIST {
	return common.MODULEHANDLELIST{
		&common.ModuleHandles{Method: "post", Path: "/user/userinfo", Handles: common.HandleArray{m.NeedLogin, m.GetUserInfo}},
		&common.ModuleHandles{Method: "post", Path: "/user/usercount", Handles: common.HandleArray{m.NeedLogin, m.GetUserCount}},
		&common.ModuleHandles{Method: "post", Path: "/user/profile", Handles: common.HandleArray{m.NeedLogin, m.UpdateProfile}},
		&common.ModuleHandles{Method: "post", Path: "/user/convertmoney", Handles: common.HandleArray{m.NeedLogin, m.ConvertMoney}},
		&common.ModuleHandles{Method: "post", Path: "/user/einfo", Handles: common.HandleArray{m.NeedLogin, m.GetExplodeState}},
	}
}

func (m *UserModule) ConvertMoney(r *gin.Context) {
	uid := r.GetInt("uid")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.ConvertMoney(uid))
}

func (m *UserModule) GetExplodeState(r *gin.Context) {
	uid := r.GetInt("uid")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.GetExplodeState(uid))
}

func (m *UserModule) UpdateProfile(r *gin.Context) {
	uid := r.GetInt("uid")
	var rq models.UpdateProfileRequest
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_USER.UpdatePorfile(uid, &rq))
}

func (m *UserModule) Update(r *gin.Context) {
	uid := r.GetInt("uid")
	data, ok := r.Get("data")
	if !ok {
		fmt.Println("no data exists")
		return
	}
	fmt.Println(data)
	models.MODEL_USER.Update(uid, data.(map[string]interface{}))
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, "ok")
}

func (m *UserModule) GetUserInfo(r *gin.Context) {
	uid := r.GetInt("uid")
	rs := models.MODEL_USER.GetBaseInfo(uid)
	rs.CashPassword = ""
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, rs)
}

func (m *UserModule) GetUserCount(r *gin.Context) {
	uid := r.GetInt("uid")
	t := m.GetInt(r, "type")
	rs := models.MODEL_USER.GetUserCount(uid, t)
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, rs)
}
