package models

import "cointrade/lib/db"

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
