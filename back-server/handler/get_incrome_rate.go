package handler

import (
	"crypto-auto-invest/model"
	"crypto-auto-invest/model/apperrors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetIncomseRate(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		log.Printf("Unable to extract user from request context for unknow reason: %v\n", c)
		err := apperrors.NewInternal()
		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	uid := user.(*model.User).UID

	ctx := c.Request.Context()
	cryptoName, ok := c.GetQuery("crypto_name")

	if !ok {
		log.Printf("Unable to extract currency type")
		err := apperrors.NewBadRequest("Need to query with currency name")
		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	strategy, ok := c.GetQuery("strategy_id")

	if !ok {
		log.Printf("Unable to extract strategy id")
		err := apperrors.NewBadRequest("Need to query with strategy id")
		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}
	strategyID, _ := strconv.Atoi(strategy)

	incomeRate, err := h.TradeService.CalIncomeRate(ctx, uid, cryptoName, strategyID)
	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": incomeRate,
	})
}
