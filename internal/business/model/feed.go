package model

import (
	"time"

	params "github.com/ssentinull/dealls-dating-service/internal/types/feed"

	"github.com/guregu/null/v5"
	"gorm.io/gorm"
)

type (
	PreferenceModel struct {
		Id        int64
		UserId    int64
		Gender    string
		MinAge    int64
		MaxAge    int64
		Location  string
		CreatedAt time.Time
		UpdatedAt null.Time
		DeletedAt gorm.DeletedAt
	}

	CreatePreferenceParams struct {
		UserId int64
		params.CreateFeedPreferenceParams
	}

	GetPreferenceParams struct {
		UserId int64
	}
)

func (p PreferenceModel) TableName() string {
	return "preferences"
}
