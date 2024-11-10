package feed_test

import (
	"context"
	"errors"
	"testing"

	"github.com/c2fo/testify/assert"
	"github.com/ssentinull/dealls-dating-service/internal/business/model"
	"github.com/ssentinull/dealls-dating-service/internal/business/usecase/feed"
	swaggerTypes "github.com/ssentinull/dealls-dating-service/internal/types"
	swaggerFeedParams "github.com/ssentinull/dealls-dating-service/internal/types/feed"

	"go.uber.org/mock/gomock"
)

func TestFeedUsecase_CreatePreference(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	defer func() {
		ctx.Done()
		ctrl.Finish()
	}()

	mockedDependency := NewMockedDependency(t, ctrl)
	feedUsecase := feed.InitFeedUsecase(
		mockedDependency.sqlTxDom,
		mockedDependency.feedDom,
		mockedDependency.userDom,
		mockedDependency.efLogger,
		mockedDependency.jsonParser,
		mockedDependency.auth,
		mockedDependency.opt,
	)

	ID := int64(1)
	createFeedParam := swaggerFeedParams.CreateFeedPreferenceParams{
		Body: &swaggerTypes.CreateFeedPreferenceRequest{
			Gender:   "MALE",
			Location: "JAKARTA",
			MinAge:   27,
			MaxAge:   30,
		},
	}

	payload := model.CreatePreferenceParams{
		UserId:                     ID,
		CreateFeedPreferenceParams: createFeedParam,
	}

	getPreferenceParam := model.GetPreferenceParams{UserId: ID}
	preference := model.PreferenceModel{
		UserId:   ID,
		Gender:   "MALE",
		MinAge:   27,
		MaxAge:   30,
		Location: "JAKARTA",
	}

	t.Run("success", func(t *testing.T) {
		mockedDependency.feedDom.EXPECT().GetPreferenceByParams(gomock.Any(), getPreferenceParam).
			Times(1).Return(model.PreferenceModel{}, nil)
		mockedDependency.feedDom.EXPECT().CreatePreference(gomock.Any(), nil, preference).
			Times(1).Return(preference, nil)

		result, err := feedUsecase.CreatePreference(ctx, payload)
		assert.NoError(t, err)
		assert.NotEmpty(t, result)
	})

	t.Run("failed - get preference by params return error", func(t *testing.T) {
		mockedDependency.feedDom.EXPECT().GetPreferenceByParams(gomock.Any(), getPreferenceParam).
			Times(1).Return(model.PreferenceModel{}, errors.New("db error"))

		result, err := feedUsecase.CreatePreference(ctx, payload)
		assert.Error(t, err)
		assert.NotEmpty(t, result)
	})

	t.Run("failed - create preference return error", func(t *testing.T) {
		mockedDependency.feedDom.EXPECT().GetPreferenceByParams(gomock.Any(), getPreferenceParam).
			Times(1).Return(model.PreferenceModel{}, nil)
		mockedDependency.feedDom.EXPECT().CreatePreference(gomock.Any(), nil, preference).
			Times(1).Return(model.PreferenceModel{}, errors.New("db error"))

		result, err := feedUsecase.CreatePreference(ctx, payload)
		assert.Error(t, err)
		assert.NotEmpty(t, result)
	})
}

func TestFeedUsecase_GetFeed(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	defer func() {
		ctx.Done()
		ctrl.Finish()
	}()

	mockedDependency := NewMockedDependency(t, ctrl)
	feedUsecase := feed.InitFeedUsecase(
		mockedDependency.sqlTxDom,
		mockedDependency.feedDom,
		mockedDependency.userDom,
		mockedDependency.efLogger,
		mockedDependency.jsonParser,
		mockedDependency.auth,
		mockedDependency.opt,
	)

	ID := int64(1)
	page := int64(1)
	size := int64(10)

	getFeedParams := swaggerFeedParams.GetFeedParams{
		PaginationPage: &page,
		PaginationSize: &size,
	}

	pagination := &swaggerTypes.Pagination{
		CurrentData: int64(1),
		CurrentPage: int64(1),
		TotalData:   int64(1),
		TotalPages:  int64(1),
	}

	params := model.GetFeedParams{
		UserId:        ID,
		GetFeedParams: getFeedParams,
	}

	getUserParam := model.GetUserParams{Id: ID}
	getPreferenceParam := model.GetPreferenceParams{UserId: ID}
	getFeedParam := model.GetFeedParams{
		UserId:        ID,
		Gender:        "FEMALE",
		Location:      "JAKARTA",
		MinAge:        27,
		MaxAge:        30,
		GetFeedParams: params.GetFeedParams,
	}

	user := model.UserModel{Id: ID}
	preference := model.PreferenceModel{
		Gender:   "FEMALE",
		Location: "JAKARTA",
		MinAge:   27,
		MaxAge:   30,
	}

	feed := []model.FeedModel{{
		Id:       1,
		Name:     "Jane Doe",
		Gender:   "FEMALE",
		Age:      28,
		Location: "JAKARTA",
	}}

	t.Run("success", func(t *testing.T) {
		mockedDependency.userDom.EXPECT().GetUserByParams(gomock.Any(), getUserParam).Times(1).Return(user, nil)
		mockedDependency.feedDom.EXPECT().GetPreferenceByParams(gomock.Any(), getPreferenceParam).Times(1).Return(preference, nil)
		mockedDependency.feedDom.EXPECT().GetFeedByParams(gomock.Any(), getFeedParam).Times(1).Return(feed, pagination, nil)

		resFeed, resPagination, err := feedUsecase.GetFeed(ctx, params)
		assert.NoError(t, err)
		assert.NotEmpty(t, resFeed)
		assert.NotNil(t, resPagination)
	})

	t.Run("failed - get user by params return error", func(t *testing.T) {
		mockedDependency.userDom.EXPECT().GetUserByParams(gomock.Any(), getUserParam).Times(1).Return(model.UserModel{}, errors.New("db error"))

		resFeed, resPagination, err := feedUsecase.GetFeed(ctx, params)
		assert.Error(t, err)
		assert.Empty(t, resFeed)
		assert.Nil(t, resPagination)
	})

	t.Run("failed - get preference by params return error", func(t *testing.T) {
		mockedDependency.userDom.EXPECT().GetUserByParams(gomock.Any(), getUserParam).Times(1).Return(user, nil)
		mockedDependency.feedDom.EXPECT().GetPreferenceByParams(gomock.Any(), getPreferenceParam).
			Times(1).Return(model.PreferenceModel{}, errors.New("db error"))

		resFeed, resPagination, err := feedUsecase.GetFeed(ctx, params)
		assert.Error(t, err)
		assert.Empty(t, resFeed)
		assert.Nil(t, resPagination)
	})

	t.Run("failed - get feed by params return error", func(t *testing.T) {
		mockedDependency.userDom.EXPECT().GetUserByParams(gomock.Any(), getUserParam).Times(1).Return(user, nil)
		mockedDependency.feedDom.EXPECT().GetPreferenceByParams(gomock.Any(), getPreferenceParam).Times(1).Return(preference, nil)
		mockedDependency.feedDom.EXPECT().GetFeedByParams(gomock.Any(), getFeedParam).
			Times(1).Return([]model.FeedModel{}, nil, errors.New("db error"))

		resFeed, resPagination, err := feedUsecase.GetFeed(ctx, params)
		assert.Error(t, err)
		assert.Empty(t, resFeed)
		assert.Nil(t, resPagination)
	})
}

func TestFeedUsecase_SwipeFeed(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	defer func() {
		ctx.Done()
		ctrl.Finish()
	}()

	mockedDependency := NewMockedDependency(t, ctrl)
	feedUsecase := feed.InitFeedUsecase(
		mockedDependency.sqlTxDom,
		mockedDependency.feedDom,
		mockedDependency.userDom,
		mockedDependency.efLogger,
		mockedDependency.jsonParser,
		mockedDependency.auth,
		mockedDependency.opt,
	)

	fromUserId := int64(1)
	toUserId := int64(2)

	fromUser := model.UserModel{
		Id:            fromUserId,
		IsPremiumUser: true,
	}
	toUser := model.UserModel{Id: toUserId}
	swipe := model.SwipeModel{
		FromUserId: fromUserId,
		ToUserId:   toUserId,
		SwipeType:  "RIGHT",
	}

	getFromUserParam := model.GetUserParams{Id: fromUserId}
	getToUserParam := model.GetUserParams{Id: toUserId}

	params := model.SwipeFeedParams{
		FromUserId: fromUserId,
		SwipeFeedParams: swaggerFeedParams.SwipeFeedParams{
			Body: &swaggerTypes.SwipeFeedRequest{
				SwipeType: "RIGHT",
				ToUserID:  toUserId,
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		mockedDependency.userDom.EXPECT().GetUserByParams(gomock.Any(), getFromUserParam).Times(1).Return(fromUser, nil)
		mockedDependency.userDom.EXPECT().GetUserByParams(gomock.Any(), getToUserParam).Times(1).Return(toUser, nil)
		mockedDependency.feedDom.EXPECT().GetSwipeByParams(gomock.Any(), gomock.Any()).Times(1).Return(model.SwipeModel{}, nil)
		mockedDependency.sqlTxDom.EXPECT().BeginTX().Times(1)
		mockedDependency.feedDom.EXPECT().CreateSwipe(gomock.Any(), gomock.Any(), swipe).Times(1).Return(swipe, nil)
		mockedDependency.feedDom.EXPECT().GetSwipeByParams(gomock.Any(), gomock.Any()).Times(1).Return(model.SwipeModel{}, nil)
		mockedDependency.sqlTxDom.EXPECT().CommitTX(gomock.Any()).Times(1)

		res, err := feedUsecase.SwipeFeed(ctx, params)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})

	t.Run("failed - get from user by params return error", func(t *testing.T) {
		mockedDependency.userDom.EXPECT().GetUserByParams(gomock.Any(), getFromUserParam).Times(1).Return(model.UserModel{}, errors.New("db error"))
		res, err := feedUsecase.SwipeFeed(ctx, params)
		assert.Error(t, err)
		assert.NotEmpty(t, res)
	})

	t.Run("failed - get to user by params return error", func(t *testing.T) {
		mockedDependency.userDom.EXPECT().GetUserByParams(gomock.Any(), getFromUserParam).Times(1).Return(fromUser, nil)
		mockedDependency.userDom.EXPECT().GetUserByParams(gomock.Any(), getToUserParam).Times(1).Return(model.UserModel{}, errors.New("db error"))

		res, err := feedUsecase.SwipeFeed(ctx, params)
		assert.Error(t, err)
		assert.NotEmpty(t, res)
	})

	t.Run("failed - get existing swipe by params return error", func(t *testing.T) {
		mockedDependency.userDom.EXPECT().GetUserByParams(gomock.Any(), getFromUserParam).Times(1).Return(fromUser, nil)
		mockedDependency.userDom.EXPECT().GetUserByParams(gomock.Any(), getToUserParam).Times(1).Return(toUser, nil)
		mockedDependency.feedDom.EXPECT().GetSwipeByParams(gomock.Any(), gomock.Any()).Times(1).Return(model.SwipeModel{}, errors.New("db error"))

		res, err := feedUsecase.SwipeFeed(ctx, params)
		assert.Error(t, err)
		assert.NotEmpty(t, res)
	})

	t.Run("failed - create swipe return error", func(t *testing.T) {
		mockedDependency.userDom.EXPECT().GetUserByParams(gomock.Any(), getFromUserParam).Times(1).Return(fromUser, nil)
		mockedDependency.userDom.EXPECT().GetUserByParams(gomock.Any(), getToUserParam).Times(1).Return(toUser, nil)
		mockedDependency.feedDom.EXPECT().GetSwipeByParams(gomock.Any(), gomock.Any()).Times(1).Return(model.SwipeModel{}, nil)
		mockedDependency.sqlTxDom.EXPECT().BeginTX().Times(1)
		mockedDependency.feedDom.EXPECT().CreateSwipe(gomock.Any(), gomock.Any(), swipe).Times(1).Return(model.SwipeModel{}, errors.New("db error"))
		mockedDependency.sqlTxDom.EXPECT().RollbackTX(gomock.Any()).Times(1)

		res, err := feedUsecase.SwipeFeed(ctx, params)
		assert.Error(t, err)
		assert.NotEmpty(t, res)
	})

	t.Run("failed - get matching swipe by params return error", func(t *testing.T) {
		mockedDependency.userDom.EXPECT().GetUserByParams(gomock.Any(), getFromUserParam).Times(1).Return(fromUser, nil)
		mockedDependency.userDom.EXPECT().GetUserByParams(gomock.Any(), getToUserParam).Times(1).Return(toUser, nil)
		mockedDependency.feedDom.EXPECT().GetSwipeByParams(gomock.Any(), gomock.Any()).Times(1).Return(model.SwipeModel{}, nil)
		mockedDependency.sqlTxDom.EXPECT().BeginTX().Times(1)
		mockedDependency.feedDom.EXPECT().CreateSwipe(gomock.Any(), gomock.Any(), swipe).Times(1).Return(swipe, nil)
		mockedDependency.feedDom.EXPECT().GetSwipeByParams(gomock.Any(), gomock.Any()).Times(1).Return(model.SwipeModel{}, errors.New("db error"))
		mockedDependency.sqlTxDom.EXPECT().RollbackTX(gomock.Any()).Times(1)

		res, err := feedUsecase.SwipeFeed(ctx, params)
		assert.Error(t, err)
		assert.NotEmpty(t, res)
	})
}
