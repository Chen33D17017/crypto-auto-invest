package repository

import (
	"context"
	"crypto-auto-invest/model"
	"crypto-auto-invest/model/apperrors"
	"log"

	"github.com/jmoiron/sqlx"
)

/*
`tid` VARCHAR(100) NOT NULL PRIMARY KEY,
  `uid` VARCHAR(100) NOT NULL,
  `from_wid` VARCHAR(100) NOT NULL,
  `from_amount` float NOT NULL,
  `to_wid` VARCHAR(100) NOT NULL,
  `to_amount` float NOT NULL,
  `timestamp` VARCHAR(255) NOT NULL
*/

const (
	queryInsertTradeLog = `INSERT INTO trades (tid, uid, from_type, from_amount, to_type, to_amount, timestamp) 
							VALUES (:tid :?, :?, :?, :?, :?, :?)`
)

type tradeRepository struct {
	DB *sqlx.DB
}

func NewTradeRepository(db *sqlx.DB) model.TradeRepository {
	return &tradeRepository{
		DB: db,
	}
}

func (r *tradeRepository) SaveTrade(ctx context.Context, t *model.Trade) error {
	_, err := r.DB.NamedExecContext(ctx, queryInsertTradeLog, *t)
	if err != nil {
		log.Printf("RESPOSITORY: Fail to Insert Trade Log: %v, err: %s", t.Tid, err.Error())
		return apperrors.NewInternal()
	}
	return nil
}
