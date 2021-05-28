package handler

import (
	"crypto-auto-invest/model"
	"crypto-auto-invest/model/apperrors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type tradeReq struct {
	CryptoName string  `json:"crypto_name" binding:"required"`
	Amount     float64 `json:"amount" binding:"required"`
	Action     string  `json:"action" binding:"required"`
	Strategy   int     `json:"strategy" binding:"required"`
}

func (h *Handler) Trade(c *gin.Context) {

	var req tradeReq

	if ok := bindData(c, &req); !ok {
		return
	}

	user, exists := c.Get("user")
	if !exists {
		log.Printf("Unable to extract user from request context for unknow reason: %v\n", c)
		err := apperrors.NewInternal()
		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	u := user.(*model.User)
	ctx := c.Request.Context()

	target, err := h.UserService.Get(ctx, u.UID)

	rst, err := h.TradeService.Trade(ctx, target, req.Amount, req.Action, req.CryptoName, req.Strategy)
	if err != nil {
		err := apperrors.NewBadRequest(err.Error())
		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": rst,
	})
}
