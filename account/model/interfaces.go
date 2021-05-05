package model

import (
	"context"

	"github.com/google/uuid"
)

type UserService interface {
	Get(ctx context.Context, uid uuid.UUID) (*User, error)
}

type UserRepository interface {
	FindById(ctx context.Context, uid uuid.UUID) (*User, error)
}
