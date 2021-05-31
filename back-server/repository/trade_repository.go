package repository

import (
	"context"
	"crypto-auto-invest/model"
	"crypto-auto-invest/model/apperrors"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

const (
	queryInsertTradeLog = `INSERT INTO orders (oid, uid, pair, action, amount, price, timestamp, fee, strategy_id) 
							VALUES (:oid, :uid, :pair, :action, :amount, :price, :timestamp, :fee, :strategy_id)`
	queryGetTradeLogs = `SELECT * FROM orders WHERE uid=? AND pair=? and strategy_id=?`

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

func (r *tradeRepository) GetOrderLogs(ctx context.Context, uid, cryptoName string, strategyID int) (*[]model.Order, error) {
	rst := &[]model.Order{}
	err := r.DB.SelectContext(ctx, rst, queryGetTradeLogs, uid, fmt.Sprintf("%s_jpy", cryptoName), strategyID)
	if err != nil {
		log.Printf("REPOSITORY: Unable to get trade logs by (uid: %s, cyptoName: %s, strategyID:%v) err: %s", uid, cryptoName, strategyID, err)
		return rst, apperrors.NewNotFound("order_log", cryptoName)
	}
	return rst, nil
}
