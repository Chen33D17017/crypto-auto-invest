package handler

import (
	"crypto-auto-invest/model"
	"crypto-auto-invest/model/apperrors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type patchReq struct {
	Name      string `json:"name" binding:"omitempty,max=50"`
	Email     string `json:"email" binding:"omitempty,email"`
	ApiKey    string `json:"api_key" binding:"omitempty"`
	ApiSecret string `json:"api_secret" binding:"omitempty"`
}

func (h *Handler) PatchUser(c *gin.Context) {
	authUser := c.MustGet("user").(*model.User)

	var req patchReq

	if ok := bindData(c, &req); !ok {
		return
	}

	// Should be returned with current imageURL
	u := &model.User{
		UID:       authUser.UID,
		Name:      req.Name,
		Email:     req.Email,
		ApiKey:    req.ApiKey,
		ApiSecret: req.ApiSecret,
	}

	ctx := c.Request.Context()
	user, err := h.UserService.PatchDetails(ctx, u)

	if err != nil {
		log.Printf("Failed to update user: %v\n", err.Error())

		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
