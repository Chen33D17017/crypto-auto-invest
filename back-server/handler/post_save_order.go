package handler

import (
	"crypto-auto-invest/model"
	"crypto-auto-invest/model/apperrors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type saveOrderReq struct {
	OrderID int64  `json:"order_id" binding:"required"`
	CryptoName    string `json:"crypto_name" binding:"required"`
}

func (h *Handler) SaveOrder(c *gin.Context) {
	var req saveOrderReq

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

	ctx := c.Request.Context()

	u, _ := h.UserService.Get(ctx, user.(*model.User).UID)

	orderID := fmt.Sprintf("%v", req.OrderID)
	err := h.TradeService.SaveOrder(ctx, u, orderID, req.CryptoName, "fixed")

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
