package handler

import (
	"crypto-auto-invest/model"
	"crypto-auto-invest/model/apperrors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) DeleteAutoTrade(c *gin.Context) {
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
	currencyName, ok := c.GetQuery("crypto_name")

	if !ok {
		log.Printf("Unable to extract currecncy type")
		err := apperrors.NewBadRequest("Need to query with currecncy type")
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

	if !ok {
		log.Printf("Unable to extract currecncy type")
		err := apperrors.NewBadRequest("Need to query with currecncy type")
		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	err := h.AutoTradeService.DeleteAutoTrade(ctx, uid, currencyName, strategyID)

	if err != nil {
		log.Printf("Unable to Delete cron: %v\n%v\n", uid, err)
		e := err.(*apperrors.Error)
		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
