package model

type Wallet struct {
	WID    string  `db:"wid" json:"-"`
	UID    string  `db:"uid" json:"-"`
	TYPE   string  `db:"type" json:"type"`
	AMOUNT float64 `db:"amount" json:"amount"`
}
