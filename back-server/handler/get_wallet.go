package handler

import (
	"crypto-auto-invest/model"
	"crypto-auto-invest/model/apperrors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetWallets(c *gin.Context) {
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
	wallets, err := h.WalletService.GetWallets(ctx, uid)

	if err != nil {
		log.Printf("Unable to find user: %v\n%v\n", uid, err)
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"wallets": wallets,
	})
}

func (h *Handler) GetWallet(c *gin.Context) {
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
	currencyType, ok := c.GetQuery("type")

	if !ok {
		log.Printf("Unable to extract currency type")
		err := apperrors.NewBadRequest("Need to query withcurrency type")
		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	wallets, err := h.WalletService.GetUserWallet(ctx, uid, currencyType)

	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"wallets": wallets,
	})
}
