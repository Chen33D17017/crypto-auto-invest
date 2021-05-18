package model

type Order struct {
	OID        string  `db:"oid"`
	UID        string  `db:"uid"`
	FromWallet string  `db:"from_wid"`
	FromAmount float64 `db:"from_amount"`
	ToWallet   string  `db:"to_wid"`
	ToAmount   float64 `db:"to_amount"`
	Timestamp  string  `db:"timestamp"`
	Type       string  `db:"type"`
}
