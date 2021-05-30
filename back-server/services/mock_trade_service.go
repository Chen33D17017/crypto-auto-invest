package services

import (
	"context"
	bm "crypto-auto-invest/bitbank/model"
	"crypto-auto-invest/model"
	"log"
)

type mocktTradeService struct {
}

func NewMockTradeService() model.TradeService {
	return &mocktTradeService{}
}

func (s *mocktTradeService) Trade(ctx context.Context, u *model.User, amount float64, action, cryptoName string, strategy int) (bm.Order, error) {
	mock := bm.Order{}
	log.Printf("%s %s %s with %v strategt: %v\n", u.Name, action, cryptoName, amount, strategy)
	return mock, nil
}

func (s *mocktTradeService) SaveOrder(ctx context.Context, u *model.User, orderID string, cryptoName string, strategy int) error {
	log.Println(u.Name + "save order")
	return nil
}

func (s *mocktTradeService) CalIncomeRate(ctx context.Context, uid string, cryptoName string, strategyID int) (float64, error) {
	return 0.0, nil
}
func (s *mocktTradeService) SendTradeRst(msg string, level string) error {
	return nil
}
