package feed

import (
	"context"

	"github.com/ssentinull/dealls-dating-service/internal/business/model"
	"github.com/ssentinull/dealls-dating-service/internal/types"
	"github.com/ssentinull/dealls-dating-service/pkg/common"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/libsql"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/libsql/utils"
	x "github.com/ssentinull/dealls-dating-service/pkg/stdlib/stacktrace"

	"gorm.io/gorm"
)

func (f *feedImpl) createPreferenceSQL(ctx context.Context, tx *gorm.DB, p model.PreferenceModel) (model.PreferenceModel, error) {
	if tx == nil {
		tx = f.sql.Leader()
	}

	var result model.PreferenceModel
	if err := tx.WithContext(ctx).Create(&p).Error; err != nil {
		f.efLogger.Error(err)
		return result, x.Wrap(err, libsql.SomethingWentWrongWithDB)
	}

	result = p

	return result, nil
}

func (f *feedImpl) getPreferenceSQL(ctx context.Context, p model.GetPreferenceParams) (model.PreferenceModel, error) {
	db := f.sql.Leader()

	var result model.PreferenceModel
	err := db.WithContext(ctx).
		Where("user_id = ?", p.UserId).
		Take(&result).
		Error
	if err != nil {
		f.efLogger.Error(err)
		return result, err
	}

	return result, nil
}

func (f *feedImpl) getFeedSQL(ctx context.Context, p model.GetFeedParams) ([]model.FeedModel, *types.Pagination, error) {

	p.PaginationPage = utils.ValidatePage(p.PaginationPage)
	p.PaginationSize = utils.ValidateSize(p.PaginationSize)

	page := common.ToValue(p.PaginationPage)
	size := common.ToValue(p.PaginationSize)
	offset := int((page - 1) * size)

	readQuery := `
	WITH potential_matches AS (
		SELECT 
			u.id,
			u.name,
			u.gender,
			u.location,
			u.created_at,
			DATE_PART('year', AGE(u.birth_date)) AS age
		FROM users u
		WHERE DATE_PART('year', AGE(u.birth_date)) <= ?
			AND DATE_PART('year', AGE(u.birth_date)) >= ?
			AND u.id != ?
			AND u.deleted_at IS NULL
	),
	swiped_today AS (
		SELECT s.to_user_id
		FROM swipes s
		WHERE s.from_user_id = ?
			AND s.created_at::DATE = CURRENT_DATE
	)
	SELECT pm.*
	FROM potential_matches pm
	LEFT JOIN swiped_today st 
		ON pm.id = st.to_user_id
	WHERE st.to_user_id IS NULL
	ORDER BY pm.id DESC
	OFFSET ?
	LIMIT ?;
	`

	result := make([]model.FeedModel, 0)
	db := f.sql.Leader().WithContext(ctx)

	err := db.Raw(readQuery, p.MaxAge, p.MinAge, p.UserId, p.UserId, offset, size).Scan(&result).Error
	if err != nil {
		f.efLogger.Error(err)
		return result, nil, x.Wrap(err, libsql.SomethingWentWrongWithDB)
	}

	countQuery := `
	WITH potential_matches AS (
		SELECT 
			u.id,
			u.name,
			u.gender,
			u.location,
			DATE_PART('year', AGE(u.birth_date)) AS age
		FROM users u
		WHERE DATE_PART('year', AGE(u.birth_date)) <= ?
			AND DATE_PART('year', AGE(u.birth_date)) >= ?
			AND u.id != ?
			AND u.deleted_at IS NULL
	),
	swiped_today AS (
		SELECT s.to_user_id
		FROM swipes s
		WHERE s.from_user_id = ?
			AND s.created_at::DATE = CURRENT_DATE
	)
	SELECT COUNT(*)
	FROM potential_matches pm
	LEFT JOIN swiped_today st 
		ON pm.id = st.to_user_id
	WHERE st.to_user_id IS NULL;
	`

	totalData := int64(0)
	err = db.Raw(countQuery, p.MaxAge, p.MinAge, p.UserId, p.UserId).Scan(&totalData).Error
	if err != nil {
		f.efLogger.Error(err)
		return result, nil, x.Wrap(err, libsql.SomethingWentWrongWithDB)
	}

	pagination := common.NewPagination(totalData, int64(len(result)), size, page)

	return result, &pagination, nil
}

func (f *feedImpl) createSwipeSQL(ctx context.Context, tx *gorm.DB, p model.SwipeModel) (model.SwipeModel, error) {
	if tx == nil {
		tx = f.sql.Leader()
	}

	var result model.SwipeModel
	if err := tx.WithContext(ctx).Create(&p).Error; err != nil {
		f.efLogger.Error(err)
		return result, x.Wrap(err, libsql.SomethingWentWrongWithDB)
	}

	result = p

	return result, nil
}

func (f *feedImpl) getSwipeSQL(ctx context.Context, p model.GetSwipeParams) (model.SwipeModel, error) {
	db := f.sql.Leader().WithContext(ctx)

	if p.FromUserId > 0 {
		db = db.Where("from_user_id = ?", p.FromUserId)
	}

	if p.ToUserId > 0 {
		db = db.Where("to_user_id = ?", p.ToUserId)
	}

	if !p.CreatedAt.IsZero() {
		db = db.Where("created_at::DATE = ?", p.CreatedAt.Format("2006-01-02"))
	}

	var result model.SwipeModel
	if err := db.Take(&result).Error; err != nil {
		f.efLogger.Error(err)
		return result, err
	}

	return result, nil
}

func (f *feedImpl) createMatchSQL(ctx context.Context, tx *gorm.DB, p model.MatchModel) (model.MatchModel, error) {
	if tx == nil {
		tx = f.sql.Leader()
	}

	var result model.MatchModel
	if err := tx.WithContext(ctx).Create(&p).Error; err != nil {
		f.efLogger.Error(err)
		return result, x.Wrap(err, libsql.SomethingWentWrongWithDB)
	}

	result = p

	return result, nil
}
