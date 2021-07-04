package services

import (
	"bytes"
	"context"
	"crypto-auto-invest/model"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type binanceTradeService struct {
	BinanceTradeRepository model.BinanceTradeRepository
	UserRepository         model.UserRepository
	Webhook                string
}

type BTSConfig struct {
	BinanceTradeRepository model.BinanceTradeRepository
	UserRepository         model.UserRepository
	Webhook                string
}

func NewBinanceTradeService(c *BTSConfig) model.BinanceTradeService {
	return &binanceTradeService{
		BinanceTradeRepository: c.BinanceTradeRepository,
		UserRepository:         c.UserRepository,
		Webhook:                c.Webhook,
	}
}

func (s *binanceTradeService) SaveOrder(ctx context.Context, uid string, symbol string, action string, avgCost float64, qty float64) (model.BinanceOrder, error) {
	cost := 0.0
	amount := 0.0
	order := model.BinanceOrder{
		UID:       uid,
		Symbol:    symbol,
		Action:    action,
		Price:     avgCost,
		Amount:    qty,
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
	}

	err := s.BinanceTradeRepository.SaveOrder(ctx, &order)
	if err != nil {
		return order, err
	}

	orders, err := s.GetOrders(ctx, uid, symbol)
	if err != nil {
		return order, err
	}

	user, err := s.UserRepository.FindByID(ctx, uid)

	for _, order := range *orders {
		if order.Action == "buy" {
			cost += order.Amount * order.Price
			amount += order.Amount
		} else {
			amount -= order.Amount
			cost -= order.Amount * order.Price
		}
		cost = normalizeFloat(cost)
		amount = normalizeFloat(amount)
	}

	incomeRate := (amount*avgCost - cost) / cost * 100

	resultMsg := fmt.Sprintf("%s %s %s amount: %v @%v, income rate now: %v%%", user.Name, action, symbol, qty, avgCost, normalizeFloat(incomeRate))
	err = s.sendResult(resultMsg)
	if err != nil {
		return order, err
	}

	return order, nil
}

func (s *binanceTradeService) GetOrders(ctx context.Context, uid string, symbol string) (*[]model.BinanceOrder, error) {
	return s.BinanceTradeRepository.GetOrders(ctx, uid, symbol)
}

func (s *binanceTradeService) sendResult(msg string) error {

	msgJSON, _ := json.Marshal(DiscordFormat{msg})
	payload := bytes.NewReader(msgJSON)

	client := &http.Client{}
	req, err := http.NewRequest("POST", s.Webhook, payload)

	if err != nil {
		return fmt.Errorf("Fail to send msg to Discord")
	}
	req.Header.Add("Content-Type", "application/json")
	client.Do(req)

	return nil
}
