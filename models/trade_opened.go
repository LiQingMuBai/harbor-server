package models

import "cointrade/lib/db"

func (m *TradeModel) GetCloseBySn(uid int, sn string) db.DB_ROW_RESULT {
	return tradeSvc.GetCloseBySn(uid, sn)
}

func (m *TradeModel) GetOpenedBySn(uid int, sn string) *OpenedInfo {
	return tradeSvc.GetOpenedBySn(uid, sn)
}

func (m *TradeModel) GetOpenedOne(uid int, coin string, tradeType int, flag int, mode int, ganggan int) *OpenedInfo {
	return tradeSvc.GetOpenedOne(uid, coin, tradeType, flag, mode, ganggan)
}

func (m *TradeModel) AddKeepOpened(delegateInfo db.DBValues) {
	tradeSvc.AddKeepOpened(delegateInfo)
}
