package model

type TradeRateReq struct {
	JPY    float64 `json:"jpy"`
	Type   string  `json:"type"`
	Amount float64 `json:"amount"`
}
