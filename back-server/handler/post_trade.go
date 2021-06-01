package handler

import (
	"crypto-auto-invest/model/apperrors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type tradeReq struct {
	UID        string  `json:"uid" binding:"required"`
	Amount     float64 `json:"amount" binding:"required"`
	Action     string  `json:"action" binding:"required"`
	CryptoName string  `json:"crypto_name" binding:"required"`
	Strategy   int     `json:"strategy"`
}

func (h *Handler) Trade(c *gin.Context) {

	var req tradeReq

	if ok := bindData(c, &req); !ok {
		return
	}

	ctx := c.Request.Context()

	user, err := h.UserService.Get(ctx, req.UID)

	_, err = h.TradeService.Trade(ctx, user, req.Amount, req.Action, req.CryptoName, 0)
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
