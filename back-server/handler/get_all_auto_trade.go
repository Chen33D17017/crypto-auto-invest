package handler

import (
	"crypto-auto-invest/model/apperrors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetAllAutoTrades(c *gin.Context) {
	ctx := c.Request.Context()
	rst, err := h.AutoTradeService.GetAllAutoTrades(ctx)
	if err != nil {
		log.Println(err.Error())
		e := err.(*apperrors.Error)
		c.JSON(e.Status(), gin.H{
			"error": err,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"data": rst,
	})
}
