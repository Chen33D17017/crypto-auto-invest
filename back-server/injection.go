package main

import (
	"context"
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
	cronJobManager := repository.NewCronJobManager(d.RedisClient)
	walletRepository := repository.NewWalletRepository(d.DB)
	tradeRepository := repository.NewTradeRepository(d.DB)
	cronRepository := repository.NewCronRepository(d.DB)
	autoTradeRepository := repository.NewAutoTradeRepository(d.DB)

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
	//infoWebhook := os.Getenv("INFO_WEBHOOK")
	//errorWebhook := os.Getenv("ERROR_WEBHOOK")
	testWebhook := os.Getenv("TEST_WEBHOOK")
	mode := os.Getenv("MODE")
	mockTradeService := services.NewMockTradeService()
	var tradeService model.TradeService
	if mode == "dev" {
		tradeService = services.NewTradeService(&services.TSConifg{
			TradeRepository:  tradeRepository,
			WalletRepository: walletRepository,
			InfoWebhook:      testWebhook,
			ErrorWebhook:     testWebhook,
			Delay:            time.Duration(time.Duration(td) * time.Second),
		})
	} else {
		tradeService = mockTradeService
	}

	// init cron job manager
	cron := cron.New()

	cron.Start()
	cronService := services.NewCronService(&services.CSConfig{
		CronRepository:   cronRepository,
		UserRepository:   userRepository,
		WalletRepository: walletRepository,
		TradeService:     tradeService,
		CronJobManager:   cronJobManager,
		Cron:             cron,
	})

	jobs, err := cronRepository.GetAllCrons()

	if err != nil {
		log.Fatalf("Fail to init cron job manager: %s\n", err)
	}

	log.Printf("Setting %v cron jobs for system\n", len(*jobs))
	for _, job := range *jobs {
		ctx := context.TODO()
		cronService.AddCronFunc(ctx, &job)
	}

	if err != nil {
		log.Fatalf("Fail to load max rate on auto trade")
	}
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
	mockWebhook := os.Getenv("MOCK_WEBHOOK")

	handler.NewHandler(&handler.Config{
		R:                router,
		UserService:      userService,
		TokenService:     tokenService,
		WalletService:    walletService,
		TradeService:     tradeService,
		CronService:      cronService,
		AutoTradeService: autoTradeService,
		BaseURL:          baseURL,
		TimeoutDuration:  time.Duration(time.Duration(ht) * time.Second),
		ServiceToken:     serviceToken,
		MockWebhook:      mockWebhook,
		MockTradeService: mockTradeService,
	})

	return router, nil
}
