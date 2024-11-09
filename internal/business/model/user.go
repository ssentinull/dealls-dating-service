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

	SignupUserParams struct {
		params.SignupUserParams
	}
)

func (b UserModel) TableName() string {
	return "users"
}
