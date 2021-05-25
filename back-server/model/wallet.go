package model

type Wallet struct {
	WID        string  `db:"wid" json:"-"`
	UID        string  `db:"uid" json:"-"`
	CryptoName string  `db:"crypto_name" json:"crypto_name"`
	Amount     float64 `db:"amount" json:"amount"`
}
