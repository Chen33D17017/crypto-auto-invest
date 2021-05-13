package services

import (
	"context"
	"crypto-auto-invest/model"
	"log"
	"time"
)

type tradeService struct {
	tradeRepository  model.TradeRepository
	walletRepository model.WalletRepository
}

type TSConifg struct {
	TradeRepository  model.TradeRepository
	WalletRepository model.WalletRepository
}

func NewTradeService(c *TSConifg) model.TradeService {
	return &tradeService{
		tradeRepository:  c.TradeRepository,
		walletRepository: c.WalletRepository,
	}
}

func (s *tradeService) Trade(ctx context.Context, uid string, trade_pair string, from_amount float64, getDelay time.Duration) error {
	panic("need to be implemented")
	// Call bitbank api
	// Callback to get transaction detail after 5 minites
}

func (s *tradeService) SaveTradeFromId(ctx context.Context, uid string, tid string) {
	// Call bitbank api to get trade detail
	// Change the value on user wallet
	// Save result with saveTrade method
}

func (s *tradeService) SaveTrade(ctx context.Context, t *model.Trade) {
	err := s.tradeRepository.SaveTrade(ctx, t)
	if err != nil {
		log.Printf("Fail to Store Trade Result with %v err: %s\n", t.Tid, err.Error())
	}
}
