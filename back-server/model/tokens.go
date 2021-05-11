package model

type RefreshToken struct {
	ID  string `json:"-"`
	UID string `json:"-"`
	SS  string `json:"refreshToken"`
}

type IDToken struct {
	SS string `json:"idToken"`
}

type TokenPair struct {
	IDToken
	RefreshToken
}
