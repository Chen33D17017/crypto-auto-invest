package handler

import (
	"crypto-auto-invest/model/apperrors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type saveOrderReq struct {
	UID        string `json:"uid" binding:"required"`
	OrderID    int64  `json:"order_id" binding:"required"`
	CryptoName string `json:"crypto_name" binding:"required"`
}

func (h *Handler) SaveOrder(c *gin.Context) {
	var req saveOrderReq

	if ok := bindData(c, &req); !ok {
		return
	}

	ctx := c.Request.Context()

	u, _ := h.UserService.Get(ctx, req.UID)

	orderID := fmt.Sprintf("%v", req.OrderID)
	err := h.TradeService.SaveOrder(ctx, u, orderID, req.CryptoName, 0)

	if err != nil {
		err := apperrors.NewBadRequest(err.Error())
		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "success",
	})
}

func (h *Handler) MockSaveOrder(c *gin.Context) {
	var req saveOrderReq

	if ok := bindData(c, &req); !ok {
		return
	}

	ctx := c.Request.Context()

	u, _ := h.UserService.Get(ctx, req.UID)

	orderID := fmt.Sprintf("%v", req.OrderID)
	err := h.MockTradeService.SaveOrder(ctx, u, orderID, req.CryptoName, 0)

	if err != nil {
		err := apperrors.NewBadRequest(err.Error())
		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "success",
	})
}
