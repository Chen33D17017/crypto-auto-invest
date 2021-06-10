package services

import (
	"bytes"
	"context"
	bm "crypto-auto-invest/bitbank/model"
	"crypto-auto-invest/model"
	"encoding/json"
	"fmt"
	"net/http"
)

type mocktTradeService struct {
	webhook string
}

func NewMockTradeService(webhook string) model.TradeService {
	return &mocktTradeService{
		webhook: webhook,
	}
}

func (s *mocktTradeService) MarketTrade(ctx context.Context, u *model.User, amount float64, action, cryptoName string, strategy int) (bm.Order, error) {
	mock := bm.Order{}
	s.sendTradeRst(fmt.Sprintf("Market: %s %s %s with %v strategy: %v\n", u.Name, action, cryptoName, amount, strategy), "info")
	return mock, nil
}

func (s *mocktTradeService) LimitTrade(ctx context.Context, u *model.User, amount float64, action, cryptoName string, price float64) (bm.Order, error) {
	mock := bm.Order{}
	s.sendTradeRst(fmt.Sprintf("Limit: %s %s %s with %v at %v\n", u.Name, action, cryptoName, amount, price), "info")
	return mock, nil
}

func (s *mocktTradeService) SaveOrder(ctx context.Context, u *model.User, orderID string, cryptoName string, strategy int, info bool) error {
	s.sendTradeRst(fmt.Sprintf("save order %s on %s strategy %v", orderID, cryptoName, strategy), "info")
	return nil
}

func (s *mocktTradeService) CalIncomeRate(ctx context.Context, uid string, cryptoName string, strategyID int) (*model.Income, error) {
	mock := &model.Income{}
	s.sendTradeRst(fmt.Sprintf("Call income rate %s with crypto %s, strategy %v", uid, cryptoName, strategyID), "info")
	return mock, nil
}

func (s *mocktTradeService) sendTradeRst(msg string, level string) error {
	msgJSON, _ := json.Marshal(DiscordFormat{msg})
	payload := bytes.NewReader(msgJSON)

	client := &http.Client{}
	req, err := http.NewRequest("POST", s.webhook, payload)

	if err != nil {
		return fmt.Errorf("Fail to send msg to Discord")
	}
	req.Header.Add("Content-Type", "application/json")
	client.Do(req)

	return nil
}
