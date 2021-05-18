package handler

import (
	"crypto-auto-invest/bitbank"
	"crypto-auto-invest/model/apperrors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetPrice(c *gin.Context) {
	currencyType, ok := c.GetQuery("type")
	if !ok {
		log.Printf("Unable to extract currency type")
		err := apperrors.NewNotFound("currency type", "")
		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	rst, err := bitbank.GetPrice(currencyType)
	if err != nil {
		log.Printf("Unable to extract currency type")
		err := apperrors.NewNotFound("currency type", currencyType)
		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": rst,
	})
}
