package repository

import (
	"context"
	"crypto-auto-invest/model"
	"crypto-auto-invest/model/apperrors"
	"log"

	"github.com/jmoiron/sqlx"
)

const (
	queryAddAutoTrade     = "INSERT INTO auto_trades(uid, crypto_id, strategy_id) VALUES(?, ?, ?);"
	queryDeleteAutoTrade  = "DELETE FROM auto_trades WHERE uid=? and crypto_id=? and strategy_id=?;"
	queryGetAutoTrades    = "SELECT * FROM auto_trades_view WHERE uid=?"
	queryGetAutoTradeUser = "SELECT * FROM auto_trades_view WHERE crypto_name=? and strategy_id=?"
	queryGetAllAutoTrade  = "SELECT * FROM auto_trades_view"
)

type autoTradeRepository struct {
	DB *sqlx.DB
}

func NewAutoTradeRepository(db *sqlx.DB) model.AutoTradeRepository {
	return &autoTradeRepository{
		DB: db,
	}
}

func (r *autoTradeRepository) AddAutoTrade(ctx context.Context, uid string, cryptoID int, strategyID int) error {
	stmt, err := r.DB.PrepareContext(ctx, queryAddAutoTrade)
	if err != nil {
		log.Printf("REPOSITORY: Unable to prepare update query: %v\n", err)
		return apperrors.NewInternal()
	}

	if _, err := stmt.ExecContext(ctx, uid, cryptoID, strategyID); err != nil {
		log.Printf("REPOSITORY: Failed to add auto trade for user: %v err: %s\n", uid, err.Error())
		return apperrors.NewInternal()
	}
	return nil
}

func (r *autoTradeRepository) DeleteAutoTrade(ctx context.Context, uid string, cryptoID int, strategyID int) error {
	stmt, err := r.DB.PrepareContext(ctx, queryDeleteAutoTrade)
	if err != nil {
		log.Printf("REPOSITORY: Unable to prepare update query: %v\n", err)
		return apperrors.NewInternal()
	}

	if _, err := stmt.ExecContext(ctx, uid, cryptoID, strategyID); err != nil {
		log.Printf("REPOSITORY: Failed to add auto trade for user: %v err: %s\n", uid, err.Error())
		return apperrors.NewInternal()
	}
	return nil
}

func (r *autoTradeRepository) GetAutoTrades(ctx context.Context, uid string) (*[]model.AutoTrade, error) {
	rst := &[]model.AutoTrade{}
	err := r.DB.SelectContext(ctx, rst, queryGetAutoTrades, uid)
	if err != nil {
		log.Printf("REPOSITORY: Unable to get auto trade setting (uid): %v err: %s", uid, err)
		return rst, apperrors.NewInternal()
	}
	return rst, nil
}

func (r *autoTradeRepository) GetAutoTradeFromStrategy(ctx context.Context, cryptoName string, strategyID int) (*[]model.AutoTrade, error) {
	rst := &[]model.AutoTrade{}
	err := r.DB.SelectContext(ctx, rst, queryGetAutoTradeUser, cryptoName, strategyID)
	if err != nil {
		log.Printf("REPOSITORY: Unable to get auto trade setting from (cryptoID: %v, strategyID: %v) err: %s", cryptoName, strategyID, err.Error())
	}
	return rst, nil
}

func (r *autoTradeRepository) GetAllAutoTrades(ctx context.Context) (*[]model.AutoTrade, error) {
	rst := &[]model.AutoTrade{}
	err := r.DB.SelectContext(ctx, rst, queryGetAllAutoTrade)
	if err != nil {
		return rst, apperrors.NewNotFound("resource", "autoTrade")
	}
	return rst, nil
}
