package feed

import (
	"context"

	"github.com/ssentinull/dealls-dating-service/internal/business/model"
	"github.com/ssentinull/dealls-dating-service/internal/types"
	x "github.com/ssentinull/dealls-dating-service/pkg/stdlib/stacktrace"

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
	result, err := f.getPreferenceSQL(ctx, p)
	if err != nil {
		err = x.WrapPassCode(err, "getPreferenceSQL")
		return result, err
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
