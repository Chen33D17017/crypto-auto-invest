package handler

import (
	"crypto-auto-invest/model"
	"crypto-auto-invest/model/apperrors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type changeWalletReq struct {
	CryptoName string  `json:"crypto_name" binding:"required,lowercase"`
	Amount     float64 `json:"amount" binding:"required"`
}

func (h *Handler) Charge(c *gin.Context) {

	var req changeWalletReq

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
	wallet, err := h.WalletService.ChangeMoney(ctx, uid, req.CryptoName, req.Amount)

	if err != nil {
		log.Printf("Fail to change Wallet value on %s to %v, err: %v\n", req.CryptoName, uid, err)
		e := err.(*apperrors.Error)
		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"wallets": wallet,
	})
}
