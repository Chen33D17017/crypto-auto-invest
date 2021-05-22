package handler

import (
	"crypto-auto-invest/model"
	"crypto-auto-invest/model/apperrors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type addWalletReq struct {
	Type string `json:"crypto_type" binding:"required,lowercase"`
}

func (h *Handler) AddWallet(c *gin.Context) {

	var req addWalletReq

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

	uid := user.(*model.User).UID

	ctx := c.Request.Context()
	wallet, err := h.WalletService.AddWallet(ctx, uid, req.Type)

	if err != nil {
		log.Printf("Fail to Add Wallet %s to %v, err: %v\n", req.Type, uid, err)
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"wallets": wallet,
	})
}
