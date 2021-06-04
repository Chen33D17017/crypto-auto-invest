package handler

import (
	"crypto-auto-invest/model/apperrors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type tradeReq struct {
	CryptoName string  `json:"crypto_name" binding:"required"`
	UID        string  `json:"uid" binding:"required"`
	Amount     float64 `json:"amount" binding:"required"`
	Action     string  `json:"action" binding:"required"`
	Strategy   int     `json:"strategy" binding:"required"`
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
		_, err = h.TradeService.LimitTrade(ctx, user, req.Amount, req.Action, req.CryptoName, req.Strategy)
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
