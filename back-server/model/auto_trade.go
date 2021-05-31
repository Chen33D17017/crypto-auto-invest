package model

type AutoTrade struct {
	ID         string `json:"-" db:"id"`
	UID        string `json:"-" db:"uid"`
	CryptoName string `json:"crypto_name" db:"crypto_name"`
}
