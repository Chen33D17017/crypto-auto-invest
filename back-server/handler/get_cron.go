package handler

import (
	"crypto-auto-invest/model"
	"crypto-auto-invest/model/apperrors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetCron(c *gin.Context) {
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

	cron, err := h.CronService.GetCron(ctx, uid, cid)

	if err != nil {
		log.Printf("Unable to find cron: %v\n%v\n", uid, err)
		e := apperrors.NewNotFound("wallet", fmt.Sprintf("%s,%s", uid, cid))
		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": cron,
	})
}

func (h *Handler) GetCrons(c *gin.Context) {
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
	crons, err := h.CronService.GetCrons(ctx, uid)

	if err != nil {
		log.Printf("Unable to find user: %v\n%v\n", uid, err)
		e := apperrors.NewNotFound("user", uid)
		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": crons,
	})
}
