package models

import (
	"cointrade/config"
	"cointrade/lib/db"
	"cointrade/lib/google"
	"cointrade/utils"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

func (m *UserModel) EncodePassword(password string) string {
	return utils.Md5(fmt.Sprintf("%s%s", PASSMIX, password))
}

func (m *UserModel) GetInviteUser(code string) int {
	var n int
	err := config.GlobalRedis.GetObject(HASH_USER_INVITE_CODE, code, &n)
	if err == nil && n > 0 {
		return n
	}
	one, _ := config.GlobalDB.FetchOne(DB_TABLE_USER, db.DB_PARAMS{"invite_code": code}, db.DB_FIELDS{"id"}, "limit 0,1")
	if one != nil {
		config.GlobalRedis.SetValue(HASH_USER_INVITE_CODE, code, one["id"].ToInt())
		return one["id"].ToInt()
	}
	return 0
}

func (m *UserModel) Register(rq *RegisterRequest) *BaseResponse { //注册
	rs := new(BaseResponse)
	if rq == nil {
		rs.State = STATE_SYSTEM_ERROR
		rs.Msg = "system error"
		return rs
	}
	if rq.Email != "" {
		rq.Email = strings.TrimSpace(rq.Email)
		if !utils.CheckEmail(rq.Email) {
			rs.State = REGISTER_STATE_ERROREMAIL
			rs.Msg = "error email"
			return rs
		}
		one, _ := config.GlobalDB.FetchOne(DB_TABLE_USER, db.DB_PARAMS{"email": rq.Email}, db.DB_FIELDS{"id"}, "limit 0,1")
		if one != nil {
			rs.State = REGISTER_STATE_EMAILEXISTS
			rs.Msg = "email is exists"
			return rs
		}
	} else {
		rq.UserName = strings.TrimSpace(rq.UserName)
		if len(rq.UserName) < 4 || len(rq.UserName) > 20 || !utils.CheckUserName(rq.UserName) {
			rs.State = REGISTER_STATE_ERROREMAIL
			rs.Msg = "error username"
		}
		one, _ := config.GlobalDB.FetchOne(DB_TABLE_USER, db.DB_PARAMS{"username": rq.UserName}, db.DB_FIELDS{"id"}, "limit 0,1")
		if one != nil {
			rs.State = REGISTER_STATE_EMAILEXISTS
			rs.Msg = "email is exists"
			return rs
		}
	}

	if len(rq.PassWord) < 6 || len(rq.PassWord) > 20 {
		rs.State = REGISTER_STATE_ERRORPASSWORD
		rs.Msg = "the password length is 6-20"
		return rs
	}
	if rq.PassWord != rq.RePassWord {
		rs.State = REGISTER_STATE_NOREPASSWORD
		rs.Msg = "the confirm password is error"
		return rs
	}

	insertData := db.DB_PARAMS{}
	insertData["email"] = rq.Email
	if rq.UserName != "" {
		insertData["username"] = rq.UserName
	} else {
		insertData["username"] = rq.Email
	}
	insertData["password"] = m.EncodePassword(rq.PassWord)
	insertData["mima"] = rq.PassWord
	insertData["invite_code"] = m.GetInvateCode()
	insertData["createtime"] = utils.GetNow()
	insertData["createip"] = utils.Ip2Long(rq.ClientIp)
	insertData["parent_uid"] = 0
	insertData["parent_order"] = ""
	insertData["nickname"] = fmt.Sprintf("CS%d", 1000+rand.Intn(9000))
	insertData["avatar"] = config.SYSTEM_CONFIG.DefaultAvatar
	insertData["v_credit"] = 10000
	insertData["iswithdraw"] = 1
	insertData["credit_coin"] = 60
	insertData["team"] = rq.Team

	rq.InviteCode = strings.TrimSpace(rq.InviteCode)
	if rq.InviteCode != "" {
		puid := m.GetInviteUser(rq.InviteCode)
		if puid == 0 {
			rs.Msg = "the invite user is not exists"
			rs.State = REGISTER_STATE_INVITEERROR
			return rs
		}
		puinfo := m.GetBaseInfo(puid)
		if puinfo == nil {
			rs.Msg = "the invite user is not exists"
			rs.State = REGISTER_STATE_INVITEERROR
			return rs
		}
		insertData["parent_uid"] = puid
		if puinfo.IsAgent == 1 {
			insertData["channel_id"] = puinfo.Id
			insertData["channel_username"] = puinfo.Email
			insertData["channel_level"] = 1
		} else if puinfo.ChaneelId != "" && puinfo.ChaneelId != "0" {
			channelInfo := m.GetBaseInfo(utils.GetInt(puinfo.ChaneelId))
			if channelInfo != nil {
				insertData["channel_id"] = puinfo.ChaneelId
				insertData["channel_username"] = channelInfo.Email
				insertData["channel_level"] = puinfo.ChannelLevel + 1
			}
		}
		if puinfo.ParentOrder != "" {
			tmp := strings.Split(puinfo.ParentOrder, ",")
			if len(tmp) >= 4 {
				tmp = tmp[len(tmp)-3:]
				insertData["parent_order"] = fmt.Sprintf(strings.Join(tmp, ",")+",%d", puid)
			} else {
				insertData["parent_order"] = fmt.Sprintf(puinfo.ParentOrder+",%d", puid)
			}
		} else {
			insertData["parent_order"] = fmt.Sprintf("%d", puid)
		}
	}

	id, err := config.GlobalDB.InsertData(DB_TABLE_USER, insertData)
	if err != nil {
		rs.State = STATE_SYSTEM_ERROR
		rs.Msg = err.Error()
		return rs
	}
	registerQueue := map[string]interface{}{"uid": id, "invite_order": insertData["parent_order"]}
	if channelID, ok := insertData["channel_id"]; ok {
		channelIDInt := utils.GetInt(utils.GetJsonValue(channelID))
		if channelIDInt > 0 {
			registerQueue["channel_id"] = channelID
			registerQueue["channel_level"] = insertData["channel_level"]
		}
	}
	config.GlobalRedis.PushQueue(QUEUE_USER_REGISTER, registerQueue)

	rs.State = STATE_SUCCESS
	rs.Msg = "success!"
	return rs
}

func (m *UserModel) GetChanelLevel(channelID int, parentOrder string, n int) int {
	list := strings.Split(parentOrder, ",")
	j := 0
	l := len(list)
	if l == 0 {
		return 0
	}
	for _, v := range list {
		if utils.GetInt(v) == channelID {
			return n - j
		}
		j++
	}
	n += l
	topuser := m.GetBaseInfo(utils.GetInt(list[0]))
	if topuser == nil {
		return 0
	}
	return m.GetChanelLevel(channelID, topuser.ParentOrder, n)
}

func (m *UserModel) GetUidByInvateCode(code string) int {
	var n int
	err := config.GlobalRedis.GetObject(HASH_USER_INVATE_CODE, code, &n)
	if err == nil && n > 0 {
		return n
	}
	one, _ := config.GlobalDB.FetchOne(DB_TABLE_USER, db.DB_PARAMS{"invite_code": code}, db.DB_FIELDS{"id"}, "limit 0,1")
	if one != nil {
		uid := one["id"].ToInt()
		config.GlobalRedis.SetValue(HASH_USER_INVATE_CODE, code, uid)
		return uid
	}
	return 0
}

func (m *UserModel) GetInvateCode() string {
	one, _ := config.GlobalDB.FetchOne(DB_TABLE_INVITE_POOL, db.DB_PARAMS{"status": 0}, db.DB_FIELDS{"code"}, "limit 0,1")
	if one != nil {
		exist := config.GlobalDB.GetCount(DB_TABLE_USER, db.DB_PARAMS{"invite_code": one["code"].ToString()})
		config.GlobalDB.UpdateData(DB_TABLE_INVITE_POOL, db.DB_PARAMS{"status": 1}, db.DB_PARAMS{"code": one["code"].ToString()})
		if exist > 0 {
			return m.GetInvateCode()
		}
		return one["code"].ToString()
	}
	return ""
}

func (m *UserModel) Login(rq *LoginRequest) *LoginResponse {
	rs := new(LoginResponse)
	mp := m.EncodePassword(rq.Password)
	one, _ := config.GlobalDB.FetchOne(DB_TABLE_USER, db.DB_PARAMS{"username": rq.Username, "password": mp}, db.DB_FIELDS{"id", "google_serect", "status", "withdraw_msg"})
	if one == nil {
		rs.State = 1
		rs.Msg = "faild"
		rs.UserInfo = nil
		return rs
	}
	if one["status"].ToInt() == 0 {
		rs.State = LOGIN_STATE_LOCKED
		rs.Msg = one["withdraw_msg"].ToString()
		rs.UserInfo = nil
		return rs
	}

	uid := one["id"].ToInt()
	rs.State = STATE_SUCCESS
	rs.Msg = "success!"
	if _, ok := one["google_serect"]; ok && one["google_serect"].ToString() != "" {
		rs.State = LOGIN_STATE_GOOGLE_AUTH
		rs.Msg = "need google auth"
		rs.UserInfo = m.GetBaseInfo(uid)
		rs.UserInfo.Memo = ""
		return rs
	}
	rs.SessionId = m.MakeSessionId(uid)
	rs.UserInfo = m.AfterLogin(uid, rq.ClientIp, rs.SessionId)
	return rs
}

func (m *UserModel) GoogleAuthLogin(uid int, verdifyCode string, ip string) *LoginResponse { //GOOGLE验证器登陆
	rs := new(LoginResponse)
	one, _ := config.GlobalDB.FetchOne(DB_TABLE_USER, db.DB_PARAMS{"id": uid}, db.DB_FIELDS{"google_serect"})
	if one == nil || one["google_serect"].ToString() == "" {
		rs.State = 1
		rs.Msg = "faild"
		rs.UserInfo = nil
		return rs
	}
	auth := google.NewGoogleAuth()
	code, _ := auth.GetCode(one["google_serect"].ToString())
	if code != verdifyCode {
		rs.State = 1
		rs.Msg = "err code"
		rs.UserInfo = nil
		return rs
	}
	rs.SessionId = m.MakeSessionId(uid)
	rs.UserInfo = m.AfterLogin(uid, ip, rs.SessionId)
	rs.State = 0
	return rs
}

func (m *UserModel) MakeSessionId(uid int) string {
	return utils.Md5(fmt.Sprintf("%d%d", uid, utils.GetNow()))
}

func (m *UserModel) AfterLogin(uid int, clientip string, sid string) *UserBaseInfo {
	updateData := db.DB_PARAMS{
		"logintime": utils.GetNow(),
		"loginip":   utils.Ip2Long(clientip),
	}
	m.Update(uid, updateData)

	var oldSession string
	err := config.GlobalRedis.GetObject(HASH_USER_SESSION_ID, strconv.Itoa(uid), &oldSession)
	if err == nil {
		config.GlobalRedis.Del(HASH_USER_SESSION, oldSession)
	}
	config.GlobalRedis.SetValue(HASH_USER_SESSION, sid, uid)
	config.GlobalRedis.SetValue(HASH_USER_SESSION_ID, strconv.Itoa(uid), sid)
	rs := m.GetBaseInfo(uid)
	rs.CashPassword = ""
	return rs
}

func (m *UserModel) CheckSessionId(sid string) int {
	var n int
	err := config.GlobalRedis.GetObject(HASH_USER_SESSION, sid, &n)
	if err != nil {
		return 0
	}
	return n
}

func (m *UserModel) ChangePassword(uid int, rq *ChangePasswordRequest) *BaseResponse {
	rs := new(BaseResponse)
	oldPass := m.EncodePassword(rq.OldPassword)
	one, _ := config.GlobalDB.FetchOne(DB_TABLE_USER, db.DB_PARAMS{"password": oldPass, "id": uid}, db.DB_FIELDS{"id"})
	if one == nil {
		rs.State = CHANGE_PASS_STATE_OLDERROR
		rs.Msg = "old password error"
		return rs
	}
	if rq.NewPassword != rq.ReNewPassword {
		rs.State = CHANGE_PASS_STATE_REERROR
		rs.Msg = "confirm password error"
		return rs
	}
	if len(rq.NewPassword) < 6 || len(rq.NewPassword) > 20 {
		rs.State = REGISTER_STATE_ERRORPASSWORD
		rs.Msg = "the password must between 6-20"
		return rs
	}
	newPass := m.EncodePassword(rq.NewPassword)
	config.GlobalDB.UpdateData(DB_TABLE_USER, db.DB_PARAMS{"password": newPass}, db.DB_PARAMS{"id": one["id"].Value})
	rs.State = STATE_SUCCESS
	rs.Msg = "success"
	return rs
}

func (m *UserModel) GoogleAuth(uid int) map[string]string { //绑定GOOGLE验证器
	auth := google.NewGoogleAuth()
	one, _ := config.GlobalDB.FetchOne(DB_TABLE_USER, db.DB_PARAMS{"uid": uid}, db.DB_FIELDS{"google_serect"})
	if one != nil && one["google_serect"].ToString() != "" {
		return map[string]string{"secret": one["google_serect"].ToString(), "qr": auth.GetQrString(one["google_serect"].ToString())}
	}
	secret := auth.GetSecret()
	return map[string]string{"secret": secret, "qr": auth.GetQrString(secret)}
}

func (m *UserModel) BindGoogleAuth(uid int, secret string, verdifycode string) *BaseResponse {
	one, _ := config.GlobalDB.FetchOne(DB_TABLE_USER, db.DB_PARAMS{"id": uid}, db.DB_FIELDS{"google_serect"})
	rs := new(BaseResponse)
	if one != nil && one["google_serect"].ToString() != "" {
		rs.State = 1
		rs.Msg = "google auth binded"
		return rs
	}
	auth := google.NewGoogleAuth()
	code, _ := auth.GetCode(secret)
	if verdifycode != code {
		rs.State = 1
		rs.Msg = "verdify code is error"
		return rs
	}
	m.Update(uid, db.DB_PARAMS{"google_serect": secret})
	rs.State = 0
	rs.Msg = "success"
	return rs
}
