package models

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
