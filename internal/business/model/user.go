package model

import (
	"time"

	params "github.com/ssentinull/dealls-dating-service/internal/types/auth"

	"github.com/guregu/null/v5"
	"gorm.io/gorm"
)

type (
	UserModel struct {
		Id                int64
		Email             string
		Password          string
		Name              string
		Gender            string
		BirthDate         time.Time
		Location          string
		ProfilePictureUrl string
		IsPremiumUser     bool
		CreatedAt         time.Time
		UpdatedAt         null.Time
		DeletedAt         gorm.DeletedAt
	}

	GetUserParams struct {
		Id    int64
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
