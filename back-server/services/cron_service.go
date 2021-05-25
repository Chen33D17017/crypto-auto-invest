package services

import (
	"context"
	"crypto-auto-invest/model"
	"fmt"

	"github.com/robfig/cron/v3"
)

// https://pkg.go.dev/github.com/robfig/cron/v3@v3.0.0?utm_source=gopls
type cronService struct {
	CronRepository   model.CronRepository
	UserRepository   model.UserRepository
	WalletRepository model.WalletRepository
	TradeService     model.TradeService
	CronJobManager   model.CronJobManager
	Cron             *cron.Cron
}

type CSConfig struct {
	CronRepository   model.CronRepository
	WalletRepository model.WalletRepository
	UserRepository   model.UserRepository
	TradeService     model.TradeService
	CronJobManager   model.CronJobManager
	Cron             *cron.Cron
}

func NewCronService(c *CSConfig) model.CronService {
	return &cronService{
		CronRepository:   c.CronRepository,
		UserRepository:   c.UserRepository,
		WalletRepository: c.WalletRepository,
		TradeService:     c.TradeService,
		CronJobManager:   c.CronJobManager,
		Cron:             c.Cron,
	}
}

// https://github.com/robfig/cron/blob/bc59245fe10efaed9d51b56900192527ed733435/cron.go#L50
// https://pkg.go.dev/github.com/robfig/cron/v3@v3.0.0?utm_source=gopls#EntryID
func (s *cronService) AddCron(ctx context.Context, cb *model.Cron) (*model.Cron, error) {
	currencyID, err := s.WalletRepository.GetCurrencyID(ctx, cb.CryptoName)
	err = s.CronRepository.AddCron(ctx, cb, currencyID)
	if err != nil {
		return nil, err
	}
	cid, err := s.CronRepository.GetCronID(ctx, cb.UID, cb.CryptoName, cb.TimePattern)
	if err != nil {
		return nil, err
	}
	cb.ID = cid

	// Add func to cron job
	err = s.AddCronFunc(ctx, cb)
	if err != nil {
		return cb, err
	}

	return cb, nil
}

func (s *cronService) GetCron(ctx context.Context, uid, cronID string) (*model.Cron, error) {
	var rst *model.Cron
	rst, err := s.CronRepository.GetCron(ctx, uid, cronID)
	if err != nil {
		return nil, err
	}
	return rst, nil
}

func (s *cronService) GetCrons(ctx context.Context, uid string) (*[]model.Cron, error) {
	rst, err := s.CronRepository.GetCrons(ctx, uid)
	if err != nil {
		return nil, err
	}
	return rst, nil
}

func (s *cronService) UpdateCron(ctx context.Context, cb *model.Cron) error {
	currencyID, err := s.WalletRepository.GetCurrencyID(ctx, cb.CryptoName)
	err = s.CronRepository.UpdateCron(ctx, cb, currencyID)
	if err != nil {
		return err
	}
	err = s.RemoveCronFunc(ctx, cb.ID)
	if err != nil {
		return err
	}

	err = s.AddCronFunc(ctx, cb)
	if err != nil {
		return err
	}

	return nil
}

func (s *cronService) DeleteCron(ctx context.Context, uid string, cronID string) error {
	err := s.CronRepository.DeleteCron(ctx, uid, cronID)
	if err != nil {
		return err
	}
	err = s.RemoveCronFunc(ctx, cronID)

	return nil
}

func (s *cronService) AddCronFunc(ctx context.Context, cb *model.Cron) error {
	u, _ := s.UserRepository.FindByID(ctx, cb.UID)
	entityID, err := s.Cron.AddFunc(cb.TimePattern, func() {
		ctxTODO := context.TODO()
		s.TradeService.Trade(ctxTODO, u, cb.Amount, "buy", cb.CryptoName, "fixed")
	})
	if err != nil {
		return err
	}

	if err := s.CronJobManager.SetCronJob(ctx, fmt.Sprintf("fixed:%s", cb.ID), int(entityID)); err != nil {
		return err
	}
	return nil
}

func (s *cronService) RemoveCronFunc(ctx context.Context, cronID string) error {
	entityID, err := s.CronJobManager.GetAndDeleteCronJob(ctx, fmt.Sprintf("fixed:%s", cronID))
	if err != nil {
		return err
	}
	s.Cron.Remove(cron.EntryID(entityID))
	return nil
}
