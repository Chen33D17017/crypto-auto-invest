package handler

import (
	"crypto-auto-invest/model"
	"crypto-auto-invest/model/apperrors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type addAutoTradeReq struct {
	Type string `json:"type" binding:"required,lowercase"`
}

func (h *Handler) AddAutoTrade(c *gin.Context) {

	var req addAutoTradeReq

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

	ctx := c.Request.Context()
	uid := user.(*model.User).UID
	err := h.AutoTradeService.AddAutoTrade(ctx, uid, req.Type)

	if err != nil {
		log.Printf("Fail to Add Auto Trade %s to %v, err: %v\n", req.Type, uid, err)
		e := err.(*apperrors.Error)
		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"uid":  uid,
		"type": req.Type,
	})
}
