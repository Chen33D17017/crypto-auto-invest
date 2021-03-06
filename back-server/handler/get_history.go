package handler

import (
	"crypto-auto-invest/bitbank"
	bm "crypto-auto-invest/bitbank/model"
	"crypto-auto-invest/model"
	"crypto-auto-invest/model/apperrors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetHistory(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		log.Printf("Unable to extract user from request context for unknow reason: %v\n", c)
		err := apperrors.NewInternal()
		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	currencyName, ok := c.GetQuery("type")
	if !ok {
		log.Printf("Unable to extract currency type")
		err := apperrors.NewBadRequest("Need to query with concurrency type")
		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	u := user.(*model.User)
	ctx := c.Request.Context()

	target, err := h.UserService.Get(ctx, u.UID)
	secret := bm.Secret{
		ApiKey:    target.ApiKey,
		ApiSecret: target.ApiSecret}
	history, err := bitbank.GetTradeHistory(secret, currencyName)

	if err != nil {
		log.Printf("bitbank err: %v\n", err)
		e := apperrors.NewBadRequest(err.Error())
		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": history,
	})
}
