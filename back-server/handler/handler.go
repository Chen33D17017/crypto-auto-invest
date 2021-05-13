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
}

type Config struct {
	R               *gin.Engine
	UserService     model.UserService
	TokenService    model.TokenService
	WalletService   model.WalletService
	BaseURL         string
	TimeoutDuration time.Duration
}

func NewHandler(c *Config) {
	h := &Handler{
		UserService:   c.UserService,
		TokenService:  c.TokenService,
		WalletService: c.WalletService,
	}
	g := c.R.Group(c.BaseURL)

	if gin.Mode() != gin.TestMode {
		g.Use(middleware.Timeout(c.TimeoutDuration, apperrors.NewServiceUnavailable()))
		g.GET("/me", middleware.AuthUser(h.TokenService), h.Me)
		g.POST("/signout", middleware.AuthUser(h.TokenService), h.Signout)
		g.PUT("/details", middleware.AuthUser(h.TokenService), h.UserDetails)
		g.PATCH("/details", middleware.AuthUser(h.TokenService), h.PatchUser)
		g.GET("/wallet", middleware.AuthUser(h.TokenService), h.GetWallet)
		g.GET("/wallets", middleware.AuthUser(h.TokenService), h.GetWallets)
		g.POST("/add_wallet", middleware.AuthUser(h.TokenService), h.AddWallet)
		g.POST("/charge", middleware.AuthUser(h.TokenService), h.Charge)
	} else {
		g.GET("/me", h.Me)
		g.POST("/signout", h.Signout)
		g.PUT("/details", h.UserDetails)
		g.PATCH("/details", h.PatchUser)
		g.GET("/wallet", h.GetWallets)
		g.GET("/wallets", h.GetWallets)
		g.POST("/add_wallet", h.AddWallet)
		g.POST("/charge", h.Charge)
	}
	g.POST("/signup", h.Signup)
	g.POST("/signin", h.Signin)
	g.POST("/tokens", h.Tokens)
	g.POST("/image", h.Image)
	g.DELETE("/image", h.DeleteImage)
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
