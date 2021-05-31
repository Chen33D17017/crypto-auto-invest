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
	infoWebhook := os.Getenv("INFO_WEBHOOK")
	errorWebhook := os.Getenv("ERROR_WEBHOOK")
	mode := os.Getenv("MODE")
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
		tradeService = services.NewMockTradeService()
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

	tradeRateApi := os.Getenv("TRADE_RATE_API")
	maxRate := os.Getenv("MAX_RATE")
	autoTradeTimePattern := os.Getenv("AUTO_TRADE_TIME")
	rate, err := strconv.ParseFloat(maxRate, 64)
	if err != nil {
		log.Fatalf("Fail to load max rate on auto trade")
	}
	autoTradeService := services.NewAutoTradeService(&services.ATSConifg{
		TradeService:        tradeService,
		WalletRepository:    walletRepository,
		UserRepository:      userRepository,
		AutoTradeRepository: autoTradeRepository,
		CronJobManager:      cronJobManager,
		Cron:                cron,
		TimePattern:         autoTradeTimePattern,
		TradeRateApi:        tradeRateApi,
		MaxRate:             rate,
	})

	settings, err := autoTradeRepository.GetAllAutoTrade()
	if err != nil {
		log.Fatalf("Fail to load auto trade setting: %s\n", err.Error())
	}
	log.Printf("Load auto buy setting number: %v\n", len(*settings))
	for _, setting := range *settings {
		ctx := context.TODO()
		err = autoTradeService.AddCronFunc(ctx, setting)
		if err != nil {
			log.Fatalf("Fail to load auto trade setting: %s\n", err.Error())
		}
	}

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
	})

	return router, nil
}
