package handler

import (
	"crypto-auto-invest/model"
	"crypto-auto-invest/model/apperrors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type addCronReq struct {
	Type        string  `json:"type" binding:"required,lowercase"`
	Amount      float64 `json:"amount" binding:"required"`
	Timepattern string  `json:"time_pattern" binding:"required"`
}

func (h *Handler) AddCron(c *gin.Context) {

	var req addCronReq

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
	cb := &model.Cron{
		UID:         user.(*model.User).UID,
		Type:        req.Type,
		Amount:      req.Amount,
		TimePattern: req.Timepattern,
	}
	cron, err := h.CronService.AddCron(ctx, cb)

	if err != nil {
		log.Printf("Fail to Add cron %s to %v, err: %v\n", req.Type, cb.UID, err)
		e := err.(*apperrors.Error)
		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": cron,
	})
}
