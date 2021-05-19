package handler

import (
	"crypto-auto-invest/model"
	"crypto-auto-invest/model/apperrors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type updateCronReq struct {
	ID          string  `json:"id" binding:"required"`
	Type        string  `json:"type" binding:"required,lowercase"`
	Amount      float64 `json:"amount" binding:"required"`
	Timepattern string  `json:"time_pattern" binding:"required"`
}

func (h *Handler) UpdateCron(c *gin.Context) {
	var req updateCronReq

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
		ID:          req.ID,
		UID:         user.(*model.User).UID,
		Type:        req.Type,
		Amount:      req.Amount,
		TimePattern: req.Timepattern,
	}
	err := h.CronService.UpdateCron(ctx, cb)

	if err != nil {
		log.Printf("Failed to update cron: %v\n", err.Error())

		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": cb,
	})
}
