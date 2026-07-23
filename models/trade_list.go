package models

func (m *TradeModel) GetDelegateList(uid int, rq *TradeListRequest) *PageBaseResponse {
	return tradeSvc.GetDelegateList(uid, rq)
}

func (m *TradeModel) GetOpendList(uid int, rq *TradeListRequest) *PageBaseResponse {
	return tradeSvc.GetOpendList(uid, rq)
}

func (m *TradeModel) GetCloseList(uid int, rq *TradeListRequest) *PageBaseResponse {
	return tradeSvc.GetCloseList(uid, rq)
}
