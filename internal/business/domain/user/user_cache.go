package user

import (
	"context"
	"fmt"
	"time"

	"github.com/ssentinull/dealls-dating-service/internal/business/model"
	x "github.com/ssentinull/dealls-dating-service/pkg/stdlib/stacktrace"

	redis "github.com/go-redis/redis/v8"
)

const (
	UserCacheKey        = "user:param:%v"
	UserCacheExpiration = 5 * time.Minute
)

func (u *userImpl) getUserCache(ctx context.Context, p model.GetUserParams) (model.UserModel, error) {
	result := model.UserModel{}
	rawKey, err := u.json.Marshal(p)
	if err != nil {
		return result, x.Wrap(err, "JSON Marshall Error")
	}

	key := fmt.Sprintf(UserCacheKey, string(rawKey))
	resultRaw, err := u.cache.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return result, err
	} else if err != nil {
		return result, x.Wrap(err, "Redis Get Cache Error")
	}

	if err := u.json.Unmarshal(resultRaw, &result); err != nil {
		return result, x.Wrap(err, "JSON Unmarshal Error")
	}

	return result, nil
}

func (u *userImpl) setUserCache(ctx context.Context, p model.GetUserParams, user model.UserModel) error {
	rawKey, err := u.json.Marshal(p)
	if err != nil {
		return x.Wrap(err, "JSON Marshall Error")
	}

	key := fmt.Sprintf(UserCacheKey, string(rawKey))
	rawJSON, err := u.json.Marshal(user)
	if err != nil {
		return x.Wrap(err, "JSON Marshall Error")
	}

	if err := u.cache.Set(ctx, key, rawJSON, UserCacheExpiration).Err(); err != nil {
		u.efLogger.Error(err)
		return x.Wrap(err, "Redis Set Cache Error")
	}

	return nil
}
