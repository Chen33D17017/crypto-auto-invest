package model

import (
	"encoding/json"
)

type AssetRst AssetRes
type AssetRes struct {
	Status int8 `json:"success"`
	Data   struct {
		Assets []Asset `json:"assets"`
	} `json:"data"`
}

type Asset struct {
	Asset           string `json:"asset"`
	AmountPrecision int    `json:"amount_precision"`
	OnhandAmount    string `json:"onhand_amount"`
	FreeAmount      string `json:"free_amount"`
}

func (w *AssetRst) UnmarshalJSON(b []byte) error {
	var tmp GeneralRes
	err := json.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}
	if tmp.Status == 0 {
		msg := tmp.Data.(map[string]interface{})["code"].(float64)
		return NewError((int64)(msg))
	}

	json.Unmarshal(b, (*AssetRes)(w))
	return nil
}
