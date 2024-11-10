package feed

import (
	"context"
	"fmt"
	"strconv"
	"time"

	x "github.com/ssentinull/dealls-dating-service/pkg/stdlib/stacktrace"

	redis "github.com/go-redis/redis/v8"
)

const (
	SwipeCountCacheKey        = "swipe:count:user_id:%d"
	SwipeCountCacheExpiration = 24 * time.Hour
)

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
