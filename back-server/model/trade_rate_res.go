package model

type TradeRateRes struct {
	Side string  `json:"action"`
	Rate float64 `json:"rate"`
}
