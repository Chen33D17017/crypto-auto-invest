package model

type CancelOrderRequest struct {
	Pair    string `json:"pair"`
	OrderID string `json:"order_id"`
}
