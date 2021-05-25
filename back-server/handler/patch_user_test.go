package handler

import (
	"bytes"
	"crypto-auto-invest/model"
	"crypto-auto-invest/model/apperrors"
	"crypto-auto-invest/model/mocks"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPatchUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Data binding error", func(t *testing.T) {
		router, mockUserService, _ := setEnv()
		rr := httptest.NewRecorder()

		reqBody, _ := json.Marshal(gin.H{
			"email": "notanemail",
		})
		request, _ := http.NewRequest(http.MethodPatch, "/details", bytes.NewBuffer(reqBody))
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		mockUserService.AssertNotCalled(t, "PatchDetails")
	})

	t.Run("Update success", func(t *testing.T) {
		router, mockUserService, ctxUser := setEnv()
		rr := httptest.NewRecorder()

		newEmail := "test@gmail.com"

		reqBody, _ := json.Marshal(gin.H{
			"email": newEmail,
		})

		request, _ := http.NewRequest(http.MethodPatch, "/details", bytes.NewBuffer(reqBody))
		request.Header.Set("Content-Type", "application/json")

		userToUpdate := &model.User{
			UID:   ctxUser.UID,
			Email: newEmail,
		}

		updateArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			userToUpdate,
		}

		rstUpdate := &model.User{
			UID:       ctxUser.UID,
			Email:     newEmail,
			ApiKey:    "apiKey",
			ApiSecret: "apiSecret",
		}

		mockUserService.
			On("PatchDetails", updateArgs...).
			Return(rstUpdate, nil)

		router.ServeHTTP(rr, request)

		respBody, _ := json.Marshal(gin.H{
			"user": rstUpdate,
		})

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockUserService.AssertCalled(t, "PatchDetails", updateArgs...)
	})

	t.Run("Update fail", func(t *testing.T) {
		router, mockUserService, ctxUser := setEnv()
		rr := httptest.NewRecorder()

		newEmail := "test@gmail.com"

		reqBody, _ := json.Marshal(gin.H{
			"email": newEmail,
		})

		request, _ := http.NewRequest(http.MethodPatch, "/details", bytes.NewBuffer(reqBody))
		request.Header.Set("Content-Type", "application/json")

		userToUpdate := &model.User{
			UID:   ctxUser.UID,
			Email: newEmail,
		}

		updateArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			userToUpdate,
		}

		mockError := apperrors.NewInternal()

		mockUserService.
			On("PatchDetails", updateArgs...).
			Return(nil, mockError)

		router.ServeHTTP(rr, request)

		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockUserService.AssertCalled(t, "PatchDetails", updateArgs...)
	})

}

func setEnv() (*gin.Engine, *mocks.MockUserService, *model.User) {
	uid, _ := uuid.NewRandom()
	ctxUser := &model.User{
		UID: uid.String(),
	}

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("user", ctxUser)
	})

	mockUserService := new(mocks.MockUserService)

	NewHandler(&Config{
		R:           router,
		UserService: mockUserService,
	})

	return router, mockUserService, ctxUser
}
