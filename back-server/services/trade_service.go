package services

import (
	"context"
	"crypto-auto-invest/bitbank"
	bm "crypto-auto-invest/bitbank/model"
	"crypto-auto-invest/model"
	"crypto-auto-invest/model/apperrors"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
)

type tradeService struct {
	TradeRepository  model.TradeRepository
	WalletRepository model.WalletRepository
	Delay            time.Duration
}

type TSConifg struct {
	TradeRepository  model.TradeRepository
	WalletRepository model.WalletRepository
	Delay            time.Duration
}

func NewTradeService(c *TSConifg) model.TradeService {
	return &tradeService{
		TradeRepository:  c.TradeRepository,
		WalletRepository: c.WalletRepository,
		Delay:            c.Delay,
	}
}

func (s *tradeService) Trade(ctx context.Context, u *model.User, amount float64, side, assetType, orderType string) (bm.Order, error) {
	secret := bm.Secret{
		ApiKey:    u.ApiKey,
		ApiSecret: u.ApiSecret,
	}
	var order bm.Order
	var err error
	switch side {
	case "buy":
		order, err = bitbank.BuyWithJPY(secret, assetType, int64(amount))
	case "sell":
		order, err = bitbank.SellToJPY(secret, assetType, amount)
	default:
		return order, apperrors.NewInternal()
	}
	if err != nil {
		log.Printf("SERVICE: Trade err with user: %s assetType: %s, Amount: %v, Side: %s\n", u.UID, assetType, amount, side)
		return order, apperrors.NewInternal()
	}

	time.AfterFunc(s.Delay, func() {
		s.SaveOrder(context.TODO(), u, fmt.Sprintf("%v", order.OrderId), assetType, orderType)
	})
	return order, nil
}

func (s *tradeService) SaveOrder(ctx context.Context, u *model.User, orderID string, assetType, orderType string) error {
	secret := bm.Secret{
		ApiKey:    u.ApiKey,
		ApiSecret: u.ApiSecret,
	}
	o, err := bitbank.GetOrderInfo(secret, assetType, orderID)
	var target model.Order
	target.UID = u.UID
	target.OID = fmt.Sprintf("%v", o.OrderId)
	amount, err := strconv.ParseFloat(o.StartAmount, 64)
	if err != nil {
		log.Printf("Fail to convert Amount")
		return apperrors.NewBadRequest("Wrong struct on order")
	}
	amount = normalizeFloat(amount)
	avgPrice, err := strconv.ParseFloat(o.AveragePrice, 64)

	if err != nil {
		log.Printf("Fail to convert AvgPrice")
		return apperrors.NewBadRequest("Wrong struct on order")
	}

	currencies := strings.Split(o.Pair, "_")
	walletFir, err1 := s.WalletRepository.GetWellet(ctx, u.UID, currencies[0])
	walletSec, err2 := s.WalletRepository.GetWellet(ctx, u.UID, currencies[1])
	if err1 != nil || err2 != nil {
		log.Printf("Wrong cuncerrency with user %v", u.UID)
		return apperrors.NewInternal()
	}

	switch o.Side {
	case "buy":
		target.FromAmount = amount * avgPrice
		target.ToAmount = amount
		target.FromWallet = walletSec.WID
		target.ToWallet = walletFir.WID
	case "sell":
		target.FromAmount = amount
		target.ToAmount = amount * float64(avgPrice)
		target.FromWallet = walletSec.WID
		target.ToWallet = walletFir.WID
	}

	target.Timestamp = time.Unix(o.OrderedAt/1000, 0).Format("2006-01-02 15:04:05")
	target.Type = orderType

	err = s.TradeRepository.SaveOrder(ctx, &target)
	if err != nil {
		log.Printf("Fail to Store Trade Result with %v err: %s\n", o.OrderId, err.Error())
		return apperrors.NewInternal()
	}
	return nil
}

func normalizeFloat(num float64) float64 {
	return math.Round(num*10000) / 10000
}
