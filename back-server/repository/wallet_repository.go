package repository

import (
	"context"
	"crypto-auto-invest/model"
	"crypto-auto-invest/model/apperrors"
	"log"

	"github.com/jmoiron/sqlx"
)

const (
	queryAddWallet     = "INSERT INTO wallets(uid, crypto_id, strategy_id, amount) VALUES(?, ?, ?, ?)"
	queryGetWalletByID = "SELECT * FROM wallets_view WHERE wid=?"
	queryGetWallet     = `SELECT * FROM wallets_view WHERE uid=? AND crypto_name=? AND strategy_id=?`
	queryGetWallets    = `SELECT * FROM wallets_view WHERE uid=? AND strategy_id=?`
	queryUpdateAmount  = `UPDATE wallets SET amount=? WHERE wid=?`
	queryGetCurrencyID = `SELECT id FROM crypto_name WHERE name=?`
	queryAddChargeLog  = `INSERT INTO charge_log(uid, crypto_id, strategy_id, amount) VALUES(?, ?, ?, ?)`
	queryGetChargeLogs = `SELECT * FROM charge_log_view WHERE uid=? AND crypto_name=? AND strategy_id=?`
)

type walletRepository struct {
	DB *sqlx.DB
}

func NewWalletRepository(db *sqlx.DB) model.WalletRepository {
	return &walletRepository{
		DB: db,
	}
}
func (r *walletRepository) AddWallet(ctx context.Context, uid string, cid int, strategyID int) error {
	stmt, err := r.DB.PrepareContext(ctx, queryAddWallet)
	if err != nil {
		log.Printf("REPOSITORY: Unable to Add Wallet: %v\n", err)
		return apperrors.NewInternal()
	}

	if _, err := stmt.ExecContext(ctx, uid, cid, strategyID, 0); err != nil {
		log.Printf("REPOSITORY: Failed to add wallet for user: %v err: %s\n", uid, err.Error())
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

func (r *walletRepository) GetWellet(ctx context.Context, uid string, currencyName string, strategyID int) (*model.Wallet, error) {
	rst := &model.Wallet{}
	err := r.DB.GetContext(ctx, rst, queryGetWallet, uid, currencyName, strategyID)
	if err != nil {
		log.Printf("REPOSITORY: Unable to get wallet by (uid, currency): (%v, %v) err: %s", uid, currencyName, err)
		return rst, apperrors.NewNotFound(uid, currencyName)
	}
	return rst, nil
}

func (r *walletRepository) GetWallets(ctx context.Context, uid string, strategyID int) (*[]model.Wallet, error) {
	rst := &[]model.Wallet{}
	err := r.DB.SelectContext(ctx, rst, queryGetWallets, uid, strategyID)
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

func (r *walletRepository) GetCurrencyID(ctx context.Context, currencyName string) (int, error) {
	var rst int
	err := r.DB.GetContext(ctx, &rst, queryGetCurrencyID, currencyName)
	if err != nil {
		log.Printf("REPOSITORY: Unable to get currency name: %s err: %s\n", currencyName, err.Error())
		return 0, apperrors.NewNotFound("currency", currencyName)
	}
	return rst, nil
}

func (r *walletRepository) AddChargeLog(ctx context.Context, uid string, cid int, strategyID int, amount float64) error {
	stmt, err := r.DB.PrepareContext(ctx, queryAddChargeLog)
	if err != nil {
		log.Printf("REPOSITORY: Unable to add charge log: %v\n", err)
		return apperrors.NewInternal()
	}

	if _, err := stmt.ExecContext(ctx, uid, cid, strategyID, amount); err != nil {
		log.Printf("REPOSITORY: Failed add charge log %v err: %s\n", uid, err.Error())
		return apperrors.NewInternal()
	}
	return nil
}

func (r *walletRepository) GetChargeLogs(ctx context.Context, uid string, cryptoName string, strategyID int) (*[]model.ChargeLog, error) {
	rst := &[]model.ChargeLog{}
	err := r.DB.SelectContext(ctx, rst, queryGetChargeLogs, uid, cryptoName, strategyID)
	if err != nil {
		log.Printf("REPOSITORY: Unable to get wallets by (uid): %v err: %s", uid, err)
		return rst, apperrors.NewNotFound("charge_log", cryptoName)
	}
	return rst, nil
}
