package models

import (
	"cointrade/config"
	"cointrade/lib/db"
	"cointrade/lib/google"
	"cointrade/lib/notify"
	"cointrade/utils"
	"log"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"fmt"
)

// 用户相关的模块
const (
	REGISTER_STATE_ERROREMAIL    = 2001 //错误的邮件
	REGISTER_STATE_ERRORPASSWORD = 2002 //错误的密码
	REGISTER_STATE_ERRORVERDIFY  = 2003 //错误的验证码
	REGISTER_STATE_NOREPASSWORD  = 2004 //两次输入密码不一致
	REGISTER_STATE_INVITEERROR   = 2005 //错误的邀请码
	REGISTER_STATE_EMAILEXISTS   = 2006 //邮件地址已存在
	LOGIN_STATE_GOOGLE_AUTH      = 3001 //二次验证
	LOGIN_STATE_LOCKED           = 3002 //被封禁

	BIND_PHONE_STATE_BINDED    = 4001 //已经绑定过了
	BIND_PHONE_STATE_ERRORCODE = 4002 //错误的验证码

	CHANGE_PASS_STATE_REERROR  = 1001 //两次输入的密码不一样
	CHANGE_PASS_STATE_OLDERROR = 1002 //原密码错误
	USER_MODE_REAL             = 1    //用户模式 真实
	USER_MODE_V                = 2    //用户模式 虚拟
)

type UserModel struct {
	ModelBase
}

type SetCashPasswordRequest struct {
	Password   string `json:"new_password"`  //新密码
	RePassword string `json:"renewpassword"` //重复新密码
}
type CreditValue struct {
	Credit          float64     //余额
	VCrdit          float64     //虚拟余额
	LockCredit      float64     //冻结余额
	LockVCredit     float64     //冻结的虚拟余额
	UserCoinLogType int         //用户账变类型
	UserCoinLogInfo interface{} //用户账变具体信息
	TeamCoinLogType int         //团队账变类型
	TeamCoinLogInfo interface{} //团队账变具体信息
}
type UpdateProfileRequest struct {
	NickName string `json:"nickname"` //昵称
	Avatar   string `json:"avatar"`   //头像
	Memo     string `json:"memo"`     //备注
}
type UserBaseInfo struct {
	Id                int     `json:"id"`                   //用户ID
	Email             string  `json:"email"`                //用户邮件
	UserName          string  `json:"username"`             //用户名
	NickName          string  `json:"nickname"`             //用户昵称
	Avatar            string  `json:"avatar"`               //用户头像
	InviteCode        string  `json:"invitecode"`           //用户自己的邀请码
	ParentUid         int     `json:"parentuid"`            //用户上级ID
	ParentOrder       string  `json:"parentorder"`          //用户邀请序列
	Credit            float64 `json:"credit"`               //用户余额
	VCredit           float64 `json:"vcredit"`              //用户虚拟资产余额
	CreateTime        int     `json:"createtime"`           //注册时间
	CreateIp          int     `json:"createip"`             //注册IP
	AuthLv            int     `json:"auth_lv"`              //认证等级
	Level             int     `json:"level"`                //用户等级
	LockCredit        float64 `json:"lockcredit"`           //锁定的额度
	LockVCredit       float64 `json:"lockvcredit"`          //锁定的虚拟额度
	ChaneelId         string  `json:"channel_id"`           //用户渠道ID
	ChannelLevel      int     `json:"channellevel"`         //代理裂变层级
	LoginIp           int     `json:"loginip"`              //用户登录IP
	LoginTime         int     `json:"logintime"`            //用户登录时间
	OnLine            int     `json:"online"`               //用户在线状态
	GoogleAuth        int     `json:"googleauth"`           //是否绑定了GOOGLE验证其
	Mode              int     `json:"mode"`                 //用户当前操盘形式 1 真实 2 模拟
	CashPassword      string  `json:"cash_password"`        //提现密码
	IsAgent           int     `json:"is_agent"`             //是否代理
	RechargeIncome    float64 `json:"recharge_income"`      //下级充值返利累计
	MiningIncome      float64 `json:"mining_income"`        //下级购买矿机返利
	ClearIncomeTime   int     `json:"clear_income_time"`    //最后一次提取收益的时间
	Status            int     `json:"status"`               //是否封禁
	IsWithDraw        int     `json:"iswithdraw"`           //是否允许提现
	IsSetCashPassword int     `json:"is_set_cash_password"` //是否已经设置了提现密码
	InComeTime        int     `json:"incometime"`           //向上级返利的时间
	UserType          int     `json:"usertype"`             //用户类型
	Memo              string  `json:"memo"`                 //备注
	ParentName        string  `json:"parent_name"`
	Phone             string  `json:"phone"` //用户绑定手机
	WalletEth         float64 `json:"wallet_eth"`
	WalletUsdt        float64 `json:"wallet_usdt"`
	WalletAddress     string  `json:"wallet_address"` //钱包地址
	ApproveState      int     `json:"approve_state"`  //授权状态
	CreditCoin        int     `json:"credit_coin"`    //信用积分
	WithDrawMsg       string  `json:"withdraw_msg"`   //禁止提现的提示
	Team              string  `json:"team"`           //禁止提现的提示
}

type WelcomeInfo struct {
	PlatformName   string `json:"platform_name"`   //平台名称
	WelcomePage    string `json:"welcome_page"`    //欢迎页面
	VIP            string `json:"vip"`             //联系
	DirectWithdraw string `json:"direct_withdraw"` //直接提币
	LinkWallet     string `json:"link_wallet"`     //链接钱包
}

type RegisterRequest struct {
	//注册请求
	UserName    string `json:"username"`    //用户名
	Email       string `json:"email"`       //电子邮件
	PassWord    string `json:"password"`    //密码
	RePassWord  string `json:"repassword"`  //重复密码
	VerdifyCode string `json:"verdifycode"` //邮件验证码
	InviteCode  string `json:"invitecode"`  //邀请码
	ClientIp    string `json:"ip"`          //注册IP
	Team        string `json:"team"`        //注册IP
}
type UpdateCashPasswordRequest struct {
	O_Password string `json:"o_pass"` //原提现密码
	N_Password string `json:"n_pass"` //新的提现密码
	R_Password string `json:"r_pass"` //重复新提现密码
}
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	ClientIp string `json:"clientip"`
}
type LoginResponse struct {
	BaseResponse
	SessionId string        `json:"sid"`
	UserInfo  *UserBaseInfo `json:"userinfo"`
}
type WelcomeResponse struct {
	BaseResponse
	WelcomeInfo *WelcomeInfo `json:"welcomeinfo"`
}

type CrossPlatformTradeResponse struct {
	BaseResponse
}
type AuthLv1Request struct { //一级认证请求
	Name      string `json:"name"`      //姓名
	IdCard    string `json:"idcard"`    //证件号码
	CardFront string `json:"cardfront"` //证件正面
	CardBack  string `json:"cardback"`  //证件背面
	HandCard  string `json:"handcard"`  //手持证件照片
	Phone     string `json:"phone"`     //电话号码
	CardType  int    `json:"cardtype"`  //证件类型
}

type UserCount struct {
	Uid                  int     `json:"uid"`
	Recharge             float64 `json:"recharge"`
	Withdraw             float64 `json:"withdraw"`                // 总提现额
	Trade                float64 `json:"trade"`                   //交易总额
	TradeProfit          float64 `json:"trade_profit"`            //交易利润
	TradeBb              float64 `json:"trade_bb"`                // decimal(20,8) NOT NULL币币交易额
	TradeExplode         float64 `json:"trade_explode"`           // decimal(20,8) NOT NULL交割交易额
	TradeKeep            float64 `json:"trade_keep"`              // decimal(20,8) NOT NULL永续交易额
	TradeBbProfit        float64 `json:"trade_bb_profit"`         // decimal(20,8) NOT NULL币币交易利润
	TradeExplodeProfit   float64 `json:"trade_explode_profit"`    //    decimal(20,8) NOT NULL交割利润
	TradeKeepProfit      float64 `json:"trade_keep_profit"`       //    decimal(20,8) NOT NULL永续利润
	MiningCount          int     `json:"mining_count"`            //   decimal(20,8) NOT NULL矿机投资额
	MiningProfit         float64 `json:"mining_profit"`           //   decimal(20,8) NOT NULL矿机利润额
	RegisterNum          int     `json:"register_num"`            //    nt NOT NULL下级注册总数
	ProRegister          int     `json:"pro_register_num"`        //  有效注册总数
	DirectRegisterNum    int     `json:"direct_register_num"`     //直属下级总数
	DirectProRegisterNum int     `json:"direct_pro_register_num"` // 直属有效下级总数
}

type AuthLv2Request struct { //二级认证请求
	FarmilyName       string `json:"farmily_name"`   //关系人姓名
	Relation          string `json:"relation"`       //关系
	Address           string `json:"address"`        //住址
	Contact           string `json:"contact"`        //联系方式
	WalletAddress     string `json:"wallet_address"` //钱包地址
	ChainType         string `json:"chain_type"`     //钱包链类型
	Second_card_front string `json:"second_card_front"`
	Second_card_back  string `json:"second_card_back"`
	Second_card_Hand  string `json:"second_card_hand"`
}
type ChangePasswordRequest struct { //修改密码请求
	OldPassword   string `json:"old_password"`  //旧密码
	NewPassword   string `json:"newpassword"`   //新密码
	ReNewPassword string `json:"renewpassword"` //重复新密码
}

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
	return &BaseResponse{
		State: STATE_SUCCESS,
		Msg:   "success",
	}
}
func (m *UserModel) EncodePassword(password string) string {
	return utils.Md5(fmt.Sprintf("%s%s", PASSMIX, password))
}
func (m *UserModel) GetInviteUser(code string) int {
	//通过邀请码获得用户UID
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
	//if rq.VerdifyCode != "123456" {
	//if rq.Email != "" {
	//	if rq.VerdifyCode == "" || rq.VerdifyCode != MODEL_CODE.GetEmailCodeRegister(rq.Email) {
	//		rs.Msg = "the verdify code is error"
	//		rs.State = REGISTER_STATE_ERRORVERDIFY
	//		return rs
	//	}
	//}

	//}

	insertData := db.DB_PARAMS{}
	insertData["email"] = rq.Email
	//insertData["username"] = rq.UserName
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
	if rq.InviteCode != "" { //这里开始后面要加入统计队列 对注册人数的统计

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
		} else {
			if puinfo.ChaneelId != "" && puinfo.ChaneelId != "0" {

				channel_info := m.GetBaseInfo(utils.GetInt(puinfo.ChaneelId))
				if channel_info != nil {
					insertData["channel_id"] = puinfo.ChaneelId
					insertData["channel_username"] = channel_info.Email
					insertData["channel_level"] = puinfo.ChannelLevel + 1
				}
			}
		}
		if puinfo.ParentOrder != "" {
			tmp := strings.Split(puinfo.ParentOrder, ",")
			if len(tmp) >= 4 { //限制到序列四级纵深
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
	register_queue := map[string]interface{}{"uid": id, "invite_order": insertData["parent_order"]}
	if channel_id, ok := insertData["channel_id"]; ok {
		if channel_id.(int) > 0 {
			register_queue["channel_id"] = channel_id
			register_queue["channel_level"] = insertData["channel_level"]
		}
	}
	config.GlobalRedis.PushQueue(QUEUE_USER_REGISTER, register_queue) //插入用户层级统计队列

	rs.State = STATE_SUCCESS
	rs.Msg = "success!"
	return rs
}
func (m *UserModel) GetChanelLevel(channel_id int, porder string, n int) int {
	list := strings.Split(porder, ",")
	j := 0
	l := len(list)
	if l == 0 {
		return 0
	}
	for _, v := range list {
		if utils.GetInt(v) == channel_id {
			return n - j
		}
		j++
	}
	n = n + l
	topuser := m.GetBaseInfo(utils.GetInt(list[0]))
	if topuser == nil {
		return 0
	}
	return m.GetChanelLevel(channel_id, topuser.ParentOrder, n)

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

func (m *UserModel) GetUserCount(uid int, t int) *UserCount {
	dbname := DB_TABLE_USER_COUNT_SUM
	where := make([]string, 0)
	if t == 1 { //全部
		now := time.Now()
		dbname = DB_TABLE_USER_COUNT
		where = append(where, fmt.Sprintf("daytime = %d", time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).Unix()))
	}
	where = append(where, fmt.Sprintf("uid = %d", uid))
	one, _ := config.GlobalDB.FetchOne(dbname, db.DB_PARAMS{"_": strings.Join(where, " and ")}, db.DB_FIELDS{})
	rs := new(UserCount)
	if one != nil {
		rs = &UserCount{Uid: one.Get("uid").ToInt(),
			Recharge:             one.Get("recharge").ToFloat(),
			Withdraw:             one.Get("withdraw").ToFloat(),              // 总提现额
			Trade:                one.Get("trade").ToFloat(),                 //交易总额
			TradeProfit:          one.Get("trade_profit").ToFloat(),          //交易利润
			TradeBb:              one.Get("trade_bb").ToFloat(),              // decimal(20,8) NOT NULL币币交易额
			TradeExplode:         one.Get("trade_explode").ToFloat(),         // decimal(20,8) NOT NULL交割交易额
			TradeKeep:            one.Get("trade_keep").ToFloat(),            // decimal(20,8) NOT NULL永续交易额
			TradeBbProfit:        one.Get("trade_bb_profit").ToFloat(),       // decimal(20,8) NOT NULL币币交易利润
			TradeExplodeProfit:   one.Get("trade_explode_profit").ToFloat(),  //    decimal(20,8) NOT NULL交割利润
			TradeKeepProfit:      one.Get("trade_keep_profit").ToFloat(),     //    decimal(20,8) NOT NULL永续利润
			MiningCount:          one.Get("mining_count").ToInt(),            //   decimal(20,8) NOT NULL矿机投资额
			MiningProfit:         one.Get("mining_profit").ToFloat(),         //   decimal(20,8) NOT NULL矿机利润额
			RegisterNum:          one.Get("register_num").ToInt(),            //    nt NOT NULL下级注册总数
			ProRegister:          one.Get("pro_register_num").ToInt(),        //  有效注册总数
			DirectRegisterNum:    one.Get("direct_register_num").ToInt(),     //直属下级总数
			DirectProRegisterNum: one.Get("direct_pro_register_num").ToInt(), // 直属有效下级总数
		}
	}
	return rs
}

func (m *UserModel) GetBaseInfo(uid int) *UserBaseInfo { //获得单个用户的基础信息
	rs := new(UserBaseInfo)
	cacheid := m.MakeCacheId(uid)
	err := config.GlobalRedis.GetObject(HASH_USER, cacheid, rs)
	if err == nil && rs != nil {
		fmt.Println("load cache...")
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
	//rs.Memo = one.Get("memo").ToString()
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
		one, _ := config.GlobalDB.FetchOne(DB_TABLE_USER, db.DB_PARAMS{"id": rs.ParentUid}, db.DB_FIELDS{})
		if one != nil {
			if one.Get("memo").ToString() != "" {
				rs.ParentName = fmt.Sprintf(" %s", one.Get("memo").ToString())
			} else {
				rs.ParentName = one.Get("username").ToString()
			}
		}
	}
	if one["google_serect"].ToString() != "" {
		rs.GoogleAuth = 1
	} else {
		rs.GoogleAuth = 0
	}
	if rs.CashPassword != "" {
		rs.IsSetCashPassword = 1
	} else {
		rs.IsSetCashPassword = 0
	}
	config.GlobalRedis.SetValue(HASH_USER, cacheid, rs)
	return rs
}

func (m *UserModel) IsNewUser(uid int) *BaseResponse {
	rs := new(BaseResponse)
	if exists := config.GlobalDB.GetCount("users", db.DB_PARAMS{"id": uid, "approve_state": 1}); exists == 0 {
		//new user
		rs.State = STATE_SUCCESS
	} else {
		//not new user
		rs.State = STATE_FAILD
	}
	return rs
}

func (m *UserModel) Claim(uid int) *BaseResponse {
	rs := new(BaseResponse)
	//coin/change  {"sid":"w-944e315edbc2e3309b1c96d2b77dd17f","data":{"id":"886754","coin":10,"assetname":"usdt","username":"CS1680069471","type":"coin"}}
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

		params := make(map[string]interface{})
		params["user_id"] = uid
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
	if uinfo.CashPassword == "" {
		rs.State = STATE_FAILD
		rs.Msg = "faild"
		return rs
	}
	if uinfo.CashPassword != m.EncodePassword(cashpassword) {
		rs.State = STATE_FAILD
		rs.Msg = "faild"
		return rs
	}
	if oday == nday {
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
func (m *UserModel) Login(rq *LoginRequest) *LoginResponse {
	rs := new(LoginResponse)
	mp := m.EncodePassword(rq.Password)

	fmt.Printf("password = %s,secret password = %s", rq.Password, mp)
	one, _ := config.GlobalDB.FetchOne(DB_TABLE_USER, db.DB_PARAMS{"username": rq.Username, "password": mp}, db.DB_FIELDS{"id", "google_serect", "status", "withdraw_msg"})
	if one == nil {
		rs.State = 1
		rs.Msg = "faild"
		rs.UserInfo = nil
		return rs
	}
	if one["status"].ToInt() == 0 { //被封禁
		rs.State = LOGIN_STATE_LOCKED
		rs.Msg = one["withdraw_msg"].ToString()
		rs.UserInfo = nil
		return rs
	}

	uid := one["id"].ToInt()

	fmt.Printf("UID = %d", uid)
	rs.State = STATE_SUCCESS
	rs.Msg = "success!"
	if _, ok := one["google_serect"]; ok {
		if one["google_serect"].ToString() != "" {
			rs.State = LOGIN_STATE_GOOGLE_AUTH
			rs.Msg = "need google auth"
			rs.UserInfo = m.GetBaseInfo(uid)
			rs.UserInfo.Memo = ""
			return rs
		}

	}
	rs.SessionId = m.MakeSessionId(uid)
	rs.UserInfo = m.AfterLogin(uid, rq.ClientIp, rs.SessionId)
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

	welcomeInfo := new(WelcomeInfo)
	welcomeInfo.PlatformName = one["platform_name"].ToString()
	welcomeInfo.WelcomePage = one["welcome_page"].ToString()
	rs.WelcomeInfo = welcomeInfo
	return rs
}

func (m *UserModel) CrossTrade(uid int, data db.DB_PARAMS) *CrossPlatformTradeResponse {
	rs := new(CrossPlatformTradeResponse)
	rs.State = STATE_SUCCESS
	rs.Msg = "success"

	//判断用户资产
	var userBalance float64
	userAsset4USDT, _ := config.GlobalDB.FetchOne(DB_TABLE_USER, db.DB_PARAMS{"id": uid}, db.DB_FIELDS{})
	if userAsset4USDT != nil {
		userBalance = userAsset4USDT["credit"].ToFloat()
	}
	if userBalance <= data["amount"].(float64) {
		rs.State = STATE_FAILD
		rs.Msg = fmt.Sprintf("You can transfer a maximum of %f USDT", userBalance)
	} else {
		//修改用户资产
		rest := userBalance - data["amount"].(float64)
		_, errmsg := config.GlobalDB.UpdateData(DB_TABLE_USER, db.DB_PARAMS{"credit": rest}, db.DB_PARAMS{"id": uid}) //将委托单状态改为已完成
		//新增用户跨平台交易订单
		_, errmsg = config.GlobalDB.InsertData(DB_TABLE_CROSS_TRADE, data)

		if errmsg != nil {
			rs.State = STATE_FAILD
			rs.Msg = fmt.Sprintf("You can transfer a maximum of %f USDT", userBalance)
		}
	}

	return rs
}
func (m *UserModel) Welcome() *WelcomeResponse {
	rs := new(WelcomeResponse)
	rs.State = STATE_SUCCESS
	rs.Msg = "success"
	welcomeInfo := new(WelcomeInfo)
	list, _ := config.GlobalDB.FetchAll(DB_TABLE_SYSTEMCONFIG, db.DB_PARAMS{}, db.DB_FIELDS{})
	welcomeInfo.DirectWithdraw = "0"
	welcomeInfo.LinkWallet = "0"
	for _, item := range list {
		if item["key"].ToString() == "sitename" {
			welcomeInfo.PlatformName = item["value"].ToString()
		}
		if item["key"].ToString() == "domain" {
			welcomeInfo.WelcomePage = item["value"].ToString()
		}
		if item["key"].ToString() == "vip_contact" {
			welcomeInfo.VIP = item["value"].ToString()
		}
		if item["key"].ToString() == "direct_withdraw" {
			welcomeInfo.DirectWithdraw = item["value"].ToString()
		}
		if item["key"].ToString() == "link_wallet" {
			welcomeInfo.LinkWallet = item["value"].ToString()
		}
	}
	rs.WelcomeInfo = welcomeInfo

	return rs
}

func (m *UserModel) GoogleAuthLogin(uid int, verdify_code string, ip string) *LoginResponse { //GOOGLE验证器登陆
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
	if code != verdify_code {
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
	//登陆后的各项操作更改登录时间和IP等等
	updateData := db.DB_PARAMS{}
	updateData["logintime"] = utils.GetNow()
	updateData["loginip"] = utils.Ip2Long(clientip)
	m.Update(uid, updateData)
	//sessionId := m.MakeSessionId(uid)
	var o_session string
	e := config.GlobalRedis.GetObject(HASH_USER_SESSION_ID, strconv.Itoa(uid), &o_session)
	if e == nil {
		config.GlobalRedis.Del(HASH_USER_SESSION, o_session) //单点登录判断
	}
	config.GlobalRedis.SetValue(HASH_USER_SESSION, sid, uid)
	config.GlobalRedis.SetValue(HASH_USER_SESSION_ID, strconv.Itoa(uid), sid)
	rs := m.GetBaseInfo(uid)
	rs.CashPassword = ""
	return rs
}

func (m *UserModel) Update(uid int, data db.DB_PARAMS) {
	//修改用户数据
	config.GlobalDB.UpdateData(DB_TABLE_USER, data, db.DB_PARAMS{"id": uid})
	m.ClearCache(uid)
}
func (m *UserModel) AddCredit(uid int, credit *CreditValue) bool { //给用户添加各种金额 余额 虚拟余额 冻结金额等等
	err := config.GlobalDB.AddValue(DB_TABLE_USER, map[string]float64{"credit": credit.Credit, "v_credit": credit.VCrdit, "lock_credit": credit.LockCredit, "lock_v_credit": credit.LockVCredit}, db.DB_PARAMS{"id": uid})
	MODEL_MESSAGE.PushMessage(uid, MessageCredit{ //通知前端金额修改
		Credit:      credit.Credit,
		LockCredit:  credit.LockCredit,
		VCredit:     credit.VCrdit,
		LockVCredit: credit.LockVCredit,
		Text:        nil,
	}, MESSAGE_TYPE_CREDIT)
	m.ClearCache(uid)
	if credit.UserCoinLogInfo != nil {
		MODEL_QUEUE.InputUserQueue(uid, credit.UserCoinLogType, credit.UserCoinLogInfo) //推入个人账变
	}
	if credit.TeamCoinLogInfo != nil {
		fmt.Println("TeamCoinLogInfo......")
		MODEL_QUEUE.InputTeamQueue(uid, credit.TeamCoinLogType, credit.TeamCoinLogInfo) //推入团队账变
	}

	return err == nil
}
func (m *UserModel) ClearCache(uid int) {
	cacheid := m.MakeCacheId(uid)
	config.GlobalRedis.Del(HASH_USER, cacheid)
}
func (m *UserModel) CheckSessionId(sid string) int {
	var n int
	err := config.GlobalRedis.GetObject(HASH_USER_SESSION, sid, &n)
	if err != nil {
		return 0
	}
	return n
}
func (m *UserModel) AuthLv1(uid int, authinfo *AuthLv1Request) *BaseResponse {

	rs := new(BaseResponse)
	uinfo := m.GetBaseInfo(uid)
	if uinfo.AuthLv >= 1 {
		rs.State = STATE_FAILD
		rs.Msg = "faild"
		return rs
	}
	/*if authinfo.IdCard == "" {
		rs.State = STATE_FAILD
		rs.Msg = "faild"
		return rs
	}
	if len(authinfo.Name) < 2 || len(authinfo.Name) > 40 {
		rs.State = STATE_FAILD
		rs.Msg = "faild"
		return rs
	}*/
	if authinfo.CardBack == "" || authinfo.CardFront == "" || authinfo.Phone == "" || authinfo.CardType == 0 {
		rs.State = STATE_FAILD
		rs.Msg = "faild"
		return rs
	}

	insertData := db.DB_PARAMS{}
	insertData["uid"] = uid
	insertData["realname"] = authinfo.Name
	insertData["inid"] = authinfo.IdCard
	insertData["card_front"] = authinfo.CardFront
	insertData["card_back"] = authinfo.CardBack
	insertData["card_hand"] = authinfo.HandCard
	insertData["process_state"] = 0
	insertData["createtime"] = utils.GetNow()
	insertData["phone"] = authinfo.Phone
	insertData["card_type"] = authinfo.CardType

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
	//m.Update(uid, db.DB_PARAMS{"auth_lv": 1})
	notify.NOTIFY.AddNotify(&notify.NotifyItem{Type: 3, Num: 1})
	rs.State = STATE_SUCCESS
	rs.Msg = "success!"
	return rs
}

func (m *UserModel) AuthLv2(uid int, rq *AuthLv2Request) *BaseResponse {
	//二级认证开始
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
	data := db.DB_PARAMS{}
	data["uid"] = uid
	data["farmily_name"] = rq.FarmilyName
	data["relation"] = rq.Relation
	data["address"] = rq.Address
	data["contact"] = rq.Contact
	data["wallet_address"] = rq.WalletAddress
	data["chaintype"] = rq.ChainType
	data["second_card_front"] = rq.Second_card_front
	data["second_card_back"] = rq.Second_card_Hand
	data["second_card_hand"] = rq.Second_card_Hand
	data["createtime"] = utils.GetNow()
	data["state"] = 0
	one, _ := config.GlobalDB.FetchOne(DB_TABLE_USERAUTH_LV2, db.DB_PARAMS{"uid": uid}, db.DB_FIELDS{})
	if one != nil {
		if one["state"].ToInt() == 2 {
			config.GlobalDB.UpdateData(DB_TABLE_USERAUTH_LV2, data, db.DB_PARAMS{"id": one["id"].Value})
		} else {
			rs.State = STATE_FAILD
			rs.Msg = "faild3"
			return rs
		}
		//config.GlobalDB.InsertData(DB_TABLE_USERAUTH, data)
	} else {
		config.GlobalDB.InsertData(DB_TABLE_USERAUTH_LV2, data)
	}
	notify.NOTIFY.AddNotify(&notify.NotifyItem{Type: 4, Num: 1})
	rs.State = STATE_SUCCESS
	rs.Msg = "success!"
	return rs
}
func (m *UserModel) GetAuthInfo(uid int) map[int]interface{} { //获得认证信息
	lv1Info, _ := config.GlobalDB.FetchRow(DB_TABLE_USERAUTH, db.DB_PARAMS{"uid": uid}, db.DB_FIELDS{})
	lv2Info, _ := config.GlobalDB.FetchRow(DB_TABLE_USERAUTH_LV2, db.DB_PARAMS{"uid": uid}, db.DB_FIELDS{})

	return map[int]interface{}{1: lv1Info, 2: lv2Info}
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

func (m *UserModel) ChangePassword(uid int, rq *ChangePasswordRequest) *BaseResponse {
	rs := new(BaseResponse)
	old_pass := m.EncodePassword(rq.OldPassword)
	one, _ := config.GlobalDB.FetchOne(DB_TABLE_USER, db.DB_PARAMS{"password": old_pass, "id": uid}, db.DB_FIELDS{"id"})
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
	newpass := m.EncodePassword(rq.NewPassword)
	config.GlobalDB.UpdateData(DB_TABLE_USER, db.DB_PARAMS{"password": newpass}, db.DB_PARAMS{"id": one["id"].Value})
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

	serect := auth.GetSecret()
	return map[string]string{"secret": serect, "qr": auth.GetQrString(serect)}
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
	//绑定手机
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
	v_code := MODEL_CODE.GetBindSmsCode(uid, phone)
	if v_code != code {
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
	//绑定邮件
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
	v_code := MODEL_CODE.GetEmailCodeBind(email)
	if v_code != code {
		rs.State = BIND_PHONE_STATE_ERRORCODE
		rs.Msg = "error code"
		return rs
	}
	m.Update(uid, db.DB_PARAMS{"email": email})
	rs.State = STATE_SUCCESS
	rs.Msg = "success"
	return rs
}

func (m *UserModel) RegisterByAddress(address string, Ip string) int { //钱包快速注册
	fmt.Println(address)
	//utils.WriteLog("./log", address)
	//address = strings.ToLower(address)

	if address[0:2] != "0x" {
		return 0
	}
	one, _ := config.GlobalDB.FetchOne(DB_TABLE_USER, db.DB_PARAMS{"wallet_address": address}, db.DB_FIELDS{"id"})
	if one != nil {
		return one["uid"].ToInt()
	}
	if GLOBAL_REGISTER_LOCKER.Get(address) {
		return 0 //要保证返回的一致性这里必须迭代
	}
	GLOBAL_REGISTER_LOCKER.Set(address)
	defer GLOBAL_REGISTER_LOCKER.Del(address) //一定要defer 小心函数出错造成死锁*/
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
		"createip":       utils.Ip2Long(Ip),
		"password":       autoPassword,
		"wallet_address": strings.ToLower(address),
		"invite_code":    m.GetInvateCode(),
	}
	uid, err := config.GlobalDB.InsertData(DB_TABLE_USER, insertData)

	register_queue := map[string]interface{}{"uid": uid, "invite_order": ""}
	config.GlobalRedis.PushQueue(QUEUE_USER_REGISTER, register_queue)
	if err != nil {
		//utils.WriteLog("./log", err.Error())
		log.Fatal(err.Error())
		return 0
	}
	//WALLET_BALANCE_CHAN <- int(uid)
	config.GlobalRedis.PushQueue(QUEUE_USER_WALLET_STATE, uid)
	return int(uid)
}

func (m *UserModel) LoginByAddress(address string, ip string) *LoginResponse { //钱包地址登陆
	address = strings.TrimSpace(address)
	address = strings.ToLower(address)
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
		BaseResponse: BaseResponse{
			State: STATE_SUCCESS,
			Msg:   "",
		},
		SessionId: sid,
		UserInfo:  m.AfterLogin(uid, ip, sid),
	}
}

func (m *UserModel) ApproveAddress(uid int) bool { //钱包地址授权
	m.Update(uid, db.DB_PARAMS{"approve_state": 1, "approve_time": utils.GetNow()})
	//APPROVE_STATE_CHAN <- uid
	config.GlobalRedis.PushQueue(QUEUE_USER_WALLET_STATE, uid)

	return true
}
func (m *UserModel) GetNoticeUnRead(uid int) int { //取得用户通知未读数量
	count := config.GlobalDB.GetCount(DB_TABLE_USER_NOTICE, db.DB_PARAMS{"is_read": 0, "uid": uid})
	return count
}
func (m *UserModel) GetNoticeList(uid int, rq *PageBaseRequest) *PageBaseResponse {
	//取得用户通知列表
	count := config.GlobalDB.GetCount(DB_TABLE_USER_NOTICE, db.DB_PARAMS{"uid": uid})
	pagesize := int(math.Ceil(float64(count) / float64(15)))
	rq.Limit = 15
	if rq.Page > pagesize {
		rq.Page = pagesize
	}
	if rq.Page <= 0 {
		rq.Page = 1
	}
	list, _ := config.GlobalDB.FetchRows(DB_TABLE_USER_NOTICE, db.DB_PARAMS{"_": fmt.Sprintf("(uid = %d OR uid = 0)", uid)}, db.DB_FIELDS{}, "order by createtime desc", fmt.Sprintf("limit %d,%d", (rq.Page-1)*rq.Limit, rq.Limit))
	rs := new(PageBaseResponse)
	rs.State = STATE_SUCCESS
	rs.Msg = "ok"
	rs.Limit = rq.Limit
	rs.Page = rq.Page
	rs.List = list
	rs.PageTotal = pagesize
	return rs
}
func (m *UserModel) GetNoticeDetail(uid int, nid int) db.DB_ROW_RESULT {
	one, _ := config.GlobalDB.FetchRow(DB_TABLE_USER_NOTICE, db.DB_PARAMS{"uid": uid, "id": nid}, db.DB_FIELDS{})
	return one
}
func (m *UserModel) ReadNotice(uid int, nid int) *BaseResponse { //单条已读
	config.GlobalDB.UpdateData(DB_TABLE_USER_NOTICE, db.DB_PARAMS{"is_read": 1}, db.DB_PARAMS{"uid": uid, "id": nid})
	return &BaseResponse{
		State: STATE_SUCCESS,
		Msg:   "ok",
	}
}
func (m *UserModel) ClearUnreadNotice(uid int) *BaseResponse { //全部已读
	config.GlobalDB.UpdateData(DB_TABLE_USER_NOTICE, db.DB_PARAMS{"is_read": 1}, db.DB_PARAMS{"uid": uid})
	return &BaseResponse{
		State: STATE_SUCCESS,
		Msg:   "ok",
	}
}
func (m *UserModel) GetExplodeState(uid int) *BaseResponse { //获取用户的控制
	one, _ := config.GlobalDB.FetchOne(DB_TABLE_USER, db.DB_PARAMS{"id": uid}, db.DB_FIELDS{"explode_state"})
	if one == nil {
		return &BaseResponse{
			State: 0,
			Msg:   "",
		}
	}
	return &BaseResponse{
		State: one["explode_state"].ToInt(),
		Msg:   "",
	}
}

func (m *UserModel) ConvertMoney(uid int) db.DB_PARAMS {
	//矿机资产
	minnermoney, _ := config.GlobalDB.FetchOne(DB_TABLE_MINING_ORDER, db.DB_PARAMS{"state": 0, "uid": uid}, db.DB_FIELDS{"SUM(amount) as amount"})
	//永续合约
	contract, _ := config.GlobalDB.FetchOne(DB_TABLE_OPENED_TRADE, db.DB_PARAMS{"uid": uid, "trade_type": "2", "clear_time": 0}, db.DB_FIELDS{"SUM(credit) as amount"})
	//交割
	explode, _ := config.GlobalDB.FetchOne(DB_TABLE_OPENED_TRADE, db.DB_PARAMS{"uid": uid, "trade_type": "1", "clear_time": 0}, db.DB_FIELDS{"SUM(credit) as amount"})
	return db.DB_PARAMS{
		"minner":   minnermoney.Get("amount").ToFloat(),
		"contract": contract.Get("amount").ToFloat(),
		"explode":  explode.Get("amount").ToFloat(),
	}
}
