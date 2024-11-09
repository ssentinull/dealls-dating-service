package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/ssentinull/dealls-dating-service/internal/business/model"
	x "github.com/ssentinull/dealls-dating-service/pkg/stdlib/stacktrace"

	"golang.org/x/crypto/bcrypt"
)

func (a *authUc) SignupUser(ctx context.Context, params model.SignupUserParams) (model.UserModel, error) {
	// TODO: check for duplicate email and return conflict status code

	encPassword, err := bcrypt.GenerateFromPassword([]byte(params.Body.Password), bcrypt.DefaultCost)
	if err != nil {
		a.efLogger.Error(err)
		return model.UserModel{}, err
	}

	book := model.UserModel{
		Email:    params.Body.Email,
		Name:     params.Body.Name,
		Password: string(encPassword),
	}

	res, err := a.userDom.CreateUser(ctx, nil, book)
	if err != nil {
		a.efLogger.Error(err)
		return model.UserModel{}, err
	}

	return res, nil
}

func (a *authUc) LoginUser(ctx context.Context, params model.LoginUserParams) (model.JWTModel, error) {
	existingUser, err := a.userDom.GetUserByParams(ctx, model.GetUserParams{Email: params.Body.Email})
	if err != nil {
		a.efLogger.Error(err)
		return model.JWTModel{}, err
	}

	if existingUser.ID <= 0 {
		return model.JWTModel{}, x.WrapWithCode(errors.New("email not found"), http.StatusNotFound, "email not found")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(params.Body.Password)); err != nil {
		return model.JWTModel{}, x.WrapWithCode(errors.New("incorrect password"), http.StatusUnauthorized, "incorrect password")
	}

	jwtToken, err := a.auth.GenerateJWTToken(existingUser)
	if err != nil {
		return model.JWTModel{}, err
	}

	return model.JWTModel{Token: jwtToken}, nil
}
