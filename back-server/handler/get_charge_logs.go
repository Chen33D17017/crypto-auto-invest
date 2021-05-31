package handler

import (
	"crypto-auto-invest/model"
	"crypto-auto-invest/model/apperrors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetChargeLogs(c *gin.Context) {
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
	rst, err := h.WalletService.GetChargeLogs(ctx, uid, cryptoName, strategyID)
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
