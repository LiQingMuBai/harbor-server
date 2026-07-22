package models

import (
	"cointrade/config"
	"cointrade/lib/db"
	"cointrade/utils"
	"fmt"
	"log"
	"strings"
	"time"
)

func (m *UserModel) UpdatePorfile(uid int, rq *UpdateProfileRequest) *BaseResponse {
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
		m.Update(uid, data)
	}
	return &BaseResponse{State: STATE_SUCCESS, Msg: "success"}
}

func (m *UserModel) GetUserCount(uid int, t int) *UserCount {
	dbname := DB_TABLE_USER_COUNT_SUM
	where := make([]string, 0)
	if t == 1 {
		now := time.Now()
		dbname = DB_TABLE_USER_COUNT
		where = append(where, fmt.Sprintf("daytime = %d", time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).Unix()))
	}
	where = append(where, fmt.Sprintf("uid = %d", uid))
	one, _ := config.GlobalDB.FetchOne(dbname, db.DB_PARAMS{"_": strings.Join(where, " and ")}, db.DB_FIELDS{})
	rs := new(UserCount)
	if one != nil {
		rs = &UserCount{
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
	return rs
}

func (m *UserModel) GetBaseInfo(uid int) *UserBaseInfo { //获得单个用户的基础信息
	rs := new(UserBaseInfo)
	cacheid := m.MakeCacheId(uid)
	err := config.GlobalRedis.GetObject(HASH_USER, cacheid, rs)
	if err == nil && rs != nil {
		return rs
	}

	one, _ := config.GlobalDB.FetchOne(DB_TABLE_USER, db.DB_PARAMS{"id": uid}, db.DB_FIELDS{})
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
		parent, _ := config.GlobalDB.FetchOne(DB_TABLE_USER, db.DB_PARAMS{"id": rs.ParentUid}, db.DB_FIELDS{})
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
	config.GlobalRedis.SetValue(HASH_USER, cacheid, rs)
	return rs
}

func (m *UserModel) IsNewUser(uid int) *BaseResponse {
	rs := new(BaseResponse)
	if exists := config.GlobalDB.GetCount("users", db.DB_PARAMS{"id": uid, "approve_state": 1}); exists == 0 {
		rs.State = STATE_SUCCESS
	} else {
		rs.State = STATE_FAILD
	}
	return rs
}

func (m *UserModel) Claim(uid int) *BaseResponse {
	rs := new(BaseResponse)
	if exists := config.GlobalDB.GetCount("users_mechanism", db.DB_PARAMS{"user_id": uid}); exists == 0 {
		log.Println("NEW USER ", uid, " ADD 10 USDT")
		MODEL_USER.AddCredit(uid, &CreditValue{
			Credit:          100,
			UserCoinLogType: 4100001,
			UserCoinLogInfo: QueueCreditLog{
				Credit:     100,
				CoinType:   "usdt",
				CreateTime: utils.GetNow(),
			},
		})

		params := map[string]interface{}{"user_id": uid}
		_, err := config.GlobalDB.InsertData("users_mechanism", params)
		log.Println("uid has not been exit,so cache uid ", uid)
		if err != nil {
			log.Println("success")
		}
	}
	rs.State = STATE_SUCCESS
	rs.Msg = "success"
	return rs
}

func (m *UserModel) ClearIncome(uid int, cashpassword string) *BaseResponse {
	rs := new(BaseResponse)
	nday := time.Now().Day()
	uinfo := m.GetBaseInfo(uid)
	oday := time.Unix(int64(uinfo.ClearIncomeTime), 0).Day()
	if uinfo.CashPassword == "" || uinfo.CashPassword != m.EncodePassword(cashpassword) || oday == nday {
		rs.State = STATE_FAILD
		rs.Msg = "faild"
		return rs
	}
	if uinfo.MiningIncome+uinfo.RechargeIncome > 0 {
		if m.AddCredit(uid, &CreditValue{
			Credit:          uinfo.MiningIncome + uinfo.RechargeIncome,
			LockCredit:      0,
			VCrdit:          0,
			LockVCredit:     0,
			UserCoinLogType: COIN_LOG_USER_CLEAR_INCOME,
			UserCoinLogInfo: QueueCreditLog{
				Credit:     uinfo.MiningIncome + uinfo.RechargeIncome,
				LockCredit: 0,
				Sn:         "",
				CreateTime: utils.GetNow(),
			},
		}) {
			config.GlobalDB.AddValue(DB_TABLE_USER, map[string]float64{"recharge_income": -1 * uinfo.RechargeIncome, "mining_income": -1 * uinfo.MiningIncome}, db.DB_PARAMS{"id": uid})
			m.ClearCache(uid)
		}
	}
	rs.State = STATE_SUCCESS
	rs.Msg = "success"
	return rs
}

func (m *UserModel) Welcome2() *WelcomeResponse {
	rs := new(WelcomeResponse)
	one, _ := config.GlobalDB.FetchOne(DB_TABLE_WELCOME, nil, db.DB_FIELDS{"id", "platform_name", "welcome_page"})
	if one == nil {
		rs.State = 1
		rs.Msg = "failed"
		rs.WelcomeInfo = nil
		return rs
	}
	rs.State = STATE_SUCCESS
	rs.Msg = "success"
	rs.WelcomeInfo = &WelcomeInfo{
		PlatformName: one["platform_name"].ToString(),
		WelcomePage:  one["welcome_page"].ToString(),
	}
	return rs
}

func (m *UserModel) CrossTrade(uid int, data db.DB_PARAMS) *CrossPlatformTradeResponse {
	rs := new(CrossPlatformTradeResponse)
	rs.State = STATE_SUCCESS
	rs.Msg = "success"

	var userBalance float64
	userAsset4USDT, _ := config.GlobalDB.FetchOne(DB_TABLE_USER, db.DB_PARAMS{"id": uid}, db.DB_FIELDS{})
	if userAsset4USDT != nil {
		userBalance = userAsset4USDT["credit"].ToFloat()
	}
	if userBalance <= data["amount"].(float64) {
		rs.State = STATE_FAILD
		rs.Msg = fmt.Sprintf("You can transfer a maximum of %f USDT", userBalance)
		return rs
	}

	rest := userBalance - data["amount"].(float64)
	_, errmsg := config.GlobalDB.UpdateData(DB_TABLE_USER, db.DB_PARAMS{"credit": rest}, db.DB_PARAMS{"id": uid})
	_, errmsg = config.GlobalDB.InsertData(DB_TABLE_CROSS_TRADE, data)
	if errmsg != nil {
		rs.State = STATE_FAILD
		rs.Msg = fmt.Sprintf("You can transfer a maximum of %f USDT", userBalance)
	}
	return rs
}

func (m *UserModel) Welcome() *WelcomeResponse {
	rs := new(WelcomeResponse)
	rs.State = STATE_SUCCESS
	rs.Msg = "success"
	welcomeInfo := &WelcomeInfo{
		DirectWithdraw: "0",
		LinkWallet:     "0",
	}
	list, _ := config.GlobalDB.FetchAll(DB_TABLE_SYSTEMCONFIG, db.DB_PARAMS{}, db.DB_FIELDS{})
	for _, item := range list {
		switch item["key"].ToString() {
		case "sitename":
			welcomeInfo.PlatformName = item["value"].ToString()
		case "domain":
			welcomeInfo.WelcomePage = item["value"].ToString()
		case "vip_contact":
			welcomeInfo.VIP = item["value"].ToString()
		case "direct_withdraw":
			welcomeInfo.DirectWithdraw = item["value"].ToString()
		case "link_wallet":
			welcomeInfo.LinkWallet = item["value"].ToString()
		}
	}
	rs.WelcomeInfo = welcomeInfo
	return rs
}

func (m *UserModel) ChangeMode(uid int) *BaseResponse {
	uinfo := m.GetBaseInfo(uid)
	if uinfo != nil {
		if uinfo.Mode == 1 {
			m.Update(uid, db.DB_PARAMS{"mode": 2})
		} else {
			m.Update(uid, db.DB_PARAMS{"mode": 1})
		}
	}
	return &BaseResponse{State: STATE_SUCCESS, Msg: ""}
}

func (m *UserModel) GetExplodeState(uid int) *BaseResponse {
	one, _ := config.GlobalDB.FetchOne(DB_TABLE_USER, db.DB_PARAMS{"id": uid}, db.DB_FIELDS{"explode_state"})
	if one == nil {
		return &BaseResponse{State: 0, Msg: ""}
	}
	return &BaseResponse{State: one["explode_state"].ToInt(), Msg: ""}
}

func (m *UserModel) ConvertMoney(uid int) db.DB_PARAMS {
	minnermoney, _ := config.GlobalDB.FetchOne(DB_TABLE_MINING_ORDER, db.DB_PARAMS{"state": 0, "uid": uid}, db.DB_FIELDS{"SUM(amount) as amount"})
	contract, _ := config.GlobalDB.FetchOne(DB_TABLE_OPENED_TRADE, db.DB_PARAMS{"uid": uid, "trade_type": "2", "clear_time": 0}, db.DB_FIELDS{"SUM(credit) as amount"})
	explode, _ := config.GlobalDB.FetchOne(DB_TABLE_OPENED_TRADE, db.DB_PARAMS{"uid": uid, "trade_type": "1", "clear_time": 0}, db.DB_FIELDS{"SUM(credit) as amount"})
	return db.DB_PARAMS{
		"minner":   minnermoney.Get("amount").ToFloat(),
		"contract": contract.Get("amount").ToFloat(),
		"explode":  explode.Get("amount").ToFloat(),
	}
}
