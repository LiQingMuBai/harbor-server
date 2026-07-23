package models

func (m *TradeModel) GetDelegateList(uid int, rq *TradeListRequest) *PageBaseResponse {
	return tradeSvc.GetDelegateList(uid, rq)
}

func (m *TradeModel) GetOpenedList(uid int, rq *TradeListRequest) *PageBaseResponse {
	return tradeSvc.GetOpenedList(uid, rq)
}

func (m *TradeModel) GetCloseList(uid int, rq *TradeListRequest) *PageBaseResponse {
	return tradeSvc.GetCloseList(uid, rq)
}
