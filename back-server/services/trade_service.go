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

// Buy unit JPY
// Sell unit crypto currency
func (s *tradeService) Trade(ctx context.Context, u *model.User, amount float64, action, cryptoName string, strategyID int) (bm.Order, error) {
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
		s.SendTradeRst(fmt.Sprintf("SERVICE: Trade err with user: %s cryptoName: %s, Amount: %v, Side: %s\n", u.UID, cryptoName, amount, action), "error")
		return order, apperrors.NewInternal()
	}

	time.AfterFunc(s.Delay, func() {
		s.SaveOrder(context.TODO(), u, fmt.Sprintf("%v", order.OrderId), cryptoName, strategyID)
	})

	time.AfterFunc(s.Delay+time.Duration(time.Second*30), func() {
		incomeRate, err := s.CalIncomeRate(context.TODO(), u.UID, cryptoName, strategyID)
		if err != nil {
			s.SendTradeRst(fmt.Sprintf("Fail to calculate %s's income rate on (cryptoName %s, strategy %v)", u.Name, cryptoName, strategyID), "error")
		}
		s.SendTradeRst(fmt.Sprintf("%s's income rate on (cryptoName %s, strategy %v): %v", u.Name, cryptoName, strategyID, incomeRate), "info")
	})
	return order, nil
}

func (s *tradeService) SaveOrder(ctx context.Context, u *model.User, orderID string, cryptoName string, strategy_id int) error {
	secret := bm.Secret{
		ApiKey:    u.ApiKey,
		ApiSecret: u.ApiSecret,
	}
	o, err := bitbank.GetOrderInfo(secret, cryptoName, orderID)
	if err != nil {
		s.SendTradeRst(fmt.Sprintf("%s fail to get order with cryptoName: %s, OrderID: %s", u.Name, cryptoName, orderID), "error")
		return apperrors.NewInternal()
	}
	var target model.Order
	target.OID = fmt.Sprintf("%v", o.OrderId)
	target.UID = u.UID
	target.Piar = o.Pair
	target.Action = o.Side

	amount, err := strconv.ParseFloat(o.StartAmount, 64)
	if err != nil {
		log.Printf("Fail to convert Amount")
		s.SendTradeRst(fmt.Sprintf("%s fail to save order with cryptoName: %s, OrderID: %s", u.Name, cryptoName, orderID), "error")
		return apperrors.NewBadRequest("Wrong struct on order")
	}
	amount = normalizeFloat(amount)
	target.Amount = amount

	avgPrice, err := strconv.ParseFloat(o.AveragePrice, 64)
	if err != nil {
		log.Printf("Fail to convert AvgPrice")
		s.SendTradeRst(fmt.Sprintf("%s fail to save order with cryptoName: %s, OrderID: %s", u.Name, cryptoName, orderID), "error")
		return apperrors.NewBadRequest("Wrong struct on order")
	}
	target.Price = avgPrice
	target.Timestamp = time.Unix(o.OrderedAt/1000, 0).Format("2006-01-02 15:04:05")

	currencies := strings.Split(o.Pair, "_")

	JPY := normalizeFloat(amount * avgPrice)
	target.Fee = normalizeFloat(amount * avgPrice * 0.0012)
	target.Strategy = strategy_id

	err = s.TradeRepository.SaveOrder(ctx, &target)
	if err != nil {
		s.SendTradeRst(fmt.Sprintf("%s fail to save order with cryptoName: %s, OrderID: %s", u.Name, cryptoName, orderID), "error")
		log.Printf("Fail to Store Trade Result with %v err: %s\n", o.OrderId, err.Error())
		return apperrors.NewInternal()
	}
	loc := time.FixedZone("UTC+9", 9*60*60)
	s.SendTradeRst(fmt.Sprintf("%s %s %v(%s, ¥%s) with ¥%v @%v",
		u.Name, o.Side, o.StartAmount, o.Pair, o.AveragePrice, JPY, time.Unix(o.OrderedAt/1000, 0).In(loc).Format(time.RFC822)), "info")

	// Money movement between wallets when orderType is auto
	if strategy_id != 0 {
		JPYWallet, err1 := s.WalletRepository.GetWellet(ctx, u.UID, currencies[0], strategy_id)
		currencyWallet, err2 := s.WalletRepository.GetWellet(ctx, u.UID, currencies[1], strategy_id)
		if err1 != nil || err2 != nil {
			log.Printf("Wrong cuncerrency with user %v", u.UID)
			s.SendTradeRst(fmt.Sprintf("%s fail to save order with cryptoName: %s, OrderID: %s", u.Name, cryptoName, orderID), "error")
			return apperrors.NewInternal()
		}
		switch o.Side {
		case "buy":
			s.WalletRepository.UpdateAmount(ctx, JPYWallet.WID, -(JPY + target.Fee))
			s.WalletRepository.UpdateAmount(ctx, currencyWallet.WID, amount)
		case "sell":
			s.WalletRepository.UpdateAmount(ctx, JPYWallet.WID, -amount)
			s.WalletRepository.UpdateAmount(ctx, currencyWallet.WID, (JPY - target.Fee))
		}
	}

	return nil
}

func (s *tradeService) CalIncomeRate(ctx context.Context, uid string, cryptoName string, strategyID int) (float64, error) {
	orders, err := s.TradeRepository.GetOrderLogs(ctx, uid, cryptoName, strategyID)
	if err != nil {
		return 0.0, err
	}

	cost := 0.0
	amount := 0.0
	for _, order := range *orders {
		if order.Action == "buy" {
			cost += order.Amount*order.Price + order.Fee
			amount += order.Amount
		} else {
			amount -= order.Amount
			cost -= (order.Amount*order.Price - order.Fee)
		}
		cost = normalizeFloat(cost)
		amount = normalizeFloat(amount)
	}
	price, err := bitbank.GetPrice(cryptoName)
	if err != nil {
		log.Printf("SERVICE: Fail to get crypto price %s err: %s\n", cryptoName, err.Error())
		return 0.0, apperrors.NewInternal()
	}
	lastPrice, err := strconv.ParseFloat(price.Last, 64)
	if err != nil {
		log.Printf("SERVICE: Fail to get crypto price %s err: %s\n", cryptoName, err.Error())
		return 0.0, apperrors.NewInternal()
	}
	incomeRate := normalizeFloat((amount * lastPrice) / cost * 100)
	return incomeRate, nil
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
