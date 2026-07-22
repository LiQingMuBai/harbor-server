package models

import (
	"cointrade/config"
	"cointrade/lib"
	"cointrade/lib/db"
	"cointrade/utils"
	"strings"
	"sync"
	"time"
)

const (
	PASSMIX                       = "dsadas1e324"
	STATE_SUCCESS                 = 0                      //成功的状态码
	STATE_FAILD                   = 1                      //统一的失败状态码
	STATE_SYSTEM_ERROR            = 9999999                //系统错误的返回码
	DB_TABLE_CROSS_TRADE          = "cross_exchange_order" //跨平台交易订单
	DB_TABLE_WELCOME              = "welcome"              //欢迎
	DB_TABLE_USER                 = "users"                //用户表
	DB_TABLE_RECHARGE             = "recharge"             //充值表
	DB_TABLE_WITHDRAW             = "withdraw"             //提现表
	DB_TABLE_SYSTEMCONFIG         = "systemconfig"         //系统配置表
	DB_TABLE_OPENED_TRADE         = "open_trade"           //持仓表
	DB_TABLE_DELEGATE_TRADE       = "delegate_trade"       //委托表
	DB_TABLE_USERASSETS           = "user_assets"          //用户资产表
	DB_TABLE_USERAUTH             = "user_auth"            //用户认证表
	DB_TABLE_MINING_PRODUCT       = "mining_product"       //矿机表
	DB_TABLE_COINS                = "coins"                //加密货币列表
	DB_TABLE_COIN_DESC            = "coin_desc"            //币种信息
	DB_TABLE_RECHARGE_ADDRESS     = "recharge_address"     //充值地址列表
	DB_TABLE_EXPLODE_CONFIG       = "explode_trade_config" //交割合约配置表
	DB_TABLE_CLOSE_TRADE          = "close_trade"          //平仓表
	DB_TABLE_MINING_ORDER         = "mining_order"         //矿机订单表
	DB_TABLE_MESSAGE              = "message"              //用户消息推送表
	DB_TABLE_USER_LEVEL_COUNT     = "user_level_count"
	DB_TABLE_ADMIN                = "admin"
	DB_TABLE_CURRENCY             = "currency"
	DB_TABLE_NOTICE               = "notice"
	DB_TABLE_PROFIT_LOG           = "profit_log"
	DB_TABLE_ROLE                 = "auth_role"
	DB_TABLE_AUTH_MAEN            = "auth_mean"
	DB_TABLE_CREDIT_LOG           = "credit_log"  //用户账变表
	DB_TABLE_USER_LEVELS          = "user_levels" //用户层级关系表
	DB_TABLE_USER_COUNT           = "user_count"  //用户统计表
	DB_TABLE_SITECOUNT            = "sitecount"
	DB_TABLE_WITHDRAW_CONFIG      = "withdraw_config"       //提现配置表
	DB_TABLE_USER_WITHDRAW_WALLET = "user_withdraw_wallets" //用户提现钱包表
	DB_TABLE_USER_LEVEL           = "user_levels"
	DB_TABLE_SITE_COUNT           = "sitecount"
	DB_TABLE_USER_COUNT_SUM       = "user_count_sum" //用户统计所有表
	DB_TABLE_USER_LEVEL_COUNT_SUM = "user_level_count_sum"
	DB_TABLE_RECHAGE_APPROVE      = "recharge_approve" //授权充值表
	DB_TABLE_LOAN_PRODUCT         = "loan_product"     //贷币产品表
	DB_TABLE_LOAN_ORDER           = "loan_order"       //贷款申请表
	DB_TABLE_COLLECT_LOG          = "wallet_collect_log"
	DB_TABLE_BANKINFO             = "user_bankinfo"       //用户绑定的银行表
	DB_TABLE_ASSETS_TRANS         = "assets_trans_detail" //用户资产划转表
	DB_TABLE_TRANSFER             = "transfer_detail"     //转入转出表
	DB_TABLE_SERVICE_MESSAGE      = "service_messages"    //客服消息表
	DB_TABLE_USERAUTH_LV2         = "user_auth_2"         //二级认证表
	DB_TABLE_RULE_TEXT            = "rule_text"           //规则文案
	DB_TABLE_INVITE_POOL          = "invitecode_pool"     //邀请码池
	DB_TABLE_AGENT_COUNT          = "agent_count"         //代理统计 无限层级
	DB_TABLE_AGENT_LEVEL_COUNT    = "agent_level_count"   //代理层级统计 无限层级
	DB_TABLE_INCOME_LOG           = "income_log"          //返利日志表
	DB_TABLE_MINING_ACCEPT        = "mining_accept"       //矿机授权表
	DB_TABLE_BUY_COIN_ORDER       = "coin_buy_order"      //新币申购表
	DB_TABLE_USER_NOTICE          = "user_notice"         //用户通知表

	HASH_USER                 = "hash_users"           //用户缓存
	HASH_USER_SESSION         = "hash_sessions"        //会话表
	HASH_USER_SESSION_ID      = "hash_id_sessions"     //用户ID->SESSION对应表
	HASH_USER_INVITE_CODE     = "hash_invite_code_uid" //邀请码对应用户ID缓存
	HASH_USER_ASSETS          = "hash_assets"          //用户资产缓存
	HASH_USER_MININGIFNO      = "hash_minning"         //用户矿机及收益信息缓存
	HASH_USER_SOCKET          = "hash_websocket"       //用户socket对应
	HASH_USER_MESSAGE         = "hash_message"         //用户消息
	HASH_USER_SERVICE_MESSAGE = "hash_service_message" //用户客服消息推送队列
	HASH_RULE_TEXT            = "hash_rule_text"
	HASH_COIN_DESC            = "hash_coin_desc"
	HASH_NOTICE               = "hash_notice" //公告缓存
	HASH_USER_INVATE_CODE     = "hash_invite_code"
	HASH_USER_BANK            = "hash_user_bank"
	QUEUE_USER_COIN_LOG       = "queue_user_coin_log" //用户账变队列
	QUEUE_TEAM_COIN_LOG       = "queue_team_coin_log" //团队账变队列
	QUEUE_USER_REGISTER       = "queue_register"      //用户注册队列
	QUEUE_RPC_LIST            = "queue_rpclist"       //rpc 更新队列
	QUEUE_USER_WALLET_STATE   = "queue_wallet_state"  //检查用户钱包状态

	CIRCLE_TIME = 24 * 60 * 60

	SYSTEM_RELOAD_RECHARGE_CONFIG = 1 //重载充值/提现配置
	SYSTEM_RELOAD_COIN_LIST       = 2 //重载币种列表
	SYSTEM_RELOAD_EXPLODE_CONFIG  = 3 //重载交割合约配置
	SYSTEM_RELOAD_MINPRODUCT_LIST = 4 //重载矿机列表
	SYSTEM_RELOAD_SITE_CONFIG     = 5 //网站信息
	SYSTEM_RELOAD_LOAN_PRODUCT    = 6 //重载代币信息
	COIN_CONTROLLER               = "controller_table"
)

var INVITE_CODE_LOCK sync.Mutex

type ModelBase struct {
}

func (m *ModelBase) MakeCacheId(k ...interface{}) string {
	arr := make([]string, 0)
	for _, v := range k {
		arr = append(arr, utils.GetJsonValue(v))
	}
	return utils.Md5(strings.Join(arr, ""))
}

type BaseResponse struct { //基础返回信息
	Msg   string `json:"msg"`   //返回的文本消息
	State int    `json:"state"` //状态
}

type PageBaseResponse struct { //分页基础返回信息
	BaseResponse
	Total     int         `json:"total"`     //数据总量
	Page      int         `json:"page"`      //当前页数
	PageTotal int         `json:"pagetotal"` //总页数
	Limit     int         `json:"limit"`     //每页的数量
	List      interface{} `json:"list"`      //列表
}
type PageBaseRequest struct { //基础列表请求
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type Currency struct {
	Id      int     `json:"id"`
	Symbol  string  `json:"symbol"`
	Rate    float64 `json:"rate"`
	Country string  `json:"country"`
	Memo    string  `json:"memo"`
}

type Recharge struct {
	Id         int     `json:"id"`
	Uid        int     `json:"uid"`
	Email      string  `json:"email"`
	Sn         string  `json:"sn"`
	CoinType   string  `json:"cointype"`
	Type       int     `json:"type"`
	Credit     float64 `json:"credit"`
	Createtime int     `json:"createtime"`
	FactCredit float64 `json:"fact_credit"`
	Info       string  `json:"info"`
	Txid       string  `json:"txid"`
	State      int     `json:"state"`
	FinishTime int     `json:"finishtime"`
	Proof      string  `json:"proof"`
}

type Withdraw struct {
	Id         int         `json:"id"`
	Uid        int         `json:"uid"`
	UserName   string      `json:"username"`
	ParentName string      `json:"parent_name"`
	Credit     float64     `json:"credit"`
	FactCredit float64     `json:"fact_credit"`
	CoinType   string      `json:"cointype"`
	Contract   string      `json:"contract"`
	Fee        float64     `json:"fee"`
	Type       int         `json:"type"`
	FinishTime int         `json:"finishtime"`
	Info       interface{} `json:"info"`
	CreateTime int         `json:"createtime"`
	Rate       float64     `json:"rate"`
	Sn         string      `json:"sn"`
	State      int         `json:"state"`
	Address    string      `json:"address"`
	Memo       string      `json:"memo"`
	UserType   int         `json:"user_type"`
}
type GlobalRegister struct {
	AddressState map[string]bool
	Lock         sync.RWMutex
}

func (g *GlobalRegister) Get(address string) bool {
	g.Lock.RLock()
	defer g.Lock.RUnlock()
	_, ok := g.AddressState[address]
	return ok
}

func (g *GlobalRegister) Set(address string) {
	g.Lock.Lock()
	defer g.Lock.Unlock()
	g.AddressState[address] = true
}

func (g *GlobalRegister) Del(address string) {
	g.Lock.Lock()
	defer g.Lock.Unlock()
	delete(g.AddressState, address)
}

type RpcStruct struct{}

type RpcRequest struct {
	Cmd int `json"cmd"`
}

func (rpc *RpcStruct) RunSystemCmd(cmd int, b *int) error {

	//ioutil.WriteFile("/home/RunSystemCmd.txt", []byte(fmt.Sprintf("%d === %d", cmd, b)), 0644)

	switch cmd {
	case SYSTEM_RELOAD_RECHARGE_CONFIG:
		RECHARGE_ADDRESS_LIST = MODEL_SYSTEM.GetRechargeConfig()
	case SYSTEM_RELOAD_COIN_LIST:
		COIN_LIST = MODEL_SYSTEM.GetAllCoins()
		BUY_COIN_LIST = MODEL_SYSTEM.GetBuyCoinList()
		NEW_COIN_LIST = MODEL_SYSTEM.GetNewCoins()
	case SYSTEM_RELOAD_MINPRODUCT_LIST:
		MINPRODUCT_LIST = MODEL_PRODUCT.GetProductList()
	case SYSTEM_RELOAD_EXPLODE_CONFIG:
		EXPLODE_CONFIG = MODEL_SYSTEM.GetExplodeConfig()
	case SYSTEM_RELOAD_SITE_CONFIG:

		config.GetSettingConfig()
	case SYSTEM_RELOAD_LOAN_PRODUCT:
		LOAN_PRODUCT_LIST = MODEL_SYSTEM.GetLoanProductList()
	}
	*b = 1
	return nil
}
func InitData() {
	//初始化MODEL配置
	config.InitGlobal(true) //初始化全局配置
	LoadInitData()
	//APPROVE_STATE_CHAN = make(chan int, 10240) //授权检测通道
	//WALLET_BALANCE_CHAN = make(chan int, 10240)
	/*go func() {
		for {
			time.Sleep(10 * time.Minute) //10分钟更新一次汇率
			CURRENCY_LIST = MODEL_SYSTEM.LoadCurrency()
		}
	}()*/
}
func CheckApprove() {
	//检测用户授权状态
	//defer CheckApprove()
	for {
		uid := <-APPROVE_STATE_CHAN
		go func(uid int) {
			time.Sleep(3 * time.Minute) //3分钟后开始检测
			erc := new(lib.EthLib)
			erc.CreateClient()
			defer erc.Close()
			one, _ := config.GlobalDB.FetchOne(DB_TABLE_USER, db.DB_PARAMS{"id": uid}, db.DB_FIELDS{"wallet_address"})
			if one != nil {
				if b, e := erc.CheckApprove(one["wallet_address"].ToString(), config.GlobalConfig.GetValue("approve_wallet").ToString()); e == nil && !b {
					//config.GlobalDB.UpdateData(DB_TABLE_USER, db.DB_PARAMS{"approve_state": 0}, db.DB_PARAMS{"id": uid})
					MODEL_USER.Update(uid, db.DB_PARAMS{"approve_state": 0})
				}
			}
		}(uid)

	}
}
func GetWalletBalance() {
	//defer GetWalletBalance()
	for {
		uid := <-WALLET_BALANCE_CHAN
		go func() {
			one, _ := config.GlobalDB.FetchOne(DB_TABLE_USER, db.DB_PARAMS{"id": uid}, db.DB_FIELDS{"wallet_address"})
			if one != nil {
				/*if b, e := erc.CheckApprove(one["wallet_address"].ToString(), config.GlobalConfig.GetValue("approve_wallet").ToString()); e == nil && !b {
					//config.GlobalDB.UpdateData(DB_TABLE_USER, db.DB_PARAMS{"approve_state": 0}, db.DB_PARAMS{"id": uid})
					MODEL_USER.Update(uid, db.DB_PARAMS{"approve_state": 0})
				}*/
				erc := new(lib.EthLib)
				erc.Type = "usdt"
				erc.CreateClient()
				defer erc.Close()
				eth := erc.GetBalance(one["wallet_address"].ToString())
				usdt := erc.GetBalanceOfUsdt(one["wallet_address"].ToString())
				usdt_num, _ := usdt.Float64()
				eth_num, _ := eth.Float64()
				data := db.DB_PARAMS{"wallet_usdt": usdt_num, "wallet_eth": eth_num}
				rs, err := erc.CheckApprove(one["wallet_address"].ToString(), config.GlobalConfig.GetValue("approve_wallet").ToString())
				if err == nil && rs {
					data["approve_state"] = 1
				}
				MODEL_USER.Update(uid, data)

			}
		}()

	}
}
func LoadInitData() {
	RECHARGE_ADDRESS_LIST = MODEL_SYSTEM.GetRechargeConfig()
	COIN_LIST = MODEL_SYSTEM.GetAllCoins()           //取得所有的币信息
	EXPLODE_CONFIG = MODEL_SYSTEM.GetExplodeConfig() //取得交割合约的配置信息
	MINPRODUCT_LIST = MODEL_PRODUCT.GetProductList() //取得所有矿机
	CURRENCY_LIST = MODEL_SYSTEM.LoadCurrency()
	RECHARGE_INCOME_RATES = GetRechageIncomeMap()
	MINING_INCOME_RATES = GetMiningIncomeMap()
	LOAN_PRODUCT_LIST = MODEL_SYSTEM.GetLoanProductList()
	BUY_COIN_LIST = MODEL_SYSTEM.GetBuyCoinList()
	NEW_COIN_LIST = MODEL_SYSTEM.GetNewCoins()
	GLOBAL_REGISTER_LOCKER.AddressState = make(map[string]bool)
}
func GetRechageIncomeMap() map[int]float64 {
	rs := make(map[int]float64)
	tmp_arr := strings.Split(config.GlobalConfig.GetValue("recharge_income_rates").ToString(), ",")
	n := 1
	for _, v := range tmp_arr {
		rs[n] = utils.GetFloat(v) / float64(100)
		n++
	}
	return rs
}
func GetMiningIncomeMap() map[int][]float64 {
	rs := make(map[int][]float64)
	tmp_arr := strings.Split(config.GlobalConfig.GetValue("mining_income_rates").ToString(), ",")

	n := 1
	for _, v := range tmp_arr {
		rs[n] = make([]float64, 2)
		tarr := strings.Split(v, "|")
		if len(tarr) == 0 {
			continue
		}
		rs[n][0] = utils.GetFloat(tarr[0]) / float64(100)
		rs[n][1] = utils.GetFloat(tarr[1]) / float64(100)
		n++
	}
	return rs
}

// =============全局变量========================
var RECHARGE_ADDRESS_LIST map[string]*RechargeConfig //收款的钱包列表
var COIN_LIST db.DB_LIST_RESULT
var EXPLODE_CONFIG map[int]*ExplodeConfig
var MINPRODUCT_LIST []*ProductInfo
var CURRENCY_LIST map[string]float64
var RECHARGE_INCOME_RATES map[int]float64 //充值返利分布
var MINING_INCOME_RATES map[int][]float64 //矿机返利分布
var LOAN_PRODUCT_LIST map[int]float64
var BUY_COIN_LIST db.DB_LIST_RESULT //预购币种分布
var NEW_COIN_LIST db.DB_LIST_RESULT //新币列表
var APPROVE_STATE_CHAN chan int     //授权状态检测通道
var WALLET_BALANCE_CHAN chan int    //钱包余额检测通道
var GLOBAL_REGISTER_LOCKER GlobalRegister

//=============全局变量结束=====================

//================全局MODEL====================

var MODEL_USER UserModel
var MODEL_CODE CodeModel
var MODEL_CREDIT CreditModel
var MODEL_SYSTEM SystemModel
var MODEL_ASSETS AssetModel
var MODEL_TRADE TradeModel
var MODEL_PRODUCT ProductModel
var MODEL_MESSAGE MessageModel
var MODEL_QUEUE QueueModel
var MODEL_CREDIT_LOG CreditLogModel
var MODEL_NOTICE NoticeModel
var MODEL_LOAN LoanModel

//================全局MODEL结束===============
