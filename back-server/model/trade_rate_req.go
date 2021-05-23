package model

type TradeRateReq struct {
	JPY    float64 `json:"jpy"`
	Type   string  `json:"crypto_name"`
	Amount float64 `json:"crypto_amount"`
}
