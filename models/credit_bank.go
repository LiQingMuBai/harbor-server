package models

import (
	"cointrade/config"
	"cointrade/lib/db"
)

func (m *CreditModel) BindBank(uid int, rq *BankInfo) *BaseResponse { //用户绑定银行卡
	if rq.Account == "" || rq.BankAddress == "" || rq.BankName == "" || rq.RealName == "" || rq.RoutNumber == "" || rq.SwiftCode == "" {
		return &BaseResponse{
			State: STATE_FAILD,
			Msg:   "the bank info is valid",
		}
	}
	data := db.DB_PARAMS{
		"uid":          uid,
		"bankname":     rq.BankName,
		"realname":     rq.RealName,
		"account":      rq.Account,
		"router_num":   rq.RoutNumber,
		"swift_code":   rq.SwiftCode,
		"bank_address": rq.BankAddress,
	}
	one, _ := config.GlobalDB.FetchOne(DB_TABLE_BANKINFO, db.DB_PARAMS{"uid": uid}, db.DB_FIELDS{})
	if one != nil {
		cacheID := m.MakeCacheId("bank", uid)
		config.GlobalDB.UpdateData(DB_TABLE_BANKINFO, data, db.DB_PARAMS{"id": one["id"].Value})
		config.GlobalRedis.Del(HASH_USER_BANK, cacheID)
	} else {
		config.GlobalDB.InsertData(DB_TABLE_BANKINFO, data)
	}
	return &BaseResponse{State: STATE_SUCCESS, Msg: "ok"}
}

func (m *CreditModel) GetBankInfo(uid int) *BankInfo {
	var rs BankInfo
	cacheID := m.MakeCacheId("bank", uid)
	err := config.GlobalRedis.GetObject(HASH_USER_BANK, cacheID, &rs)
	if err == nil && rs.Account != "" {
		return &rs
	}
	one, _ := config.GlobalDB.FetchOne(DB_TABLE_BANKINFO, db.DB_PARAMS{"uid": uid}, db.DB_FIELDS{})
	if one != nil {
		rs.Account = one["account"].ToString()
		rs.BankAddress = one["bank_address"].ToString()
		rs.BankName = one["bankname"].ToString()
		rs.RealName = one["realname"].ToString()
		rs.RoutNumber = one["router_num"].ToString()
		rs.SwiftCode = one["swift_code"].ToString()
		config.GlobalRedis.SetValue(HASH_USER_BANK, cacheID, rs)
		return &rs
	}
	return nil
}
