package services

import (
	"context"
	"crypto-auto-invest/model"
)

type cronService struct {
	CronRepository model.CronRepository
}

type CSConfig struct {
	CronRepository model.CronRepository
}

func NewCronService(c *CSConfig) model.CronService {
	return &cronService{
		CronRepository: c.CronRepository,
	}
}

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
	return nil
}

func (s *cronService) DeleteCron(ctx context.Context, uid string, cronID string) error {
	err := s.CronRepository.DeleteCron(ctx, uid, cronID)
	if err != nil {
		return err
	}
	return nil
}
