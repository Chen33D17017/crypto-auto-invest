package handler

import (
	"crypto-auto-invest/model/apperrors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type tradeReq struct {
	CryptoName string  `json:"crypto_name" binding:"required"`
	Rate       float64 `json:"rate" binding:"required"`
	Action     string  `json:"action" binding:"required"`
	Strategy   int     `json:"strategy" binding:"required"`
}

func (h *Handler) Trade(c *gin.Context) {

	var req tradeReq

	if ok := bindData(c, &req); !ok {
		return
	}

	ctx := c.Request.Context()

	err := h.AutoTradeService.TradeWithRate(ctx, req.CryptoName, req.Action, req.Rate, req.Strategy)
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
