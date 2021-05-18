package model

type Secret struct {
	ApiKey    string `json:"key"`
	ApiSecret string `json:"secret"`
}
