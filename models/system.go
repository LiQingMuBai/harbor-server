package models

const (
	COIN_BUY_STATE_NOTENGOUGH = 100001 //剩余数量不足
	COIN_BUY_STATE_NOMONEY    = 100002 //余额不足
)

type ExplodeConfig struct {
	//交割全局配置
	Time     int     `json:"time"`
	Winrate  float64 `json:"winrate"`
	Loserate float64 `json:"loserate"`
	Minprice float64 `json:"minprice"`
}
type SystemModel struct {
	ModelBase
}
type RechargeContractConfig struct {
	Contract    string  `json:"contract"`     //合约名称
	Min         float64 `json:"min"`          //充值最小金额
	WithDrawMin float64 `json:"withdraw_min"` //提现最小金额
	Address     string  `json:"address"`      //充值地址
}
type RechargeConfig struct {
	CoinType  string                    `json:"cointype"`  //币种
	Logo      string                    `json:"logo"`      //logo
	Contracts []*RechargeContractConfig `json:"contracts"` //合约类型
}

type KlineControlConfig struct {
	StartTime   int     `json:"starttime"`    //开始时间
	EndTime     int     `json:"endtime"`      //结束时间
	TargetPrice float64 `json:"target_price"` //目标价格
}
type CoinKlineConfig struct { //Kline控制参数
	MaxPrice   float64 `json:"max_price"`   //最高价格
	MinPrice   float64 `json:"min_price"`   //最低价格
	Heart      float64 `json:"heart"`       //每一跳的震荡幅度最大幅度 在震荡幅度内随机
	UpRate     int     `json:"up_rate"`     //看涨几率 1-100
	HighRate   float64 `json:"high_rate"`   //高价幅度 1-5 最合适 百分比
	LowRate    float64 `json:"low_rate"`    //低价幅度 1-5 最合适 百分比
	BaseAmount int     `json:"base_amount"` //购买量基础数值 在数值内随机
}

var PERIOD_LIST []string = []string{"1min", "5min", "15min", "30min", "60min", "4hour", "1day", "1mon", "1week", "1year"}
