package handler

import (
	"crypto-auto-invest/model"
	"crypto-auto-invest/model/apperrors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) DeleteCron(c *gin.Context) {
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
	cid, ok := c.GetQuery("id")

	if !ok {
		log.Printf("Unable to extract cron id")
		err := apperrors.NewBadRequest("Need to query with cron id")
		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	err := h.CronService.DeleteCron(ctx, uid, cid)

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
