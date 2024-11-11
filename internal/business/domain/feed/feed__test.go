package feed_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/c2fo/testify/assert"
	"github.com/ssentinull/dealls-dating-service/internal/business/domain/feed"
	"github.com/ssentinull/dealls-dating-service/internal/business/model"
	"go.uber.org/mock/gomock"
)

func TestFeedDomain_CreatePreference(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	defer func() {
		ctx.Done()
		ctrl.Finish()
	}()

	mockedDependency := NewMockedDependency(t, ctrl)
	feedDomain := feed.InitFeedDomain(
		mockedDependency.efLogger,
		mockedDependency.json,
		mockedDependency.sql,
		mockedDependency.cache,
		mockedDependency.opt,
	)

	ID := 1
	userID := int64(2)
	preference := model.PreferenceModel{
		UserId: userID,
		Gender: "FEMALE",
		MinAge: 27,
		MaxAge: 30,
	}

	createPreferenceSQL := `INSERT INTO "preferences" ("user_id","gender","min_age","max_age","location","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "id"`

	t.Run("success", func(t *testing.T) {
		mockedDependency.sql.MockLeader().ExpectBegin()
		mockedDependency.sql.MockLeader().ExpectQuery(regexp.QuoteMeta(createPreferenceSQL)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(ID))
		mockedDependency.sql.MockLeader().ExpectCommit()

		res, err := feedDomain.CreatePreference(ctx, nil, preference)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
		assert.Nil(t, mockedDependency.sql.MockLeader().ExpectationsWereMet())
	})

	t.Run("failed", func(t *testing.T) {
		mockedDependency.sql.MockLeader().ExpectBegin()
		mockedDependency.sql.MockLeader().ExpectQuery(regexp.QuoteMeta(createPreferenceSQL)).WillReturnError(errors.New("db error"))
		mockedDependency.sql.MockLeader().ExpectRollback()

		res, err := feedDomain.CreatePreference(ctx, nil, preference)
		assert.NotNil(t, err)
		assert.NotEmpty(t, res)
		assert.Nil(t, mockedDependency.sql.MockLeader().ExpectationsWereMet())
	})
}

func TestFeedDomain_GetPreferenceByParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	defer func() {
		ctx.Done()
		ctrl.Finish()
	}()

	mockedDependency := NewMockedDependency(t, ctrl)
	feedDomain := feed.InitFeedDomain(
		mockedDependency.efLogger,
		mockedDependency.json,
		mockedDependency.sql,
		mockedDependency.cache,
		mockedDependency.opt,
	)

	ID := int64(1)
	userID := int64(2)
	param := model.GetPreferenceParams{UserId: userID}

	paramRaw, err := json.Marshal(param)
	assert.NoError(t, err)

	cacheKey := fmt.Sprintf(feed.PreferenceCacheKey, string(paramRaw))
	getPreferenceByParamSQL := `SELECT * FROM "preferences" WHERE user_id = $1 AND "preferences"."deleted_at" IS NULL LIMIT 1`

	t.Run("success", func(t *testing.T) {
		mockedDependency.json.EXPECT().Marshal(param).Times(1).Return(paramRaw, nil)
		mockedDependency.cache.Mock().ExpectGet(cacheKey).SetErr(errors.New("redis error"))
		mockedDependency.sql.MockLeader().ExpectQuery(regexp.QuoteMeta(getPreferenceByParamSQL)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}).AddRow(ID, userID))

		res, err := feedDomain.GetPreferenceByParams(ctx, param)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
		assert.Nil(t, mockedDependency.sql.MockLeader().ExpectationsWereMet())
	})

	t.Run("failed", func(t *testing.T) {
		mockedDependency.json.EXPECT().Marshal(param).Times(1).Return(paramRaw, nil)
		mockedDependency.cache.Mock().ExpectGet(cacheKey).SetErr(errors.New("redis error"))
		mockedDependency.sql.MockLeader().ExpectQuery(regexp.QuoteMeta(getPreferenceByParamSQL)).WillReturnError(errors.New("db error"))

		activity, err := feedDomain.GetPreferenceByParams(ctx, param)
		assert.NotNil(t, err)
		assert.NotEmpty(t, activity)
		assert.Nil(t, mockedDependency.sql.MockLeader().ExpectationsWereMet())
	})
}

func TestFeedDomain_GetFeedByParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	defer func() {
		ctx.Done()
		ctrl.Finish()
	}()

	mockedDependency := NewMockedDependency(t, ctrl)
	feedDomain := feed.InitFeedDomain(
		mockedDependency.efLogger,
		mockedDependency.json,
		mockedDependency.sql,
		mockedDependency.cache,
		mockedDependency.opt,
	)

	ID := int64(1)
	userID := int64(2)

	param := model.GetFeedParams{
		Gender:   "FEMALE",
		Location: "JAKARTA",
		MinAge:   27,
		MaxAge:   30,
	}

	getFeedByParamSQL := `WITH potential_matches AS ( SELECT u.id, u.name, u.gender, u.location, u.created_at, DATE_PART('year', AGE(u.birth_date)) AS age FROM users u WHERE DATE_PART('year', AGE(u.birth_date)) <= $1 AND DATE_PART('year', AGE(u.birth_date)) >= $2 AND u.gender = $3 AND u.location = $4 AND u.id != $5 AND u.deleted_at IS NULL ), swiped_today AS ( SELECT s.to_user_id FROM swipes s WHERE s.from_user_id = $6 AND s.created_at::DATE = CURRENT_DATE ) SELECT pm.* FROM potential_matches pm LEFT JOIN swiped_today st ON pm.id = st.to_user_id WHERE st.to_user_id IS NULL ORDER BY pm.id DESC OFFSET $7 LIMIT $8;`
	countFeedByParamSQL := `WITH potential_matches AS ( SELECT u.id, u.name, u.gender, u.location, DATE_PART('year', AGE(u.birth_date)) AS age FROM users u WHERE DATE_PART('year', AGE(u.birth_date)) <= $1 AND DATE_PART('year', AGE(u.birth_date)) >= $2 AND u.gender = $3 AND u.location = $4 AND u.id != $5 AND u.deleted_at IS NULL ), swiped_today AS ( SELECT s.to_user_id FROM swipes s WHERE s.from_user_id = $6 AND s.created_at::DATE = CURRENT_DATE ) SELECT COUNT(*) FROM potential_matches pm LEFT JOIN swiped_today st ON pm.id = st.to_user_id WHERE st.to_user_id IS NULL;`

	t.Run("success", func(t *testing.T) {
		mockedDependency.sql.MockLeader().ExpectQuery(regexp.QuoteMeta(getFeedByParamSQL)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}).AddRow(ID, userID))
		mockedDependency.sql.MockLeader().ExpectQuery(regexp.QuoteMeta(countFeedByParamSQL)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		resFeed, resPagination, err := feedDomain.GetFeedByParams(ctx, param)
		assert.NoError(t, err)
		assert.NotEmpty(t, resFeed)
		assert.NotNil(t, resPagination)
		assert.Nil(t, mockedDependency.sql.MockLeader().ExpectationsWereMet())
	})

	t.Run("failed - get feed query return error", func(t *testing.T) {
		mockedDependency.sql.MockLeader().ExpectQuery(regexp.QuoteMeta(getFeedByParamSQL)).
			WillReturnError(errors.New("db error"))

		resFeed, resPagination, err := feedDomain.GetFeedByParams(ctx, param)
		assert.NotNil(t, err)
		assert.Empty(t, resFeed)
		assert.Nil(t, resPagination)
		assert.Nil(t, mockedDependency.sql.MockLeader().ExpectationsWereMet())
	})

	t.Run("failed - count feed query return error", func(t *testing.T) {
		mockedDependency.sql.MockLeader().ExpectQuery(regexp.QuoteMeta(getFeedByParamSQL)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}).AddRow(ID, userID))
		mockedDependency.sql.MockLeader().ExpectQuery(regexp.QuoteMeta(countFeedByParamSQL)).WillReturnError(errors.New("db error"))

		resFeed, resPagination, err := feedDomain.GetFeedByParams(ctx, param)
		assert.NotNil(t, err)
		assert.Empty(t, resFeed)
		assert.Nil(t, resPagination)
		assert.Nil(t, mockedDependency.sql.MockLeader().ExpectationsWereMet())
	})
}

func TestFeedDomain_CreateSwipe(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	defer func() {
		ctx.Done()
		ctrl.Finish()
	}()

	mockedDependency := NewMockedDependency(t, ctrl)
	feedDomain := feed.InitFeedDomain(
		mockedDependency.efLogger,
		mockedDependency.json,
		mockedDependency.sql,
		mockedDependency.cache,
		mockedDependency.opt,
	)

	ID := 1
	fromUserID := int64(2)
	toUserID := int64(3)

	swipe := model.SwipeModel{
		FromUserId: fromUserID,
		ToUserId:   toUserID,
		SwipeType:  "RIGHT",
	}

	createSwipeSQL := `INSERT INTO "swipes" ("from_user_id","to_user_id","swipe_type","created_at") VALUES ($1,$2,$3,$4) RETURNING "id"`

	t.Run("success", func(t *testing.T) {
		mockedDependency.sql.MockLeader().ExpectBegin()
		mockedDependency.sql.MockLeader().ExpectQuery(regexp.QuoteMeta(createSwipeSQL)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(ID))
		mockedDependency.sql.MockLeader().ExpectCommit()

		res, err := feedDomain.CreateSwipe(ctx, nil, swipe)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
		assert.Nil(t, mockedDependency.sql.MockLeader().ExpectationsWereMet())
	})

	t.Run("failed", func(t *testing.T) {
		mockedDependency.sql.MockLeader().ExpectBegin()
		mockedDependency.sql.MockLeader().ExpectQuery(regexp.QuoteMeta(createSwipeSQL)).WillReturnError(errors.New("db error"))
		mockedDependency.sql.MockLeader().ExpectRollback()

		res, err := feedDomain.CreateSwipe(ctx, nil, swipe)
		assert.NotNil(t, err)
		assert.NotEmpty(t, res)
		assert.Nil(t, mockedDependency.sql.MockLeader().ExpectationsWereMet())
	})
}

func TestFeedDomain_GetSwipeByParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	defer func() {
		ctx.Done()
		ctrl.Finish()
	}()

	mockedDependency := NewMockedDependency(t, ctrl)
	feedDomain := feed.InitFeedDomain(
		mockedDependency.efLogger,
		mockedDependency.json,
		mockedDependency.sql,
		mockedDependency.cache,
		mockedDependency.opt,
	)

	ID := int64(1)
	fromUserID := int64(2)
	toUserID := int64(3)

	param := model.GetSwipeParams{
		FromUserId: fromUserID,
		ToUserId:   toUserID,
		SwipeType:  "RIGHT",
	}

	getSwipeByParamsSQL := `SELECT * FROM "swipes" WHERE from_user_id = $1 AND to_user_id = $2 LIMIT 1`

	t.Run("success", func(t *testing.T) {
		mockedDependency.sql.MockLeader().ExpectQuery(regexp.QuoteMeta(getSwipeByParamsSQL)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(ID))

		res, err := feedDomain.GetSwipeByParams(ctx, param)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
		assert.Nil(t, mockedDependency.sql.MockLeader().ExpectationsWereMet())
	})

	t.Run("failed", func(t *testing.T) {
		mockedDependency.sql.MockLeader().ExpectQuery(regexp.QuoteMeta(getSwipeByParamsSQL)).WillReturnError(errors.New("db error"))

		res, err := feedDomain.GetSwipeByParams(ctx, param)
		assert.NotNil(t, err)
		assert.NotEmpty(t, res)
		assert.Nil(t, mockedDependency.sql.MockLeader().ExpectationsWereMet())
	})
}

func TestFeedDomain_GetSwipeCountByUserId(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	defer func() {
		ctx.Done()
		ctrl.Finish()
	}()

	mockedDependency := NewMockedDependency(t, ctrl)
	feedDomain := feed.InitFeedDomain(
		mockedDependency.efLogger,
		mockedDependency.json,
		mockedDependency.sql,
		mockedDependency.cache,
		mockedDependency.opt,
	)

	count := int64(1)
	userID := int64(2)
	cacheKey := fmt.Sprintf(feed.SwipeCountCacheKey, userID)

	t.Run("success", func(t *testing.T) {
		mockedDependency.cache.Mock().ExpectGet(cacheKey).SetVal(strconv.FormatInt(count, 10))

		res, err := feedDomain.GetSwipeCountByUserId(ctx, userID)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
		assert.Nil(t, mockedDependency.sql.MockLeader().ExpectationsWereMet())
	})

	t.Run("failed", func(t *testing.T) {
		mockedDependency.cache.Mock().ExpectGet(cacheKey).SetErr(errors.New("redis error"))

		res, err := feedDomain.GetSwipeCountByUserId(ctx, userID)
		assert.NotNil(t, err)
		assert.Empty(t, res)
		assert.Nil(t, mockedDependency.sql.MockLeader().ExpectationsWereMet())
	})
}

func TestFeedDomain_SetSwipeCountByUserId(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	defer func() {
		ctx.Done()
		ctrl.Finish()
	}()

	mockedDependency := NewMockedDependency(t, ctrl)
	feedDomain := feed.InitFeedDomain(
		mockedDependency.efLogger,
		mockedDependency.json,
		mockedDependency.sql,
		mockedDependency.cache,
		mockedDependency.opt,
	)

	count := int64(1)
	countStr := strconv.FormatInt(count, 10)
	userID := int64(2)
	cacheKey := fmt.Sprintf(feed.SwipeCountCacheKey, userID)

	t.Run("success", func(t *testing.T) {
		mockedDependency.cache.Mock().ExpectSet(cacheKey, countStr, feed.SwipeCountCacheExpiration).SetVal(countStr)

		err := feedDomain.SetSwipeCountByUserId(ctx, userID, count)
		assert.NoError(t, err)
		assert.Nil(t, mockedDependency.sql.MockLeader().ExpectationsWereMet())
	})

	t.Run("failed", func(t *testing.T) {
		mockedDependency.cache.Mock().ExpectSet(cacheKey, countStr, feed.SwipeCountCacheExpiration).SetErr(errors.New("redis error"))

		err := feedDomain.SetSwipeCountByUserId(ctx, userID, count)
		assert.NotNil(t, err)
		assert.Nil(t, mockedDependency.sql.MockLeader().ExpectationsWereMet())
	})
}

func TestFeedDomain_CreateMatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	defer func() {
		ctx.Done()
		ctrl.Finish()
	}()

	mockedDependency := NewMockedDependency(t, ctrl)
	feedDomain := feed.InitFeedDomain(
		mockedDependency.efLogger,
		mockedDependency.json,
		mockedDependency.sql,
		mockedDependency.cache,
		mockedDependency.opt,
	)

	ID := 1
	myUserID := int64(2)
	matchedUserID := int64(3)

	match := model.MatchModel{
		MyUserId:      myUserID,
		MatchedUserId: matchedUserID,
	}

	createMatchSQL := `INSERT INTO "matches" ("my_user_id","matched_user_id","created_at") VALUES ($1,$2,$3) RETURNING "id"`

	t.Run("success", func(t *testing.T) {
		mockedDependency.sql.MockLeader().ExpectBegin()
		mockedDependency.sql.MockLeader().ExpectQuery(regexp.QuoteMeta(createMatchSQL)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(ID))
		mockedDependency.sql.MockLeader().ExpectCommit()

		res, err := feedDomain.CreateMatch(ctx, nil, match)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
		assert.Nil(t, mockedDependency.sql.MockLeader().ExpectationsWereMet())
	})

	t.Run("failed", func(t *testing.T) {
		mockedDependency.sql.MockLeader().ExpectBegin()
		mockedDependency.sql.MockLeader().ExpectQuery(regexp.QuoteMeta(createMatchSQL)).WillReturnError(errors.New("db error"))
		mockedDependency.sql.MockLeader().ExpectRollback()

		res, err := feedDomain.CreateMatch(ctx, nil, match)
		assert.NotNil(t, err)
		assert.NotEmpty(t, res)
		assert.Nil(t, mockedDependency.sql.MockLeader().ExpectationsWereMet())
	})
}
