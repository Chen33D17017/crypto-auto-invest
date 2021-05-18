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
	UserService   model.UserService
	TokenService  model.TokenService
	WalletService model.WalletService
	TradeService  model.TradeService
	Delay         time.Duration
}

type Config struct {
	R               *gin.Engine
	UserService     model.UserService
	TokenService    model.TokenService
	WalletService   model.WalletService
	TradeService    model.TradeService
	BaseURL         string
	TimeoutDuration time.Duration
	Delay           time.Duration
}

func NewHandler(c *Config) {
	h := &Handler{
		UserService:   c.UserService,
		TokenService:  c.TokenService,
		WalletService: c.WalletService,
		TradeService:  c.TradeService,
		Delay:         c.Delay,
	}
	g_user := c.R.Group(c.BaseURL)
	g_price := c.R.Group("/api/bitbank")
	g_crypto := c.R.Group("/api/crypto")

	if gin.Mode() != gin.TestMode {
		g_user.Use(middleware.Timeout(c.TimeoutDuration, apperrors.NewServiceUnavailable()))

		g_user.GET("/me", middleware.AuthUser(h.TokenService), h.Me)
		g_user.POST("/signout", middleware.AuthUser(h.TokenService), h.Signout)
		g_user.PUT("/details", middleware.AuthUser(h.TokenService), h.UserDetails)
		g_user.PATCH("/details", middleware.AuthUser(h.TokenService), h.PatchUser)
		g_user.GET("/wallet", middleware.AuthUser(h.TokenService), h.GetWallet)
		g_user.GET("/wallets", middleware.AuthUser(h.TokenService), h.GetWallets)
		g_user.POST("/add_wallet", middleware.AuthUser(h.TokenService), h.AddWallet)
		g_user.POST("/charge", middleware.AuthUser(h.TokenService), h.Charge)

		g_price.GET("/assets", middleware.AuthUser(h.TokenService), h.GetAssets)
		g_price.GET("/trade", middleware.AuthUser(h.TokenService), h.GetTrade)
		g_price.GET("/historys", middleware.AuthUser(h.TokenService), h.GetHistory)

		g_crypto.POST("/trade", middleware.AuthUser(h.TokenService), h.Trade)
		g_crypto.POST("/order", middleware.AuthUser(h.TokenService), h.SaveOrder)
	} else {
		g_user.GET("/me", h.Me)
		g_user.POST("/signout", h.Signout)
		g_user.PUT("/details", h.UserDetails)
		g_user.PATCH("/details", h.PatchUser)
		g_user.GET("/wallet", h.GetWallets)
		g_user.GET("/wallets", h.GetWallets)
		g_user.POST("/add_wallet", h.AddWallet)
		g_user.POST("/charge", h.Charge)
	}
	g_user.POST("/signup", h.Signup)
	g_user.POST("/signin", h.Signin)
	g_user.POST("/tokens", h.Tokens)
	g_user.POST("/image", h.Image)
	g_user.DELETE("/image", h.DeleteImage)

	g_price.GET("/price", h.GetPrice)
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
