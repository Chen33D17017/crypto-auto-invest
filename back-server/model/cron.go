package model

type Cron struct {
	ID          string  `db:"id" json:"id"`
	UID         string  `db:"uid" json:"-"`
	CryptoName  string  `db:"crypto_name" json:"crypto_name"`
	Amount      float64 `db:"amount" json:"amount"`
	TimePattern string  `db:"time_pattern" json:"time_pattern"`
}
