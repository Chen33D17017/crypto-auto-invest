package middleware

import (
	"crypto-auto-invest/model/apperrors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHeader struct {
	ServiceToken string `header:"Authorization"`
}

func AuthService(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := AuthHeader{}

		// bind Authorization Header to h and check for validation errors
		if err := c.ShouldBindHeader(&h); err != nil {
			if errs, ok := err.(validator.ValidationErrors); ok {
				// we used this type in bind_data to extract desired fields from errs
				// you might consider extracting it
				var invalidArgs []invalidArgument

				for _, err := range errs {
					invalidArgs = append(invalidArgs, invalidArgument{
						err.Field(),
						err.Value().(string),
						err.Tag(),
						err.Param(),
					})
				}

				err := apperrors.NewBadRequest("Invalid request parameters. See invalidArgs")

				c.JSON(err.Status(), gin.H{
					"error":       err,
					"invalidArgs": invalidArgs,
				})
				c.Abort()
				return
			}

			// otherwise error type is unknown
			err := apperrors.NewInternal()
			c.JSON(err.Status(), gin.H{
				"error": err,
			})
			c.Abort()
			return
		}

		secretHeader := strings.Split(h.ServiceToken, "Bearer ")

		if len(secretHeader) < 2 {
			err := apperrors.NewAuthorization("Must provide Authorization header with format `Bearer {token}`")

			c.JSON(err.Status(), gin.H{
				"error": err,
			})
			c.Abort()
			return
		}

		if secret != secretHeader[1] {
			err := apperrors.NewAuthorization("Provided token is invalid")
			c.JSON(err.Status(), gin.H{
				"error": err,
			})
			c.Abort()
			return
		}
	}
}
