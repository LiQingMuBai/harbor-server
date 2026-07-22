package models

// 交易模块
type TradeModel struct {
	ModelBase
}

const (
	OPEN_TYPE_KEEP               = 1 //永续
	OPEN_TYPE_EXPLODE            = 2 //交割
	OPEN_TYPE_BB                 = 3 //币币交易
	PRICE_TYPE_LIMIT             = 1 //价格类型 限价
	PRICE_TYPE_MARKET            = 2 //价格类型 市场价格
	DIRECT_TYPE_BIG              = 1 //买涨
	DIRECT_TYPE_SMALL            = 2 //买空
	DELEGATE_TYPE_BUY            = 1 //开仓
	DELEGATE_TYPE_SELL           = 2 //平仓
	TRADE_BUY_PREFIX             = "B"
	TRADE_SELL_PREFIX            = "S"
	DELEGATE_STATE_NOCOIN        = 40001 //持仓不足
	DELEGATE_STATE_CREDIT        = 40002 //余额不足
	DELEGATE_STATE_CLOSETIME     = 40003 //交割合约下平仓时间间隔不合法
	DELEGATE_STATE_NOASSET       = 40004 //没有这个资产 币币交易出售时
	DELEGATE_STATE_MIN           = 40005 //低于交易最小值
	DELEGATE_STATE_TRADE_CLOSED  = 40006 //此类型的交易已关闭
	DELEGATE_STATE_GANGGAN_ERROR = 40007 //杠杆倍率错误
)

type TradeDelegateRequest struct { //委托请求
	OpenType              int     `json:"opentype"`           //开单类型
	DelegateType          int     `json:"closeoropen"`        //委托类型 开仓 还是平仓
	Pair                  string  `json:"pair"`               //交易对
	Coin                  string  `json:"coin"`               //币种
	PriceType             int     `json:"pricetype"`          //价格类型 当前市场价格
	DirectType            int     `json:"directtype"`         //买卖方向
	GangGan               int     `json:"ganggan"`            //杠杆倍率
	Amount                float64 `json:"amount"`             //购买量 永续合约按张数算 1张=1000U 交割合约按 USDT数量算 币币交易按购买的COIN当前汇率或者限价来算
	Price                 float64 `json:"price"`              //限价模式下用户提交的价格
	CloseTime             int     `json:"closetime"`          //平仓时间间隔 交割合约下存在
	StopUpPrice           float64 `json:"stop_up_price"`      //止盈价格 为0为不指定
	StopDownPrice         float64 `json:"stop_down_price"`    //止损价格 为0为不指定
	StopUpDelegatePrice   float64 `json:"stop_up_delegate"`   //止盈委托价格 为0为跟随市场价
	StopDownDelegatePrice float64 `json:"stop_down_delegate"` //止损委托价格 为0为跟随市场价
	Sn                    string  `json:"sn"`                 //订单号 杠杆手动平仓时指定 其他不需要 如果指定了SN其他类型的平仓会返回错误
}

type OpenedInfo struct { //持仓信息结构
	Id            int     `json:"id"`  //系统Id
	Uid           int     `json:"uid"` //用户ID
	Sn            string  `json:"sn"`  //交易单号
	UserType      int     `json:"user_type"`
	TradeType     int     `json:"tradetype"`      //交易类型 交易类型 1 永续合约 2 交割合约 3币币交易
	Flag          int     `json:"flag"`           //方向 多/空
	OpenPrice     float64 `json:"openprice"`      //开仓价格 永续合约为成本价
	ClosePrice    float64 `json:"closeprice"`     //平仓价格
	CoinId        int     `json:"coinid"`         //币系统ID
	CoinPair      string  `json:"pair"`           //币的交易对
	CoinSymbol    string  `json:"symbol"`         //币唯一标识
	CloseTime     int     `json:"closetime"`      //交割合约的下平仓时间间隔
	CloseRealTime int     `json:"close_realtime"` //实际的平仓时间 真实的时间线
	ClearTime     int     `json:"cleartime"`      //结算时间
	CreateTime    int     `json:"createtime"`     //开仓时间
	Ganggan       int     `json:"gangan"`         //杠杆倍率
	Credit        float64 `json:"credit"`         //总投入额
	Profit        float64 `json:"profit"`         //产生的利润
	WinRate       float64 `json:"winrate"`        //交割合约下赢的比例
	LoseRate      float64 `json:"loserate"`       //交割合约下输的比列
	Num           float64 `json:"num"`            //币的总量
	LockNum       float64 `json:"lock_num"`
	UserName      string  `json:"username"`
	Mode          int     `json:"mode"` //交易模式
}

type CloseTrade struct {
	Id         int     `json:"id"`
	Uid        int     `json:"uid"`
	Sn         string  `json:"sn"`
	CoinSymbol string  `json:"coin_symbol"`
	TradeType  string  `json:"trade_type"`
	Flag       int     `json:"flag"`
	Amount     float64 `json:"amount"`
	ClosePrice float64 `json:"close_price"`
	CreateTime int     `json:"createtime"`
	Num        float64 `json:"num"`
	Mode       int     `json:"mode"`
	AllPrice   float64 `json:"allprice"`
	Oprice     float64 `json:"o_price"`

	UserType int     `json:"user_type"`
	Profit   float64 `json:"profit"`
	UserName string  `json:"username"`
}

type DelegateInfo struct {
	Id  int `json:"id"`
	Uid int `json:"uid"`

	UserType     int     `json:"user_type"`
	Sn           string  `json:"sn"`
	DelegameType int     `json:"delegate_type"`
	TradeType    int     `json:"trade_type"`
	Flag         int     `json:"flag"`
	Fee          float64 `json:"fee"`
	Price        float64 `json:"price"`
	CoinId       int     `json:"coinid"`
	CoinPair     string  `json:"coinpair"`
	CoinSymbol   string  `json:"coin_symbol"`
	CloseTime    int     `json:"close_time"`
	Createtime   int     `json:"createtime"`
	Credit       float64 `json:"credit"`
	Num          float64 `json:"num"`
	State        int     `json:"state"`
	Mode         int     `json:"mode"`
	ChangeTime   int     `json:"changetime"`
	UserName     string  `json:"username"`
}
type TradeListRequest struct {
	//委托列表请求
	PageBaseRequest
	TradeType    int    `json:"tradetype"`
	Flag         int    `json:"flag"`
	State        int    `json:"state"`
	Coin         string `json:"coin"`
	DelegateType int    `json:"delegate_type"`
	Ganggan      int    `json:"ganggan"` //是否为杠杆订单 1 为只获取杠杆订单 0 为不指定
}
