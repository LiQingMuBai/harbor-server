package models

import "cointrade/lib/db"

func (m *TradeModel) GetCloseBySn(uid int, sn string) db.DB_ROW_RESULT {
	return tradeSvc.GetCloseBySn(uid, sn)
}

func (m *TradeModel) GetOpendBySn(uid int, sn string) *OpenedInfo {
	return tradeSvc.GetOpendBySn(uid, sn)
}

func (m *TradeModel) GetOpendOne(uid int, coin string, tradeType int, flag int, mode int, ganggan int) *OpenedInfo {
	return tradeSvc.GetOpendOne(uid, coin, tradeType, flag, mode, ganggan)
}

func (m *TradeModel) AddKeepOpend(delegateInfo db.DBValues) {
	tradeSvc.AddKeepOpend(delegateInfo)
}
