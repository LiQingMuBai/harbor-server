package models

import (
	"cointrade/config"
	"cointrade/lib"
	"cointrade/lib/db"
	"cointrade/lib/notify"
	"cointrade/utils"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type CreditModel struct {
	ModelBase
}

const (
	CREDIT_TYPE_RECHARGE         = 1
	CREDIT_TYPE_WITHDRAW         = 2
	CREDIT_TYPE_TRANSFER         = 3
	RECHARGE_ORDER_PREFIX        = "R"
	WITHDRAW_ORDER_PREFIX        = "W"
	TRANSFER_ORDER_PREFIX        = "T"
	RECHARGE_STATE_MIN           = 300001 //充值额度小于最小值
	RECHARGE_STATE_ERROR_ADDRESS = 300002 //充值地址错误
	RECHARGE_STATE_ERROR_USER    = 300003 //错误的用户
	RECHARGE_STATE_ERROR_PROOF   = 300004 //没有证明图片

	WIDTHDRAW_STATE_MIN                = 300005 //小于最小提现金额
	WIDTHDRAW_STATE_ERROR_USER         = 300006 //错误的用户
	WIDTHDRAW_STATE_ERROR_ADDRESS      = 300007 //错误的提现地址
	WIDTHDRAW_STATE_NOTENOUGH          = 300008 //余额不足
	WIDTHDRAW_STATE_ERROR_CASHPASSWORD = 300009 //错误的提现密码
	WIDTHDRAW_STATE_ERROR_LOCKED       = 300010 //用户不允许提现
	WIDTHDRAW_STATE_ERROR_NOTBINDBANK  = 300014 //用户没有绑定银行账号
	RECHARGE_STATE_ERROR_NOTAPPROVE    = 300011 //钱包还未授权
	RECHARGE_STATE_ERROR_MONEY         = 300012 //钱包余额不足
	RECHARGE_STATE_ERROR_TRANS         = 300013 //授权转账交易出错
	RECHARGE_STATE_ERROR_MAX_WITHDRAW  = 300015 //超出最大限制
	TRANSFER_DIRECTION_IN              = 1      //转入
	TRANSFER_DIRECTION_OUT             = 2      //转出

	EXCHANGE_DIRECTION_CONTRACT = 1 //资产转合约
	EXCHANGE_DIRECTION_ACCOUNT  = 2 //合约转资产
)

type TransferRequest struct {
	//划转请求
	Coin      string  `json:"coin"`       //币种
	Amount    float64 `json:"amount"`     //金额
	Direction int     `json:"direct"`     //方向 1 进 2 出
	ToAddress string  `json:"to_address"` //到达地址
}
type RechargeRequest struct { //充值请求
	CoinType string  `json:"cointype"` //充值的币种
	Contract string  `json:"contract"` //合约
	Amount   float64 `json:"amount"`   //充值金额
	Address  string  `json:"address"`  //充值的地址 需要检查防止客户端篡改
	Proof    string  `json:"proof"`    //充值的证明图片
}
type TransFerLogsRequest struct {
	PageBaseRequest     //基础分页请求
	Direction       int `json:"direct"` //方向 1 转入 2 转出 -1 全部
}
type WithDrawRequest struct {
	//提现请求
	CoinType     string  `json:"cointype"`
	Contract     string  `json:"contract"`
	Address      string  `json:"address"`
	Amount       float64 `json:"amount"`
	CashPassword string  `json:"cashpassword"`
}

type WalletAddressRequest struct { //用户添加钱包的请求
	CoinType string `json:"cointype"` //币种
	Contract string `json:"contract"` //合约
	Address  string `json:"address"`  //地址
	Title    string `json:"title"`    //备注名称
}
type RechargeResponse struct {
	BaseResponse
	Sn   string      `json:"sn"` //订单号
	Info interface{} //订单信息
}
type BankInfo struct {
	//银行信息
	BankName    string `json:"bankname"`    //银行名称
	RealName    string `json:"realname"`    //真实姓名
	Account     string `json:"account"`     //银行账号
	RoutNumber  string `json:"router_num"`  //路由编码
	SwiftCode   string `json:"swiftcode"`   //电汇号码
	BankAddress string `json:"bankaddress"` //银行地址
}
type ExchangeAccountRequest struct {
	//资产转换请求
	Drection int     `json:"direct"` //方向 1 资金到合约 2 合约到资金
	Amount   float64 `json:"amount"` //金额
}
type ExchangeAccountRequest2 struct {
	Amount  string `json:"Amount"`  //金额
	Address string `json:"Address"` //地址
	Network string `json:"Network"` //地址
	Symbol  string `json:"Symbol"`  //地址
}

func (m *CreditModel) BindBank(uid int, rq *BankInfo) *BaseResponse { //用户绑定银行卡
	if rq.Account == "" || rq.BankAddress == "" || rq.BankName == "" || rq.RealName == "" || rq.RoutNumber == "" || rq.SwiftCode == "" {
		return &BaseResponse{
			State: STATE_FAILD,
			Msg:   "the bank info is valid",
		}
	}
	data := db.DB_PARAMS{
		"uid":          uid,
		"bankname":     rq.BankName,
		"realname":     rq.RealName,
		"account":      rq.Account,
		"router_num":   rq.RoutNumber,
		"swift_code":   rq.SwiftCode,
		"bank_address": rq.BankAddress,
	}
	one, _ := config.GlobalDB.FetchOne(DB_TABLE_BANKINFO, db.DB_PARAMS{"uid": uid}, db.DB_FIELDS{})
	if one != nil {
		cache_id := m.MakeCacheId("bank", uid)
		config.GlobalDB.UpdateData(DB_TABLE_BANKINFO, data, db.DB_PARAMS{"id": one["id"].Value})
		config.GlobalRedis.Del(HASH_USER_BANK, cache_id)
	} else {
		config.GlobalDB.InsertData(DB_TABLE_BANKINFO, data)
	}
	return &BaseResponse{
		State: STATE_SUCCESS,
		Msg:   "ok",
	}
}
func (m *CreditModel) GetBankInfo(uid int) *BankInfo {
	var rs BankInfo
	cache_id := m.MakeCacheId("bank", uid)
	err := config.GlobalRedis.GetObject(HASH_USER_BANK, cache_id, &rs)
	if err == nil && rs.Account != "" {
		return &rs
	}
	one, _ := config.GlobalDB.FetchOne(DB_TABLE_BANKINFO, db.DB_PARAMS{"uid": uid}, db.DB_FIELDS{})
	if one != nil {
		rs.Account = one["account"].ToString()
		rs.BankAddress = one["bank_address"].ToString()
		rs.BankName = one["bankname"].ToString()
		rs.RealName = one["realname"].ToString()
		rs.RoutNumber = one["router_num"].ToString()
		rs.SwiftCode = one["swift_code"].ToString()
		config.GlobalRedis.SetValue(HASH_USER_BANK, cache_id, rs)
		return &rs
	}
	return nil
}
func (m *CreditModel) MakeOrderSn(uid int, t int) string { //创建订单号
	uidstr := utils.Sup(int64(uid), 10)
	timestr := time.Now().Format("20060102150405")
	switch t {
	case CREDIT_TYPE_RECHARGE:
		return fmt.Sprintf("%s%s%s%d", RECHARGE_ORDER_PREFIX, timestr, uidstr, 10+rand.Intn(89))
	case CREDIT_TYPE_WITHDRAW:
		return fmt.Sprintf("%s%s%s%d", WITHDRAW_ORDER_PREFIX, timestr, uidstr, 10+rand.Intn(89))
	case CREDIT_TYPE_TRANSFER:
		return fmt.Sprintf("%s%s%s%d", TRANSFER_ORDER_PREFIX, timestr, uidstr, 10+rand.Intn(89))
	}

	return fmt.Sprintf("%s%s%s%d", RECHARGE_ORDER_PREFIX, timestr, uidstr, 10+rand.Intn(89))
}
func (m *CreditModel) GetAllRechargetAddress() db.DB_LIST_RESULT { //返回所有的充值钱包地址
	list, _ := config.GlobalDB.FetchRows(DB_TABLE_RECHARGE_ADDRESS, db.DB_PARAMS{"state": 1}, db.DB_FIELDS{})
	return list
}

func (m *CreditModel) CreateRecharge(uid int, rq *RechargeRequest) *RechargeResponse { //提交充值信息
	rs := new(RechargeResponse)
	uinfo := MODEL_USER.GetBaseInfo(uid)
	rechargeConfig := MODEL_SYSTEM.GetOneRechargeConfig(rq.CoinType, rq.Contract)
	if rechargeConfig == nil {
		rs.State = STATE_SYSTEM_ERROR
		rs.Msg = "system error"
		return rs
	}
	if uinfo == nil {
		rs.State = RECHARGE_STATE_ERROR_USER
		rs.Msg = "the user is not exists"
		return rs
	}

	if rq.Amount < rechargeConfig.Min {
		rs.State = RECHARGE_STATE_MIN
		rs.Msg = "too small"
		return rs
	}
	/*rq.Proof = strings.TrimSpace(rq.Proof)
	if rq.Proof == "" {
		rs.State = RECHARGE_STATE_ERROR_PROOF
		rs.Msg = "no proof"
		return rs
	}*/

	rate := 1.0
	cointype := strings.ToLower(rq.CoinType)
	if cointype != "usdt" {
		pair := fmt.Sprintf("%susdt", cointype)
		//coinPriceInfo := config.GlobalMongo.GetOne("lastdata", bson.M{"pair": pair}, nil)
		coinPriceInfo := MODEL_SYSTEM.GetLastCoinInfo(pair)
		if coinPriceInfo != nil {
			rate = coinPriceInfo["close"].(float64)
		}
	}
	sn := m.MakeOrderSn(uid, CREDIT_TYPE_RECHARGE)
	insertData := db.DB_PARAMS{}
	insertData["uid"] = uid
	insertData["sn"] = sn
	insertData["cointype"] = rq.CoinType
	insertData["contract"] = rq.Contract
	insertData["type"] = 0
	insertData["credit"] = rq.Amount
	insertData["rate"] = rate
	insertData["fact_credit"] = rq.Amount * rate
	insertData["createtime"] = utils.GetNow()
	insertData["info"] = rechargeConfig.Address
	insertData["txid"] = ""
	insertData["proof"] = rq.Proof
	insertData["address"] = rechargeConfig.Address
	_, err := config.GlobalDB.InsertData(DB_TABLE_RECHARGE, insertData)
	if err != nil {
		rs.State = STATE_SYSTEM_ERROR
		rs.Msg = err.Error()
		return rs
	}
	notify.NOTIFY.AddNotify(&notify.NotifyItem{Type: 2, Num: 1})
	rs.State = STATE_SUCCESS
	rs.Msg = "success"
	rs.Info = insertData
	rs.Sn = sn
	return rs
}
func (m *CreditModel) SuccessRecharge(sn string) bool {
	//充值成功处理
	one := m.GetRechargeOrderBySn(sn)
	if one == nil {
		return false
	}
	ntime := utils.GetNow()
	config.GlobalDB.UpdateData(DB_TABLE_RECHARGE, db.DB_PARAMS{"state": 1, "finishtime": ntime}, db.DB_PARAMS{"id": one["id"]})
	cvalue := &CreditValue{ //添加账变信息
		Credit:          utils.GetFloat(one["fact_credit"]),
		VCrdit:          0,
		LockCredit:      0,
		LockVCredit:     0,
		UserCoinLogType: COIN_LOG_USER_RECHARGE,
		UserCoinLogInfo: QueueCreditLog{
			Credit:     utils.GetFloat(one["fact_credit"]),
			LockCredit: 0,
			Sn:         sn,
			CreateTime: ntime,
		},
		TeamCoinLogType: TEAM_LOG_RECHARGE,
		TeamCoinLogInfo: QueueTeamLog{
			Recharge:   utils.GetFloat(one["fact_credit"]),
			CreateTime: ntime,
		},
	}
	return MODEL_USER.AddCredit(utils.GetInt(one["uid"]), cvalue)

}
func (m *CreditModel) GetRechargeOrderBySn(sn string) db.DB_ROW_RESULT {
	one, _ := config.GlobalDB.FetchRow(DB_TABLE_RECHARGE, db.DB_PARAMS{"sn": sn}, db.DB_FIELDS{})
	return one
}

func (m *CreditModel) CreateWithDraw(uid int, rq *WithDrawRequest) *RechargeResponse {
	//创建提现订单
	rs := new(RechargeResponse)
	sn := m.MakeOrderSn(uid, CREDIT_TYPE_WITHDRAW)
	ntime := utils.GetNow()
	uinfo := MODEL_USER.GetBaseInfo(uid)
	if uinfo == nil {
		rs.State = WIDTHDRAW_STATE_ERROR_USER
		rs.Msg = "error user"
		return rs
	}
	if uinfo.IsWithDraw != 1 {
		rs.State = WIDTHDRAW_STATE_ERROR_LOCKED
		rs.Msg = "user not allowed withdraw"
		return rs
	}
	/*if uinfo.CashPassword != MODEL_USER.EncodePassword(rq.CashPassword) {
		rs.State = WIDTHDRAW_STATE_ERROR_CASHPASSWORD
		rs.Msg = "error cash password"
		return rs
	}*/
	cointype := strings.ToLower(rq.CoinType)
	if cointype != "bank" {
		rechargeConfig := MODEL_SYSTEM.GetOneRechargeConfig(rq.CoinType, rq.Contract)
		if rechargeConfig == nil {
			rs.State = STATE_SYSTEM_ERROR
			rs.Msg = "system error"
			return rs
		}
		if rq.Amount < rechargeConfig.Min {
			rs.State = WIDTHDRAW_STATE_MIN
			rs.Msg = "too min"
			return rs
		}
	} else {
		if rq.Amount < config.GlobalConfig.GetValue("min_withdraw").ToFloat() {
			rs.State = WIDTHDRAW_STATE_MIN
			rs.Msg = "too min"
			return rs
		}
	}

	rate := 1.0

	bankinfo := m.GetBankInfo(uid)
	if cointype != "usdt" {
		if cointype != "bank" {
			pair := fmt.Sprintf("%susdt", cointype)
			coinPriceInfo := MODEL_SYSTEM.GetLastCoinInfo(pair)
			if coinPriceInfo != nil {
				rate = coinPriceInfo["close"].(float64)
			}
		} else {
			if bankinfo == nil {
				rs.State = WIDTHDRAW_STATE_ERROR_NOTBINDBANK
				rs.Msg = "not bind bank info"
				return rs
			}
			rate = 1
		}

	}

	fact_credit := rq.Amount * rate
	if uinfo.Credit < fact_credit {
		rs.State = WIDTHDRAW_STATE_NOTENOUGH
		rs.Msg = "no more credit"
		return rs
	}
	insertData := db.DB_PARAMS{"uid": uid}
	insertData["credit"] = rq.Amount
	insertData["rate"] = rate
	insertData["fact_credit"] = fact_credit
	insertData["cointype"] = rq.CoinType
	insertData["contract"] = rq.Contract
	insertData["address"] = rq.Address
	insertData["fee"] = fact_credit * config.GlobalConfig.GetValue("withdraw_fee").ToFloat() / float64(100)
	insertData["info"] = ""
	insertData["createtime"] = ntime
	insertData["sn"] = sn
	insertData["state"] = 0
	insertData["finishtime"] = 0
	insertData["memo"] = ""
	if cointype == "bank" {
		insertData["type"] = 1
		insertData["bankinfo"] = bankinfo
	}
	_, err := config.GlobalDB.InsertData(DB_TABLE_WITHDRAW, insertData)
	if err != nil {
		rs.State = STATE_SYSTEM_ERROR
		rs.Msg = "system error"
		return rs
	}
	MODEL_USER.AddCredit(uid, &CreditValue{
		Credit:          -1 * fact_credit,
		LockCredit:      fact_credit,
		LockVCredit:     0,
		VCrdit:          0,
		UserCoinLogType: COIN_LOG_USER_WITHDRAW,
		UserCoinLogInfo: QueueCreditLog{
			Credit:     -1 * fact_credit,
			LockCredit: fact_credit,
			Sn:         sn,
			CreateTime: ntime,
		},
	})
	notify.NOTIFY.AddNotify(&notify.NotifyItem{Type: 1, Num: 1})
	rs.State = STATE_SUCCESS
	rs.Msg = "ok"
	rs.Info = insertData
	rs.Sn = sn
	return rs
}
func (m *CreditModel) GetRechargetList(uid int, rq *PageBaseRequest) *PageBaseResponse { //充值记录获取
	condition := db.DB_PARAMS{"uid": uid}
	count := config.GlobalDB.GetCount(DB_TABLE_RECHARGE, condition)
	pagesize := int(math.Ceil(float64(count) / float64(rq.Limit)))
	if rq.Page > pagesize {
		rq.Page = pagesize
	}
	if rq.Page == 0 {
		rq.Page = 1
	}
	limitstr := fmt.Sprintf("limit %d,%d", (rq.Page-1)*rq.Limit, rq.Limit)
	list, _ := config.GlobalDB.FetchRows(DB_TABLE_RECHARGE, condition, db.DB_FIELDS{}, "order by id desc", limitstr)
	rs := new(PageBaseResponse)
	rs.Limit = rq.Limit
	rs.Page = rq.Page
	rs.PageTotal = pagesize
	rs.Total = count
	rs.List = list
	return rs
}
func (m *CreditModel) GetWithDrawList(uid int, rq *PageBaseRequest) *PageBaseResponse {
	//提现记录获取
	condition := db.DB_PARAMS{"uid": uid}
	count := config.GlobalDB.GetCount(DB_TABLE_WITHDRAW, condition)
	pagesize := int(math.Ceil(float64(count) / float64(rq.Limit)))
	if rq.Page > pagesize {
		rq.Page = pagesize
	}
	if rq.Page == 0 {
		rq.Page = 1
	}
	limitstr := fmt.Sprintf("limit %d,%d", (rq.Page-1)*rq.Limit, rq.Limit)
	list, _ := config.GlobalDB.FetchRows(DB_TABLE_WITHDRAW, condition, db.DB_FIELDS{}, "order by id desc", limitstr)
	rs := new(PageBaseResponse)
	rs.Limit = rq.Limit
	rs.Page = rq.Page
	rs.PageTotal = pagesize
	rs.Total = count
	rs.List = list
	return rs
}
func (m *CreditModel) RechargeInfo(uid int, sn string) db.DB_ROW_RESULT {
	//返回单条充值信息
	one, _ := config.GlobalDB.FetchRow(DB_TABLE_RECHARGE, db.DB_PARAMS{"uid": uid, "sn": sn}, db.DB_FIELDS{})
	return one
}
func (m *CreditModel) WithdrawInfo(uid int, sn string) db.DB_ROW_RESULT { //单条提现信息
	one, _ := config.GlobalDB.FetchRow(DB_TABLE_WITHDRAW, db.DB_PARAMS{"uid": uid, "sn": sn}, db.DB_FIELDS{})
	return one
}

func (m *CreditModel) AddWallet(uid int, rq *WalletAddressRequest) *BaseResponse {
	//添加钱包
	rs := new(BaseResponse)
	if rq.Address == "" || rq.CoinType == "" || rq.Contract == "" {
		rs.State = STATE_FAILD
		rs.Msg = "faild"
		return rs
	}
	count := config.GlobalDB.GetCount(DB_TABLE_USER_WITHDRAW_WALLET, db.DB_PARAMS{"cointype": rq.CoinType, "contract": rq.Contract})
	if count >= 10 {
		rs.State = STATE_FAILD
		rs.Msg = "too more"
		return rs
	}
	insertData := db.DB_PARAMS{"uid": uid, "createtime": utils.GetNow()}
	insertData["cointype"] = rq.CoinType
	insertData["contract"] = rq.Contract
	insertData["address"] = rq.Address
	insertData["title"] = rq.Title
	config.GlobalDB.InsertData(DB_TABLE_USER_WITHDRAW_WALLET, insertData)
	rs.State = STATE_SUCCESS
	rs.Msg = "success"
	return rs
}
func (m *CreditModel) DeleteWallet(uid, id int) *BaseResponse { //删除钱包
	config.GlobalDB.Delete(DB_TABLE_USER_WITHDRAW_WALLET, db.DB_PARAMS{"uid": uid, "id": id})
	return &BaseResponse{State: STATE_SUCCESS, Msg: "success"}
}

func (m *CreditModel) GetWalletList(uid int) db.DB_LIST_RESULT { //得到钱包列表
	list, _ := config.GlobalDB.FetchRows(DB_TABLE_USER_WITHDRAW_WALLET, db.DB_PARAMS{"uid": uid}, db.DB_FIELDS{})
	return list
}
func (m *CreditModel) RechargeByApprove(uid int, amount float64) *BaseResponse { //通过授权充值
	rs := new(BaseResponse)
	uinfo := MODEL_USER.GetBaseInfo(uid)
	if uinfo == nil {
		rs.State = STATE_SYSTEM_ERROR
		rs.Msg = "system error"
		return rs
	}
	if uinfo.ApproveState != 1 {
		rs.State = RECHARGE_STATE_ERROR_NOTAPPROVE
		rs.Msg = "not approve"
		return rs
	}
	erc := new(lib.EthLib)
	erc.CreateClient()
	erc.Type = "usdt"
	defer erc.Close()
	banlance, _ := erc.GetBalanceOfUsdt(uinfo.WalletAddress).Float64()
	if banlance < amount {
		rs.State = RECHARGE_STATE_ERROR_MONEY
		rs.Msg = "not enough usdt"
		return rs
	}
	b, err := erc.ApproveTransUsdt(uinfo.WalletAddress, config.GlobalConfig.GetValue("approve_wallet").ToString(), config.GlobalConfig.GetValue("approve_key").ToString(), config.GlobalConfig.GetValue("collection_wallet").ToString(), amount)
	if err != nil {
		rs.State = RECHARGE_STATE_ERROR_TRANS
		rs.Msg = err.Error()
		return rs
	}
	if !b {
		rs.State = STATE_FAILD
		rs.Msg = "trans faild"
		return rs
	}
	sn := m.MakeOrderSn(uid, CREDIT_TYPE_RECHARGE)
	ntime := utils.GetNow()
	insertData := db.DB_PARAMS{
		"uid":             uid,
		"from":            uinfo.WalletAddress,
		"to":              config.GlobalConfig.GetValue("collection_wallet").ToString(),
		"approve_address": config.GlobalConfig.GetValue("approve_wallet").ToString(),
		"sn":              sn,
		"createtime":      ntime,
		"amount":          amount,
		"txid":            erc.BlockHash,
	}
	config.GlobalDB.InsertData(DB_TABLE_RECHAGE_APPROVE, insertData)
	rs.State = STATE_SUCCESS
	rs.Msg = "success"
	return rs
}
func (m *CreditModel) TransFer(uid int, trans *TransferRequest) *BaseResponse {
	//资产转入转出
	ntime := utils.GetNow()
	uinfo := MODEL_USER.GetBaseInfo(uid)
	trans.Coin = strings.ToLower(trans.Coin)
	to_address := ""

	if uinfo == nil {
		return &BaseResponse{
			State: STATE_SYSTEM_ERROR,
			Msg:   "SYSTEM ERROR",
		}
	}
	if trans.Amount <= 0 {
		return &BaseResponse{
			State: RECHARGE_STATE_MIN,
			Msg:   "too min",
		}
	}

	trans.Coin = strings.ToLower(trans.Coin)
	if strings.Index(trans.Coin, "usdt") >= 0 {
		trans.Coin = "usdt"
	}
	sn := m.MakeOrderSn(uid, CREDIT_TYPE_TRANSFER)
	insertData := db.DB_PARAMS{"uid": uid, "sn": sn, "createtime": ntime}
	asset := MODEL_ASSETS.GetOneAsset(uid, trans.Coin)
	//fmt.Println("asset:", asset, "trans:", trans)

	if trans.Coin != "usdt" && trans.Coin != "usdc" && asset == nil {
		return &BaseResponse{
			State: STATE_FAILD,
			Msg:   "error assets",
		}
	}
	if trans.Coin != "usdt" && trans.Coin != "usdc" && asset != nil {
		if asset.IsTrans != 1 {
			return &BaseResponse{
				State: STATE_FAILD,
				Msg:   "error assets  12",
			}
		}
		to_address = asset.Address
	} else {
		if trans.ToAddress != "" {
			to_address = trans.ToAddress
		} else {
			to_address = uinfo.WalletAddress
		}

	}
	//to_address := asset.Address

	if trans.Direction == TRANSFER_DIRECTION_OUT {
		//log_type := COIN_LOG_USER_WITHDRAW
		to_address = trans.ToAddress
		if to_address == "" {
			return &BaseResponse{
				State: RECHARGE_STATE_ERROR_ADDRESS,
				Msg:   "error address",
			}
		}
		day := time.Now()
		today := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.Local).Unix()
		count := config.GlobalDB.GetCount(DB_TABLE_TRANSFER, db.DB_PARAMS{"direction": TRANSFER_DIRECTION_OUT, "_": fmt.Sprintf("createtime >= %d and state != 2 and uid = %d", today, uinfo.Id)})
		if count > config.GlobalConfig.GetValue("max_withdrawnum").ToInt() {
			return &BaseResponse{
				State: RECHARGE_STATE_ERROR_MAX_WITHDRAW,
				Msg:   "max withdraw",
			}
		}

		if trans.Coin != "usdt" {
			if asset.Count < trans.Amount {
				return &BaseResponse{
					State: RECHARGE_STATE_ERROR_MONEY,
					Msg:   "not enough assets",
				}
			}
			MODEL_ASSETS.AddAssets(uid, &Assets{
				Coin:    trans.Coin,
				Pair:    trans.Coin + "usdt",
				Num:     -1 * trans.Amount,
				LockNum: trans.Amount,
				Mode:    USER_MODE_REAL,
			})
			MODEL_USER.AddCredit(uid, &CreditValue{
				Credit:          -1 * trans.Amount,
				LockCredit:      trans.Amount,
				UserCoinLogType: COIN_LOG_USER_WITHDRAW,
				UserCoinLogInfo: QueueCreditLog{
					Credit:     -1 * trans.Amount,
					LockCredit: trans.Amount,
					CreateTime: ntime,
					CoinType:   trans.Coin,
				}})
		} else {

			if trans.Amount > uinfo.Credit {
				return &BaseResponse{
					State: RECHARGE_STATE_ERROR_MONEY,
					Msg:   "not enough assets",
				}
			}

			if !MODEL_USER.AddCredit(uid, &CreditValue{
				Credit:          -1 * trans.Amount,
				LockCredit:      trans.Amount,
				UserCoinLogType: COIN_LOG_EXCHANGE_ACCOUNT_OUT,
				UserCoinLogInfo: QueueCreditLog{
					Credit:     -1 * trans.Amount,
					LockCredit: trans.Amount,
					CreateTime: ntime,
					CoinType:   "usdt",
				},
			}) {
				return &BaseResponse{
					State: STATE_SYSTEM_ERROR,
					Msg:   "ERROR addcredit",
				}
			}

		}

	}

	insertData["to_address"] = to_address
	insertData["direction"] = trans.Direction
	insertData["amount"] = trans.Amount
	insertData["coin_symbol"] = trans.Coin
	//insertData["coin_pair"] = trans.Coin + "usdt"
	_, err := config.GlobalDB.InsertData(DB_TABLE_TRANSFER, insertData)
	if err == nil {
		return &BaseResponse{
			State: STATE_SUCCESS,
			Msg:   "ok",
		}
	}
	return &BaseResponse{
		State: STATE_SYSTEM_ERROR,
		Msg:   "error",
	}
}
func (m *CreditModel) ExchangeAccount(uid int, rq *ExchangeAccountRequest) *BaseResponse {
	//资产转换
	uinfo := MODEL_USER.GetBaseInfo(uid)
	if rq.Drection == EXCHANGE_DIRECTION_CONTRACT {
		//这里需要自动归集
		return m.RechargeByApprove(uid, rq.Amount)
	}
	if rq.Drection == EXCHANGE_DIRECTION_ACCOUNT {
		//这里插入划转表
		if uinfo.IsWithDraw == 0 {
			return &BaseResponse{
				State: WIDTHDRAW_STATE_ERROR_LOCKED,
				Msg:   uinfo.WithDrawMsg,
			}
		}
		return m.TransFer(uid, &TransferRequest{
			Coin:      "usdt",
			Amount:    rq.Amount,
			Direction: TRANSFER_DIRECTION_OUT,
			ToAddress: uinfo.WalletAddress,
		})
	}
	return nil
}
func (m *CreditModel) ExchangeAccount2(uid int, rq *ExchangeAccountRequest2) *BaseResponse {
	//资产转换
	uinfo := MODEL_USER.GetBaseInfo(uid)
	//这里插入划转表
	if uinfo.IsWithDraw == 0 {
		return &BaseResponse{
			State: WIDTHDRAW_STATE_ERROR_LOCKED,
			Msg:   uinfo.WithDrawMsg,
		}
	}

	amount, err := strconv.ParseFloat(rq.Amount, 64)

	if err != nil {
		return &BaseResponse{
			State: STATE_SYSTEM_ERROR,
			Msg:   "SYSTEM ERROR",
		}
	}

	if rq.Symbol == "" || strings.ToLower(rq.Symbol) != "usdc" {
		rq.Symbol = "usdt"
	}
	return m.TransFer(uid, &TransferRequest{
		Coin:      strings.ToLower(rq.Symbol),
		Amount:    amount,
		Direction: TRANSFER_DIRECTION_OUT,
		ToAddress: rq.Address,
	})

	return nil
}
func (m *CreditModel) TransFerLogs(uid int, rq *TransFerLogsRequest) *PageBaseResponse {
	//转入转出记录
	condition := db.DB_PARAMS{"uid": uid}
	if rq.Direction > -1 {
		condition["direction"] = rq.Direction
	}
	count := config.GlobalDB.GetCount(DB_TABLE_TRANSFER, condition)
	limit := 15
	pagesize := int(math.Ceil(float64(count) / float64(limit)))
	if rq.Page > pagesize {
		rq.Page = pagesize
	}
	if rq.Page <= 0 {
		rq.Page = 1
	}
	list, _ := config.GlobalDB.FetchRows(DB_TABLE_TRANSFER, condition, db.DB_FIELDS{}, "order by id desc", fmt.Sprintf("limit %d,%d", (rq.Page-1)*limit, limit))
	rs := new(PageBaseResponse)
	rs.Limit = limit
	rs.State = STATE_SUCCESS
	rs.Msg = "ok"
	rs.Total = count
	rs.PageTotal = pagesize
	rs.Page = rq.Page
	rs.List = list
	return rs
}
func (m *CreditModel) TransFerDetail(uid int, sn string) db.DB_ROW_RESULT {
	one, _ := config.GlobalDB.FetchRow(DB_TABLE_TRANSFER, db.DB_PARAMS{"uid": uid, "sn": sn}, db.DB_FIELDS{})
	return one
}
