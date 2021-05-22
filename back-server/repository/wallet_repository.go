package repository

import (
	"context"
	"crypto-auto-invest/model"
	"crypto-auto-invest/model/apperrors"
	"log"

	"github.com/jmoiron/sqlx"
)

const (
	queryAddWallet     = "INSERT INTO wallets(uid, type, amount) VALUES(?, ?, ?)"
	queryGetWalletByID = "SELECT * FROM wallets_view WHERE wid=?"
	queryGetWallet     = `SELECT * FROM wallets_view WHERE uid=? AND type=?`
	queryGetWallets    = `SELECT * FROM wallets_view WHERE uid=?`
	queryUpdateAmount  = `UPDATE wallets SET amount=? WHERE wid=?`
	queryGetCurrencyID = `SELECT id FROM currency_type WHERE`
)

type walletRepository struct {
	DB *sqlx.DB
}

func NewWalletRepository(db *sqlx.DB) model.WalletRepository {
	return &walletRepository{
		DB: db,
	}
}
func (r *walletRepository) AddWallet(ctx context.Context, uid string, currencyName string) error {
	stmt, err := r.DB.PrepareContext(ctx, queryAddWallet)
	if err != nil {
		log.Printf("REPOSITORY: Unable to Add Wallet: %v\n", err)
		return apperrors.NewInternal()
	}

	if _, err := stmt.ExecContext(ctx, uid, currencyName, 0); err != nil {
		log.Printf("REPOSITORY: Failed to update details for user: %v err: %s\n", uid, err.Error())
		return apperrors.NewInternal()
	}
	return nil
}

func (r *walletRepository) GetWalletByID(ctx context.Context, wid string) (*model.Wallet, error) {
	var rst *model.Wallet
	err := r.DB.GetContext(ctx, rst, queryGetWalletByID, wid)
	if err != nil {
		log.Printf("REPOSITORY: Unable to get wallet by id: %v err: %s", wid, err)
		return rst, apperrors.NewNotFound("wallet", wid)
	}
	return rst, nil
}

func (r *walletRepository) GetWellet(ctx context.Context, uid string, currencyType string) (*model.Wallet, error) {
	rst := &model.Wallet{}
	err := r.DB.GetContext(ctx, rst, queryGetWallet, uid, currencyType)
	if err != nil {
		log.Printf("REPOSITORY: Unable to get wallet by (uid, currency): (%v, %v) err: %s", uid, currencyType, err)
		return rst, apperrors.NewNotFound(uid, currencyType)
	}
	return rst, nil
}

func (r *walletRepository) GetWallets(ctx context.Context, uid string) (*[]model.Wallet, error) {
	rst := &[]model.Wallet{}
	err := r.DB.SelectContext(ctx, rst, queryGetWallets, uid)
	if err != nil {
		log.Printf("REPOSITORY: Unable to get wallets by (uid): %v err: %s", uid, err)
		return rst, apperrors.NewNotFound("uid", uid)
	}
	return rst, nil
}

func (r *walletRepository) UpdateAmount(ctx context.Context, wid string, amount float64) error {
	stmt, err := r.DB.PrepareContext(ctx, queryUpdateAmount)
	if err != nil {
		log.Printf("REPOSITORY: Unable to prepare update query: %v\n", err)
		return apperrors.NewInternal()
	}

	if _, err := stmt.ExecContext(ctx, amount, wid); err != nil {
		log.Printf("REPOSITORY: Failed to update wallet for wid: %v err: %s\n", wid, err.Error())
		return apperrors.NewInternal()
	}
	return nil
}

