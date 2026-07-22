package adminmodel

import (
	"cointrade/lib/db"
	"cointrade/models"
)

type LoginRequest struct {
	UserName   string `json:"username"`
	Password   string `json:"password"`
	ClientIp   string `json:"client_ip"`
	VerifyCode string `json:"verifycode"`
}

type AdminInfo struct {
	Id            int    `json:"id"`
	UserName      string `json:"username"`
	Password      string `json:"password"`
	Avatar        string `json:"avatar"`
	LastLoginTime string `json:"last_login_time"`
	LastLoginIp   string `json:"last_login_ip"`
	RoleId        int    `json:"role_id"`
	Token         string `json:"token"`
}

type AdminResponse struct {
	State int         `json:"state"`
	Data  interface{} `json:"data"`
}
type LoginResponse struct {
	models.BaseResponse
	Data interface{} `json:"data"`
}

type Notice struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Pos     string `json:"pos"`
	Pic     string `json:"pic"`
	Lang    string `json:"lang"`
	PubTime int    `json:"pubtime"`
	Content string `json:"content"`
}

type Role struct {
	Id       int    `json:"id"`
	RoleName string `json:"role_name"`
	IsSuper  int    `json:"is_surper"`
	RoleIds  string `json:"role_ids"`
}

type Coin struct {
	Id              int     `json:"id"`
	Name            string  `json:"name"`           //'名称',
	Symbol          string  `json:"symbol"`         //'简称 全部小写',
	Pair            string  `json:"pair"`           //'交易对 全部小写',
	Logo            string  `json:"logo"`           //'LOGO',
	Desc            string  `json:"desc"`           // '简介',
	Sort            int     `json:"sort"`           //排序
	OpenCoinCoin    int     `json:"open_coin2coin"` //'币币交易开启',
	OpenTrade       int     `json:"open_trade"`     //'秒合约交易开启',
	IsNative        int     `json:"isnative"`       //'是否是市场币 0 为自发币',
	Dnum            int     `json:"dnum"`           //'价格小数保留位数',
	Address         string  `json:"address"`
	Cnum            int     `json:"cnum"`            //'数量小数保留位数',
	BasePrice       float64 `json:"baseprice"`       //'基础价格 自有币需要基础价格开始浮动',
	MinPriceFloat   float64 `json:"min_price_float"` //'最小价格浮动 自发币有效 跟随价格小数保留位数',
	MaxPriceFloat   float64 `json:"max_price_float"` //'最大价格浮动 自发币有效 跟随价格小数保留位数',
	MaxFloat        float64 `json:"max_float"`       // '最大波动百分比',
	Vpair           string  `json:"vpair"`           //'相对行情字段 自发币参照物',
	Is_F            int     `json:"is_f"`
	IsMarket        int     `json:"is_market"`
	Fprice          float64 `json:"f_price"`
	IsIn            int     `json:"is_in"`
	IsOut           int     `json:"is_out"`
	AllAmount       int     `json:"all_amount"`
	IsNew           int     `json:"is_new"` //是否新币
	PubTime         int     `json:"pubtime"`
	ControlPriceMin float64 `json:"contorl_price_min"`
	ControlPriceMax float64 `json:"contorl_price_max"`
	KLineConfig     string  `json:"kline_config"` //K线控制
}

type RechargeAddress struct {
	Id          int     `json:"id"`
	CoinType    string  `json:"cointype"`     //'币种',
	Contract    string  `json:"contract"`     // '合约',
	Logo        string  `json:"logo"`         // 'LOGO',
	Address     string  `json:"address"`      //'钱包地址',
	State       int     `json:"state"`        // '是否开启',
	Min         float64 `json:"min"`          // '最小充值',
	WithdrawMin float64 `json:"withdraw_min"` //'最小提现'
}

type SiteCount struct {
	Withdraw           float64 `json:"withdraw"`
	Withdraw_num       int     `json:"withdraw_num"`
	Recharge           float64 `json:"recharge"`
	Register_num       int     `json:"register_num"`
	Pro_num            int     `json:"pro_num"`
	Trade              float64 `json:"trade"`
	Trade_profit       float64 `json:"trade_profit"`
	Minning_count      float64 `json:"minning_count"`
	First_recharge     float64 `json:"first_recharge"`
	Minning_profit     float64 `json:"minning_profit"`
	First_recharge_num int     `json:"first_recharge_num"`
	Close_num          int     `json:"close_num"`
	Open_num           int     `json:"open_num"`
	Daytime            int     `json:"daytime"`
}

type CustomMsg struct {
	Uid      int    `json:"uid"`
	Msg      string `json:"msg"`
	Type     string `json:"type"`
	SnId     string `json:"sn_id"`
	IsNotice int    `json:"is_notice"`
}
type UserNoticeMsg struct {
	Id      string `json:"id"`
	Uid     string `json:"uid"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Type    int    `json:"type"`
}
type P map[string]interface{}

func (s *P) Ts() *db.DBValue {
	return &db.DBValue{Value: s}
}

const (
	PARAM_ERROR = 4001
	LOGIN_ERROR = 4002 //登录失败
	LOGOUT      = 50008
	SUCCESS     = 2000
	ERROR       = 2001
)

const (
	HASH_BY_LOGIN_ADMIN = "admin_session"
	OPERATION_PASSWORD  = "666888"
	HASH_AGENT_SESSION  = "agent_session"
)

var SYSTEM_MODEL SystemModel
var MODEL_USER UserModel
var MODEL_TRADE TradeModel
