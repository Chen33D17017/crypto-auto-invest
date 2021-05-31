package repository

import (
	"context"
	"crypto-auto-invest/model"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

const (
	queryInsertTradeLog = `INSERT INTO orders (oid, uid, pair, action, amount, price, timestamp, fee, strategy_id) 
							VALUES (:oid, :uid, :pair, :action, :amount, :price, :timestamp, :fee, :strategy_id)`
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
		errString := fmt.Sprintf("RESPOSITORY: Fail to Insert Trade Log: %v, err: %s\n", t.OID, err.Error())
		log.Printf(errString)
		return fmt.Errorf(errString)
	}
	return nil
}
