package model

type OrderRequest struct {
	Pair   string `json:"pair"`
	Amount string `json:"amount"`
	Side   string `json:"side"`
	Type   string `json:"type"`
}
