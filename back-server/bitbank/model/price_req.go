package model

import (
	"encoding/json"
)

type PriceRst PriceRes
type PriceRes struct {
	Status int8  `json:"success"`
	Data   Price `json:"data"`
}
type Price struct {
	Sell      string `json:"sell"`
	Buy       string `json:"buy"`
	High      string `json:"high"`
	Low       string `json:"low"`
	Last      string `json:"last"`
	Vol       string `json:"vol"`
	Timestamp int64  `json:"timestamp"`
}

func (w *PriceRst) UnmarshalJSON(b []byte) error {
	var tmp GeneralRes
	err := json.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}
	if tmp.Status == 0 {
		msg := tmp.Data.(map[string]interface{})["code"].(int64)
		return NewError(msg)
	}

	json.Unmarshal(b, (*PriceRes)(w))
	return nil
}
