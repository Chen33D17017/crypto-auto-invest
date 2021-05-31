package handler

import (
	"crypto-auto-invest/model/apperrors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetAutoTradeInfo(c *gin.Context) {

	cryptoName, ok := c.GetQuery("crypto_name")
	if !ok {
		err := apperrors.NewBadRequest("Need to query with crypto_name")
		c.JSON(err.Status(), gin.H{
			"error": err,
		})
	}

	strategy, ok := c.GetQuery("strategy_id")
	if !ok {
		err := apperrors.NewBadRequest("Need to query with strategy id")
		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}
	strategyID, err := strconv.Atoi(strategy)
	if err != nil {
		err := apperrors.NewBadRequest("Wrong format on strategy id")
		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	ctx := c.Request.Context()
	rst, err := h.AutoTradeService.GetAutoTradesFromStrategy(ctx, cryptoName, strategyID)
	if err != nil {
		e := err.(*apperrors.Error)
		c.JSON(e.Status(), gin.H{
			"error": err,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"data": rst,
	})
}
