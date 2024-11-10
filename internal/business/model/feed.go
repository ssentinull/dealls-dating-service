package model

import (
	"time"

	params "github.com/ssentinull/dealls-dating-service/internal/types/feed"

	"github.com/guregu/null/v5"
	"gorm.io/gorm"
)

const (
	SwipeTypeRight = "RIGHT"
	SwipeTypeLeft  = "LEFT"
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

	FeedModel struct {
		Id       int64
		Name     string
		Gender   string
		Age      int64
		Location string
	}

	GetFeedParams struct {
		UserId   int64
		Gender   string
		Location string
		MinAge   int64
		MaxAge   int64
		params.GetFeedParams
	}

	SwipeModel struct {
		Id         int64
		FromUserId int64
		ToUserId   int64
		SwipeType  string
		CreatedAt  time.Time
	}

	GetSwipeParams struct {
		FromUserId int64
		ToUserId   int64
		CreatedAt  time.Time
		SwipeType  string
	}

	SwipeFeedParams struct {
		FromUserId int64
		params.SwipeFeedParams
	}

	MatchModel struct {
		Id            int64
		MyUserId      int64
		MatchedUserId int64
		CreatedAt     time.Time
	}
)

func (p PreferenceModel) TableName() string {
	return "preferences"
}

func (s SwipeModel) TableName() string {
	return "swipes"
}

func (m MatchModel) TableName() string {
	return "matches"
}
