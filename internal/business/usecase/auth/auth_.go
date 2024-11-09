package auth

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/ssentinull/dealls-dating-service/internal/business/model"
	x "github.com/ssentinull/dealls-dating-service/pkg/stdlib/stacktrace"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (a *authUc) SignupUser(ctx context.Context, params model.SignupUserParams) (model.UserModel, error) {
	birthDate, err := time.Parse("02-01-2006", params.Body.BirthDate)
	if err != nil {
		a.efLogger.Error(err)
		return model.UserModel{}, x.WrapWithCode(err, http.StatusBadRequest, "birth_date field layout is invalid")
	}

	existingUser, err := a.userDom.GetUserByParams(ctx, model.GetUserParams{Email: params.Body.Email})
	if err != nil && x.GetCause(err) != gorm.ErrRecordNotFound {
		a.efLogger.Error(err)
		return model.UserModel{}, err
	}

	if existingUser.Id > 0 {
		err = errors.New("email is already in use")
		a.efLogger.Error(err)
		return model.UserModel{}, x.WrapWithCode(err, http.StatusConflict, "email is already in use")
	}

	encPassword, err := bcrypt.GenerateFromPassword([]byte(params.Body.Password), bcrypt.DefaultCost)
	if err != nil {
		a.efLogger.Error(err)
		return model.UserModel{}, err
	}

	user := model.UserModel{
		Email:     params.Body.Email,
		Name:      params.Body.Name,
		Password:  string(encPassword),
		Gender:    string(params.Body.Gender),
		BirthDate: birthDate,
		Location:  string(params.Body.Location),
	}

	res, err := a.userDom.CreateUser(ctx, nil, user)
	if err != nil {
		a.efLogger.Error(err)
		return model.UserModel{}, err
	}

	return res, nil
}

func (a *authUc) LoginUser(ctx context.Context, params model.LoginUserParams) (model.JWTModel, error) {
	existingUser, err := a.userDom.GetUserByParams(ctx, model.GetUserParams{Email: params.Body.Email})
	if err != nil {
		err = errors.New("email not found")
		a.efLogger.Error(err)
		return model.JWTModel{}, x.WrapWithCode(err, http.StatusNotFound, "email not found")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(params.Body.Password)); err != nil {
		a.efLogger.Error(err)
		return model.JWTModel{}, x.WrapWithCode(errors.New("incorrect password"), http.StatusUnauthorized, "incorrect password")
	}

	jwtToken, err := a.auth.GenerateJWTToken(existingUser)
	if err != nil {
		a.efLogger.Error(err)
		return model.JWTModel{}, err
	}

	return model.JWTModel{Token: jwtToken}, nil
}
