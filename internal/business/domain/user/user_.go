package user

import (
	"context"

	"github.com/ssentinull/dealls-dating-service/internal/business/model"
	x "github.com/ssentinull/dealls-dating-service/pkg/stdlib/stacktrace"

	redis "github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func (u *userImpl) CreateUser(ctx context.Context, tx *gorm.DB, p model.UserModel) (model.UserModel, error) {
	result, err := u.createUserSQL(ctx, tx, p)
	if err != nil {
		err = x.WrapPassCode(err, "createUserSQL")
		return result, err
	}

	return result, nil
}

func (u *userImpl) GetUserByParams(ctx context.Context, p model.GetUserParams) (model.UserModel, error) {
	result, err := u.getUserCache(ctx, p)
	if err == redis.Nil {
		result, err := u.getUserSQL(ctx, p)
		if err != nil {
			err = x.WrapPassCode(err, "getUserSQL")
			return result, err
		}

		if err = u.setUserCache(ctx, p, result); err != nil {
			u.efLogger.Warn(err)
		}

		return result, nil
	} else if err != nil {
		result, err := u.getUserSQL(ctx, p)
		if err != nil {
			err = x.WrapPassCode(err, "getUserSQL")
			return result, err
		}

		return result, nil
	}

	return result, nil
}
