package models

// 用户资产相关包
const (
	EXCHANGE_STATE_NOCOIN      = 600001 //没有这个币种
	EXCHANGE_STATE_NOTENNOUGH  = 600002 //没有足够的余额支持兑换
	EXCHANGE_STATE_TOOMIN      = 600003 //交易量过小
	EXCHANGE_STATE_NOT_TRANS   = 600004 //未允许交易
	ASSETS_TRANS_TYPE_IN       = 1      //转入
	ASSETS_TRANS_TYPE_OUT      = 2      //转出
	ASSETS_TRANS_TYPE_CONTRACT = 3      //转入交易账户
)

type AssetModel struct {
	ModelBase
}
type AssetInfo struct {
	Symbol        string  `json:"symbol"`          //币种唯一标识
	CoinId        int     `json:"coinid"`          //币的系统ID
	Pair          string  `json:"pair"`            //交易对
	Count         float64 `json:"count"`           //拥有的数量
	O_Price       float64 `json:"o_price"`         //成本单价
	LockCount     float64 `json:"lockcount"`       //锁定数量
	Address       string  `json:"address"`         //用户私有地址
	IsTrans       int     `json:"istrans"`         //是否允许兑换
	TransOpenTime int     `json:"trans_open_time"` //交易开启时间
}
type Assets struct {
	Coin          string  //币种
	Pair          string  //交易对
	Num           float64 //数量
	LockNum       float64 //锁定数量
	Price         float64 //开仓价格
	Mode          int     //模式
	IsTrans       int     //是否可以交易划转
	OpenTransTime int     //交易划转权限开启时间
}
type ExchangeRequest struct { //兑换请求
	From   string  `json:"from"`   //来源币种
	To     string  `json:"to"`     //兑换币种
	Amount float64 `json:"amount"` //兑换数量
}
type AssetsTransRequest struct { //划转请求
	Coin      string  `json:"coin"`       //币种
	Type      int     `json:"type"`       //划转类型
	Amount    float64 `json:"amount"`     //金额
	ToAddress string  `json:"to_address"` //到达地址
}
type QuickExchangeRequest struct {
	//闪兑请求
	Coin   string  `json:"coin"`   //币种
	Amount float64 `json:"amount"` //金额
}
