package mapper

import (
	"github.com/ssentinull/dealls-dating-service/internal/business/model"
	"github.com/ssentinull/dealls-dating-service/internal/types"
)

func MapUserModelToUserType(src model.UserModel) *types.User {
	return &types.User{
		ID:    src.ID,
		Email: src.Email,
	}
}

func MapJWTModelToJWTType(src model.JWTModel) *types.JWT {
	return &types.JWT{
		Token: src.Token,
	}
}
