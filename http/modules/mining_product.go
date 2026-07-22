package modules

import (
	"cointrade/http/common"
	"cointrade/models"

	"github.com/gin-gonic/gin"
)

func (m *MiningModule) productRoutes() common.MODULEHANDLELIST {
	return common.MODULEHANDLELIST{
		&common.ModuleHandles{Method: "post", Path: "/mining/list", Handles: common.HandleArray{m.NeedLogin, m.ListProducts}},
		&common.ModuleHandles{Method: "post", Path: "/mining/detail", Handles: common.HandleArray{m.NeedLogin, m.GetProductDetail}},
	}
}

func (m *MiningModule) ListProducts(r *gin.Context) {
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MINPRODUCT_LIST)
}

func (m *MiningModule) GetProductDetail(r *gin.Context) {
	pid := m.GetInt(r, "pid")
	m.SendResponse(r, common.HTTP_CODE_SUCCESS, models.MODEL_PRODUCT.GetProductInfo(pid))
}
