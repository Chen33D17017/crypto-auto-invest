package services

import (
	"bytes"
	"context"
	"crypto-auto-invest/model"
	"crypto-auto-invest/model/apperrors"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/robfig/cron/v3"
)

type autoTradeService struct {
	TradeService        model.TradeService
	WalletRepository    model.WalletRepository
	UserRepository      model.UserRepository
	AutoTradeRepository model.AutoTradeRepository
	CronJobManager      model.CronJobManager
	Cron                *cron.Cron
	TimePattern         string
	TradeRateApi        string
	MaxRate             float64
}

type ATSConifg struct {
	WalletRepository    model.WalletRepository
	UserRepository      model.UserRepository
	AutoTradeRepository model.AutoTradeRepository
	TradeService        model.TradeService
	CronJobManager      model.CronJobManager
	Cron                *cron.Cron
	TimePattern         string
	TradeRateApi        string
	MaxRate             float64
}

func NewAutoTradeService(c *ATSConifg) model.AutoTradeService {
	return &autoTradeService{
		TradeService:        c.TradeService,
		WalletRepository:    c.WalletRepository,
		UserRepository:      c.UserRepository,
		AutoTradeRepository: c.AutoTradeRepository,
		CronJobManager:      c.CronJobManager,
		Cron:                c.Cron,
		TimePattern:         c.TimePattern,
		TradeRateApi:        c.TradeRateApi,
		MaxRate:             c.MaxRate,
	}
}

func (s *autoTradeService) AddAutoTrade(ctx context.Context, uid, cryptoName string) error {
	cID, err := s.WalletRepository.GetCurrencyID(ctx, cryptoName)
	if err != nil {
		return err
	}
	err = s.AutoTradeRepository.AddAutoTrade(ctx, uid, cID)
	if err != nil {
		return err
	}

	setting, err := s.AutoTradeRepository.GetAutoTrade(ctx, uid, cryptoName)
	if err != nil {
		return err
	}
	err = s.AddCronFunc(ctx, *setting)
	if err != nil {
		return err
	}
	return nil
}

func (s *autoTradeService) DeleteAutoTrade(ctx context.Context, uid, cryptoName string) error {
	setting, err := s.AutoTradeRepository.GetAutoTrade(ctx, uid, cryptoName)
	if err != nil {
		return err
	}

	err = s.RemoveCronFunc(ctx, setting.ID)
	if err != nil {
		return err
	}

	cID, err := s.WalletRepository.GetCurrencyID(ctx, cryptoName)
	if err != nil {
		return err
	}

	err = s.AutoTradeRepository.DeleteAutoTrade(ctx, uid, cID)
	if err != nil {
		return err
	}

	return nil
}

func (s *autoTradeService) GetAutoTrades(ctx context.Context, uid string) (*[]model.AutoTrade, error) {
	return s.AutoTradeRepository.GetAutoTrades(ctx, uid)
}

// Get Trade rate from other service
func (s *autoTradeService) GetTradeRate(reqBody model.TradeRateReq) (model.TradeRateRes, error) {
	var rst model.TradeRateRes

	payload, _ := json.Marshal(reqBody)
	payloadReader := bytes.NewReader(payload)

	client := &http.Client{}
	req, err := http.NewRequest("POST", s.TradeRateApi, payloadReader)

	if err != nil {
		errMsg := fmt.Sprintf("Fail to buiild request for getting trade rate: %s", err.Error())
		log.Println(errMsg)
		s.TradeService.SendTradeRst(errMsg, "error")
		return rst, err
	}

	resp, err := client.Do(req)
	if err != nil {
		errMsg := fmt.Sprintf("Fail to request for trade rate: %s", err.Error())
		log.Println(errMsg)
		s.TradeService.SendTradeRst(errMsg, "error")
		return rst, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&rst)
	if err != nil {
		errMsg := fmt.Sprintf("Trade rate response err: %s", err.Error())
		log.Println(errMsg)
		s.TradeService.SendTradeRst(errMsg, "error")
	}

	if rst.Rate == 0 || rst.Side == "none" {
		return rst, nil
	}

	if rst.Side != "buy" && rst.Side != "sell" {
		errMsg := fmt.Sprintf("Trade rate response err: side is not buy or sell")
		log.Println(errMsg, "error")
		return rst, err
	}

	if rst.Rate > s.MaxRate || rst.Rate < 0 {
		errMsg := fmt.Sprintf("Trade rate response err: rate is bigger than %v or less than 0", s.MaxRate)
		log.Println(errMsg)
		return rst, err
	}
	return rst, nil
}

func (s *autoTradeService) TradeWithRate(ctx context.Context, cryptoName string, action string, rate float64, strategy int) error {

	uids, err := s.AutoTradeRepository.GetAutoTradeUser(ctx, cryptoName)
	if err != nil {
		return apperrors.NewBadRequest("Fail to find auto trade user")
	}
	for _, uid := range *uids {
		err := s.tradeWithRate(ctx, uid, cryptoName, action, rate, strategy)
		if err != nil {
			s.TradeService.SendTradeRst(fmt.Sprintf("Fail to trade user %s with strategy %v err: %s", uid, strategy, err.Error()), "error")
		}
	}
	return nil
}

func (s *autoTradeService) tradeWithRate(ctx context.Context, uid string, cryptoName string, action string, rate float64, strategy int) error {
	u, err := s.UserRepository.FindByID(ctx, uid)
	if err != nil {
		return fmt.Errorf("autoTradeWithRate: %s", err.Error())
	}

	jpyWallet, err := s.WalletRepository.GetWellet(ctx, uid, "jpy")
	if err != nil {
		return fmt.Errorf("AutoTrade get jpy wallet err: %s", err.Error())
	}

	wallet, err := s.WalletRepository.GetWellet(ctx, uid, cryptoName)
	if err != nil {
		return fmt.Errorf("AutoTrade: get %s wallet err: %s\n", cryptoName, err.Error())
	}

	if rate > 0 {
		switch action {
		case "buy":
			_, err = s.TradeService.Trade(ctx, u, jpyWallet.Amount*rate, "buy", cryptoName, strategy)
		case "sell":
			_, err = s.TradeService.Trade(ctx, u, wallet.Amount*rate, "sell", cryptoName, strategy)
		default:
			return nil
		}
	}

	return nil
}

func (s *autoTradeService) AutoTrade(uid string, cryptoName string) error {
	ctx := context.TODO()
	u, err := s.UserRepository.FindByID(ctx, uid)
	if err != nil {
		return fmt.Errorf("AutoTrade find user err: %s", err.Error())
	}

	jpyWallet, err := s.WalletRepository.GetWellet(ctx, uid, "jpy")
	if err != nil {
		return fmt.Errorf("AutoTrade get jpy wallet err: %s", err.Error())
	}

	wallet, err := s.WalletRepository.GetWellet(ctx, uid, cryptoName)
	if err != nil {
		return fmt.Errorf("AutoTrade: get %s wallet err: %s\n", cryptoName, err.Error())
	}

	req := model.TradeRateReq{
		JPY:        jpyWallet.Amount,
		CryptoName: cryptoName,
		Amount:     wallet.Amount,
	}
	resp, err := s.GetTradeRate(req)
	if err != nil {
		return fmt.Errorf("AutoTrade: request err: %s\n", err.Error())
	}

	if resp.Rate > 0 {
		switch resp.Side {
		case "buy":
			_, err = s.TradeService.Trade(ctx, u, jpyWallet.Amount*resp.Rate, "buy", cryptoName, resp.Strategy)
		case "sell":
			_, err = s.TradeService.Trade(ctx, u, wallet.Amount*resp.Rate, "sell", cryptoName, resp.Strategy)
		default:
			return nil
		}
	}
	return nil
}

func (s *autoTradeService) AddCronFunc(ctx context.Context, setting model.AutoTrade) error {
	entityID, err := s.Cron.AddFunc(s.TimePattern, func() {
		s.AutoTrade(setting.UID, setting.CryptoName)
	})
	if err != nil {
		return err
	}

	if err := s.CronJobManager.SetCronJob(ctx, fmt.Sprintf("auto:%s", setting.ID), int(entityID)); err != nil {
		return err
	}
	return nil
}

func (s *autoTradeService) RemoveCronFunc(ctx context.Context, autoTradeID string) error {
	entityID, err := s.CronJobManager.GetAndDeleteCronJob(ctx, fmt.Sprintf("auto:%s", autoTradeID))
	if err != nil {
		return err
	}
	s.Cron.Remove(cron.EntryID(entityID))
	return nil
}
