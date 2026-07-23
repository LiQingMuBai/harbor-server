package assets

const (
	EXCHANGE_STATE_NOCOIN      = 600001
	EXCHANGE_STATE_NOTENNOUGH  = 600002
	EXCHANGE_STATE_TOOMIN      = 600003
	EXCHANGE_STATE_NOT_TRANS   = 600004
	ASSETS_TRANS_TYPE_IN       = 1
	ASSETS_TRANS_TYPE_OUT      = 2
	ASSETS_TRANS_TYPE_CONTRACT = 3
)

type AssetInfo struct {
	Symbol        string  `json:"symbol"`
	CoinId        int     `json:"coinid"`
	Pair          string  `json:"pair"`
	Count         float64 `json:"count"`
	O_Price       float64 `json:"o_price"`
	LockCount     float64 `json:"lockcount"`
	Address       string  `json:"address"`
	IsTrans       int     `json:"istrans"`
	TransOpenTime int     `json:"trans_open_time"`
}

type Assets struct {
	Coin          string
	Pair          string
	Num           float64
	LockNum       float64
	Price         float64
	Mode          int
	IsTrans       int
	OpenTransTime int
}

type ExchangeRequest struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

type AssetsTransRequest struct {
	Coin      string  `json:"coin"`
	Type      int     `json:"type"`
	Amount    float64 `json:"amount"`
	ToAddress string  `json:"to_address"`
}

type QuickExchangeRequest struct {
	Coin   string  `json:"coin"`
	Amount float64 `json:"amount"`
}
