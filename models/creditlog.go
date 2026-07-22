package models

// 账变信息包
const (
	COIN_LOG_USER_RECHARGE         = 5100001 //用户充值
	COIN_LOG_USER_WITHDRAW         = 5100002 //用户提现
	COIN_LOG_USER_PROFIT           = 5100003 //用户矿机收益
	COIN_LOG_USER_CLOSE            = 5100004 //用户平仓
	COIN_LOG_USER_DELEGATE         = 5100005 //用户委托
	COIN_LOG_USER_DELEGATE_SUCCESS = 5100010 //用户委托成功
	COIN_LOG_USER_CANCLE           = 5100007 //用户撤单
	COIN_LOG_USER_BUY_MINING       = 5100008 //用户购买矿机
	COIN_LOG_USER_MINING_PROFIT    = 5100009 //用户挖矿获利

	COIN_LOG_USER_CLEAR_INCOME    = 5100011 //用户提取下级返利收入
	COIN_LOG_USER_WITHDRAW_FAILD  = 5100012 //用户提现失败
	COIN_LOG_USER_EXCHANGE        = 5100013 //用户兑换
	COIN_LOG_BB_TRADE             = 5100014 //币币交易
	COIN_LOG_EXPLODE_TRADE        = 5100015 //交割交易
	COIN_LOG_KEEP_TRADE           = 5100016 //永续交易
	COIN_LOG_BACKEND              = 5100017 //后台划转
	COIN_LOG_USER_MINING_BACK     = 5100018 //用户矿机本金返还
	COIN_LOG_ASSETS_EXCHANGE      = 5100019 //币种兑换
	COIN_LOG_MINING_UNLOCK        = 5100020 //矿机解锁
	COIN_LOG_BUY_COIN             = 5100021 //用户申购新币
	COIN_LOG_EXCHANGE_ACCOUNT_IN  = 5100022 //用户资产转换 资金账户到合约账户
	COIN_LOG_EXCHANGE_ACCOUNT_OUT = 5100023 //用户资产转换 合约账户到资金账户
	COIN_LOG_LOAN_BACK            = 5100025 //归还贷款
	COIN_LOG_LORA_IN              = 5100024 //贷款成功
	COIN_LOG_KEEP_BREAK           = 5100026 //杠杆穿仓
	COIN_LOG_USER_REVERVATION     = 5100027 // 预约扣除金额
	TEAM_LOG_RECHARGE             = 5200001 //用户充值-团队
	TEAM_LOG_WITHDRAW             = 5200002 //提现-团队
	TEAM_LOG_MINING               = 5200003 //挖矿总投入
	TEAM_LOG_MINING_PROFIT        = 5200004 //挖矿总收益
	TEAM_LOG_TRADE                = 5200005 //交易总投入
	TEAM_LOG_TRADE_PROFIT         = 5200006 //交易总收益
	COIN_LOG_SPOT_BACK            = 5200028 //现货申购退还

	LOG_TIMETYPE_ALL   = 0 //全部
	LOG_TIMETYPE_DAY   = 1 //当日
	LOG_TIMETYPE_MONTH = 2 //当月

	INCOME_TYPE_RECHARGE      = 1 //下级充值返利
	INCOME_TYPE_MINING_BUY    = 2 //下级购买矿机返利
	INCOME_TYPE_MINING_PROFIT = 3 //下级矿机收入返利
)

type CreditLogModel struct {
	ModelBase
}
type CreditLogInfo struct { //用户账变添加结构
	Uid        int     `json:"uid"`
	Credit     float64 `json:"credit"`
	LockCredit float64 `json:"lockcredit"`
	Mode       int     `josn:"mode"`
	Sn         string  `josn:"sn"`
	Type       int     `json:"type"`
	CoinType   string  `json:"cointype"`
	Createtime int     `json:"credittime"`
}
type QueueCreditLog struct { //用户账变入队结构
	Credit     float64 `json:"credit"`
	LockCredit float64 `json:"lockcredit"`
	CoinType   string  `json:"cointype"`
	Sn         string  `json:"sn"`
	CreateTime int     `json:"createtime"`
}
type QueueTeamLog struct { //团队统计入队结构
	Recharge            float64 `json:"recharge"`             //个人充值量
	WithDraw            float64 `json:"withdraw"`             //个人提现量
	Trade               float64 `json:"trade"`                //个人交易量
	TradeProfit         float64 `json:"trade_profit"`         //个人交易获利
	MiningCount         float64 `json:"mining_count"`         //个人矿机投入额度
	MiningProfit        float64 `json:"mining_profit"`        //个人矿机获利
	CreateTime          int     `json:"createtime"`           //发生时间
	TradeBB             float64 `json:"trade_bb"`             //币币交易
	TradeExplode        float64 `json:"trade_explode"`        //交割交易
	TradeKeep           float64 `json:"trade_keep"`           //永续合约
	TradeBB_Profit      float64 `json:"trade_bb_profit"`      //币币利润
	TradeExplode_Profit float64 `json:"trade_explode_profit"` //交割利润
	TradeKeep_Profit    float64 `json:"trade_keep_profit"`    //永续利润
	//IsFirstRecharge int     `json:"is_first_recharge"` //是否首充
	/*RegisterNum          int     `json:"register_num"`            //下级注册数量
	ProRegisterNum       int     `json:"pro_register_num"`        //下级有效注册量
	DirectRegisterNum    int     `json:"direct_register_num"`     //直属注册量
	DirectProRegisterNum int     `json:"direct_pro_register_num"` //直属有效注册
	TeamRecharge         float64 `json:"team_recharge"`           //团队充值总额
	TeamWithDraw         float64 `json:"team_withdraw"`           //团队提现总额
	TeamMining           float64 `json:"team_mining"`             //团队挖矿投入总额
	TeamMiningProfit     float64 `json:"team_mining_profit"`      //团队挖矿总利润
	TeamTrade            float64 `json:"team_trade"`              //团队交易总投入
	TeamTradeProfit      float64 `json:"team_trade_profit"`       //团队交易总获利*/
}
type CoinLogRequest struct {
	PageBaseRequest
	Type int `json:"type"`
}
