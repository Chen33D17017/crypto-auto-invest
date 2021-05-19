package services

import (
	"bytes"
	"context"
	"crypto-auto-invest/bitbank"
	bm "crypto-auto-invest/bitbank/model"
	"crypto-auto-invest/model"
	"crypto-auto-invest/model/apperrors"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type DiscordFormat struct {
	Msg string `json:"content"`
}
type tradeService struct {
	TradeRepository  model.TradeRepository
	WalletRepository model.WalletRepository
	InfoWebhook      string
	ErrorWebhook     string
}

type TSConifg struct {
	TradeRepository  model.TradeRepository
	WalletRepository model.WalletRepository
	InfoWebhook      string
	ErrorWebhook     string
}

func NewTradeService(c *TSConifg) model.TradeService {
	return &tradeService{
		TradeRepository:  c.TradeRepository,
		WalletRepository: c.WalletRepository,
		InfoWebhook:      c.InfoWebhook,
		ErrorWebhook:     c.ErrorWebhook,
	}
}

func (s *tradeService) Trade(ctx context.Context, u *model.User, amount float64, getDelay time.Duration, side, assetType, orderType string) (bm.Order, error) {
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
		s.SendTradeRst(fmt.Sprintf("SERVICE: Trade err with user: %s assetType: %s, Amount: %v, Side: %s\n", u.UID, assetType, amount, side), "error")
		return order, apperrors.NewInternal()
	}

	time.AfterFunc(getDelay, func() {
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
	if err != nil {
		s.SendTradeRst(fmt.Sprintf("%s fail to get order with assertType: %s, OrderID: %s", u.Name, assetType, orderID), "error")
		return apperrors.NewInternal()
	}
	var target model.Order
	target.UID = u.UID
	target.OID = fmt.Sprintf("%v", o.OrderId)
	amount, err := strconv.ParseFloat(o.StartAmount, 64)
	if err != nil {
		log.Printf("Fail to convert Amount")
		s.SendTradeRst(fmt.Sprintf("%s fail to save order with assertType: %s, OrderID: %s", u.Name, assetType, orderID), "error")
		return apperrors.NewBadRequest("Wrong struct on order")
	}
	amount = normalizeFloat(amount)
	avgPrice, err := strconv.ParseFloat(o.AveragePrice, 64)

	if err != nil {
		log.Printf("Fail to convert AvgPrice")
		s.SendTradeRst(fmt.Sprintf("%s fail to save order with assertType: %s, OrderID: %s", u.Name, assetType, orderID), "error")
		return apperrors.NewBadRequest("Wrong struct on order")
	}

	currencies := strings.Split(o.Pair, "_")
	walletFir, err1 := s.WalletRepository.GetWellet(ctx, u.UID, currencies[0])
	walletSec, err2 := s.WalletRepository.GetWellet(ctx, u.UID, currencies[1])
	if err1 != nil || err2 != nil {
		log.Printf("Wrong cuncerrency with user %v", u.UID)
		s.SendTradeRst(fmt.Sprintf("%s fail to save order with assertType: %s, OrderID: %s", u.Name, assetType, orderID), "error")
		return apperrors.NewInternal()
	}
	JPY := math.Round(amount * avgPrice * 1.0012)
	// TODO: adding fee on JPY
	switch o.Side {
	case "buy":
		target.FromAmount = JPY
		target.ToAmount = amount
		target.FromWallet = walletSec.WID
		target.ToWallet = walletFir.WID
	case "sell":
		target.FromAmount = amount
		target.ToAmount = JPY
		target.FromWallet = walletSec.WID
		target.ToWallet = walletFir.WID
	}

	target.Timestamp = time.Unix(o.OrderedAt/1000, 0).Format("2006-01-02 15:04:05")
	target.Type = orderType

	err = s.TradeRepository.SaveOrder(ctx, &target)
	if err != nil {
		s.SendTradeRst(fmt.Sprintf("%s fail to save order with assertType: %s, OrderID: %s", u.Name, assetType, orderID), "error")
		log.Printf("Fail to Store Trade Result with %v err: %s\n", o.OrderId, err.Error())
		return apperrors.NewInternal()
	}
	loc := time.FixedZone("UTC+9", 9*60*60)
	s.SendTradeRst(fmt.Sprintf("%s %s %v(%s, ¥%s) cost ¥%v @%v",
		u.Name, o.Side, o.StartAmount, o.Pair, o.AveragePrice, JPY, time.Unix(o.OrderedAt/1000, 0).In(loc).Format(time.RFC822)), "info")
	return nil
}

func (s *tradeService) SendTradeRst(msg string, level string) error {
	var url string
	switch level {
	case "info":
		url = s.InfoWebhook
	case "error":
		url = s.ErrorWebhook
	default:
		return fmt.Errorf("Wrong type of level")
	}
	msgJSON, _ := json.Marshal(DiscordFormat{msg})
	payload := bytes.NewReader(msgJSON)

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, payload)

	if err != nil {
		return fmt.Errorf("Fail to send msg to Discord")
	}
	req.Header.Add("Content-Type", "application/json")
	client.Do(req)

	return nil
}

func normalizeFloat(num float64) float64 {
	return math.Round(num*10000) / 10000
}
