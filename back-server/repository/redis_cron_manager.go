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

func (r *cronJobManager) SetCronJob(ctx context.Context, key string, value int) error {
	if err := r.Redis.Set(ctx, key, value, 0).Err(); err != nil {
		log.Printf("REPOSITORY: Could not set job entity to redis for cronID/entityID: %s/%v: %s\n", key, value, err)
		return apperrors.NewInternal()
	}
	return nil
}

func (r *cronJobManager) GetCronJob(ctx context.Context, key string) (int, error) {
	rst, err := r.Redis.Get(ctx, key).Int()
	if err != nil {
		log.Printf("REPOSITORY: Could not get job entity from redis for cronID: %s\n", key)
		return rst, apperrors.NewInternal()
	}

	return rst, nil
}

func (r *cronJobManager) DeleteCronJob(ctx context.Context, key string) error {
	result := r.Redis.Del(ctx, key)

	if err := result.Err(); err != nil {
		log.Printf("REPOSITORY: Could not delete cronID to redis for cronID: %s %s\n", key, err)
		return apperrors.NewInternal()
	}

	if result.Val() < 1 {
		log.Printf("REPOSITORY: CronID %s to redis does not exist\n", key)
		return apperrors.NewInternal()
	}

	return nil
}

func (r *cronJobManager) GetAndDeleteCronJob(ctx context.Context, key string) (int, error) {
	entityID, err := r.GetCronJob(ctx, key)
	if err != nil {
		return 0, err
	}
	err = r.DeleteCronJob(ctx, key)
	if err != nil {
		return 0, err
	}
	return entityID, nil
}
