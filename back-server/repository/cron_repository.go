package repository

import (
	"context"
	"crypto-auto-invest/model"
	"crypto-auto-invest/model/apperrors"
	"log"

	"github.com/jmoiron/sqlx"
)

const (
	queryAddCron    = "INSERT INTO crons(uid, type, amount, time_pattern) VALUES(?, ?, ?, ?);"
	queryGetCron    = "SELECT * FROM crons WHERE id=? and uid=?;"
	queryGetCrons   = `SELECT * FROM crons WHERE uid=?;`
	queryUpdateCron = `UPDATE crons SET amount=?, type=?, time_pattern=? WHERE id=? and uid=?;`
	queryDeleteCron = `DELETE FROM crons WHERE id=? and uid=?;`
	queryGetCronID  = `SELECT * FROM crons WHERE uid=? and type=? and time_pattern=?`
)

type cronRepository struct {
	DB *sqlx.DB
}

func NewCronRepository(db *sqlx.DB) model.CronRepository {
	return &cronRepository{
		DB: db,
	}
}

func (r *cronRepository) AddCron(ctx context.Context, cb *model.Cron) error {
	stmt, err := r.DB.PrepareContext(ctx, queryAddCron)
	if err != nil {
		log.Printf("REPOSITORY: Unable to Add Cron: %v\n", err)
		return apperrors.NewInternal()
	}

	if _, err := stmt.ExecContext(ctx, cb.UID, cb.Type, cb.Amount, cb.TimePattern); err != nil {
		log.Printf("REPOSITORY: Failed to update details for user: %v err: %s\n", cb.UID, err.Error())
		return apperrors.NewInternal()
	}
	return nil
}

func (r *cronRepository) GetCron(ctx context.Context, uid string, cronID string) (*model.Cron, error) {
	rst := &model.Cron{}
	err := r.DB.GetContext(ctx, rst, queryGetCron, cronID, uid)
	if err != nil {
		log.Printf("REPOSITORY: Unable to get cron by (uid, cronID): (%v, %v) err: %s", uid, cronID, err)
		return rst, apperrors.NewInternal()
	}
	return rst, nil
}

func (r *cronRepository) GetCrons(ctx context.Context, uid string) (*[]model.Cron, error) {
	rst := &[]model.Cron{}
	err := r.DB.SelectContext(ctx, rst, queryGetCrons, uid)
	if err != nil {
		log.Printf("REPOSITORY: Unable to get crons by (uid): %v err: %s", uid, err)
		return rst, apperrors.NewInternal()
	}
	return rst, nil
}

func (r *cronRepository) UpdateCron(ctx context.Context, cb *model.Cron) error {
	stmt, err := r.DB.PrepareContext(ctx, queryUpdateCron)
	if err != nil {
		log.Printf("REPOSITORY: Unable to prepare update query: %v\n", err)
		return apperrors.NewInternal()
	}

	if _, err := stmt.ExecContext(ctx, cb.Amount, cb.Type, cb.TimePattern, cb.ID, cb.UID); err != nil {
		log.Printf("REPOSITORY: Failed to update wallet for wid: %v err: %s\n", cb.ID, err.Error())
		return apperrors.NewInternal()
	}
	return nil
}

func (r *cronRepository) DeleteCron(ctx context.Context, userID string, cronID string) error {
	stmt, err := r.DB.PrepareContext(ctx, queryDeleteCron)
	if err != nil {
		log.Printf("REPOSITORY: Unable to prepare update query: %v\n", err)
		return apperrors.NewInternal()
	}

	if _, err := stmt.ExecContext(ctx, cronID, userID); err != nil {
		log.Printf("REPOSITORY: Failed to delete cron on id: %v err: %s\n", cronID, err.Error())
		return apperrors.NewInternal()
	}
	return nil
}

func (r *cronRepository) GetCronID(ctx context.Context, uid, cryptoType, timePattern string) (string, error) {
	rst := &model.Cron{}
	err := r.DB.GetContext(ctx, rst, queryGetCronID, uid, cryptoType, timePattern)
	if err != nil {
		log.Printf("REPOSITORY: Unable to get cron by (uid, type, pattern): (%s, %s, %s) err: %s", uid, cryptoType, timePattern, err)
		return rst.ID, apperrors.NewInternal()
	}
	return rst.ID, nil
}
