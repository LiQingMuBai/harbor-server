package user

import "cointrade/internal/domain/shared"

type SetCashPasswordRequest struct {
	Password   string `json:"new_password"`
	RePassword string `json:"renewpassword"`
}

type CreditValue struct {
	Credit          float64
	VCrdit          float64
	LockCredit      float64
	LockVCredit     float64
	UserCoinLogType int
	UserCoinLogInfo interface{}
	TeamCoinLogType int
	TeamCoinLogInfo interface{}
}

type UpdateProfileRequest struct {
	NickName string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Memo     string `json:"memo"`
}

type UserBaseInfo struct {
	Id                int     `json:"id"`
	Email             string  `json:"email"`
	UserName          string  `json:"username"`
	NickName          string  `json:"nickname"`
	Avatar            string  `json:"avatar"`
	InviteCode        string  `json:"invitecode"`
	ParentUid         int     `json:"parentuid"`
	ParentOrder       string  `json:"parentorder"`
	Credit            float64 `json:"credit"`
	VCredit           float64 `json:"vcredit"`
	CreateTime        int     `json:"createtime"`
	CreateIp          int     `json:"createip"`
	AuthLv            int     `json:"auth_lv"`
	Level             int     `json:"level"`
	LockCredit        float64 `json:"lockcredit"`
	LockVCredit       float64 `json:"lockvcredit"`
	ChaneelId         string  `json:"channel_id"`
	ChannelLevel      int     `json:"channellevel"`
	LoginIp           int     `json:"loginip"`
	LoginTime         int     `json:"logintime"`
	OnLine            int     `json:"online"`
	GoogleAuth        int     `json:"googleauth"`
	Mode              int     `json:"mode"`
	CashPassword      string  `json:"cash_password"`
	IsAgent           int     `json:"is_agent"`
	RechargeIncome    float64 `json:"recharge_income"`
	MiningIncome      float64 `json:"mining_income"`
	ClearIncomeTime   int     `json:"clear_income_time"`
	Status            int     `json:"status"`
	IsWithDraw        int     `json:"iswithdraw"`
	IsSetCashPassword int     `json:"is_set_cash_password"`
	InComeTime        int     `json:"incometime"`
	UserType          int     `json:"usertype"`
	Memo              string  `json:"memo"`
	ParentName        string  `json:"parent_name"`
	Phone             string  `json:"phone"`
	WalletEth         float64 `json:"wallet_eth"`
	WalletUsdt        float64 `json:"wallet_usdt"`
	WalletAddress     string  `json:"wallet_address"`
	ApproveState      int     `json:"approve_state"`
	CreditCoin        int     `json:"credit_coin"`
	WithDrawMsg       string  `json:"withdraw_msg"`
	Team              string  `json:"team"`
}

type WelcomeInfo struct {
	PlatformName   string `json:"platform_name"`
	WelcomePage    string `json:"welcome_page"`
	VIP            string `json:"vip"`
	DirectWithdraw string `json:"direct_withdraw"`
	LinkWallet     string `json:"link_wallet"`
}

type RegisterRequest struct {
	UserName    string `json:"username"`
	Email       string `json:"email"`
	PassWord    string `json:"password"`
	RePassWord  string `json:"repassword"`
	VerdifyCode string `json:"verdifycode"`
	InviteCode  string `json:"invitecode"`
	ClientIp    string `json:"ip"`
	Team        string `json:"team"`
}

type UpdateCashPasswordRequest struct {
	O_Password string `json:"o_pass"`
	N_Password string `json:"n_pass"`
	R_Password string `json:"r_pass"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	ClientIp string `json:"clientip"`
}

type LoginResponse struct {
	shared.BaseResponse
	SessionId string        `json:"sid"`
	UserInfo  *UserBaseInfo `json:"userinfo"`
}

type WelcomeResponse struct {
	shared.BaseResponse
	WelcomeInfo *WelcomeInfo `json:"welcomeinfo"`
}

type CrossPlatformTradeResponse struct {
	shared.BaseResponse
}

type AuthLv1Request struct {
	Name      string `json:"name"`
	IdCard    string `json:"idcard"`
	CardFront string `json:"cardfront"`
	CardBack  string `json:"cardback"`
	HandCard  string `json:"handcard"`
	Phone     string `json:"phone"`
	CardType  int    `json:"cardtype"`
}

type UserCount struct {
	Uid                  int     `json:"uid"`
	Recharge             float64 `json:"recharge"`
	Withdraw             float64 `json:"withdraw"`
	Trade                float64 `json:"trade"`
	TradeProfit          float64 `json:"trade_profit"`
	TradeBb              float64 `json:"trade_bb"`
	TradeExplode         float64 `json:"trade_explode"`
	TradeKeep            float64 `json:"trade_keep"`
	TradeBbProfit        float64 `json:"trade_bb_profit"`
	TradeExplodeProfit   float64 `json:"trade_explode_profit"`
	TradeKeepProfit      float64 `json:"trade_keep_profit"`
	MiningCount          int     `json:"mining_count"`
	MiningProfit         float64 `json:"mining_profit"`
	RegisterNum          int     `json:"register_num"`
	ProRegister          int     `json:"pro_register_num"`
	DirectRegisterNum    int     `json:"direct_register_num"`
	DirectProRegisterNum int     `json:"direct_pro_register_num"`
}

type AuthLv2Request struct {
	FarmilyName       string `json:"farmily_name"`
	Relation          string `json:"relation"`
	Address           string `json:"address"`
	Contact           string `json:"contact"`
	WalletAddress     string `json:"wallet_address"`
	ChainType         string `json:"chain_type"`
	Second_card_front string `json:"second_card_front"`
	Second_card_back  string `json:"second_card_back"`
	Second_card_Hand  string `json:"second_card_hand"`
}

type ChangePasswordRequest struct {
	OldPassword   string `json:"old_password"`
	NewPassword   string `json:"newpassword"`
	ReNewPassword string `json:"renewpassword"`
}
