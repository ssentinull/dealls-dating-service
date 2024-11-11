package feed

import (
	"context"

	"github.com/ssentinull/dealls-dating-service/internal/business/model"
	"github.com/ssentinull/dealls-dating-service/internal/types"
	x "github.com/ssentinull/dealls-dating-service/pkg/stdlib/stacktrace"

	redis "github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func (f *feedImpl) CreatePreference(ctx context.Context, tx *gorm.DB, p model.PreferenceModel) (model.PreferenceModel, error) {
	result, err := f.createPreferenceSQL(ctx, tx, p)
	if err != nil {
		err = x.WrapPassCode(err, "createPreferenceSQL")
		return result, err
	}

	return result, nil
}

func (f *feedImpl) GetPreferenceByParams(ctx context.Context, p model.GetPreferenceParams) (model.PreferenceModel, error) {
	result, err := f.getPreferenceCache(ctx, p)
	if err == redis.Nil {
		result, err := f.getPreferenceSQL(ctx, p)
		if err != nil {
			err = x.WrapPassCode(err, "getPreferenceSQL")
			return result, err
		}

		if err = f.setPreferenceCache(ctx, p, result); err != nil {
			f.efLogger.Warn(err)
		}

		return result, nil
	} else if err != nil {
		result, err := f.getPreferenceSQL(ctx, p)
		if err != nil {
			err = x.WrapPassCode(err, "getPreferenceSQL")
			return result, err
		}

		return result, nil
	}

	return result, nil
}

func (f *feedImpl) GetFeedByParams(ctx context.Context, p model.GetFeedParams) ([]model.FeedModel, *types.Pagination, error) {
	result, pagination, err := f.getFeedSQL(ctx, p)
	if err != nil {
		err = x.WrapPassCode(err, "getFeedSQL")
		return []model.FeedModel{}, nil, err
	}

	return result, pagination, nil
}

func (f *feedImpl) CreateSwipe(ctx context.Context, tx *gorm.DB, p model.SwipeModel) (model.SwipeModel, error) {
	result, err := f.createSwipeSQL(ctx, tx, p)
	if err != nil {
		err = x.WrapPassCode(err, "createSwipeSQL")
		return result, err
	}

	return result, nil
}

func (f *feedImpl) GetSwipeByParams(ctx context.Context, p model.GetSwipeParams) (model.SwipeModel, error) {
	result, err := f.getSwipeSQL(ctx, p)
	if err != nil {
		err = x.WrapPassCode(err, "getSwipeSQL")
		return model.SwipeModel{}, err
	}

	return result, nil
}

func (f *feedImpl) GetSwipeCountByUserId(ctx context.Context, userId int64) (int64, error) {
	result, err := f.getSwipeCountCache(ctx, userId)
	if err != nil {
		err = x.WrapPassCode(err, "getSwipeCountCache")
		return 0, err
	}

	return result, nil
}

func (f *feedImpl) SetSwipeCountByUserId(ctx context.Context, userId, swipeCount int64) error {
	if err := f.setSwipeCountCache(ctx, userId, swipeCount); err != nil {
		err = x.WrapPassCode(err, "setSwipeCountCache")
		return err
	}

	return nil
}

func (f *feedImpl) CreateMatch(ctx context.Context, tx *gorm.DB, p model.MatchModel) (model.MatchModel, error) {
	result, err := f.createMatchSQL(ctx, tx, p)
	if err != nil {
		err = x.WrapPassCode(err, "createMatchSQL")
		return result, err
	}

	return result, nil
}
