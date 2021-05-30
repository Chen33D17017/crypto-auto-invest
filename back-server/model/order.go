package model

type Order struct {
	OID       string  `db:"oid"`
	UID       string  `db:"uid"`
	Piar      string  `db:"pair"`
	Action    string  `db:"action"`
	Amount    float64 `db:"amount"`
	Price     float64 `db:"price"`
	Timestamp string  `db:"timestamp"`
	Fee       float64 `db:"fee"`
	Strategy  int     `db:"strategy"`
}
