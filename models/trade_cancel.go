package models

func (m *TradeModel) CancleDelegate(uid int, sn string) *BaseResponse {
	return tradeSvc.CancleDelegate(uid, sn)
}
