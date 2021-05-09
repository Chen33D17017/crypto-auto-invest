package mocks

import (
	"context"
	"crypto-auto-invest/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockTokenService struct {
	mock.Mock
}

func (m *MockTokenService) NewPairFromUser(ctx context.Context, u *model.User, prevTokenID string) (*model.TokenPair, error) {
	ret := m.Called(ctx, u, prevTokenID)

	var r0 *model.TokenPair
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*model.TokenPair)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

func (m *MockTokenService) ValidateIDToken(tokenString string) (*model.User, error) {
	ret := m.Called(tokenString)

	// first value passed to "Return"
	var r0 *model.User
	if ret.Get(0) != nil {
		// we can just return this if we know we won't be passing function to "Return"
		r0 = ret.Get(0).(*model.User)
	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

func (m *MockTokenService) ValidateRefreshToken(refreshTokenString string) (*model.RefreshToken, error) {
	ret := m.Called(refreshTokenString)

	var r0 *model.RefreshToken
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*model.RefreshToken)
	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

func (m *MockTokenService) Signout(ctx context.Context, uid uuid.UUID) error {
	ret := m.Called(ctx, uid)
	var r0 error

	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}
