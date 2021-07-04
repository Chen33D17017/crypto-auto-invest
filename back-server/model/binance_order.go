package model

type BinanceOrder struct {
	ID        string  `db:"id" json:"-"`
	UID       string  `db:"uid" json:"-"`
	Symbol    string  `db:"symbol" json:"symbol"`
	Action    string  `db:"action" json:"action"`
	Amount    float64 `db:"amount" json:"amount"`
	Price     float64 `db:"price" json:"price"`
	Timestamp string  `db:"timestamp" json:"timestamp"`
}
