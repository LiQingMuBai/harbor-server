package product

import "cointrade/internal/domain/shared"

const (
	BUY_STATE_CREDIT    = 400001
	BUY_STATE_MIN       = 400002
	BUY_STATE_MAX       = 400003
	BUY_STATE_PER_LIMIT = 400003
	BUY_V_GETED         = 400004
	BUY_STATE_LOCKED    = 40005

	MTYPE_V = 2
	MTYPE_R = 1
)

type ProductProfile struct {
	Algorithm   string `json:"algorithm"`
	MathPower   string `json:"mathpower"`
	GPowW       string `json:"gpoww"`
	Factory     string `json:"factory"`
	WallW       string `json:"wallw"`
	Size        string `json:"size"`
	Weight      string `json:"weight"`
	Temperature string `json:"temperature"`
	Humidity    string `json:"humidity"`
	Chan        string `json:"chan"`
}

type ProductInfo struct {
	Id       int             `json:"id"`
	Name     string          `json:"name"`
	Type     int             `json:"type"`
	Rate     float64         `json:"rate"`
	RateMin  float64         `json:"rate_min"`
	Profit   float64         `json:"profit"`
	Circle   int             `json:"circle"`
	Price    float64         `json:"price"`
	Logo     string          `json:"logo"`
	Desc     string          `json:"desc"`
	Profile  *ProductProfile `json:"profile"`
	PerLimit int             `json:"per_limit"`
	IsOpen   int             `json:"isopen"`
	Min      float64         `json:"min"`
	Max      float64         `json:"max"`
	UserMin  float64         `json:"user_min"`
	IsPublic int             `json:"is_public"`
}

type BuyRequest struct {
	Pid    int     `json:"pid"`
	Amount float64 `json:"amount"`
}

type OrderListRequest struct {
	shared.PageBaseRequest
	State int `json:"state"`
}
