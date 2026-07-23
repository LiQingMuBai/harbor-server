package system

const (
	COIN_BUY_STATE_NOTENGOUGH = 100001
	COIN_BUY_STATE_NOMONEY    = 100002
)

type ExplodeConfig struct {
	Time     int     `json:"time"`
	Winrate  float64 `json:"winrate"`
	Loserate float64 `json:"loserate"`
	Minprice float64 `json:"minprice"`
}

type RechargeContractConfig struct {
	Contract    string  `json:"contract"`
	Min         float64 `json:"min"`
	WithDrawMin float64 `json:"withdraw_min"`
	Address     string  `json:"address"`
}

type RechargeConfig struct {
	CoinType  string                    `json:"cointype"`
	Logo      string                    `json:"logo"`
	Contracts []*RechargeContractConfig `json:"contracts"`
}

type KlineControlConfig struct {
	StartTime   int     `json:"starttime"`
	EndTime     int     `json:"endtime"`
	TargetPrice float64 `json:"target_price"`
}

type CoinKlineConfig struct {
	MaxPrice   float64 `json:"max_price"`
	MinPrice   float64 `json:"min_price"`
	Heart      float64 `json:"heart"`
	UpRate     int     `json:"up_rate"`
	HighRate   float64 `json:"high_rate"`
	LowRate    float64 `json:"low_rate"`
	BaseAmount int     `json:"base_amount"`
}

var PERIOD_LIST = []string{"1min", "5min", "15min", "30min", "60min", "4hour", "1day", "1mon", "1week", "1year"}
