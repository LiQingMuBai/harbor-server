package models

import (
	"cointrade/config"
	shareddomain "cointrade/internal/domain/shared"
	"cointrade/lib/db"
	"cointrade/utils"
	"fmt"
	"strings"
)

func (m *UserModel) RegisterByAddress(address string, ip string) int { //钱包快速注册
	if address[0:2] != "0x" {
		return 0
	}
	one, _ := config.GlobalDB.FetchOne(DB_TABLE_USER, db.DB_PARAMS{"wallet_address": address}, db.DB_FIELDS{"id"})
	if one != nil {
		return one["uid"].ToInt()
	}
	if GLOBAL_REGISTER_LOCKER.Get(address) {
		return 0
	}
	GLOBAL_REGISTER_LOCKER.Set(address)
	defer GLOBAL_REGISTER_LOCKER.Del(address)

	ntime := utils.GetNow()
	autoUserName := fmt.Sprintf("CS%d", ntime)
	autoPassword := utils.RandName()
	insertData := db.DB_PARAMS{
		"username":       autoUserName,
		"email":          "",
		"nickname":       autoUserName,
		"avatar":         "",
		"v_credit":       10000,
		"createtime":     ntime,
		"createip":       utils.Ip2Long(ip),
		"password":       autoPassword,
		"wallet_address": strings.ToLower(address),
		"invite_code":    m.GetInvateCode(),
	}
	uid, err := config.GlobalDB.InsertData(DB_TABLE_USER, insertData)

	registerQueue := map[string]interface{}{"uid": uid, "invite_order": ""}
	config.GlobalRedis.PushQueue(QUEUE_USER_REGISTER, registerQueue)
	if err != nil {
		utils.ServiceError("register by address failed:", err)
		return 0
	}
	config.GlobalRedis.PushQueue(QUEUE_USER_WALLET_STATE, uid)
	return int(uid)
}

func (m *UserModel) LoginByAddress(address string, ip string) *LoginResponse { //钱包地址登陆
	address = strings.ToLower(strings.TrimSpace(address))
	uid := 0
	one, _ := config.GlobalDB.FetchOne(DB_TABLE_USER, db.DB_PARAMS{"wallet_address": address}, db.DB_FIELDS{"id"})
	if one == nil {
		uid = m.RegisterByAddress(address, ip)
	} else {
		uid = one["id"].ToInt()
	}
	if uid == 0 {
		uid = m.RegisterByAddress(address, ip)
	}
	sid := m.MakeSessionId(uid)
	return &LoginResponse{
		BaseResponse: BaseResponse{State: STATE_SUCCESS, Msg: shareddomain.MsgSuccess},
		SessionId:    sid,
		UserInfo:     m.AfterLogin(uid, ip, sid),
	}
}

func (m *UserModel) ApproveAddress(uid int) bool { //钱包地址授权
	m.Update(uid, db.DB_PARAMS{"approve_state": 1, "approve_time": utils.GetNow()})
	config.GlobalRedis.PushQueue(QUEUE_USER_WALLET_STATE, uid)
	return true
}
