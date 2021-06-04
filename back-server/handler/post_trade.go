package handler

import (
	"crypto-auto-invest/model/apperrors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type tradeReq struct {
	UID        string  `json:"uid" binding:"required"`
	Amount     float64 `json:"amount" binding:"required"`
	Action     string  `json:"action" binding:"required"`
	CryptoName string  `json:"crypto_name" binding:"required"`
	Price      float64 `json:"price"`
	Strategy   int     `json:"strategy"`
	Type       string  `json:"type" binding:"required"`
}

func (h *Handler) Trade(c *gin.Context) {
	var req tradeReq
	if ok := bindData(c, &req); !ok {
		return
	}

	ctx := c.Request.Context()

	user, err := h.UserService.Get(ctx, req.UID)

	switch req.Type {
	case "market":
		_, err = h.TradeService.MarketTrade(ctx, user, req.Amount, req.Action, req.CryptoName, req.Strategy)

	case "limit":
		if req.Price == 0 {
			err := apperrors.NewBadRequest(err.Error())
			c.JSON(err.Status(), gin.H{
				"error": err,
			})
			return
		}
		_, err = h.TradeService.LimitTrade(ctx, user, req.Amount, req.Action, req.CryptoName, req.Price)
	}

	if err != nil {
		log.Println(err.(*apperrors.Error).Internal)
		err := apperrors.NewBadRequest(err.Error())
		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": "success",
	})
}

func (h *Handler) MockTrade(c *gin.Context) {
	var req tradeReq
	if ok := bindData(c, &req); !ok {
		return
	}

	ctx := c.Request.Context()

	user, err := h.UserService.Get(ctx, req.UID)

	switch req.Type {
	case "market":
		_, err = h.MockTradeService.MarketTrade(ctx, user, req.Amount, req.Action, req.CryptoName, req.Strategy)

	case "limit":
		_, err = h.MockTradeService.LimitTrade(ctx, user, req.Amount, req.Action, req.CryptoName, req.Price)
	}
	if err != nil {
		err := apperrors.NewBadRequest(err.Error())
		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": "success",
	})
}
