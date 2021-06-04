package model

type LimitOrderRequest struct {
	Pair     string `json:"pair"`
	Amount   string `json:"amount"`
	Price    string `json:"price"`
	Side     string `json:"side"`
	Type     string `json:"type"`
	PostOnly bool   `json:"post_only"`
}
