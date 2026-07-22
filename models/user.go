package models

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
