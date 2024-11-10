package auth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ssentinull/dealls-dating-service/internal/business/model"
	"github.com/ssentinull/dealls-dating-service/internal/business/usecase/auth"
	swaggerTypes "github.com/ssentinull/dealls-dating-service/internal/types"
	swaggerAuthParams "github.com/ssentinull/dealls-dating-service/internal/types/auth"

	"github.com/c2fo/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAuthUsecase_SignupUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	defer func() {
		ctx.Done()
		ctrl.Finish()
	}()

	mockedDependency := NewMockedDependency(t, ctrl)
	authUsecase := auth.InitAuthUsecase(
		mockedDependency.userDom,
		mockedDependency.efLogger,
		mockedDependency.jsonParser,
		mockedDependency.auth,
		mockedDependency.opt,
	)

	ID := int64(1)
	email := "john.doe@mail.com"

	signupUserParam := swaggerAuthParams.SignupUserParams{
		Body: &swaggerTypes.SignupUserRequest{
			Email:     email,
			Password:  "12345",
			Name:      "John Doe",
			Gender:    "FEMALE",
			BirthDate: "01-01-1991",
			Location:  "JAKARTA",
		},
	}

	payload := model.SignupUserParams{SignupUserParams: signupUserParam}
	getUserParam := model.GetUserParams{Email: email}
	user := model.UserModel{Id: ID}

	t.Run("success", func(t *testing.T) {
		mockedDependency.userDom.EXPECT().GetUserByParams(gomock.Any(), getUserParam).Times(1).Return(model.UserModel{}, nil)
		mockedDependency.userDom.EXPECT().CreateUser(gomock.Any(), nil, gomock.Any()).Times(1).Return(user, nil)

		result, err := authUsecase.SignupUser(ctx, payload)
		assert.NoError(t, err)
		assert.NotEmpty(t, result)
	})

	t.Run("failed - invalid birth date layout", func(t *testing.T) {
		signupUserParam := swaggerAuthParams.SignupUserParams{
			Body: &swaggerTypes.SignupUserRequest{BirthDate: "ABCDE"},
		}

		invalidPayload := model.SignupUserParams{SignupUserParams: signupUserParam}

		result, err := authUsecase.SignupUser(ctx, invalidPayload)
		assert.Error(t, err)
		assert.NotEmpty(t, result)
	})

	t.Run("failed - get user by params return error", func(t *testing.T) {
		mockedDependency.userDom.EXPECT().GetUserByParams(gomock.Any(), getUserParam).
			Times(1).Return(model.UserModel{}, errors.New("db error"))

		result, err := authUsecase.SignupUser(ctx, payload)
		assert.Error(t, err)
		assert.NotEmpty(t, result)
	})

	t.Run("failed - email is already in use", func(t *testing.T) {
		existingUser := model.UserModel{
			Id:    int64(1),
			Email: email,
		}

		mockedDependency.userDom.EXPECT().GetUserByParams(gomock.Any(), getUserParam).
			Times(1).Return(existingUser, nil)

		result, err := authUsecase.SignupUser(ctx, payload)
		assert.Error(t, err)
		assert.NotEmpty(t, result)
	})

	t.Run("failed - create user return error", func(t *testing.T) {
		mockedDependency.userDom.EXPECT().GetUserByParams(gomock.Any(), getUserParam).
			Times(1).Return(model.UserModel{}, nil)
		mockedDependency.userDom.EXPECT().CreateUser(gomock.Any(), nil, gomock.Any()).
			Times(1).Return(model.UserModel{}, errors.New("db error"))

		result, err := authUsecase.SignupUser(ctx, payload)
		assert.Error(t, err)
		assert.NotEmpty(t, result)
	})
}

func TestAuthUsecase_LoginUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	defer func() {
		ctx.Done()
		ctrl.Finish()
	}()

	mockedDependency := NewMockedDependency(t, ctrl)
	authUsecase := auth.InitAuthUsecase(
		mockedDependency.userDom,
		mockedDependency.efLogger,
		mockedDependency.jsonParser,
		mockedDependency.auth,
		mockedDependency.opt,
	)

	ID := int64(1)
	email := "john.doe@mail.com"
	password := "12345"
	hashedPassword := "$2a$10$F43bitSBfOFurnMDHyhKF.8o8v5oPOfO5CETCzNZrByF1LNX2eCsy"
	jwtToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzEyMzYwNTksImlkIj"

	loginUserParam := swaggerAuthParams.LoginUserParams{
		Body: &swaggerTypes.LoginUserRequest{
			Email:    email,
			Password: password,
		},
	}

	payload := model.LoginUserParams{LoginUserParams: loginUserParam}
	getUserParam := model.GetUserParams{Email: email}
	user := model.UserModel{
		Id:       ID,
		Email:    email,
		Password: hashedPassword,
	}

	t.Run("success", func(t *testing.T) {
		mockedDependency.userDom.EXPECT().GetUserByParams(gomock.Any(), getUserParam).
			Times(1).Return(user, nil)
		mockedDependency.auth.EXPECT().GenerateJWTToken(user).Times(1).Return(jwtToken, nil)

		result, err := authUsecase.LoginUser(ctx, payload)
		assert.NoError(t, err)
		assert.NotEmpty(t, result)
	})

	t.Run("failed - get user by param return error", func(t *testing.T) {
		mockedDependency.userDom.EXPECT().GetUserByParams(gomock.Any(), getUserParam).
			Times(1).Return(model.UserModel{}, errors.New("db error"))

		result, err := authUsecase.LoginUser(ctx, payload)
		assert.Error(t, err)
		assert.NotEmpty(t, result)
	})

	t.Run("failed - password does not match", func(t *testing.T) {
		user := model.UserModel{
			Id:       ID,
			Email:    email,
			Password: "randompassword",
		}

		mockedDependency.userDom.EXPECT().GetUserByParams(gomock.Any(), getUserParam).
			Times(1).Return(user, nil)

		result, err := authUsecase.LoginUser(ctx, payload)
		assert.Error(t, err)
		assert.NotEmpty(t, result)
	})

	t.Run("failed - generate jwt token return error", func(t *testing.T) {
		mockedDependency.userDom.EXPECT().GetUserByParams(gomock.Any(), getUserParam).
			Times(1).Return(user, nil)
		mockedDependency.auth.EXPECT().GenerateJWTToken(user).Times(1).Return("", errors.New("db error"))

		result, err := authUsecase.LoginUser(ctx, payload)
		assert.Error(t, err)
		assert.NotEmpty(t, result)
	})
}
