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
	Delay            time.Duration
	InfoWebhook      string
	ErrorWebhook     string
}

type TSConifg struct {
	TradeRepository  model.TradeRepository
	WalletRepository model.WalletRepository
	Delay            time.Duration
	InfoWebhook      string
	ErrorWebhook     string
}

func NewTradeService(c *TSConifg) model.TradeService {
	return &tradeService{
		TradeRepository:  c.TradeRepository,
		WalletRepository: c.WalletRepository,
		Delay:            c.Delay,
		InfoWebhook:      c.InfoWebhook,
		ErrorWebhook:     c.ErrorWebhook,
	}
}

func (s *tradeService) MarketTrade(ctx context.Context, u *model.User, amount float64, action, cryptoName string, strategyID int) (bm.Order, error) {
	secret := bm.Secret{
		ApiKey:    u.ApiKey,
		ApiSecret: u.ApiSecret,
	}
	var order bm.Order
	var err error
	switch action {
	case "buy":
		order, err = bitbank.BuyWithJPY(secret, cryptoName, int64(amount))
	case "sell":
		order, err = bitbank.SellToJPY(secret, cryptoName, amount)
	default:
		return order, apperrors.NewInternal()
	}
	if err != nil {
		s.sendTradeRst(fmt.Sprintf("SERVICE: Trade err with user: %s cryptoName: %s, Amount: %v, Side: %s\n", u.UID, cryptoName, amount, action), "error", true)
		return order, apperrors.NewInternal()
	}

	time.AfterFunc(s.Delay, func() {
		s.SaveOrder(context.TODO(), u, fmt.Sprintf("%v", order.OrderId), cryptoName, strategyID, false)
	})

	time.AfterFunc(s.Delay+time.Duration(time.Second*30), func() {
		incomeRate, err := s.CalIncomeRate(context.TODO(), u.UID, cryptoName, strategyID)
		if err != nil {
			s.sendTradeRst(fmt.Sprintf("Fail to calculate %s's income rate on (cryptoName %s, strategy %v)", u.Name, cryptoName, strategyID), "error", true)
		}
		s.sendTradeRst(fmt.Sprintf("%s's income rate on (cryptoName %s, strategy %v): %s", u.Name, cryptoName, strategyID, incomeRate.IncomeRate), "info", true)
	})
	return order, nil
}

func (s *tradeService) LimitTrade(ctx context.Context, u *model.User, amount float64, action, cryptoName string, price float64) (bm.Order, error) {
	var order bm.Order
	secret := bm.Secret{
		ApiKey:    u.ApiKey,
		ApiSecret: u.ApiSecret,
	}
	order, err := bitbank.MakeTrade(secret, cryptoName, action, amount, "limit", price, true)
	if err != nil {
		return order, apperrors.NewInternalWithReason(fmt.Sprintf("SERVICE MakeOrder: %s", err.Error()))
	}

	return order, nil
}

func (s *tradeService) SaveOrder(ctx context.Context, u *model.User, orderID string, cryptoName string, strategyID int, info bool) error {
	secret := bm.Secret{
		ApiKey:    u.ApiKey,
		ApiSecret: u.ApiSecret,
	}
	o, err := bitbank.GetOrderInfo(secret, cryptoName, orderID)
	if err != nil {
		s.sendTradeRst(fmt.Sprintf("%s fail to get order with cryptoName: %s, OrderID: %s", u.Name, cryptoName, orderID), "error", info)
		return apperrors.NewInternalWithReason(fmt.Sprintf("SERVICE SaveOrder: %s", err.Error()))
	}

	// if the order is not be fully filled yet
	if o.Status != "FULLY_FILLED" {
		time.AfterFunc(time.Duration(time.Second*30), func() {
			s.SaveOrder(context.TODO(), u, orderID, cryptoName, strategyID, false)
		})
		return nil
	}

	var target model.Order
	target.OID = fmt.Sprintf("%v", o.OrderId)
	target.UID = u.UID
	target.Piar = o.Pair
	target.Action = o.Side

	amount, err := strconv.ParseFloat(o.StartAmount, 64)
	if err != nil {
		s.sendTradeRst(fmt.Sprintf("%s fail to save order with cryptoName: %s, OrderID: %s", u.Name, cryptoName, orderID), "error", info)
		return apperrors.NewBadRequest("Fail to convert Amount")
	}
	amount = normalizeFloat(amount)
	target.Amount = amount

	avgPrice, err := strconv.ParseFloat(o.AveragePrice, 64)
	if err != nil {
		log.Printf("Fail to convert AvgPrice")
		s.sendTradeRst(fmt.Sprintf("%s fail to save order with cryptoName: %s, OrderID: %s", u.Name, cryptoName, orderID), "error", info)
		return apperrors.NewBadRequest("Wrong struct on order")
	}
	target.Price = avgPrice
	target.Timestamp = time.Unix(o.OrderedAt/1000, 0).Format("2006-01-02 15:04:05")

	currencies := strings.Split(o.Pair, "_")

	JPY := normalizeFloat(amount * avgPrice)
	if o.Type == "limit" {
		target.Fee = -normalizeFloat(amount * avgPrice * 0.0002)
	} else {
		target.Fee = normalizeFloat(amount * avgPrice * 0.0012)
	}
	target.Strategy = strategyID

	err = s.TradeRepository.SaveOrder(ctx, &target)
	if err != nil {
		s.sendTradeRst(fmt.Sprintf("%s fail to save order with cryptoName: %s, OrderID: %s", u.Name, cryptoName, orderID), "error", info)
		return apperrors.NewInternalWithReason(fmt.Sprintf("SERVICE SaveOrder: %s", err.Error()))
	}

	loc := time.FixedZone("UTC+9", 9*60*60)
	s.sendTradeRst(fmt.Sprintf("%s %s %v(%s, ??%s) with ??%v @%v",
		u.Name, o.Side, o.StartAmount, o.Pair, o.AveragePrice, JPY, time.Unix(o.OrderedAt/1000, 0).In(loc).Format(time.RFC822)), "info", info)

	// Money movement between sub wallets when strategy is not 0
	if strategyID != 0 {
		cryptoWallet, err1 := s.WalletRepository.GetWellet(ctx, u.UID, currencies[0], strategyID)
		JPYwallet, err2 := s.WalletRepository.GetWellet(ctx, u.UID, currencies[1], strategyID)
		if err1 != nil || err2 != nil {
			log.Printf("Wrong cuncerrency with user %v", u.UID)
			s.sendTradeRst(fmt.Sprintf("%s fail to save order with cryptoName: %s, OrderID: %s", u.Name, cryptoName, orderID), "error", info)
			return apperrors.NewInternal()
		}
		switch o.Side {
		case "buy":
			s.WalletRepository.UpdateAmount(ctx, JPYwallet.WID, JPYwallet.Amount-(JPY+target.Fee))
			s.WalletRepository.UpdateAmount(ctx, cryptoWallet.WID, cryptoWallet.Amount+amount)
		case "sell":
			s.WalletRepository.UpdateAmount(ctx, cryptoWallet.WID, cryptoWallet.Amount-amount)
			s.WalletRepository.UpdateAmount(ctx, JPYwallet.WID, JPYwallet.Amount+(JPY-target.Fee))
		}
	}

	return nil
}

func (s *tradeService) CalIncomeRate(ctx context.Context, uid string, cryptoName string, strategyID int) (*model.Income, error) {
	rst := &model.Income{}
	orders, err := s.TradeRepository.GetOrderLogs(ctx, uid, cryptoName, strategyID)
	if err != nil {
		return rst, err
	}

	cost := 0.0
	amount := 0.0
	JPY := 0.0
	for _, order := range *orders {
		if order.Action == "buy" {
			cost += order.Amount*order.Price + order.Fee
			amount += order.Amount
		} else {
			amount -= order.Amount
			JPY += (order.Amount*order.Price - order.Fee)
		}
		cost = normalizeFloat(cost)
		amount = normalizeFloat(amount)
		JPY = normalizeFloat(JPY)
	}
	price, err := bitbank.GetPrice(cryptoName)
	if err != nil {
		log.Printf("SERVICE: Fail to get crypto price %s err: %s\n", cryptoName, err.Error())
		return rst, apperrors.NewInternal()
	}
	lastPrice, err := strconv.ParseFloat(price.Last, 64)
	if err != nil {
		log.Printf("SERVICE: Fail to get crypto price %s err: %s\n", cryptoName, err.Error())
		return rst, apperrors.NewInternal()
	}
	incomeRate := normalizeFloat((amount*lastPrice + JPY - cost) / cost * 100)
	rst.CryptoName = cryptoName
	rst.Strategy = strategyID
	rst.Amount = amount
	rst.Cost = cost
	rst.JPY = JPY
	rst.IncomeRate = fmt.Sprintf("%v%%", incomeRate)
	if strategyID != 0 {
		chargeAmount := 0.0
		chargeJPY := 0.0
		cryptoWallet, _ := s.WalletRepository.GetWellet(ctx, uid, cryptoName, strategyID)
		JPYWallet, _ := s.WalletRepository.GetWellet(ctx, uid, "jpy", strategyID)
		if err != nil {
			return rst, err
		}
		logs, err := s.WalletRepository.GetChargeLogs(ctx, uid, cryptoName, strategyID)
		if err != nil {
			return rst, err
		}
		for _, log := range *logs {
			chargeAmount += log.Amount
		}

		logs, err = s.WalletRepository.GetChargeLogs(ctx, uid, "jpy", strategyID)
		if err != nil {
			return rst, err
		}
		for _, log := range *logs {
			chargeJPY += log.Amount
		}
		JPYincome := JPYWallet.Amount - chargeJPY
		cryptoIncmome := cryptoWallet.Amount - chargeAmount
		rst.Deposit = chargeAmount*lastPrice + chargeJPY
		rst.DepositIncomeRate = fmt.Sprintf("%v%%", normalizeFloat((cryptoIncmome*lastPrice+JPYincome)/rst.Deposit*100))
	} else {
		rst.Deposit = cost
		rst.DepositIncomeRate = fmt.Sprintf("%v%%", incomeRate)
	}
	return rst, nil
}

func (s *tradeService) sendTradeRst(msg string, level string, info bool) error {
	var url string
	if !info {
		return nil
	}
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
