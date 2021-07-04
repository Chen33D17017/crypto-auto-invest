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
	AddWallet(ctx context.Context, uid string, cryptoName string, strategyID int) (*Wallet, error)
	GetUserWallet(ctx context.Context, uid string, cryptoName string, strategyID int) (*Wallet, error)
	GetWallets(ctx context.Context, uid string, strategyID int) (*[]Wallet, error)
	ChangeMoney(ctx context.Context, uid string, cryptoName string, amount float64, strategyID int) (*Wallet, error)
	GetChargeLogs(ctx context.Context, uid string, cryptoName string, strategyID int) (*[]ChargeLog, error)
}

type TradeService interface {
	MarketTrade(ctx context.Context, u *User, amount float64, action, cryptoName string, strategy int) (bm.Order, error)
	LimitTrade(ctx context.Context, u *User, amount float64, action, cryptoName string, price float64) (bm.Order, error)
	SaveOrder(ctx context.Context, u *User, orderID string, cryptoName string, strategy int, info bool) error
	CalIncomeRate(ctx context.Context, uid string, cryptoName string, strategyID int) (*Income, error)
}

type AutoTradeService interface {
	AddAutoTrade(ctx context.Context, uid, cryptoName string, strategyID int) error
	DeleteAutoTrade(ctx context.Context, uid, cryptoName string, strategyID int) error
	GetAutoTrades(ctx context.Context, uid string) (*[]AutoTrade, error)
	GetAutoTradesFromStrategy(ctx context.Context, cryptoName string, strategyID int) ([]AutoTradeRes, error)
	GetAllAutoTrades(ctx context.Context) (*[]AutoTrade, error)
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
	AddWallet(ctx context.Context, uid string, cid int, strategyID int) error
	GetWalletByID(ctx context.Context, wid string) (*Wallet, error)
	GetWellet(ctx context.Context, uid string, cryptoName string, strategyID int) (*Wallet, error)
	GetWallets(ctx context.Context, uid string, strategyID int) (*[]Wallet, error)
	UpdateAmount(ctx context.Context, wid string, amount float64) error
	GetCurrencyID(ctx context.Context, cryptoName string) (int, error)
	AddChargeLog(ctx context.Context, uid string, cid int, strategyID int, amount float64) error
	GetChargeLogs(ctx context.Context, uid string, cryptoName string, strategyID int) (*[]ChargeLog, error)
}

type TradeRepository interface {
	SaveOrder(ctx context.Context, t *Order) error
	GetOrderLogs(ctx context.Context, uid, cryptoName string, strategyID int) (*[]Order, error)
}

type AutoTradeRepository interface {
	AddAutoTrade(ctx context.Context, uid string, int, strategyID int) error
	DeleteAutoTrade(ctx context.Context, uid string, int, strategyID int) error
	GetAutoTrades(ctx context.Context, uid string) (*[]AutoTrade, error)
	GetAutoTradeFromStrategy(ctx context.Context, cryptoName string, strategyID int) (*[]AutoTrade, error)
	GetAllAutoTrades(ctx context.Context) (*[]AutoTrade, error)
}

type CronJobManager interface {
	SetCronJob(ctx context.Context, cronID string, entityID int) error
	GetCronJob(ctx context.Context, cronID string) (int, error)
	DeleteCronJob(ctx context.Context, cronID string) error
	GetAndDeleteCronJob(ctx context.Context, cronID string) (int, error)
}

type BinanceTradeService interface {
	SaveOrder(ctx context.Context, uid string, symbol string, action string, avgCost float64, qty float64) (BinanceOrder, error)
	GetOrders(ctx context.Context, uid string, symbol string) (*[]BinanceOrder, error)
}

type BinanceTradeRepository interface {
	SaveOrder(ctx context.Context, order *BinanceOrder) error
	GetOrders(ctx context.Context, uid string, symbol string) (*[]BinanceOrder, error)
}
