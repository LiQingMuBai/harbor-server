package models

import (
	productdomain "cointrade/internal/domain/product"
	productrepo "cointrade/internal/product/repo"
	productservice "cointrade/internal/product/service"
	"cointrade/lib/db"
)

// 质押挖矿 挖矿没有虚拟模式
type ProductModel struct {
	ModelBase
}

const (
	BUY_STATE_CREDIT    = productdomain.BUY_STATE_CREDIT
	BUY_STATE_MIN       = productdomain.BUY_STATE_MIN
	BUY_STATE_MAX       = productdomain.BUY_STATE_MAX
	BUY_STATE_PER_LIMIT = productdomain.BUY_STATE_PER_LIMIT
	BUY_V_GETED         = productdomain.BUY_V_GETED
	BUY_STATE_LOCKED    = productdomain.BUY_STATE_LOCKED

	MTYPE_V = productdomain.MTYPE_V
	MTYPE_R = productdomain.MTYPE_R
)

type ProductProfile = productdomain.ProductProfile

type ProductInfo = productdomain.ProductInfo

type BuyRequest = productdomain.BuyRequest

type OrderListRequest = productdomain.OrderListRequest

type productUserGateway struct{}

func (productUserGateway) GetBaseInfo(uid int) *UserBaseInfo {
	return MODEL_USER.GetBaseInfo(uid)
}

func (productUserGateway) AddCredit(uid int, value *CreditValue) bool {
	return MODEL_USER.AddCredit(uid, value)
}

type cachedProductCatalog struct{}

func (cachedProductCatalog) GetProductList() []*ProductInfo {
	return MINPRODUCT_LIST
}

var productSvc = productservice.New(
	productrepo.NewDBRepository(),
	productUserGateway{},
	cachedProductCatalog{},
)

func (m *ProductModel) GetAcceptPids(uid int) map[int]int { //获取用户的授权
	return productSvc.GetAcceptPids(uid)
}
func (m *ProductModel) CheckProcutAccept(uid int, pid int) bool {
	return productSvc.CheckProductAccept(uid, pid)
}
func (m *ProductModel) MakeSn(uid int) string {
	return productSvc.MakeSn(uid)
}
func (m *ProductModel) GetProductList() []*ProductInfo {
	return productSvc.GetProductList()
}
func (m *ProductModel) Buy(uid int, rq *BuyRequest) *BaseResponse {
	return productSvc.Buy(uid, rq)
}
func (m *ProductModel) GetROrder(uid int) []int { //获取用户正在预约中的产品ID
	return productSvc.GetROrder(uid)
}

func (m *ProductModel) GetProductInfo(pid int) *ProductInfo {
	return productSvc.GetProductInfo(pid)
}
func (m *ProductModel) GetOrderList(uid int, rq OrderListRequest) *PageBaseResponse { //取得矿机列表
	return productSvc.GetOrderList(uid, rq)
}
func (m *ProductModel) Unlock(uid int, sn string) *BaseResponse {
	return productSvc.Unlock(uid, sn)
}

func (m *ProductModel) GetOrderInfo(uid int, sn string) db.DB_ROW_RESULT { //获得单个矿机订单信息
	one := productSvc.GetOrderInfo(uid, sn)
	if one == nil {
		return nil
	}
	rs := make(db.DB_ROW_RESULT)
	for key, value := range one {
		rs[key] = value.ToString()
	}
	return rs
}
func (m *ProductModel) GetOrderCount(uid int) map[string]interface{} {
	return productSvc.GetOrderCount(uid)
}
