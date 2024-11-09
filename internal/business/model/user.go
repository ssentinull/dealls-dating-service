package model

import (
	params "github.com/ssentinull/dealls-dating-service/internal/types/auth"
)

type (
	UserModel struct {
		ID       int64
		Email    string
		Password string
		Name     string
	}

	GetUserParams struct {
		Email string
	}

	SignupUserParams struct {
		params.SignupUserParams
	}

	LoginUserParams struct {
		params.LoginUserParams
	}
)

func (b UserModel) TableName() string {
	return "users"
}
