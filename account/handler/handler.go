package handler

import (
	"account-tutorial/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	UserService  model.UserService
	TokenService model.TokenService
}

type Config struct {
	R            *gin.Engine
	UserService  model.UserService
	TokenService model.TokenService
	BaseURL      string
}

func NewHandler(c *Config) {
	h := &Handler{
		UserService:  c.UserService,
		TokenService: c.TokenService,
	}
	g := c.R.Group(c.BaseURL)
	g.GET("/me", h.Me)
	g.POST("/signup", h.Signup)
	g.POST("/signin", h.Signin)
	g.POST("/tokens", h.Tokens)
	g.POST("/image", h.Image)
	g.DELETE("/image", h.DeleteImage)
	g.PUT("details", h.Details)
}

func (h *Handler) Signin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Signin",
	})
}

func (h *Handler) Tokens(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "it's token",
	})
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

func (h *Handler) Details(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Details",
	})
}
