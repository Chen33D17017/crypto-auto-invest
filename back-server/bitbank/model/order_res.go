package model

import (
	"encoding/json"
)

type OrderRst OrderRes
type OrderRes struct {
	Status int8  `json:"success"`
	Data   Order `json:"data"`
}
type Order struct {
	OrderId         int64  `json:"order_id"`
	Pair            string `json:"pair"`
	Side            string `json:"side"`
	Type            string `json:"type"`
	StartAmount     string `json:"start_amount"`
	RemainingAmount string `json:"remaining_amount"`
	ExecutedAmount  string `json:"executed_amount"`
	Price           string `json:"Price"`
	AveragePrice    string `json:"average_price"`
	OrderedAt       int64  `json:"ordered_at"`
	Status          string `json:"status"`
}

func (w *OrderRst) UnmarshalJSON(b []byte) error {
	var tmp GeneralRes
	err := json.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}
	if tmp.Status == 0 {
		msg := tmp.Data.(map[string]interface{})["code"].(float64)
		return NewError((int64)(msg))
	}

	json.Unmarshal(b, (*OrderRes)(w))

	return nil
}
