package repository

import (
	"context"
	"crypto-auto-invest/model"
	"crypto-auto-invest/model/apperrors"
	"log"

	"github.com/jmoiron/sqlx"
)

const (
	queryAddAutoTrade    = "INSERT INTO auto_trades(uid, type_id) VALUES(?, ?);"
	queryDeleteAutoTrade = "DELETE FROM auto_trades WHERE uid=? and type_id=?;"
	queryGetAutoTrades   = "SELECT * FROM auto_trades_view WHERE uid=?"
)

type autoTradeRepository struct {
	DB *sqlx.DB
}

func NewAutoTradeRepository(db *sqlx.DB) model.AutoTradeRepository {
	return &autoTradeRepository{
		DB: db,
	}
}

func (r *autoTradeRepository) AddAutoTrade(ctx context.Context, uid string, type_id int) error {
	stmt, err := r.DB.PrepareContext(ctx, queryAddAutoTrade)
	if err != nil {
		log.Printf("REPOSITORY: Unable to prepare update query: %v\n", err)
		return apperrors.NewInternal()
	}

	if _, err := stmt.ExecContext(ctx, uid, type_id); err != nil {
		log.Printf("REPOSITORY: Failed to add auto trade for user: %v err: %s\n", uid, err.Error())
		return apperrors.NewInternal()
	}
	return nil
}

func (r *autoTradeRepository) DeleteAutoTrade(ctx context.Context, uid string, type_id int) error {
	stmt, err := r.DB.PrepareContext(ctx, queryDeleteAutoTrade)
	if err != nil {
		log.Printf("REPOSITORY: Unable to prepare update query: %v\n", err)
		return apperrors.NewInternal()
	}

	if _, err := stmt.ExecContext(ctx, uid, type_id); err != nil {
		log.Printf("REPOSITORY: Failed to add auto trade for user: %v err: %s\n", uid, err.Error())
		return apperrors.NewInternal()
	}
	return nil
}

func (r *autoTradeRepository) GetAutoTrades(ctx context.Context, uid string) (*[]model.AutoTrade, error) {
	rst := &[]model.AutoTrade{}
	err := r.DB.SelectContext(ctx, rst, queryGetAutoTrades, uid)
	if err != nil {
		log.Printf("REPOSITORY: Unable to get crons by (uid): %v err: %s", uid, err)
		return rst, apperrors.NewInternal()
	}
	return rst, nil
}
