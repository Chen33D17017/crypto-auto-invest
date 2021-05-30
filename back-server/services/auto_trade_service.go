package services

import (
	"context"
	"crypto-auto-invest/model"
	"crypto-auto-invest/model/apperrors"
	"log"
)

type autoTradeService struct {
	WalletRepository    model.WalletRepository
	UserRepository      model.UserRepository
	AutoTradeRepository model.AutoTradeRepository
}

type ATSConifg struct {
	WalletRepository    model.WalletRepository
	AutoTradeRepository model.AutoTradeRepository
	UserRepository      model.UserRepository
}

func NewAutoTradeService(c *ATSConifg) model.AutoTradeService {
	return &autoTradeService{
		WalletRepository:    c.WalletRepository,
		UserRepository:      c.UserRepository,
		AutoTradeRepository: c.AutoTradeRepository,
	}
}

func (s *autoTradeService) AddAutoTrade(ctx context.Context, uid, cryptoName string, strategyID int) error {
	cID, err := s.WalletRepository.GetCurrencyID(ctx, cryptoName)
	if err != nil {
		return err
	}
	err = s.checkAndAddWallet(ctx, uid, cryptoName, strategyID)
	if err != nil {
		return err
	}
	err = s.checkAndAddWallet(ctx, uid, "jpy", strategyID)
	if err != nil {
		return err
	}

	err = s.AutoTradeRepository.AddAutoTrade(ctx, uid, cID, strategyID)
	if err != nil {
		return err
	}
	return nil
}

func (s *autoTradeService) DeleteAutoTrade(ctx context.Context, uid, cryptoName string, strategyID int) error {
	cID, err := s.WalletRepository.GetCurrencyID(ctx, cryptoName)
	if err != nil {
		return err
	}

	err = s.AutoTradeRepository.DeleteAutoTrade(ctx, uid, cID, strategyID)
	if err != nil {
		return err
	}

	return nil
}

func (s *autoTradeService) GetAutoTrades(ctx context.Context, uid string) (*[]model.AutoTrade, error) {
	return s.AutoTradeRepository.GetAutoTrades(ctx, uid)
}

func (s *autoTradeService) GetAutoTradesFromStrategy(ctx context.Context, cryptoName string, strategyID int) ([]model.AutoTradeRes, error) {
	rst := []model.AutoTradeRes{}
	settings, err := s.AutoTradeRepository.GetAutoTradeFromStrategy(ctx, cryptoName, strategyID)
	if err != nil {
		log.Printf("GetAutoTradesFromStrategy (strategyID: %v) err: %s\n", strategyID, err.Error())
		return rst, nil
	}
	for _, setting := range *settings {
		tmp := model.AutoTradeRes{}
		tmp.UID = setting.UID
		tmp.CryptoName = setting.CryptoName
		w, err := s.WalletRepository.GetWellet(ctx, setting.UID, setting.CryptoName, strategyID)
		if err != nil {
			log.Printf("GetAutoTradesFromStrategy (uid:%s, cryptoName%s, strategyID:%v)err: %s", setting.UID, setting.CryptoName, strategyID, err.Error())
			continue
		}
		tmp.Amount = w.Amount
		jpyw, err := s.WalletRepository.GetWellet(ctx, setting.UID, "jpy", strategyID)
		if err != nil {
			log.Printf("GetAutoTradesFromStrategy (uid:%s, cryptoName%s, strategyID:%v)err: %s\n", setting.UID, "jpy", strategyID, err.Error())
			continue
		}
		tmp.JPY = jpyw.Amount
		rst = append(rst, tmp)
	}
	return rst, nil
}

func (s *autoTradeService) checkAndAddWallet(ctx context.Context, uid string, cryptoName string, strategyID int) error {
	if _, err := s.WalletRepository.GetWellet(ctx, uid, cryptoName, strategyID); err != nil {
		cid, err := s.WalletRepository.GetCurrencyID(ctx, cryptoName)
		if err != nil {
			return apperrors.NewBadRequest("Wrong crypto name")
		}
		s.WalletRepository.AddWallet(ctx, uid, cid, strategyID)
	}
	return nil
}
