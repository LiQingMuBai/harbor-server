package module

import (
	"cointrade/http/common"
	"cointrade/models"

	"github.com/gin-gonic/gin"
)

func (m *UserModule) bankRoutes() common.MODULEHANDLELIST {
	return common.MODULEHANDLELIST{
		&common.ModuleHandles{Method: "post", Path: "/user/bank/bind", Handles: common.HandleArray{m.NeedLogin, m.BindBank}},
		&common.ModuleHandles{Method: "post", Path: "/user/bank/info", Handles: common.HandleArray{m.NeedLogin, m.GetBank}},
	}
}

func (m *UserModule) BindBank(r *gin.Context) {
	uid := r.GetInt("uid")
	var rq models.BankInfo
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_CREDIT.BindBank(uid, &rq))
}

func (m *UserModule) GetBank(r *gin.Context) {
	uid := r.GetInt("uid")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_CREDIT.GetBankInfo(uid))
}
