package feed

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/ssentinull/dealls-dating-service/internal/business/model"
	x "github.com/ssentinull/dealls-dating-service/pkg/stdlib/stacktrace"

	redis "github.com/go-redis/redis/v8"
)

const (
	PreferenceCacheKey = "preference:param:%v"
	SwipeCountCacheKey = "swipe:count:user_id:%d"

	PreferenceCacheExpiration = 5 * time.Minute
	SwipeCountCacheExpiration = 24 * time.Hour
)

func (f *feedImpl) getPreferenceCache(ctx context.Context, p model.GetPreferenceParams) (model.PreferenceModel, error) {
	result := model.PreferenceModel{}
	rawKey, err := f.json.Marshal(p)
	if err != nil {
		return result, x.Wrap(err, "JSON Marshall Error")
	}

	key := fmt.Sprintf(PreferenceCacheKey, string(rawKey))
	resultRaw, err := f.cache.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return result, err
	} else if err != nil {
		return result, x.Wrap(err, "Redis Get Cache Error")
	}

	if err := f.json.Unmarshal(resultRaw, &result); err != nil {
		return result, x.Wrap(err, "JSON Unmarshal Error")
	}

	return result, nil
}

func (f *feedImpl) setPreferenceCache(ctx context.Context, p model.GetPreferenceParams, preference model.PreferenceModel) error {
	rawKey, err := f.json.Marshal(p)
	if err != nil {
		return x.Wrap(err, "JSON Marshall Error")
	}

	key := fmt.Sprintf(PreferenceCacheKey, string(rawKey))
	rawJSON, err := f.json.Marshal(preference)
	if err != nil {
		return x.Wrap(err, "JSON Marshall Error")
	}

	if err := f.cache.Set(ctx, key, rawJSON, PreferenceCacheExpiration).Err(); err != nil {
		f.efLogger.Error(err)
		return x.Wrap(err, "Redis Set Cache Error")
	}

	return nil
}

func (f *feedImpl) getSwipeCountCache(ctx context.Context, userId int64) (int64, error) {
	result := int64(0)
	key := fmt.Sprintf(SwipeCountCacheKey, userId)

	resultStr, err := f.cache.Get(ctx, key).Result()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		f.efLogger.Error(err)
		return result, x.Wrap(err, "Redis Get Cache Error")
	}

	result, err = strconv.ParseInt(resultStr, 10, 64)
	if err != nil {
		f.efLogger.Error(err)
		return result, x.Wrap(err, "Parse Int Error")
	}

	return result, nil
}

func (f *feedImpl) setSwipeCountCache(ctx context.Context, userId, swipeCount int64) error {
	key := fmt.Sprintf(SwipeCountCacheKey, userId)
	swipeCountStr := strconv.FormatInt(swipeCount, 10)

	if err := f.cache.Set(ctx, key, swipeCountStr, SwipeCountCacheExpiration).Err(); err != nil {
		f.efLogger.Error(err)
		return x.Wrap(err, "Redis Set Cache Error")
	}

	return nil
}
