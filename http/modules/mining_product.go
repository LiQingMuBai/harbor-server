package modules

import (
	"cointrade/http/common"
	"cointrade/models"

	"github.com/gin-gonic/gin"
)

func (m *MingingModule) productRoutes() common.MODULEHANDLELIST {
	return common.MODULEHANDLELIST{
		&common.ModuleHandles{Method: "post", Path: "/mining/list", Handles: common.HandleArray{m.NeedLogin, m.ProductList}},
		&common.ModuleHandles{Method: "post", Path: "/mining/detail", Handles: common.HandleArray{m.NeedLogin, m.Detail}},
	}
}

func (m *MingingModule) ProductList(r *gin.Context) {
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MINPRODUCT_LIST)
}

func (m *MingingModule) Detail(r *gin.Context) {
	pid := m.GetInt(r, "pid")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_PRODUCT.GetProductInfo(pid))
}
