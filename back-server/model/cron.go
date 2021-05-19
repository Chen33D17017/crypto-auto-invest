package model

type Cron struct {
	ID          string  `db:"id" json:"id"`
	UID         string  `db:"uid" json:"-"`
	Type        string  `db:"type" json:"type"`
	Amount      float64 `db:"amount" json:"amount"`
	TimePattern string  `db:"time_pattern" json:"time_pattern"`
}
