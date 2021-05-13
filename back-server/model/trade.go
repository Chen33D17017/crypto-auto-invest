package model

type Trade struct {
	Tid        string `db:"tid" json:"tid"`
	Uid        string `db:"uid" json:"-"`
	FromType   string `db:"from_type" json:"from_type"`
	FromAmount float64 `db:"from_amount" json:"from_amount"`
	ToType     string `db:"to_type" json:"to_type"`
	ToAmount   float64 `db:"to_amount" json:"to_amount"`
	Timestamp  string `db:"time" json:"time"`
}
