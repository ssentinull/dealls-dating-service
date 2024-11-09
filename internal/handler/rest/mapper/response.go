package mapper

import (
	"github.com/ssentinull/dealls-dating-service/internal/business/model"
	"github.com/ssentinull/dealls-dating-service/internal/types"
)

func MapUserModelToUserType(src model.UserModel) *types.User {
	res := &types.User{
		ID:    src.ID,
		Email: src.Email,
	}

	return res
}
