package shared

// BaseResponse is the shared base response payload.
type BaseResponse struct {
	Msg   string `json:"msg"`
	State int    `json:"state"`
}

// PageBaseResponse is the shared paged response payload.
type PageBaseResponse struct {
	BaseResponse
	Total     int         `json:"total"`
	Page      int         `json:"page"`
	PageTotal int         `json:"pagetotal"`
	Limit     int         `json:"limit"`
	List      interface{} `json:"list"`
}

// PageBaseRequest is the shared paged request payload.
type PageBaseRequest struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type Currency struct {
	Id      int     `json:"id"`
	Symbol  string  `json:"symbol"`
	Rate    float64 `json:"rate"`
	Country string  `json:"country"`
	Memo    string  `json:"memo"`
}

type Recharge struct {
	Id         int     `json:"id"`
	Uid        int     `json:"uid"`
	Email      string  `json:"email"`
	Sn         string  `json:"sn"`
	CoinType   string  `json:"cointype"`
	Type       int     `json:"type"`
	Credit     float64 `json:"credit"`
	Createtime int     `json:"createtime"`
	FactCredit float64 `json:"fact_credit"`
	Info       string  `json:"info"`
	Txid       string  `json:"txid"`
	State      int     `json:"state"`
	FinishTime int     `json:"finishtime"`
	Proof      string  `json:"proof"`
}

type Withdraw struct {
	Id         int         `json:"id"`
	Uid        int         `json:"uid"`
	UserName   string      `json:"username"`
	ParentName string      `json:"parent_name"`
	Credit     float64     `json:"credit"`
	FactCredit float64     `json:"fact_credit"`
	CoinType   string      `json:"cointype"`
	Contract   string      `json:"contract"`
	Fee        float64     `json:"fee"`
	Type       int         `json:"type"`
	FinishTime int         `json:"finishtime"`
	Info       interface{} `json:"info"`
	CreateTime int         `json:"createtime"`
	Rate       float64     `json:"rate"`
	Sn         string      `json:"sn"`
	State      int         `json:"state"`
	Address    string      `json:"address"`
	Memo       string      `json:"memo"`
	UserType   int         `json:"user_type"`
}
