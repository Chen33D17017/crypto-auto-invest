package model

import (
	"encoding/json"
)

type TradeRst TradeRes
type TradeRes struct {
	Status int8 `json:"success"`
	Data   struct {
		Trades []Trade `json:"trades"`
	} `json:"data"`
}
type Trade struct {
	TradeId        int64  `json:"trade_id"`
	Pair           string `json:"pair"`
	OrderId        int64  `json:"order_id"`
	Side           string `json:"side"`
	Type           string `json:"type"`
	Amount         string `json:"amount"`
	Price          string `json:"price"`
	MakerTaker     string `json:"maker_taker"`
	FeeAmountBase  string `json:"fee_amount_base"`
	FeeAmountQuote string `json:"fee_amount_quote"`
	ExecutedAt     int64  `json:"executed_at"`
}

func (w *TradeRst) UnmarshalJSON(b []byte) error {
	var tmp GeneralRes
	err := json.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}
	if tmp.Status == 0 {
		msg := tmp.Data.(map[string]interface{})["code"].(float64)
		return NewError((int64)(msg))
	}

	json.Unmarshal(b, (*TradeRes)(w))
	return nil
}

type Trades []Trade

func (trs Trades) Len() int {
	return len(trs)
}

func (trs Trades) Less(i, j int) bool {
	return trs[i].ExecutedAt < trs[j].ExecutedAt
}

func (trs Trades) Swap(i, j int) {
	trs[i], trs[j] = trs[j], trs[i]
}
