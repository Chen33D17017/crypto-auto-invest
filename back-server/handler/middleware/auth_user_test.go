package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"crypto-auto-invest/model"
	"crypto-auto-invest/model/apperrors"
	"crypto-auto-invest/model/mocks"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAuthUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockTokenService := new(mocks.MockTokenService)

	uid, _ := uuid.NewRandom()
	u := &model.User{
		UID:   uid.String(),
		Email: "bob@bob.com",
	}

	// Since we mock tokenService, we need not
	// create actual JWTs
	validTokenHeader := "validTokenString"
	invalidTokenHeader := "invalidTokenString"
	invalidTokenErr := apperrors.NewAuthorization("Unable to verify user from idToken")

	mockTokenService.On("ValidateIDToken", validTokenHeader).Return(u, nil)
	mockTokenService.On("ValidateIDToken", invalidTokenHeader).Return(nil, invalidTokenErr)

	t.Run("Adds a user to context", func(t *testing.T) {
		rr := httptest.NewRecorder()

		_, r := gin.CreateTestContext(rr)

		var contextUser *model.User

		r.GET("/me", AuthUser(mockTokenService), func(c *gin.Context) {
			contextKeyVal, _ := c.Get("user")
			contextUser = contextKeyVal.(*model.User)
		})

		request, _ := http.NewRequest(http.MethodGet, "/me", http.NoBody)

		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", validTokenHeader))
		r.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, u, contextUser)

		mockTokenService.AssertCalled(t, "ValidateIDToken", validTokenHeader)
	})

	t.Run("Invalid Token", func(t *testing.T) {
		rr := httptest.NewRecorder()

		// creates a test context and gin engine
		_, r := gin.CreateTestContext(rr)

		r.GET("/me", AuthUser(mockTokenService))

		request, _ := http.NewRequest(http.MethodGet, "/me", http.NoBody)

		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", invalidTokenHeader))
		r.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		mockTokenService.AssertCalled(t, "ValidateIDToken", invalidTokenHeader)
	})

	t.Run("Missing Authorization Header", func(t *testing.T) {
		rr := httptest.NewRecorder()

		// creates a test context and gin engine
		_, r := gin.CreateTestContext(rr)

		r.GET("/me", AuthUser(mockTokenService))

		request, _ := http.NewRequest(http.MethodGet, "/me", http.NoBody)

		r.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		mockTokenService.AssertNotCalled(t, "ValidateIDToken")
	})
}
