package model

type User struct {
	UID      string `db:"uid" json:"uid"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"-"`
	Name     string `db:"name" json:"name"`
	ImageURL string `db:"image_url" json:"imageUrl"`
}
