package services

import (
	"bytes"
	"context"
	"crypto-auto-invest/model"
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

func (s *autoTradeService) AddAutoTrade(ctx context.Context, uid, currencyName string) error {
	cID, err := s.WalletRepository.GetCurrencyID(ctx, currencyName)
	if err != nil {
		return err
	}
	err = s.AutoTradeRepository.AddAutoTrade(ctx, uid, cID)
	if err != nil {
		return err
	}

	setting, err := s.AutoTradeRepository.GetAutoTrade(ctx, uid, currencyName)
	if err != nil {
		return err
	}
	err = s.AddCronFunc(ctx, *setting)
	if err != nil {
		return err
	}
	return nil
}

func (s *autoTradeService) DeleteAutoTrade(ctx context.Context, uid, currencyName string) error {
	cID, err := s.WalletRepository.GetCurrencyID(ctx, currencyName)
	if err != nil {
		return err
	}
	err = s.AutoTradeRepository.DeleteAutoTrade(ctx, uid, cID)
	if err != nil {
		return err
	}
	setting, err := s.AutoTradeRepository.GetAutoTrade(ctx, uid, currencyName)
	if err != nil {
		return err
	}

	err = s.RemoveCronFunc(ctx, setting.ID)
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

func (s *autoTradeService) AutoTrade(uid string, currencyName string) error {
	ctx := context.TODO()
	u, err := s.UserRepository.FindByID(ctx, uid)
	if err != nil {
		return fmt.Errorf("AutoTrade find user err: %s", err.Error())
	}

	jpyWallet, err := s.WalletRepository.GetWellet(ctx, uid, "jpy")
	if err != nil {
		return fmt.Errorf("AutoTrade get jpy wallet err: %s", err.Error())
	}

	wallet, err := s.WalletRepository.GetWellet(ctx, uid, currencyName)
	if err != nil {
		return fmt.Errorf("AutoTrade: get %s wallet err: %s\n", currencyName, err.Error())
	}

	req := model.TradeRateReq{
		JPY:    jpyWallet.Amount,
		Type:   currencyName,
		Amount: wallet.Amount,
	}
	resp, err := s.GetTradeRate(req)
	if err != nil {
		return fmt.Errorf("AutoTrade: request err: %s\n", err.Error())
	}

	if resp.Rate > 0 {
		switch resp.Side {
		case "buy":
			_, err = s.TradeService.Trade(ctx, u, jpyWallet.Amount*resp.Rate, "buy", currencyName, "auto")
		case "sell":
			_, err = s.TradeService.Trade(ctx, u, wallet.Amount*resp.Rate, "sell", currencyName, "auto")
		default:
			return nil
		}
	}
	return nil
}

func (s *autoTradeService) AddCronFunc(ctx context.Context, setting model.AutoTrade) error {
	entityID, err := s.Cron.AddFunc(s.TimePattern, func() {
		s.AutoTrade(setting.UID, setting.Type)
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
