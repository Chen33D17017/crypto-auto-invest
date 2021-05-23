package model

type AutoTrade struct {
	ID   string `json:"-" db:"id"`
	UID  string `json:"-" db:"uid"`
	Type string `json:"type" db:"type"`
}
