package repository

import (
	"context"
	"crypto-auto-invest/model"
	"crypto-auto-invest/model/apperrors"
	"log"

	"github.com/go-redis/redis/v8"
)

//https://edward-cernera.medium.com/intro-to-redis-with-docker-compose-8d53962336cb

type cronJobManager struct {
	Redis *redis.Client
}

func NewCronJobManager(redisClient *redis.Client) model.CronJobManager {
	return &cronJobManager{
		Redis: redisClient,
	}
}

// https://medium.com/easyread/unit-test-redis-in-golang-c22b5589ea37In

func (r *cronJobManager) SetCronJob(ctx context.Context, cronID string, entityID int) error {
	if err := r.Redis.Set(ctx, cronID, entityID, 0).Err(); err != nil {
		log.Printf("REPOSITORY: Could not set job entity to redis for cronID/entityID: %s/%v: %s\n", cronID, entityID, err)
		return apperrors.NewInternal()
	}
	return nil
}

func (r *cronJobManager) GetCronJob(ctx context.Context, cronID string) (int, error) {
	rst, err := r.Redis.Get(ctx, cronID).Int()
	if err != nil {
		log.Printf("REPOSITORY: Could not get job entity from redis for cronID: %s\n", cronID)
		return rst, apperrors.NewInternal()
	}

	return rst, nil
}

func (r *cronJobManager) DeleteCronJob(ctx context.Context, cronID string) error {
	result := r.Redis.Del(ctx, cronID)

	if err := result.Err(); err != nil {
		log.Printf("REPOSITORY: Could not delete cronID to redis for cronID: %s %s\n", cronID, err)
		return apperrors.NewInternal()
	}

	if result.Val() < 1 {
		log.Printf("REPOSITORY: CronID %s to redis does not exist\n", cronID)
		return apperrors.NewInternal()
	}

	return nil
}

func (r *cronJobManager) GetAndDeleteCronJob(ctx context.Context, cronID string) (int, error) {
	entityID, err := r.GetCronJob(ctx, cronID)
	if err != nil {
		return 0, err
	}
	err = r.DeleteCronJob(ctx, cronID)
	if err != nil {
		return 0, err
	}
	return entityID, nil
}
