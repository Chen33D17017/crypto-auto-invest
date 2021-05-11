package handler

import (
	"crypto-auto-invest/model"
	"crypto-auto-invest/model/apperrors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type detailsReq struct {
	Name    string `json:"name" binding:"omitempty,max=50"`
	Email   string `json:"email" binding:"required,email"`
	Website string `json:"website" binding:"omitempty,url"`
}

// Details handler
func (h *Handler) UserDetails(c *gin.Context) {
	authUser := c.MustGet("user").(*model.User)

	var req detailsReq

	if ok := bindData(c, &req); !ok {
		return
	}

	// Should be returned with current imageURL
	u := &model.User{
		UID:     authUser.UID,
		Name:    req.Name,
		Email:   req.Email,
	}

	ctx := c.Request.Context()
	err := h.UserService.UpdateDetails(ctx, u)

	if err != nil {
		log.Printf("Failed to update user: %v\n", err.Error())

		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": u,
	})
}
