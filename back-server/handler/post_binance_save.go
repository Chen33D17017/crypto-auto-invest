package handler

import (
	"crypto-auto-invest/model/apperrors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type binanceReq struct {
	UID     string  `json:"uid" binding:"required"`
	Symbol  string  `json:"symbol" binding:"required"`
	Action  string  `json:"action" binding:"required"`
	AvgCost float64 `json:"average_cost" binding:"required"`
	Qty     float64 `json:"qty" binding:"required"`
}

func (h *Handler) BinanceSave(c *gin.Context) {

	var req binanceReq

	if ok := bindData(c, &req); !ok {
		return
	}

	ctx := c.Request.Context()
	if req.Action != "sell" && req.Action != "buy" {
		e := apperrors.NewBadRequest("wrong type of cation")
		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return

	}
	order, err := h.BinanceTradeService.SaveOrder(ctx, req.UID, req.Symbol, req.Action, req.AvgCost, req.Qty)

	if err != nil {
		log.Printf("Fail to save order from Binance, err: %v\n", err)
		e := apperrors.NewInternal()
		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": order,
	})
}
