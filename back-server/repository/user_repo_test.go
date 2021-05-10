package repository

import (
	"context"
	"crypto-auto-invest/model"
	"crypto-auto-invest/model/apperrors"
	"fmt"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

// https://github.com/jmoiron/sqlx/issues/204

func TestCreate(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	defer mockDB.Close()
	assert.NoError(t, err)

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	ur := NewUserRepository(sqlxDB)
	ctx := context.TODO()

	t.Run("Success to create", func(t *testing.T) {
		uid, _ := uuid.NewRandom()
		mu := &model.User{
			UID:      uid,
			Email:    "bob@bob.com",
			Name:     "Bobby Bobson",
			Password: "password",
		}

		mock.ExpectExec("INSERT INTO users").WithArgs(mu.Email, mu.Password).WillReturnResult(sqlmock.NewResult(0, 1))
		err = ur.Create(ctx, mu)
		assert.NoError(t, err)
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	t.Run("Duplicate email", func(t *testing.T) {
		duplicateErr := mysql.MySQLError{
			Number:  1062,
			Message: "Duplicate",
		}
		uid, _ := uuid.NewRandom()
		mu := &model.User{
			UID:      uid,
			Email:    "bob@bob.com",
			Name:     "Bobby Bobson",
			Password: "password",
		}
		mock.ExpectExec("^INSERT (.+)").WillReturnError(&duplicateErr)
		err = ur.Create(ctx, mu)
		log.Println(err)
		apperror, ok := err.(*apperrors.Error)
		assert.True(t, ok)

		assert.Equal(t, apperrors.Conflict, apperror.Type)
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	t.Run("Unexpected Error", func(t *testing.T) {
		uid, _ := uuid.NewRandom()
		mu := &model.User{
			UID:      uid,
			Email:    "bob@bob.com",
			Name:     "Bobby Bobson",
			Password: "password",
		}
		mock.ExpectExec("^INSERT (.+)").WillReturnError(fmt.Errorf("Unexpected Error"))
		err = ur.Create(ctx, mu)
		log.Println(err)
		apperror, ok := err.(*apperrors.Error)
		assert.True(t, ok)

		assert.Equal(t, apperrors.Internal, apperror.Type)
	})
}
