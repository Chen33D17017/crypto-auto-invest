package repository

import (
	"context"
	"crypto-auto-invest/model"
	"crypto-auto-invest/model/apperrors"
	"log"

	"github.com/jmoiron/sqlx"
)

const (
	queryInsertTradeLog = `INSERT INTO orders (oid, uid, from_wid, from_amount, to_wid, to_amount, timestamp, fee, strategy) 
							VALUES (:oid, :uid, :from_wid, :from_amount, :to_wid, :to_amount, :timestamp, :fee, :strategy)`
)

type tradeRepository struct {
	DB *sqlx.DB
}

func NewTradeRepository(db *sqlx.DB) model.TradeRepository {
	return &tradeRepository{
		DB: db,
	}
}

func (r *tradeRepository) SaveOrder(ctx context.Context, t *model.Order) error {
	_, err := r.DB.NamedExecContext(ctx, queryInsertTradeLog, *t)
	if err != nil {
		log.Printf("RESPOSITORY: Fail to Insert Trade Log: %v, err: %s", t.OID, err.Error())
		return apperrors.NewInternal()
	}
	return nil
}
