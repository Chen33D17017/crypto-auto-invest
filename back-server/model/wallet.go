package model

type Wallet struct {
	WID    string  `db:"wid" json:"-"`
	UID    string  `db:"uid" json:"-"`
	Type   string  `db:"type" json:"type"`
	Amount float64 `db:"amount" json:"amount"`
}
