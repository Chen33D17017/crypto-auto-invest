package handler

import (
	"crypto-auto-invest/handler/middleware"
	"crypto-auto-invest/model"
	"crypto-auto-invest/model/apperrors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	UserService      model.UserService
	TokenService     model.TokenService
	WalletService    model.WalletService
	TradeService     model.TradeService
	CronService      model.CronService
	AutoTradeService model.AutoTradeService
	MockTradeService model.TradeService
}

type Config struct {
	R                *gin.Engine
	UserService      model.UserService
	TokenService     model.TokenService
	WalletService    model.WalletService
	TradeService     model.TradeService
	CronService      model.CronService
	AutoTradeService model.AutoTradeService
	BaseURL          string
	TimeoutDuration  time.Duration
	ServiceToken     string
	MockTradeService model.TradeService
}

func NewHandler(c *Config) {
	h := &Handler{
		UserService:      c.UserService,
		TokenService:     c.TokenService,
		WalletService:    c.WalletService,
		TradeService:     c.TradeService,
		CronService:      c.CronService,
		AutoTradeService: c.AutoTradeService,
		MockTradeService: c.MockTradeService,
	}
	g_user := c.R.Group(c.BaseURL)
	g_price := c.R.Group("/api/bitbank")
	g_crypto := c.R.Group("/api/crypto")
	g_admin := c.R.Group("/api/admin")
	g_mock := c.R.Group("/api/mock")

	if gin.Mode() != gin.TestMode {
		g_user.Use(middleware.Timeout(c.TimeoutDuration, apperrors.NewServiceUnavailable()))

		g_user.GET("/me", middleware.AuthUser(h.TokenService), h.Me)
		g_user.POST("/signout", middleware.AuthUser(h.TokenService), h.Signout)
		g_user.PUT("/details", middleware.AuthUser(h.TokenService), h.UserDetails)
		g_user.PATCH("/details", middleware.AuthUser(h.TokenService), h.PatchUser)
		g_user.GET("/wallet", middleware.AuthUser(h.TokenService), h.GetWallet)
		g_user.GET("/wallets", middleware.AuthUser(h.TokenService), h.GetWallets)
		g_user.POST("/charge", middleware.AuthUser(h.TokenService), h.Charge)
		g_user.GET("/charge", middleware.AuthUser(h.TokenService), h.GetChargeLogs)
		g_user.POST("/cron", middleware.AuthUser(h.TokenService), h.AddCron)
		g_user.GET("/cron", middleware.AuthUser(h.TokenService), h.GetCron)
		g_user.GET("/crons", middleware.AuthUser(h.TokenService), h.GetCrons)
		g_user.PUT("/cron", middleware.AuthUser(h.TokenService), h.UpdateCron)
		g_user.DELETE("/cron", middleware.AuthUser(h.TokenService), h.DeleteCron)

		g_price.GET("/assets", middleware.AuthUser(h.TokenService), h.GetAssets)
		g_price.GET("/trade", middleware.AuthUser(h.TokenService), h.GetTrade)
		g_price.GET("/historys", middleware.AuthUser(h.TokenService), h.GetHistory)

		g_crypto.POST("/auto_trade", middleware.AuthUser(h.TokenService), h.AddAutoTrade)
		g_crypto.DELETE("/auto_trade", middleware.AuthUser(h.TokenService), h.DeleteAutoTrade)
		g_crypto.GET("/auto_trades", middleware.AuthUser(h.TokenService), h.GetAutoTrades)
		g_crypto.GET("/income", middleware.AuthUser(h.TokenService), h.GetIncomeRate)
	} else {
		g_user.GET("/me", h.Me)
		g_user.POST("/signout", h.Signout)
		g_user.PUT("/details", h.UserDetails)
		g_user.PATCH("/details", h.PatchUser)
		g_user.GET("/wallet", h.GetWallets)
		g_user.GET("/wallets", h.GetWallets)
		g_user.POST("/charge", h.Charge)
	}
	g_user.POST("/signup", h.Signup)
	g_user.POST("/signin", h.Signin)
	g_user.POST("/tokens", h.Tokens)
	g_user.POST("/image", h.Image)
	g_user.DELETE("/image", h.DeleteImage)

	g_admin.POST("/trade", middleware.AuthService(c.ServiceToken), h.Trade)
	g_mock.POST("/trade", middleware.AuthService(c.ServiceToken), h.MockTrade)
	g_admin.GET("/auto_trade", middleware.AuthService(c.ServiceToken), h.GetAutoTradeInfo)
	g_admin.GET("/all_auto_trade", middleware.AuthService(c.ServiceToken), h.GetAllAutoTrades)
	g_admin.POST("/order", middleware.AuthService(c.ServiceToken), h.SaveOrder)
	g_mock.POST("/order", middleware.AuthService(c.ServiceToken), h.SaveOrder)
}

func (h *Handler) Image(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Image",
	})
}

func (h *Handler) DeleteImage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Delete Image",
	})
}
