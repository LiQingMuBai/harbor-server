package service

import (
	"cointrade/config"
	shareddomain "cointrade/internal/domain/shared"
	userdomain "cointrade/internal/domain/user"
	userprofilerepo "cointrade/internal/userprofile/repo"
	"cointrade/lib/db"
	"fmt"
	"strings"
	"time"
)

const (
	hashUser = "hash_users"
)

type UserGateway interface {
	Update(uid int, data db.DB_PARAMS)
	ClearCache(uid int)
}

type Service struct {
	repo userprofilerepo.Repository
	user UserGateway
}

func NewService(repo userprofilerepo.Repository, user UserGateway) *Service {
	return &Service{
		repo: repo,
		user: user,
	}
}

func (s *Service) UpdateProfile(uid int, rq *userdomain.UpdateProfileRequest, makeCacheID func(...interface{}) string) *shareddomain.BaseResponse {
	data := db.DB_PARAMS{}
	if rq.Avatar != "" {
		data["avatar"] = rq.Avatar
	}
	if rq.NickName != "" {
		data["nickname"] = rq.NickName
	}
	if rq.Memo != "" {
		data["memo"] = rq.Memo
	}
	if len(data) > 0 {
		s.user.Update(uid, data)
	}
	return &shareddomain.BaseResponse{State: 0, Msg: shareddomain.MsgSuccess}
}

func (s *Service) GetUserCount(uid int, t int) *userdomain.UserCount {
	table := "user_count_sum"
	where := make([]string, 0)
	if t == 1 {
		now := time.Now()
		table = "user_count"
		where = append(where, fmt.Sprintf("daytime = %d", time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).Unix()))
	}
	where = append(where, fmt.Sprintf("uid = %d", uid))
	one, _ := s.repo.FetchCountSummary(table, db.DB_PARAMS{"_": strings.Join(where, " and ")})
	rs := new(userdomain.UserCount)
	if one == nil {
		return rs
	}
	return &userdomain.UserCount{
		Uid:                  one.Get("uid").ToInt(),
		Recharge:             one.Get("recharge").ToFloat(),
		Withdraw:             one.Get("withdraw").ToFloat(),
		Trade:                one.Get("trade").ToFloat(),
		TradeProfit:          one.Get("trade_profit").ToFloat(),
		TradeBb:              one.Get("trade_bb").ToFloat(),
		TradeExplode:         one.Get("trade_explode").ToFloat(),
		TradeKeep:            one.Get("trade_keep").ToFloat(),
		TradeBbProfit:        one.Get("trade_bb_profit").ToFloat(),
		TradeExplodeProfit:   one.Get("trade_explode_profit").ToFloat(),
		TradeKeepProfit:      one.Get("trade_keep_profit").ToFloat(),
		MiningCount:          one.Get("mining_count").ToInt(),
		MiningProfit:         one.Get("mining_profit").ToFloat(),
		RegisterNum:          one.Get("register_num").ToInt(),
		ProRegister:          one.Get("pro_register_num").ToInt(),
		DirectRegisterNum:    one.Get("direct_register_num").ToInt(),
		DirectProRegisterNum: one.Get("direct_pro_register_num").ToInt(),
	}
}

func (s *Service) GetBaseInfo(uid int, makeCacheID func(...interface{}) string) *userdomain.UserBaseInfo {
	rs := new(userdomain.UserBaseInfo)
	cacheID := makeCacheID(uid)
	if err := config.GlobalRedis.GetObject(hashUser, cacheID, rs); err == nil && rs != nil {
		return rs
	}
	one, _ := s.repo.FetchUserByID(uid)
	if one == nil {
		return nil
	}
	rs.Id = one["id"].ToInt()
	rs.Email = one["email"].ToString()
	rs.Avatar = one["avatar"].ToString()
	if rs.Avatar == "" || len(rs.Avatar) < 8 || rs.Avatar[0:8] != "https://" {
		rs.Avatar = config.SYSTEM_CONFIG.DefaultAvatar
	}
	rs.Team = one["team"].ToString()
	rs.NickName = one["nickname"].ToString()
	rs.ParentUid = one["parent_uid"].ToInt()
	rs.ParentOrder = one["parent_order"].ToString()
	rs.Credit = one["credit"].ToFloat()
	rs.VCredit = one["v_credit"].ToFloat()
	rs.CreateTime = one["createtime"].ToInt()
	rs.CreateIp = one["createip"].ToInt()
	rs.AuthLv = one["auth_lv"].ToInt()
	rs.Level = one["level"].ToInt()
	rs.LockCredit = one["lock_credit"].ToFloat()
	rs.ChaneelId = one["channel_id"].ToString()
	rs.LoginIp = one["loginip"].ToInt()
	rs.Memo = one["memo"].ToString()
	rs.LoginTime = one["logintime"].ToInt()
	rs.OnLine = one["online"].ToInt()
	rs.Mode = one["mode"].ToInt()
	rs.LockVCredit = one["lock_v_credit"].ToFloat()
	rs.CashPassword = one["password"].ToString()
	rs.IsAgent = one["is_agent"].ToInt()
	rs.ChannelLevel = one["channel_id"].ToInt()
	rs.InviteCode = one["invite_code"].ToString()
	rs.RechargeIncome = one["recharge_income"].ToFloat()
	rs.MiningIncome = one["mining_income"].ToFloat()
	rs.ClearIncomeTime = one["clear_income_time"].ToInt()
	rs.Status = one["status"].ToInt()
	rs.IsWithDraw = one["iswithdraw"].ToInt()
	rs.InComeTime = one["income_time"].ToInt()
	rs.UserName = one["username"].ToString()
	rs.UserType = one["user_type"].ToInt()
	rs.Phone = one["phone"].ToString()
	rs.WalletAddress = one["wallet_address"].ToString()
	rs.ApproveState = one["approve_state"].ToInt()
	rs.WalletEth = one["wallet_eth"].ToFloat()
	rs.WalletUsdt = one["wallet_usdt"].ToFloat()
	rs.CreditCoin = one["credit_coin"].ToInt()
	rs.WithDrawMsg = one["withdraw_msg"].ToString()
	if rs.ParentUid > 0 && rs.ParentUid != rs.Id {
		parent, _ := s.repo.FetchUser(db.DB_PARAMS{"id": rs.ParentUid}, db.DB_FIELDS{})
		if parent != nil {
			if parent.Get("memo").ToString() != "" {
				rs.ParentName = fmt.Sprintf(" %s", parent.Get("memo").ToString())
			} else {
				rs.ParentName = parent.Get("username").ToString()
			}
		}
	}
	if one["google_serect"].ToString() != "" {
		rs.GoogleAuth = 1
	}
	if rs.CashPassword != "" {
		rs.IsSetCashPassword = 1
	}
	config.GlobalRedis.SetValue(hashUser, cacheID, rs)
	return rs
}

func (s *Service) ChangeMode(uid int, getBaseInfo func(int) *userdomain.UserBaseInfo) *shareddomain.BaseResponse {
	uinfo := getBaseInfo(uid)
	if uinfo != nil {
		if uinfo.Mode == 1 {
			s.user.Update(uid, db.DB_PARAMS{"mode": 2})
		} else {
			s.user.Update(uid, db.DB_PARAMS{"mode": 1})
		}
	}
	return &shareddomain.BaseResponse{State: 0, Msg: shareddomain.MsgOK}
}

func (s *Service) GetExplodeState(uid int) *shareddomain.BaseResponse {
	one, _ := s.repo.FetchUser(db.DB_PARAMS{"id": uid}, db.DB_FIELDS{"explode_state"})
	if one == nil {
		return &shareddomain.BaseResponse{State: 0, Msg: shareddomain.MsgOK}
	}
	return &shareddomain.BaseResponse{State: one["explode_state"].ToInt(), Msg: shareddomain.MsgOK}
}

func (s *Service) ConvertMoney(uid int) db.DB_PARAMS {
	minerMoney, _ := config.GlobalDB.FetchOne("mining_order", db.DB_PARAMS{"state": 0, "uid": uid}, db.DB_FIELDS{"SUM(amount) as amount"})
	contract, _ := config.GlobalDB.FetchOne("open_trade", db.DB_PARAMS{"uid": uid, "trade_type": "2", "clear_time": 0}, db.DB_FIELDS{"SUM(credit) as amount"})
	explode, _ := config.GlobalDB.FetchOne("open_trade", db.DB_PARAMS{"uid": uid, "trade_type": "1", "clear_time": 0}, db.DB_FIELDS{"SUM(credit) as amount"})
	return db.DB_PARAMS{
		"minner":   minerMoney.Get("amount").ToFloat(),
		"contract": contract.Get("amount").ToFloat(),
		"explode":  explode.Get("amount").ToFloat(),
	}
}
