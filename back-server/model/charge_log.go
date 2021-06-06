package model

type ChargeLog struct {
	ID         int     `json:"-" db:"id"`
	UID        string  `json:"uid" db:"uid"`
	StrategyID int     `json:"strategy_id" db:"strategy_id"`
	CryptoName string  `json:"crypto_name" db:"crypto_name"`
	Amount     float64 `json:"amount" db:"amount"`
}
