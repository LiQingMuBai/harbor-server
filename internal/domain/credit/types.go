package credit

import "cointrade/internal/domain/shared"

const (
	CREDIT_TYPE_RECHARGE         = 1
	CREDIT_TYPE_WITHDRAW         = 2
	CREDIT_TYPE_TRANSFER         = 3
	RECHARGE_ORDER_PREFIX        = "R"
	WITHDRAW_ORDER_PREFIX        = "W"
	TRANSFER_ORDER_PREFIX        = "T"
	RECHARGE_STATE_MIN           = 300001
	RECHARGE_STATE_ERROR_ADDRESS = 300002
	RECHARGE_STATE_ERROR_USER    = 300003
	RECHARGE_STATE_ERROR_PROOF   = 300004

	WIDTHDRAW_STATE_MIN                = 300005
	WIDTHDRAW_STATE_ERROR_USER         = 300006
	WIDTHDRAW_STATE_ERROR_ADDRESS      = 300007
	WIDTHDRAW_STATE_NOTENOUGH          = 300008
	WIDTHDRAW_STATE_ERROR_CASHPASSWORD = 300009
	WIDTHDRAW_STATE_ERROR_LOCKED       = 300010
	WIDTHDRAW_STATE_ERROR_NOTBINDBANK  = 300014
	RECHARGE_STATE_ERROR_NOTAPPROVE    = 300011
	RECHARGE_STATE_ERROR_MONEY         = 300012
	RECHARGE_STATE_ERROR_TRANS         = 300013
	RECHARGE_STATE_ERROR_MAX_WITHDRAW  = 300015
	TRANSFER_DIRECTION_IN              = 1
	TRANSFER_DIRECTION_OUT             = 2

	EXCHANGE_DIRECTION_CONTRACT = 1
	EXCHANGE_DIRECTION_ACCOUNT  = 2
)

type TransferRequest struct {
	Coin      string  `json:"coin"`
	Amount    float64 `json:"amount"`
	Direction int     `json:"direct"`
	ToAddress string  `json:"to_address"`
}

type RechargeRequest struct {
	CoinType string  `json:"cointype"`
	Contract string  `json:"contract"`
	Amount   float64 `json:"amount"`
	Address  string  `json:"address"`
	Proof    string  `json:"proof"`
}

type TransFerLogsRequest struct {
	shared.PageBaseRequest
	Direction int `json:"direct"`
}

type WithDrawRequest struct {
	CoinType     string  `json:"cointype"`
	Contract     string  `json:"contract"`
	Address      string  `json:"address"`
	Amount       float64 `json:"amount"`
	CashPassword string  `json:"cashpassword"`
}

type WalletAddressRequest struct {
	CoinType string `json:"cointype"`
	Contract string `json:"contract"`
	Address  string `json:"address"`
	Title    string `json:"title"`
}

type RechargeResponse struct {
	shared.BaseResponse
	Sn   string `json:"sn"`
	Info interface{}
}

type BankInfo struct {
	BankName    string `json:"bankname"`
	RealName    string `json:"realname"`
	Account     string `json:"account"`
	RoutNumber  string `json:"router_num"`
	SwiftCode   string `json:"swiftcode"`
	BankAddress string `json:"bankaddress"`
}

type ExchangeAccountRequest struct {
	Drection int     `json:"direct"`
	Amount   float64 `json:"amount"`
}

type ExchangeAccountRequest2 struct {
	Amount  string `json:"Amount"`
	Address string `json:"Address"`
	Network string `json:"Network"`
	Symbol  string `json:"Symbol"`
}
