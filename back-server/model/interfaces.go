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
	Trade(ctx context.Context, u *User, amount float64, action, cryptoName string, strategy int) (bm.Order, error)
	SaveOrder(ctx context.Context, u *User, orderID string, cryptoName string, strategy int) error
	SendTradeRst(msg string, level string) error
}

type CronService interface {
	AddCron(ctx context.Context, cb *Cron) (*Cron, error)
	GetCron(ctx context.Context, uid, cronID string) (*Cron, error)
	GetCrons(ctx context.Context, uid string) (*[]Cron, error)
	UpdateCron(ctx context.Context, cb *Cron) error
	DeleteCron(ctx context.Context, uid string, cronID string) error
	AddCronFunc(ctx context.Context, cb *Cron) error
	RemoveCronFunc(ctx context.Context, cronID string) error
}

type AutoTradeService interface {
	AddAutoTrade(ctx context.Context, uid, currencyName string) error
	DeleteAutoTrade(ctx context.Context, uid, currencyName string) error
	GetAutoTrades(ctx context.Context, uid string) (*[]AutoTrade, error)
	GetTradeRate(reqBody TradeRateReq) (TradeRateRes, error)
	AutoTrade(uid string, currencyName string) error
	AddCronFunc(ctx context.Context, setting AutoTrade) error
	RemoveCronFunc(ctx context.Context, autoTradeID string) error
	TradeWithRate(ctx context.Context, currencyName string, action string, rate float64, strategy int) error
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
	AddWallet(ctx context.Context, uid string, cid int) error
	GetWalletByID(ctx context.Context, wid string) (*Wallet, error)
	GetWellet(ctx context.Context, uid string, currencyName string) (*Wallet, error)
	GetWallets(ctx context.Context, uid string) (*[]Wallet, error)
	UpdateAmount(ctx context.Context, wid string, amount float64) error
	GetCurrencyID(ctx context.Context, currencyName string) (int, error)
}

type TradeRepository interface {
	SaveOrder(ctx context.Context, t *Order) error
}

type CronRepository interface {
	AddCron(ctx context.Context, cb *Cron, currencyID int) error
	GetCron(ctx context.Context, uid string, cronID string) (*Cron, error)
	GetCrons(ctx context.Context, uid string) (*[]Cron, error)
	UpdateCron(ctx context.Context, cb *Cron, currencyID int) error
	DeleteCron(ctx context.Context, userID string, cronID string) error
	GetCronID(ctx context.Context, uid, cryptoName, timePattern string) (string, error)
	GetAllCrons() (*[]Cron, error)
}

type AutoTradeRepository interface {
	AddAutoTrade(ctx context.Context, uid string, type_id int) error
	DeleteAutoTrade(ctx context.Context, uid string, type_id int) error
	GetAutoTrades(ctx context.Context, uid string) (*[]AutoTrade, error)
	GetAutoTrade(ctx context.Context, uid, currencyName string) (*AutoTrade, error)
	GetAllAutoTrade() (*[]AutoTrade, error)
	GetAutoTradeUser(ctx context.Context, currencyName string) (*[]string, error)
}

type CronJobManager interface {
	SetCronJob(ctx context.Context, cronID string, entityID int) error
	GetCronJob(ctx context.Context, cronID string) (int, error)
	DeleteCronJob(ctx context.Context, cronID string) error
	GetAndDeleteCronJob(ctx context.Context, cronID string) (int, error)
}
