package handler

import (
	"account-tutorial/model"
	"account-tutorial/model/apperrors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type signinReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,gte=6,lte=30"`
}

func (h *Handler) Signin(c *gin.Context) {
	var req signinReq

	if ok := bindData(c, &req); !ok {
		return
	}

	u := &model.User{
		Email:    req.Email,
		Password: req.Password,
	}

	ctx := c.Request.Context()
	err := h.UserService.Signin(ctx, u)

	if err != nil {
		log.Printf("Failed to sign in user: %v\n", err.Error())

		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	tokens, err := h.TokenService.NewPairFromUser(ctx, u, "")

	if err != nil {
		log.Printf("Failed to create token for user: %v\n", err.Error())

		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return 
	}

	c.JSON(http.StatusOK, gin.H{
		"tokens": tokens,
	})
}
