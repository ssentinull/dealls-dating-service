package mapper

import (
	"github.com/ssentinull/dealls-dating-service/internal/business/model"
	"github.com/ssentinull/dealls-dating-service/internal/types"
)

func MapUserModelToUserType(src model.UserModel) *types.User {
	return &types.User{
		ID:        src.Id,
		Email:     src.Email,
		Name:      src.Name,
		Gender:    src.Gender,
		BirthDate: src.BirthDate.Format("02-01-2006"),
		Location:  src.Location,
		CreatedAt: types.CreatedAt(src.CreatedAt),
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
