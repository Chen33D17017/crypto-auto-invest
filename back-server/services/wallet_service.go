package services

import (
	"context"
	"crypto-auto-invest/model"
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

func (w *walletService) AddWallet(ctx context.Context, uid string, currencyName string) (*model.Wallet, error) {
	var rst *model.Wallet
	err := w.WalletRepository.AddWallet(ctx, uid, currencyName)
	if err != nil {
		return nil, err
	}

	rst, err = w.WalletRepository.GetWellet(ctx, uid, currencyName)
	if err != nil {
		return nil, err
	}
	return rst, nil
}

func (w *walletService) GetUserWallet(ctx context.Context, uid string, currencyName string) (*model.Wallet, error) {
	var rst *model.Wallet
	rst, err := w.WalletRepository.GetWellet(ctx, uid, currencyName)
	if err != nil {
		return nil, err
	}
	return rst, nil
}

func (w *walletService) GetWallets(ctx context.Context, uid string) (*[]model.Wallet, error) {
	rst, err := w.WalletRepository.GetWallets(ctx, uid)
	if err != nil {
		return nil, err
	}
	return rst, nil
}

func (w *walletService) ChangeMoney(ctx context.Context, uid string, currencyName string, amount float64) (*model.Wallet, error) {
	var rst *model.Wallet
	rst, err := w.WalletRepository.GetWellet(ctx, uid, currencyName)
	if err != nil {
		return nil, err
	}

	rst.AMOUNT += amount
	if rst.AMOUNT < 0 {
		return nil, err
	}
	err = w.WalletRepository.UpdateAmount(ctx, rst.WID, rst.AMOUNT)
	if err != nil {
		return nil, err
	}
	return rst, nil
}
