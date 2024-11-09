package auth

import (
	"context"

	"github.com/ssentinull/dealls-dating-service/internal/business/model"

	"golang.org/x/crypto/bcrypt"
)

func (a *authUc) SignupUser(ctx context.Context, params model.SignupUserParams) (model.UserModel, error) {
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
