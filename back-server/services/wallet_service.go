package services

import (
	"context"
	"crypto-auto-invest/model"
	"crypto-auto-invest/model/apperrors"
)

type walletService struct {
	WalletRepository model.WalletRepository
}

type WAConfig struct {
	WalletRepository model.WalletRepository
}

func NewWalletService(c *WAConfig) model.WalletService {
	return &walletService{
		WalletRepository: c.WalletRepository,
	}
}

func (w *walletService) AddWallet(ctx context.Context, uid string, cryptoName string, strategyID int) (*model.Wallet, error) {
	var rst *model.Wallet
	cid, err := w.WalletRepository.GetCurrencyID(ctx, cryptoName)
	if err != nil {
		return nil, err
	}
	err = w.WalletRepository.AddWallet(ctx, uid, cid, strategyID)
	if err != nil {
		return nil, err
	}

	rst, err = w.WalletRepository.GetWellet(ctx, uid, cryptoName, strategyID)
	if err != nil {
		return nil, err
	}
	return rst, nil
}

func (w *walletService) GetUserWallet(ctx context.Context, uid string, cryptoName string, strategyID int) (*model.Wallet, error) {
	var rst *model.Wallet
	rst, err := w.WalletRepository.GetWellet(ctx, uid, cryptoName, strategyID)
	if err != nil {
		return nil, err
	}
	return rst, nil
}

func (w *walletService) GetWallets(ctx context.Context, uid string, strategyID int) (*[]model.Wallet, error) {
	rst, err := w.WalletRepository.GetWallets(ctx, uid, strategyID)
	if err != nil {
		return nil, err
	}
	return rst, nil
}

func (w *walletService) ChangeMoney(ctx context.Context, uid string, cryptoName string, amount float64, strategyID int) (*model.Wallet, error) {
	var rst *model.Wallet
	cid, err := w.WalletRepository.GetCurrencyID(ctx, cryptoName)
	if err != nil {
		return rst, apperrors.NewBadRequest("Unknow crypto name")
	}
	err = w.WalletRepository.AddChargeLog(ctx, uid, cid, strategyID, amount)
	rst, err = w.WalletRepository.GetWellet(ctx, uid, cryptoName, strategyID)
	if err != nil {
		return nil, err
	}

	rst.Amount += amount
	if rst.Amount < 0 {
		return nil, err
	}
	err = w.WalletRepository.UpdateAmount(ctx, rst.WID, rst.Amount)
	if err != nil {
		return nil, err
	}
	return rst, nil
}

func (w *walletService) GetChargeLogs(ctx context.Context, uid string, cryptoName string, strategyID int) (*[]model.ChargeLog, error) {
	return w.WalletRepository.GetChargeLogs(ctx, uid, cryptoName, strategyID)
}
