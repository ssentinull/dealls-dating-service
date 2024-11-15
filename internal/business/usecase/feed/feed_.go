package feed

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/ssentinull/dealls-dating-service/internal/business/model"
	"github.com/ssentinull/dealls-dating-service/internal/types"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/libsql"
	x "github.com/ssentinull/dealls-dating-service/pkg/stdlib/stacktrace"

	"gorm.io/gorm"
)

func (f *feedUc) CreatePreference(ctx context.Context, params model.CreatePreferenceParams) (model.PreferenceModel, error) {
	existingPreference, err := f.feedDom.GetPreferenceByParams(ctx, model.GetPreferenceParams{UserId: params.UserId})
	if err != nil && x.GetCause(err) != gorm.ErrRecordNotFound {
		f.efLogger.Error(err)
		return model.PreferenceModel{}, err
	}

	if existingPreference.Id > 0 {
		err = errors.New("preference already exists")
		f.efLogger.Error(err)
		return model.PreferenceModel{}, x.WrapWithCode(err, http.StatusConflict, "preference already exists")
	}

	res, err := f.feedDom.CreatePreference(ctx, nil, model.PreferenceModel{
		UserId:   params.UserId,
		Gender:   string(params.Body.Gender),
		MinAge:   params.Body.MinAge,
		MaxAge:   params.Body.MaxAge,
		Location: string(params.Body.Location),
	})
	if err != nil {
		f.efLogger.Error(err)
		return model.PreferenceModel{}, err
	}

	return res, nil
}

func (f *feedUc) GetFeed(ctx context.Context, params model.GetFeedParams) ([]model.FeedModel, *types.Pagination, error) {
	if _, err := f.userDom.GetUserByParams(ctx, model.GetUserParams{Id: params.UserId}); err != nil {
		f.efLogger.Error(err)
		return []model.FeedModel{}, nil, x.WrapWithCode(err, http.StatusNotFound, "user not found")
	}

	preference, err := f.feedDom.GetPreferenceByParams(ctx, model.GetPreferenceParams{UserId: params.UserId})
	if err != nil {
		f.efLogger.Error(err)
		if x.GetCause(err) == gorm.ErrRecordNotFound {
			err = errors.New("create feed preference first")
			return []model.FeedModel{}, nil, x.WrapWithCode(err, http.StatusUnprocessableEntity, "create feed preference first")
		}
		return []model.FeedModel{}, nil, err
	}

	getFeedParams := model.GetFeedParams{
		UserId:        params.UserId,
		Gender:        preference.Gender,
		Location:      preference.Location,
		MinAge:        preference.MinAge,
		MaxAge:        preference.MaxAge,
		GetFeedParams: params.GetFeedParams,
	}

	feed, pagination, err := f.feedDom.GetFeedByParams(ctx, getFeedParams)
	if err != nil {
		f.efLogger.Error(err)
		return []model.FeedModel{}, nil, err
	}

	return feed, pagination, nil
}

func (f *feedUc) SwipeFeed(ctx context.Context, params model.SwipeFeedParams) (model.SwipeModel, error) {
	user, err := f.userDom.GetUserByParams(ctx, model.GetUserParams{Id: params.FromUserId})
	if err != nil {
		f.efLogger.Error(err)
		return model.SwipeModel{}, x.WrapWithCode(err, http.StatusNotFound, "from_user not found")
	}

	if _, err := f.userDom.GetUserByParams(ctx, model.GetUserParams{Id: params.Body.ToUserID}); err != nil {
		f.efLogger.Error(err)
		return model.SwipeModel{}, x.WrapWithCode(err, http.StatusNotFound, "to_user not found")
	}

	swipeCount := int64(0)
	if !user.IsPremiumUser {
		swipeCount, err = f.feedDom.GetSwipeCountByUserId(ctx, params.FromUserId)
		if err != nil {
			f.efLogger.Error(err)
			return model.SwipeModel{}, err
		}

		if swipeCount >= f.opt.DailySwipeThreshold {
			err = errors.New("user has reached daily swipe threshold")
			return model.SwipeModel{}, x.WrapWithCode(err, http.StatusTooManyRequests, "user has reached daily swipe threshold")
		}
	}

	existingSwipe, err := f.feedDom.GetSwipeByParams(ctx, model.GetSwipeParams{
		FromUserId: params.FromUserId,
		ToUserId:   params.Body.ToUserID,
		CreatedAt:  time.Now(),
	})
	if err != nil && x.GetCause(err) != gorm.ErrRecordNotFound {
		f.efLogger.Error(err)
		return model.SwipeModel{}, err
	}

	if existingSwipe.Id > 0 {
		err = errors.New("user has already swiped this person")
		f.efLogger.Error(err)
		return model.SwipeModel{}, x.WrapWithCode(err, http.StatusConflict, "user has already swiped this person")
	}

	tx := f.sqlTxDom.BeginTX()
	defer func(tx *gorm.DB) {
		if err != nil {
			if errRollback := f.sqlTxDom.RollbackTX(tx); errRollback != nil {
				err = x.WrapWithCode(errRollback, x.GetCode(errRollback), libsql.ErrorTxRollback)
			}
		}
	}(tx)

	res, err := f.feedDom.CreateSwipe(ctx, tx, model.SwipeModel{
		FromUserId: params.FromUserId,
		ToUserId:   params.Body.ToUserID,
		SwipeType:  string(params.Body.SwipeType),
	})
	if err != nil {
		f.efLogger.Error(err)
		return model.SwipeModel{}, err
	}

	matchingSwipe, err := f.feedDom.GetSwipeByParams(ctx, model.GetSwipeParams{
		FromUserId: params.Body.ToUserID,
		ToUserId:   params.FromUserId,
		SwipeType:  model.SwipeTypeRight,
		CreatedAt:  time.Now(),
	})
	if err != nil && x.GetCause(err) != gorm.ErrRecordNotFound {
		f.efLogger.Error(err)
		return model.SwipeModel{}, err
	}

	if matchingSwipe.Id > 0 {
		_, err = f.feedDom.CreateMatch(ctx, tx, model.MatchModel{
			MyUserId:      params.FromUserId,
			MatchedUserId: params.Body.ToUserID,
		})
		if err != nil {
			f.efLogger.Error(err)
			return model.SwipeModel{}, err
		}

		_, err = f.feedDom.CreateMatch(ctx, tx, model.MatchModel{
			MyUserId:      params.Body.ToUserID,
			MatchedUserId: params.FromUserId,
		})
		if err != nil {
			f.efLogger.Error(err)
			return model.SwipeModel{}, err
		}
	}

	if err = f.sqlTxDom.CommitTX(tx); err != nil {
		f.efLogger.Error(err)
		return model.SwipeModel{}, err
	}

	if !user.IsPremiumUser {
		swipeCount++
		if err = f.feedDom.SetSwipeCountByUserId(ctx, params.FromUserId, swipeCount); err != nil {
			f.efLogger.Warn(err)
		}
	}

	return res, nil
}
