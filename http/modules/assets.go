package modules

import (
	"cointrade/http/common"
	"cointrade/models"

	"github.com/gin-gonic/gin"
)

//资产相关
type AssetModule struct {
	common.ModuleBase
}

func (m *AssetModule) ModuleList() common.MODULEHANDLELIST {
	//每个MODULE必须要实现的MODULE
	return common.MODULEHANDLELIST{
		&common.ModuleHandles{Method: "post", Path: "/assets/exchange", Handles: common.HandleArray{m.NeedLogin, m.Exchange}},
		&common.ModuleHandles{Method: "post", Path: "/assets/list", Handles: common.HandleArray{m.NeedLogin, m.GetList}},
		&common.ModuleHandles{Method: "post", Path: "/assets/quickexchange", Handles: common.HandleArray{m.NeedLogin, m.QuickExchange}}, //闪兑
	}
}
func (m *AssetModule) Exchange(r *gin.Context) {
	//兑换
	uid := r.GetInt("uid")
	var rq models.ExchangeRequest
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_ASSETS.Exchange(uid, rq.From, rq.To, rq.Amount))
}
func (m *AssetModule) GetList(r *gin.Context) {
	uid := r.GetInt("uid")
	uinfo := models.MODEL_USER.GetBaseInfo(uid)
	if uinfo == nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_ASSETS.GetAllAssets(uid, uinfo.Mode))
}
func (m *AssetModule) QuickExchange(r *gin.Context) {
	uid := r.GetInt("uid")
	var rq models.QuickExchangeRequest
	err := m.ConvertObject(r, &rq)
	if err != nil {
		m.SendResponse(r, common.HTTP_CODE_ERRORPARAM, nil)
		return
	}
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_ASSETS.QuickExchange(uid, &rq))
}
