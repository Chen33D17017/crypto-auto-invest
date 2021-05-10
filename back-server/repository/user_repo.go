package repository

import (
	"context"
	"crypto-auto-invest/model"
	"crypto-auto-invest/model/apperrors"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type userRepository struct {
	DB *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) model.UserRepository {
	return &userRepository{
		DB: db,
	}
}

func (r *userRepository) Create(ctx context.Context, u *model.User) error {
	query := "INSERT INTO users (email, password) VALUES (?, ?)"

	// Deal with mysql error:
	// https://stackoverflow.com/questions/47009068/how-to-get-the-mysql-error-type-in-golang
	// https://github.com/go-sql-driver/mysql/blob/a8b7ed4454a6a4f98f85d3ad558cd6d97cec6959/errors.go#L19
	// https://dev.mysql.com/doc/mysql-errors/8.0/en/server-error-reference.html
	// https://github.com/VividCortex/mysqlerr

	_, err := r.DB.ExecContext(ctx, query, u.Email, u.Password)
	if err != nil {
		log.Println(err)
		if err, ok := err.(*mysql.MySQLError); ok && err.Number == 1062 {
			log.Printf("Could not create a user with email: %v. Reason: %v\n", u.Email, err.Message)
			return apperrors.NewConflict("email", u.Email)
		}

		log.Printf("Could not create a user with email: %v. Reason: %v\n", u.Email, err)
		return apperrors.NewInternal()
	}
	return nil
}

func (r *userRepository) FindByID(ctx context.Context, uid string) (*model.User, error) {
	user := &model.User{}

	query := "SELECT * FROM users WHERE uid=?"

	// we need to actually check errors as it could be something other than not found
	if err := r.DB.Get(user, query, uid); err != nil {
		return user, apperrors.NewNotFound("uid", uid)
	}

	return user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}

	query := "SELECT * FROM users WHERE email=?"

	if err := r.DB.GetContext(ctx, user, query, email); err != nil {
		log.Printf("Unable to get user with email address: %v. Err: %v\n", email, err)
		return user, apperrors.NewNotFound("email", email)
	}

	return user, nil
}

func (r *userRepository) Update(ctx context.Context, u *model.User) error {
	query := `
	UPDATE users 
	SET 
		name=:name, 
		email=:email, 
		website=:website
	WHERE
		uid=:uid;
	`

	nstmt, err := r.DB.PrepareNamedContext(ctx, query)

	if err != nil {
		log.Printf("Unable to prepare user update query: %v\n", err)
		return apperrors.NewInternal()
	}

	if err := nstmt.GetContext(ctx, u, u); err != nil {
		log.Printf("Failed to update details for user: %v\n", u)
		return apperrors.NewInternal()
	}

	return nil
}
