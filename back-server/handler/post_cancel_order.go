package handler

import (
	"crypto-auto-invest/bitbank"
	bm "crypto-auto-invest/bitbank/model"
	"crypto-auto-invest/model/apperrors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type cancelOrderReq struct {
	UID        string `json:"uid" binding:"required"`
	OrderID    string `json:"order_id" binding:"required"`
	CryptoName string `json:"crypto_name" binding:"required"`
}

func (h *Handler) CancelOrder(c *gin.Context) {
	var req cancelOrderReq
	if ok := bindData(c, &req); !ok {
		return
	}

	ctx := c.Request.Context()

	user, err := h.UserService.Get(ctx, req.UID)

	secret := bm.Secret{
		ApiKey:    user.ApiKey,
		ApiSecret: user.ApiSecret,
	}
	order, err := bitbank.CancelOrder(secret, req.CryptoName, req.OrderID)

	if err != nil {
		log.Println(err.Error())
		err := apperrors.NewInternal()
		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": order,
	})
}
