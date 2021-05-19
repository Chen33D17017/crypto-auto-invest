package services

import (
	"context"
	"crypto-auto-invest/model"

	"github.com/robfig/cron/v3"
)

// https://pkg.go.dev/github.com/robfig/cron/v3@v3.0.0?utm_source=gopls
type cronService struct {
	CronRepository model.CronRepository
	TradeService   model.TradeService
	UserRepository model.UserRepository
	CronJobManager model.CronJobManager
	Cron           *cron.Cron
}

type CSConfig struct {
	CronRepository model.CronRepository
	TradeService   model.TradeService
	UserRepository model.UserRepository
	CronJobManager model.CronJobManager
	Cron           *cron.Cron
}

func NewCronService(c *CSConfig) model.CronService {
	return &cronService{
		CronRepository: c.CronRepository,
		TradeService:   c.TradeService,
		UserRepository: c.UserRepository,
		CronJobManager: c.CronJobManager,
		Cron:           c.Cron,
	}
}

// https://github.com/robfig/cron/blob/bc59245fe10efaed9d51b56900192527ed733435/cron.go#L50
// https://pkg.go.dev/github.com/robfig/cron/v3@v3.0.0?utm_source=gopls#EntryID
func (s *cronService) AddCron(ctx context.Context, cb *model.Cron) (*model.Cron, error) {
	err := s.CronRepository.AddCron(ctx, cb)
	if err != nil {
		return nil, err
	}
	cid, err := s.CronRepository.GetCronID(ctx, cb.UID, cb.Type, cb.TimePattern)
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
	err := s.CronRepository.UpdateCron(ctx, cb)
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
		s.TradeService.Trade(ctxTODO, u, cb.Amount, "buy", cb.Type, "fixed")
	})
	if err != nil {
		return err
	}

	// log.Printf("cronID: %v with entityID: %v\n", cb.ID, entityID)

	if err := s.CronJobManager.SetCronJob(ctx, cb.ID, int(entityID)); err != nil {
		return err
	}
	return nil
}

func (s *cronService) RemoveCronFunc(ctx context.Context, cronID string) error {
	entityID, err := s.CronJobManager.GetAndDeleteCronJob(ctx, cronID)
	if err != nil {
		return err
	}
	s.Cron.Remove(cron.EntryID(entityID))
	return nil
}
