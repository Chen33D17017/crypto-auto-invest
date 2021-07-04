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
	queryInserOrder = `INSERT INTO binance_orders (uid, symbol, action, amount, price, timestamp) 
	VALUES (:uid, :symbol, :action, :amount, :price, :timestamp)`
	queryBinanceGetOrders = "SELECT * FROM binance_orders WHERE uid=? and symbol=?"
)

type binanceTradeRepository struct {
	DB *sqlx.DB
}

func NewBinanceTradeRepository(db *sqlx.DB) model.BinanceTradeRepository {
	return &binanceTradeRepository{
		DB: db,
	}
}

func (r *binanceTradeRepository) SaveOrder(ctx context.Context, o *model.BinanceOrder) error {
	_, err := r.DB.NamedExecContext(ctx, queryInserOrder, *o)
	if err != nil {
		errString := fmt.Sprintf("RESPOSITORY: Fail to Insert Binance Log: %v, err: %s\n", o.UID, err.Error())
		log.Printf(errString)
		return fmt.Errorf(errString)
	}
	return nil
}

func (r *binanceTradeRepository) GetOrders(ctx context.Context, uid string, symbol string) (*[]model.BinanceOrder, error) {
	rst := &[]model.BinanceOrder{}
	err := r.DB.SelectContext(ctx, rst, queryBinanceGetOrders, uid, symbol)
	if err != nil {
		log.Printf("REPOSITORY: Unable to get order logs by (uid: %s, cyptoName: %s) err: %s", uid, symbol, err)
		return rst, apperrors.NewNotFound("order_log", symbol)
	}
	return rst, nil
}
