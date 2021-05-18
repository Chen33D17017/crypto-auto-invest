package model

import (
	"context"
	"time"

	bm "crypto-auto-invest/bitbank/model"
)

type UserService interface {
	Get(ctx context.Context, uid string) (*User, error)
	Signup(ctx context.Context, u *User) error
	Signin(ctx context.Context, u *User) error
	UpdateDetails(ctx context.Context, u *User) error
	PatchDetails(ctx context.Context, u *User) (*User, error)
}

type TokenService interface {
	NewPairFromUser(ctx context.Context, u *User, prevTokenID string) (*TokenPair, error)
	ValidateIDToken(tokenString string) (*User, error)
	ValidateRefreshToken(refreshToken string) (*RefreshToken, error)
	Signout(ctx context.Context, uid string) error
}

type WalletService interface {
	AddWallet(ctx context.Context, uid string, currencyName string) (*Wallet, error)
	GetUserWallet(ctx context.Context, uid string, currencyName string) (*Wallet, error)
	GetWallets(ctx context.Context, uid string) (*[]Wallet, error)
	ChangeMoney(ctx context.Context, uid string, currencyName string, amount float64) (*Wallet, error)
}

type TradeService interface {
	Trade(ctx context.Context, u *User, amount float64, getDelay time.Duration, side, assetType, orderType string) (bm.Order, error)
	SaveOrder(ctx context.Context, u *User, orderID string, assetType, orderType string) error
}

type UserRepository interface {
	FindByID(ctx context.Context, uid string) (*User, error)
	Create(ctx context.Context, u *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, u *User) error
	Patch(ctx context.Context, u *User) error
}

type TokenRepository interface {
	SetRefreshToken(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) error
	DeleteRefreshToken(ctx context.Context, userID string, prevTokenID string) error
	DeleteUserRefreshTokens(ctx context.Context, userID string) error
}

type WalletRepository interface {
	AddWallet(ctx context.Context, uid string, currencyName string) error
	GetWalletByID(ctx context.Context, wid string) (*Wallet, error)
	GetWellet(ctx context.Context, uid string, currencyType string) (*Wallet, error)
	GetWallets(ctx context.Context, uid string) (*[]Wallet, error)
	UpdateAmount(ctx context.Context, wid string, amount float64) error
}

type TradeRepository interface {
	SaveOrder(ctx context.Context, t *Order) error
}
