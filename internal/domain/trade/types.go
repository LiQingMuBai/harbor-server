package trade

import "cointrade/internal/domain/shared"

const (
	OPEN_TYPE_KEEP               = 1
	OPEN_TYPE_EXPLODE            = 2
	OPEN_TYPE_BB                 = 3
	PRICE_TYPE_LIMIT             = 1
	PRICE_TYPE_MARKET            = 2
	DIRECT_TYPE_BIG              = 1
	DIRECT_TYPE_SMALL            = 2
	DELEGATE_TYPE_BUY            = 1
	DELEGATE_TYPE_SELL           = 2
	TRADE_BUY_PREFIX             = "B"
	TRADE_SELL_PREFIX            = "S"
	DELEGATE_STATE_NOCOIN        = 40001
	DELEGATE_STATE_CREDIT        = 40002
	DELEGATE_STATE_CLOSETIME     = 40003
	DELEGATE_STATE_NOASSET       = 40004
	DELEGATE_STATE_MIN           = 40005
	DELEGATE_STATE_TRADE_CLOSED  = 40006
	DELEGATE_STATE_GANGGAN_ERROR = 40007
)

type TradeDelegateRequest struct {
	OpenType              int     `json:"opentype"`
	DelegateType          int     `json:"closeoropen"`
	Pair                  string  `json:"pair"`
	Coin                  string  `json:"coin"`
	PriceType             int     `json:"pricetype"`
	DirectType            int     `json:"directtype"`
	GangGan               int     `json:"ganggan"`
	Amount                float64 `json:"amount"`
	Price                 float64 `json:"price"`
	CloseTime             int     `json:"closetime"`
	StopUpPrice           float64 `json:"stop_up_price"`
	StopDownPrice         float64 `json:"stop_down_price"`
	StopUpDelegatePrice   float64 `json:"stop_up_delegate"`
	StopDownDelegatePrice float64 `json:"stop_down_delegate"`
	Sn                    string  `json:"sn"`
}

type OpenedInfo struct {
	Id            int     `json:"id"`
	Uid           int     `json:"uid"`
	Sn            string  `json:"sn"`
	UserType      int     `json:"user_type"`
	TradeType     int     `json:"tradetype"`
	Flag          int     `json:"flag"`
	OpenPrice     float64 `json:"openprice"`
	ClosePrice    float64 `json:"closeprice"`
	CoinId        int     `json:"coinid"`
	CoinPair      string  `json:"pair"`
	CoinSymbol    string  `json:"symbol"`
	CloseTime     int     `json:"closetime"`
	CloseRealTime int     `json:"close_realtime"`
	ClearTime     int     `json:"cleartime"`
	CreateTime    int     `json:"createtime"`
	Ganggan       int     `json:"gangan"`
	Credit        float64 `json:"credit"`
	Profit        float64 `json:"profit"`
	WinRate       float64 `json:"winrate"`
	LoseRate      float64 `json:"loserate"`
	Num           float64 `json:"num"`
	LockNum       float64 `json:"lock_num"`
	UserName      string  `json:"username"`
	Mode          int     `json:"mode"`
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
	UserType   int     `json:"user_type"`
	Profit     float64 `json:"profit"`
	UserName   string  `json:"username"`
}

type DelegateInfo struct {
	Id           int     `json:"id"`
	Uid          int     `json:"uid"`
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
	shared.PageBaseRequest
	TradeType    int    `json:"tradetype"`
	Flag         int    `json:"flag"`
	State        int    `json:"state"`
	Coin         string `json:"coin"`
	DelegateType int    `json:"delegate_type"`
	Ganggan      int    `json:"ganggan"`
}
