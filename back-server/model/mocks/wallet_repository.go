package mocks

import (
	"context"
	"crypto-auto-invest/model"

	"github.com/stretchr/testify/mock"
)

type mockWalletRepository struct {
	mock.Mock
}

func (m *mockWalletRepository) AddWallet(ctx context.Context, uid string, currencyName string) error {
	ret := m.Called(ctx, uid, currencyName)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}
func (m *mockWalletRepository) GetWalletByID(ctx context.Context, wid string) (*model.Wallet, error) {
	ret := m.Called(ctx, wid)

	var r0 *model.Wallet
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*model.Wallet)
	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

func (m *mockWalletRepository) GetWellet(ctx context.Context, uid string, currencyName string) (*model.Wallet, error) {
	ret := m.Called(ctx, uid, currencyName)

	var r0 *model.Wallet
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*model.Wallet)
	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

func (m *mockWalletRepository) GetWallets(ctx context.Context, uid string) (*[]model.Wallet, error) {
	ret := m.Called(ctx, uid)

	var r0 *[]model.Wallet
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*[]model.Wallet)
	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

func (m *mockWalletRepository) UpdateAmount(ctx context.Context, wid string, amount float64) error {
	ret := m.Called(ctx, wid)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}
