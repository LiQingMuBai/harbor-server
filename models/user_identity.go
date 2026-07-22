package models

import (
	"cointrade/config"
	"cointrade/lib/db"
	"cointrade/lib/notify"
	"cointrade/utils"
)

func (m *UserModel) AuthLv1(uid int, authinfo *AuthLv1Request) *BaseResponse {
	rs := new(BaseResponse)
	uinfo := m.GetBaseInfo(uid)
	if uinfo.AuthLv >= 1 {
		rs.State = STATE_FAILD
		rs.Msg = "faild"
		return rs
	}
	if authinfo.CardBack == "" || authinfo.CardFront == "" || authinfo.Phone == "" || authinfo.CardType == 0 {
		rs.State = STATE_FAILD
		rs.Msg = "faild"
		return rs
	}

	insertData := db.DB_PARAMS{
		"uid":           uid,
		"realname":      authinfo.Name,
		"inid":          authinfo.IdCard,
		"card_front":    authinfo.CardFront,
		"card_back":     authinfo.CardBack,
		"card_hand":     authinfo.HandCard,
		"process_state": 0,
		"createtime":    utils.GetNow(),
		"phone":         authinfo.Phone,
		"card_type":     authinfo.CardType,
	}

	one, _ := config.GlobalDB.FetchOne(DB_TABLE_USERAUTH, db.DB_PARAMS{"uid": uid}, db.DB_FIELDS{})
	if one != nil {
		if one["process_state"].ToInt() == 2 {
			config.GlobalDB.UpdateData(DB_TABLE_USERAUTH, insertData, db.DB_PARAMS{"uid": uid})
		} else {
			rs.State = STATE_FAILD
			rs.Msg = "faild"
			return rs
		}
	} else {
		config.GlobalDB.InsertData(DB_TABLE_USERAUTH, insertData)
	}
	notify.NOTIFY.AddNotify(&notify.NotifyItem{Type: 3, Num: 1})
	rs.State = STATE_SUCCESS
	rs.Msg = "success!"
	return rs
}

func (m *UserModel) AuthLv2(uid int, rq *AuthLv2Request) *BaseResponse {
	rs := new(BaseResponse)
	uinfo := m.GetBaseInfo(uid)
	if uinfo.AuthLv >= 2 || uinfo.AuthLv == 0 {
		rs.State = STATE_FAILD
		rs.Msg = "faild1"
		return rs
	}
	if rq.FarmilyName == "" || rq.WalletAddress == "" || rq.Address == "" {
		rs.State = STATE_FAILD
		rs.Msg = "faild2"
		return rs
	}
	data := db.DB_PARAMS{
		"uid":               uid,
		"farmily_name":      rq.FarmilyName,
		"relation":          rq.Relation,
		"address":           rq.Address,
		"contact":           rq.Contact,
		"wallet_address":    rq.WalletAddress,
		"chaintype":         rq.ChainType,
		"second_card_front": rq.Second_card_front,
		"second_card_back":  rq.Second_card_Hand,
		"second_card_hand":  rq.Second_card_Hand,
		"createtime":        utils.GetNow(),
		"state":             0,
	}
	one, _ := config.GlobalDB.FetchOne(DB_TABLE_USERAUTH_LV2, db.DB_PARAMS{"uid": uid}, db.DB_FIELDS{})
	if one != nil {
		if one["state"].ToInt() == 2 {
			config.GlobalDB.UpdateData(DB_TABLE_USERAUTH_LV2, data, db.DB_PARAMS{"id": one["id"].Value})
		} else {
			rs.State = STATE_FAILD
			rs.Msg = "faild3"
			return rs
		}
	} else {
		config.GlobalDB.InsertData(DB_TABLE_USERAUTH_LV2, data)
	}
	notify.NOTIFY.AddNotify(&notify.NotifyItem{Type: 4, Num: 1})
	rs.State = STATE_SUCCESS
	rs.Msg = "success!"
	return rs
}

func (m *UserModel) GetAuthInfo(uid int) map[int]interface{} {
	lv1Info, _ := config.GlobalDB.FetchRow(DB_TABLE_USERAUTH, db.DB_PARAMS{"uid": uid}, db.DB_FIELDS{})
	lv2Info, _ := config.GlobalDB.FetchRow(DB_TABLE_USERAUTH_LV2, db.DB_PARAMS{"uid": uid}, db.DB_FIELDS{})
	return map[int]interface{}{1: lv1Info, 2: lv2Info}
}

func (m *UserModel) ChangeCashPassword(uid int, rq *SetCashPasswordRequest) *BaseResponse {
	rs := new(BaseResponse)
	if rq.Password == "" {
		rs.State = 1
		rs.Msg = "the password is volid"
		return rs
	}
	if len(rq.Password) < 6 || len(rq.Password) > 20 {
		rs.State = 1
		rs.Msg = "the password length is must between 6-20"
		return rs
	}
	if rq.Password != rq.RePassword {
		rs.State = 1
		rs.Msg = "confirm password error"
		return rs
	}
	m.Update(uid, db.DB_PARAMS{"cash_password": rq.Password})
	rs.State = 0
	rs.Msg = "success"
	return rs
}

func (m *UserModel) UpdateCashPassword(uid int, rq *UpdateCashPasswordRequest) *BaseResponse {
	uinfo := m.GetBaseInfo(uid)
	rs := new(BaseResponse)
	if uinfo == nil {
		rs.State = STATE_FAILD
		rs.Msg = "user is not exists"
		return rs
	}
	if rq.O_Password != uinfo.CashPassword {
		rs.State = CHANGE_PASS_STATE_OLDERROR
		rs.Msg = "cash password is error"
		return rs
	}
	if len(rq.N_Password) < 6 || len(rq.N_Password) > 20 {
		rs.State = REGISTER_STATE_ERRORPASSWORD
		rs.Msg = "password length is must between 6-20"
		return rs
	}
	if rq.N_Password != rq.R_Password {
		rs.State = REGISTER_STATE_ERRORPASSWORD
		rs.Msg = "confirm password is error"
		return rs
	}
	m.Update(uid, db.DB_PARAMS{"cash_password": rq.N_Password})
	rs.State = STATE_SUCCESS
	rs.Msg = "success"
	return rs
}

func (m *UserModel) BindPhone(uid int, phone string, code string) *BaseResponse {
	rs := new(BaseResponse)
	uinfo := m.GetBaseInfo(uid)
	if uinfo.Phone != "" {
		rs.State = BIND_PHONE_STATE_BINDED
		rs.Msg = "binded"
		return rs
	}
	one, _ := config.GlobalDB.FetchOne(DB_TABLE_USER, db.DB_PARAMS{"phone": phone}, db.DB_FIELDS{"id"})
	if one != nil {
		rs.State = BIND_PHONE_STATE_BINDED
		rs.Msg = "phone num binded"
		return rs
	}
	vCode := MODEL_CODE.GetBindSmsCode(uid, phone)
	if vCode != code {
		rs.State = BIND_PHONE_STATE_ERRORCODE
		rs.Msg = "error code"
		return rs
	}
	m.Update(uid, db.DB_PARAMS{"phone": phone})
	rs.State = STATE_SUCCESS
	rs.Msg = "success"
	return rs
}

func (m *UserModel) BindEmail(uid int, email string, code string) *BaseResponse {
	rs := new(BaseResponse)
	uinfo := m.GetBaseInfo(uid)
	if uinfo.Email != "" && uinfo.Email != "0" {
		rs.State = BIND_PHONE_STATE_BINDED
		rs.Msg = "binded"
		return rs
	}
	one, _ := config.GlobalDB.FetchOne(DB_TABLE_USER, db.DB_PARAMS{"email": email}, db.DB_FIELDS{"id"})
	if one != nil {
		rs.State = BIND_PHONE_STATE_BINDED
		rs.Msg = "phone num binded"
		return rs
	}
	vCode := MODEL_CODE.GetEmailCodeBind(email)
	if vCode != code {
		rs.State = BIND_PHONE_STATE_ERRORCODE
		rs.Msg = "error code"
		return rs
	}
	m.Update(uid, db.DB_PARAMS{"email": email})
	rs.State = STATE_SUCCESS
	rs.Msg = "success"
	return rs
}
