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

func (s *mocktTradeService) Trade(ctx context.Context, u *model.User, amount float64, side, assetType, orderType string) (bm.Order, error) {
	mock := bm.Order{}
	log.Printf("%s %s %s with %v (%s)\n", u.Name, side, assetType, amount, orderType)
	return mock, nil
}

func (s *mocktTradeService) SaveOrder(ctx context.Context, u *model.User, orderID string, assetType, orderType string) error {
	log.Println(u.Name + "save order")
	return nil
}

func (s *mocktTradeService) CalIncomeRate(ctx context.Context, uid string, cryptoName string, strategyID int) (*model.Income, error) {
	mock := &model.Income{}
	return mock, nil
}
func (s *mocktTradeService) SendTradeRst(msg string, level string) error {
	return nil
}
