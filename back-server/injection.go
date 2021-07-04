package main

import (
	"crypto-auto-invest/handler"
	"crypto-auto-invest/model"
	"crypto-auto-invest/repository"
	"crypto-auto-invest/services"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

func inject(d *dataSources) (*gin.Engine, error) {
	log.Println("Injecting data sources")

	/*
	 * repository layer
	 */
	userRepository := repository.NewUserRepository(d.DB)
	tokenRepository := repository.NewTokenRepository(d.RedisClient)
	walletRepository := repository.NewWalletRepository(d.DB)
	tradeRepository := repository.NewTradeRepository(d.DB)
	autoTradeRepository := repository.NewAutoTradeRepository(d.DB)
	binanceTradeRepository := repository.NewBinanceTradeRepository(d.DB)

	/*
	 * service layer
	 */
	walletService := services.NewWalletService(&services.WAConfig{
		WalletRepository: walletRepository,
	})

	userService := services.NewUserService(&services.USConfig{
		UserRepository: userRepository,
		WalletService:  walletService,
	})

	tradeDelay := os.Getenv("DELAY")
	td, err := strconv.ParseInt(tradeDelay, 0, 64)
	infoWebhook := os.Getenv("INFO_WEBHOOK")
	errorWebhook := os.Getenv("ERROR_WEBHOOK")
	mode := os.Getenv("MODE")
	mockWebhook := os.Getenv("MOCK_WEBHOOK")
	mockTradeService := services.NewMockTradeService(mockWebhook)
	var tradeService model.TradeService
	if mode == "prod" {
		tradeService = services.NewTradeService(&services.TSConifg{
			TradeRepository:  tradeRepository,
			WalletRepository: walletRepository,
			InfoWebhook:      infoWebhook,
			ErrorWebhook:     errorWebhook,
			Delay:            time.Duration(time.Duration(td) * time.Second),
		})
	} else {
		tradeService = mockTradeService
	}

	// init cron job manager
	cron := cron.New()

	cron.Start()

	autoTradeService := services.NewAutoTradeService(&services.ATSConifg{
		WalletRepository:    walletRepository,
		UserRepository:      userRepository,
		AutoTradeRepository: autoTradeRepository,
	})

	// load rsa keys
	privKeyFile := os.Getenv("PRIV_KEY_FILE")
	priv, err := ioutil.ReadFile(privKeyFile)

	if err != nil {
		return nil, fmt.Errorf("could not read private key pem file: %w", err)
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(priv)

	if err != nil {
		return nil, fmt.Errorf("could not parse private key: %w", err)
	}

	pubKeyFile := os.Getenv("PUB_KEY_FILE")
	pub, err := ioutil.ReadFile(pubKeyFile)

	if err != nil {
		return nil, fmt.Errorf("could not read public key pem file: %w", err)
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pub)

	if err != nil {
		return nil, fmt.Errorf("could not parse public key: %w", err)
	}

	// load refresh token secret from env variable
	refreshSecret := os.Getenv("REFRESH_SECRET")

	// load expiration lengths from env variables and parse as int
	idTokenExp := os.Getenv("ID_TOKEN_EXP")
	refreshTokenExp := os.Getenv("REFRESH_TOKEN_EXP")

	idExp, err := strconv.ParseInt(idTokenExp, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse ID_TOKEN_EXP as int: %w", err)
	}

	refreshExp, err := strconv.ParseInt(refreshTokenExp, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse REFRESH_TOKEN_EXP as int: %w", err)
	}

	tokenService := services.NewTokenService(&services.TSConfig{
		TokenRepository:       tokenRepository,
		PrivKey:               privKey,
		PubKey:                pubKey,
		RefreshSecret:         refreshSecret,
		IDExpirationSecs:      idExp,
		RefreshExpirationSecs: refreshExp,
	})

	binanceTradeService := services.NewBinanceTradeService(&services.BTSConfig{
		BinanceTradeRepository: binanceTradeRepository,
		UserRepository:         userRepository,
		Webhook:                infoWebhook,
	})

	// initialize gin.Engine
	router := gin.Default()

	// read in ACCOUNT_API_URL
	baseURL := os.Getenv("ACCOUNT_API_URL")

	handlerTimeout := os.Getenv("HANDLER_TIMEOUT")
	ht, err := strconv.ParseInt(handlerTimeout, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse HANDLER_TIMEOUT as int: %w", err)
	}

	serviceToken := os.Getenv("HEADER_SECRET")

	handler.NewHandler(&handler.Config{
		R:                   router,
		UserService:         userService,
		TokenService:        tokenService,
		WalletService:       walletService,
		TradeService:        tradeService,
		AutoTradeService:    autoTradeService,
		BaseURL:             baseURL,
		TimeoutDuration:     time.Duration(time.Duration(ht) * time.Second),
		ServiceToken:        serviceToken,
		MockTradeService:    mockTradeService,
		BinanceTradeService: binanceTradeService,
	})

	return router, nil
}
