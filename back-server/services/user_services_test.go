package services

import (
	"context"
	"crypto-auto-invest/model"
	"crypto-auto-invest/model/apperrors"
	"crypto-auto-invest/model/mocks"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGet(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockUserResp := &model.User{
			UID:   uid.String(),
			Email: "bob@bob.com",
			Name:  "Bobby Bobson",
		}

		mockUserRepository := new(mocks.MockUserRepository)
		us := NewUserService(&USConfig{
			UserRepository: mockUserRepository,
		})
		mockUserRepository.On("FindByID", mock.Anything, uid.String()).Return(mockUserResp, nil)

		ctx := context.TODO()
		u, err := us.Get(ctx, uid.String())

		assert.NoError(t, err)
		assert.Equal(t, u, mockUserResp)
		mockUserRepository.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockUserRepository := new(mocks.MockUserRepository)
		us := NewUserService(&USConfig{
			UserRepository: mockUserRepository,
		})

		mockUserRepository.On("FindByID", mock.Anything, uid.String()).Return(nil, fmt.Errorf("Some error down the call chain"))

		ctx := context.TODO()
		u, err := us.Get(ctx, uid.String())

		assert.Nil(t, u)
		assert.Error(t, err)
		mockUserRepository.AssertExpectations(t)
	})
}

func TestSignup(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockUser := &model.User{
			UID:       uid.String(),
			Email:     "test@gmail.com",
			Name:      "testUser",
			ApiKey:    "apiKey",
			ApiSecret: "apiSecret",
		}

		mockUserRepository := new(mocks.MockUserRepository)
		us := NewUserService(&USConfig{
			UserRepository: mockUserRepository,
		})

		// We can use Run method to modify the user when the Create method is called.
		//  We can then chain on a Return method to return no error
		mockUserRepository.
			On("Create", mock.AnythingOfType("*context.emptyCtx"), mockUser).
			Run(func(args mock.Arguments) {
				userArg := args.Get(1).(*model.User) // arg 0 is context, arg 1 is *User
				userArg.UID = uid.String()
			}).Return(nil)

		ctx := context.TODO()
		err := us.Signup(ctx, mockUser)

		assert.NoError(t, err)

		// assert user now has a userID
		assert.Equal(t, uid.String(), mockUser.UID)

		mockUserRepository.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockUser := &model.User{
			Email:    "bob@bob.com",
			Password: "howdyhoneighbor!",
		}

		mockUserRepository := new(mocks.MockUserRepository)
		us := NewUserService(&USConfig{
			UserRepository: mockUserRepository,
		})

		mockErr := apperrors.NewConflict("email", mockUser.Email)

		// We can use Run method to modify the user when the Create method is called.
		//  We can then chain on a Return method to return no error
		mockUserRepository.
			On("Create", mock.AnythingOfType("*context.emptyCtx"), mockUser).
			Return(mockErr)

		ctx := context.TODO()
		err := us.Signup(ctx, mockUser)

		// assert error is error we response with in mock
		assert.EqualError(t, err, mockErr.Error())

		mockUserRepository.AssertExpectations(t)
	})
}

func TestUpdateDetails(t *testing.T) {
	mockUserRepository := new(mocks.MockUserRepository)
	us := NewUserService(&USConfig{
		UserRepository: mockUserRepository,
	})

	t.Run("Success", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockUser := &model.User{
			UID:       uid.String(),
			Email:     "new@bob.com",
			Name:      "A New Bob!",
			ApiKey:    "apiKey",
			ApiSecret: "apiSecret",
		}

		mockArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mockUser,
		}

		mockUserRepository.
			On("Update", mockArgs...).Return(nil)

		ctx := context.TODO()
		err := us.UpdateDetails(ctx, mockUser)

		assert.NoError(t, err)
		mockUserRepository.AssertCalled(t, "Update", mockArgs...)
	})

	t.Run("Failure", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockUser := &model.User{
			UID: uid.String(),
		}

		mockArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mockUser,
		}

		mockError := apperrors.NewInternal()

		mockUserRepository.
			On("Update", mockArgs...).Return(mockError)

		ctx := context.TODO()
		err := us.UpdateDetails(ctx, mockUser)
		assert.Error(t, err)

		apperror, ok := err.(*apperrors.Error)
		assert.True(t, ok)
		assert.Equal(t, apperrors.Internal, apperror.Type)

		mockUserRepository.AssertCalled(t, "Update", mockArgs...)
	})
}

func TestPatchDetails(t *testing.T) {
	mockUserRepository := new(mocks.MockUserRepository)
	us := NewUserService(&USConfig{
		UserRepository: mockUserRepository,
	})

	t.Run("Success", func(t *testing.T) {

		uid, _ := uuid.NewRandom()

		mockReq := &model.User{
			UID:   uid.String(),
			Email: "test2@gmail.com",
		}

		mockUser := &model.User{
			UID:       uid.String(),
			Email:     "test2@gmail.com",
			Name:      "testUser",
			ApiKey:    "apiKey",
			ApiSecret: "apiSecret",
		}

		mockArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mockReq,
		}

		mockGetArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			uid.String(),
		}

		mockUserRepository.On("Patch", mockArgs...).Return(nil)

		mockUserRepository.On("FindByID", mockGetArgs...).Return(mockUser, nil)

		ctx := context.TODO()
		u, err := us.PatchDetails(ctx, mockReq)

		assert.NoError(t, err)
		assert.Equal(t, mockUser, u)

		mockUserRepository.AssertCalled(t, "Patch", mockArgs...)
		mockUserRepository.AssertCalled(t, "FindByID", mockGetArgs...)
	})

	t.Run("Fail to Patch User Detail", func(t *testing.T) {

		uid, _ := uuid.NewRandom()

		mockReq := &model.User{
			UID:   uid.String(),
			Email: "test2@gmail.com",
		}

		mockUser := &model.User{
			UID:       uid.String(),
			Email:     "test2@gmail.com",
			Name:      "testUser",
			ApiKey:    "apiKey",
			ApiSecret: "apiSecret",
		}

		mockArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mockReq,
		}

		mockGetArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			uid.String(),
		}

		mockUserRepository.
			On("Patch", mockArgs...).Return(fmt.Errorf("Fail to Patch"))
		mockUserRepository.On("FindByID", mockGetArgs...).Return(mockUser, nil)

		ctx := context.TODO()
		u, err := us.PatchDetails(ctx, mockReq)
		assert.Nil(t, u)
		assert.Error(t, err)

		mockUserRepository.AssertCalled(t, "Patch", mockArgs...)
		mockUserRepository.AssertNotCalled(t, "FindByID", mockGetArgs...)
	})

	t.Run("Fail to Get User Detail", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockReq := &model.User{
			UID:   uid.String(),
			Email: "test2@gmail.com",
		}

		mockArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mockReq,
		}

		mockGetArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			uid.String(),
		}

		mockUserRepository.On("Patch", mockArgs...).Return(nil)
		mockUserRepository.On("FindByID", mockGetArgs...).Return(nil, fmt.Errorf("Some error down the call chain"))

		ctx := context.TODO()
		u, err := us.PatchDetails(ctx, mockReq)
		assert.Nil(t, u)
		assert.Error(t, err)

		mockUserRepository.AssertCalled(t, "Patch", mockArgs...)
		mockUserRepository.AssertCalled(t, "FindByID", mockGetArgs...)
	})
}
