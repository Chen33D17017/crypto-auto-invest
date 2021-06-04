package model

type AutoTrade struct {
	ID         string `json:"-" db:"id"`
	UID        string `json:"uid" db:"uid"`
	CryptoName string `json:"crypto_name" db:"crypto_name"`
	StrategyID int    `json:"strategy_id" db:"strategy_id"`
}
