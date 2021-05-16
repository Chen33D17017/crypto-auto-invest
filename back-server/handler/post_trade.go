package handler

import (
	"crypto-auto-invest/bitbank"
	bm "crypto-auto-invest/bitbank/model"
	"crypto-auto-invest/model"
	"crypto-auto-invest/model/apperrors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type tradeReq struct {
	Type   string  `json:"type" binding:"required"`
	Amount float64 `json:"amount" binding:"required"`
	Side   string  `json:"side" binding:"required"`
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
	secret := bm.Secret{
		ApiKey:    target.ApiKey,
		ApiSecret: target.ApiSecret}

	var rst bm.Order

	switch req.Side {
	case "buy":
		rst, err = bitbank.BuyWithJPY(secret, req.Type, int64(req.Amount))
	case "sell":
		rst, err = bitbank.SellToJPY(secret, req.Type, req.Amount)
	default:
		err := apperrors.NewBadRequest(fmt.Sprintf("Trade with wrong type of side: %v", req.Type))
		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

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
