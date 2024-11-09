package feed

import (
	"context"

	"github.com/ssentinull/dealls-dating-service/internal/business/model"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/libsql"
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
