package services

import (
	"context"
	"crypto-auto-invest/model"
	"crypto-auto-invest/model/apperrors"
	"log"
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
		log.Printf("SERVICE: Fail to add wallet to user with (uid, currencyName): (%s, %s)\n", uid, currencyName)
		return nil, apperrors.NewInternal()
	}

	rst, err = w.WalletRepository.GetWellet(ctx, uid, currencyName)
	if err != nil {
		log.Printf("SERVICE: Fail to get wallet from (uid, currencyName): (%s, %s)\n", uid, currencyName)
		return nil, apperrors.NewInternal()
	}
	return rst, nil
}

func (w *walletService) GetUserWallet(ctx context.Context, uid string, currencyName string) (*model.Wallet, error) {
	var rst *model.Wallet
	rst, err := w.WalletRepository.GetWellet(ctx, uid, currencyName)
	if err != nil {
		log.Printf("SERVICE: Fail to get wallet from (uid, currencyName): (%s, %s)\n", uid, currencyName)
		return nil, apperrors.NewInternal()
	}
	return rst, nil
}

func (w *walletService) GetWallets(ctx context.Context, uid string) (*[]model.Wallet, error) {
	rst, err := w.WalletRepository.GetWallets(ctx, uid)
	if err != nil {
		log.Printf("SERVICE: Fail to get wallets from (uid): (%s)\n", uid)
		return nil, apperrors.NewInternal()
	}
	return rst, nil
}

func (w *walletService) ChangeMoney(ctx context.Context, uid string, currencyName string, amount float64) (*model.Wallet, error) {
	var rst *model.Wallet
	rst, err := w.WalletRepository.GetWellet(ctx, uid, currencyName)
	if err != nil {
		log.Printf("SERVICE: Fail to get wallet from (uid, currencyName): (%s, %s)\n", uid, currencyName)
		return nil, apperrors.NewInternal()
	}

	rst.AMOUNT += amount
	if rst.AMOUNT < 0 {
		log.Printf("SERVICE: User balance not enough: %s, %v\n", uid, rst.AMOUNT)
		return nil, apperrors.NewBadRequest("Balance is not enough")
	}
	err = w.WalletRepository.UpdateAmount(ctx, rst.WID, rst.AMOUNT)
	if err != nil {
		log.Printf("SERVICE: Fail to update wallet from (uid, currencyName): (%s, %s)\n", uid, currencyName)
		return nil, apperrors.NewInternal()
	}
	return rst, nil
}
