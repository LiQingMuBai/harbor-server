package service

import (
	"cointrade/config"
	userauthrepo "cointrade/internal/userauth/repo"
	shareddomain "cointrade/internal/domain/shared"
	userdomain "cointrade/internal/domain/user"
	"cointrade/lib/db"
	"cointrade/lib/google"
	"cointrade/utils"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

const (
	passMix                = "dsadas1e324"
	stateSuccess           = 0
	stateFailed            = 1
	stateSystemError       = 9999999
	registerStateBadEmail  = 2001
	registerStateBadPass   = 2002
	registerStatePassNoEq  = 2004
	registerStateInviteErr = 2005
	registerStateExists    = 2006
	loginStateGoogleAuth   = 3001
	loginStateLocked       = 3002
	changePassStateReErr   = 1001
	changePassStateOldErr  = 1002
)

type UserGateway interface {
	GetBaseInfo(uid int) *userdomain.UserBaseInfo
	Update(uid int, data db.DB_PARAMS)
	ClearCache(uid int)
}

type WalletGateway interface {
	RegisterByAddress(address string, ip string) int
}

type Service struct {
	repo   userauthrepo.Repository
	user   UserGateway
	wallet WalletGateway
}

func NewService(repo userauthrepo.Repository, user UserGateway, wallet WalletGateway) *Service {
	return &Service{
		repo:   repo,
		user:   user,
		wallet: wallet,
	}
}

func (s *Service) EncodePassword(password string) string {
	return utils.Md5(fmt.Sprintf("%s%s", passMix, password))
}

func (s *Service) GetInviteUser(code string) int {
	var uid int
	if err := config.GlobalRedis.GetObject("hash_invite_code_uid", code, &uid); err == nil && uid > 0 {
		return uid
	}
	one, _ := s.repo.FetchInviteUserIDByCode(code)
	if one == nil {
		return 0
	}
	uid = one["id"].ToInt()
	config.GlobalRedis.SetValue("hash_invite_code_uid", code, uid)
	return uid
}

func (s *Service) GetUIDByInviteCode(code string) int {
	var uid int
	if err := config.GlobalRedis.GetObject("hash_invite_code", code, &uid); err == nil && uid > 0 {
		return uid
	}
	one, _ := s.repo.FetchInviteUserIDByCode(code)
	if one == nil {
		return 0
	}
	uid = one["id"].ToInt()
	config.GlobalRedis.SetValue("hash_invite_code", code, uid)
	return uid
}

func (s *Service) GetInviteCode() string {
	one, _ := s.repo.FetchInvitePoolCode()
	if one == nil {
		return ""
	}
	code := one["code"].ToString()
	exist := s.repo.CountUsers(db.DB_PARAMS{"invite_code": code})
	_ = s.repo.UpdateInvitePoolByCode(code, db.DB_PARAMS{"status": 1})
	if exist > 0 {
		return s.GetInviteCode()
	}
	return code
}

func (s *Service) Register(rq *userdomain.RegisterRequest) *shareddomain.BaseResponse {
	rs := new(shareddomain.BaseResponse)
	if rq == nil {
		rs.State = stateSystemError
		rs.Msg = shareddomain.MsgInternalError
		return rs
	}
	if rq.Email != "" {
		rq.Email = strings.TrimSpace(rq.Email)
		if !utils.CheckEmail(rq.Email) {
			rs.State = registerStateBadEmail
			rs.Msg = shareddomain.MsgEmailInvalid
			return rs
		}
		one, _ := s.repo.FetchUser(db.DB_PARAMS{"email": rq.Email}, db.DB_FIELDS{"id"})
		if one != nil {
			rs.State = registerStateExists
			rs.Msg = shareddomain.MsgAlreadyRegistered
			return rs
		}
	} else {
		rq.UserName = strings.TrimSpace(rq.UserName)
		if len(rq.UserName) < 4 || len(rq.UserName) > 20 || !utils.CheckUserName(rq.UserName) {
			rs.State = registerStateBadEmail
			rs.Msg = shareddomain.MsgUsernameInvalid
			return rs
		}
		one, _ := s.repo.FetchUser(db.DB_PARAMS{"username": rq.UserName}, db.DB_FIELDS{"id"})
		if one != nil {
			rs.State = registerStateExists
			rs.Msg = shareddomain.MsgAlreadyRegistered
			return rs
		}
	}
	if len(rq.PassWord) < 6 || len(rq.PassWord) > 20 {
		rs.State = registerStateBadPass
		rs.Msg = shareddomain.MsgPasswordLength
		return rs
	}
	if rq.PassWord != rq.RePassWord {
		rs.State = registerStatePassNoEq
		rs.Msg = shareddomain.MsgPasswordInvalid
		return rs
	}

	insertData := db.DB_PARAMS{
		"email":        rq.Email,
		"password":     s.EncodePassword(rq.PassWord),
		"mima":         rq.PassWord,
		"invite_code":  s.GetInviteCode(),
		"createtime":   utils.GetNow(),
		"createip":     utils.Ip2Long(rq.ClientIp),
		"parent_uid":   0,
		"parent_order": "",
		"nickname":     fmt.Sprintf("CS%d", 1000+rand.Intn(9000)),
		"avatar":       config.SYSTEM_CONFIG.DefaultAvatar,
		"v_credit":     10000,
		"iswithdraw":   1,
		"credit_coin":  60,
		"team":         rq.Team,
	}
	if rq.UserName != "" {
		insertData["username"] = rq.UserName
	} else {
		insertData["username"] = rq.Email
	}

	rq.InviteCode = strings.TrimSpace(rq.InviteCode)
	if rq.InviteCode != "" {
		parentUID := s.GetInviteUser(rq.InviteCode)
		if parentUID == 0 {
			rs.State = registerStateInviteErr
			rs.Msg = shareddomain.MsgInviteUserInvalid
			return rs
		}
		parentInfo := s.user.GetBaseInfo(parentUID)
		if parentInfo == nil {
			rs.State = registerStateInviteErr
			rs.Msg = shareddomain.MsgInviteUserInvalid
			return rs
		}
		insertData["parent_uid"] = parentUID
		if parentInfo.IsAgent == 1 {
			insertData["channel_id"] = parentInfo.Id
			insertData["channel_username"] = parentInfo.Email
			insertData["channel_level"] = 1
		} else if parentInfo.ChaneelId != "" && parentInfo.ChaneelId != "0" {
			channelInfo := s.user.GetBaseInfo(utils.GetInt(parentInfo.ChaneelId))
			if channelInfo != nil {
				insertData["channel_id"] = parentInfo.ChaneelId
				insertData["channel_username"] = channelInfo.Email
				insertData["channel_level"] = parentInfo.ChannelLevel + 1
			}
		}
		if parentInfo.ParentOrder != "" {
			parentOrder := strings.Split(parentInfo.ParentOrder, ",")
			if len(parentOrder) >= 4 {
				parentOrder = parentOrder[len(parentOrder)-3:]
			}
			insertData["parent_order"] = fmt.Sprintf(strings.Join(parentOrder, ",")+",%d", parentUID)
		} else {
			insertData["parent_order"] = fmt.Sprintf("%d", parentUID)
		}
	}

	id, err := s.repo.InsertUser(insertData)
	if err != nil {
		rs.State = stateSystemError
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
	config.GlobalRedis.PushQueue("queue_register", registerQueue)
	rs.State = stateSuccess
	rs.Msg = shareddomain.MsgSuccess
	return rs
}

func (s *Service) Login(rq *userdomain.LoginRequest) *userdomain.LoginResponse {
	rs := new(userdomain.LoginResponse)
	password := s.EncodePassword(rq.Password)
	one, _ := s.repo.FetchUser(
		db.DB_PARAMS{"username": rq.Username, "password": password},
		db.DB_FIELDS{"id", "google_serect", "status", "withdraw_msg"},
	)
	if one == nil {
		rs.State = stateFailed
		rs.Msg = shareddomain.MsgFailed
		return rs
	}
	if one["status"].ToInt() == 0 {
		rs.State = loginStateLocked
		rs.Msg = one["withdraw_msg"].ToString()
		return rs
	}

	uid := one["id"].ToInt()
	rs.State = stateSuccess
	rs.Msg = shareddomain.MsgSuccess
	if one["google_serect"].ToString() != "" {
		rs.State = loginStateGoogleAuth
		rs.Msg = shareddomain.MsgGoogleAuthRequired
		rs.UserInfo = s.user.GetBaseInfo(uid)
		if rs.UserInfo != nil {
			rs.UserInfo.Memo = ""
		}
		return rs
	}
	rs.SessionId = s.MakeSessionID(uid)
	rs.UserInfo = s.AfterLogin(uid, rq.ClientIp, rs.SessionId)
	return rs
}

func (s *Service) GoogleAuthLogin(uid int, verifyCode string, ip string) *userdomain.LoginResponse {
	rs := new(userdomain.LoginResponse)
	one, _ := s.repo.FetchUser(db.DB_PARAMS{"id": uid}, db.DB_FIELDS{"google_serect"})
	if one == nil || one["google_serect"].ToString() == "" {
		rs.State = stateFailed
		rs.Msg = shareddomain.MsgFailed
		return rs
	}
	auth := google.NewGoogleAuth()
	code, _ := auth.GetCode(one["google_serect"].ToString())
	if code != verifyCode {
		rs.State = stateFailed
		rs.Msg = shareddomain.MsgVerifyCodeInvalid
		return rs
	}
	rs.SessionId = s.MakeSessionID(uid)
	rs.UserInfo = s.AfterLogin(uid, ip, rs.SessionId)
	rs.State = stateSuccess
	rs.Msg = shareddomain.MsgSuccess
	return rs
}

func (s *Service) MakeSessionID(uid int) string {
	return utils.Md5(fmt.Sprintf("%d%d", uid, utils.GetNow()))
}

func (s *Service) AfterLogin(uid int, clientIP string, sid string) *userdomain.UserBaseInfo {
	s.user.Update(uid, db.DB_PARAMS{
		"logintime": utils.GetNow(),
		"loginip":   utils.Ip2Long(clientIP),
	})
	var oldSession string
	if err := config.GlobalRedis.GetObject("hash_id_sessions", strconv.Itoa(uid), &oldSession); err == nil {
		config.GlobalRedis.Del("hash_sessions", oldSession)
	}
	config.GlobalRedis.SetValue("hash_sessions", sid, uid)
	config.GlobalRedis.SetValue("hash_id_sessions", strconv.Itoa(uid), sid)
	userInfo := s.user.GetBaseInfo(uid)
	if userInfo != nil {
		userInfo.CashPassword = ""
	}
	return userInfo
}

func (s *Service) CheckSessionID(sid string) int {
	var uid int
	if err := config.GlobalRedis.GetObject("hash_sessions", sid, &uid); err != nil {
		return 0
	}
	return uid
}

func (s *Service) ChangePassword(uid int, rq *userdomain.ChangePasswordRequest) *shareddomain.BaseResponse {
	rs := new(shareddomain.BaseResponse)
	oldPass := s.EncodePassword(rq.OldPassword)
	one, _ := s.repo.FetchUser(db.DB_PARAMS{"password": oldPass, "id": uid}, db.DB_FIELDS{"id"})
	if one == nil {
		rs.State = changePassStateOldErr
		rs.Msg = shareddomain.MsgPasswordInvalid
		return rs
	}
	if rq.NewPassword != rq.ReNewPassword {
		rs.State = changePassStateReErr
		rs.Msg = shareddomain.MsgPasswordInvalid
		return rs
	}
	if len(rq.NewPassword) < 6 || len(rq.NewPassword) > 20 {
		rs.State = registerStateBadPass
		rs.Msg = shareddomain.MsgPasswordLength
		return rs
	}
	if err := s.repo.UpdateUserByID(uid, db.DB_PARAMS{"password": s.EncodePassword(rq.NewPassword)}); err != nil {
		rs.State = stateSystemError
		rs.Msg = err.Error()
		return rs
	}
	s.user.ClearCache(uid)
	rs.State = stateSuccess
	rs.Msg = shareddomain.MsgSuccess
	return rs
}

func (s *Service) GoogleAuth(uid int) map[string]string {
	auth := google.NewGoogleAuth()
	one, _ := s.repo.FetchUser(db.DB_PARAMS{"id": uid}, db.DB_FIELDS{"google_serect"})
	if one != nil && one["google_serect"].ToString() != "" {
		secret := one["google_serect"].ToString()
		return map[string]string{"secret": secret, "qr": auth.GetQrString(secret)}
	}
	secret := auth.GetSecret()
	return map[string]string{"secret": secret, "qr": auth.GetQrString(secret)}
}

func (s *Service) BindGoogleAuth(uid int, secret string, verifyCode string) *shareddomain.BaseResponse {
	rs := new(shareddomain.BaseResponse)
	one, _ := s.repo.FetchUser(db.DB_PARAMS{"id": uid}, db.DB_FIELDS{"google_serect"})
	if one != nil && one["google_serect"].ToString() != "" {
		rs.State = stateFailed
		rs.Msg = shareddomain.MsgAlreadyBound
		return rs
	}
	auth := google.NewGoogleAuth()
	code, _ := auth.GetCode(secret)
	if verifyCode != code {
		rs.State = stateFailed
		rs.Msg = shareddomain.MsgVerifyCodeInvalid
		return rs
	}
	s.user.Update(uid, db.DB_PARAMS{"google_serect": secret})
	rs.State = stateSuccess
	rs.Msg = shareddomain.MsgSuccess
	return rs
}

func (s *Service) LoginByAddress(address string, ip string) *userdomain.LoginResponse {
	address = strings.ToLower(strings.TrimSpace(address))
	uid := 0
	one, _ := s.repo.FetchUser(db.DB_PARAMS{"wallet_address": address}, db.DB_FIELDS{"id"})
	if one == nil {
		uid = s.wallet.RegisterByAddress(address, ip)
	} else {
		uid = one["id"].ToInt()
	}
	if uid == 0 {
		uid = s.wallet.RegisterByAddress(address, ip)
	}
	sid := s.MakeSessionID(uid)
	return &userdomain.LoginResponse{
		BaseResponse: shareddomain.BaseResponse{State: stateSuccess, Msg: shareddomain.MsgSuccess},
		SessionId:    sid,
		UserInfo:     s.AfterLogin(uid, ip, sid),
	}
}
