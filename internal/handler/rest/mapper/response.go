package mapper

import (
	"github.com/go-openapi/strfmt"
	"github.com/ssentinull/dealls-dating-service/internal/business/model"
	"github.com/ssentinull/dealls-dating-service/internal/types"
)

func MapUserModelToUserType(src model.UserModel) *types.User {
	return &types.User{
		ID:                src.Id,
		Email:             src.Email,
		Name:              src.Name,
		Gender:            types.Gender(src.Gender),
		BirthDate:         src.BirthDate.Format("02-01-2006"),
		Location:          src.Location,
		ProfilePictureURL: strfmt.URI(src.ProfilePictureUrl),
		CreatedAt:         types.CreatedAt(src.CreatedAt),
	}
}

func MapJWTModelToJWTType(src model.JWTModel) *types.JWT {
	return &types.JWT{
		Token: src.Token,
	}
}

func MapPreferenceModelToPreferenceType(src model.PreferenceModel) *types.Preference {
	return &types.Preference{
		ID:        src.Id,
		UserID:    src.UserId,
		Gender:    src.Gender,
		MinAge:    src.MinAge,
		MaxAge:    src.MaxAge,
		Location:  src.Location,
		CreatedAt: types.CreatedAt(src.CreatedAt),
	}
}

func MapFeedModelToFeedType(src model.FeedModel) *types.Feed {
	return &types.Feed{
		ID:                src.Id,
		Name:              src.Name,
		Gender:            src.Gender,
		Age:               src.Age,
		Location:          src.Location,
		ProfilePictureURL: strfmt.URI(src.ProfilePictureUrl),
	}
}

func MapSwipeModelToSwipeType(src model.SwipeModel) *types.Swipe {
	return &types.Swipe{
		ID:         src.Id,
		FromUserID: src.FromUserId,
		ToUserID:   src.ToUserId,
		SwipeType:  types.SwipeType(src.SwipeType),
		CreatedAt:  types.CreatedAt(src.CreatedAt),
	}
}
