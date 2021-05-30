package model

type AutoTradeRes struct{
	UID string `json:"uid"`
	CryptoName string `json:"crypto_name"`
	Amount float64 `json:"amount"`
	JPY float64 `json:"JPY"`
}