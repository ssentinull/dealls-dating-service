package feed

import (
	"context"
	"errors"
	"net/http"

	"github.com/ssentinull/dealls-dating-service/internal/business/model"
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

	preference := model.PreferenceModel{
		UserId:   params.UserId,
		Gender:   string(params.Body.Gender),
		MinAge:   params.Body.MinAge,
		MaxAge:   params.Body.MaxAge,
		Location: string(params.Body.Location),
	}

	res, err := f.feedDom.CreatePreference(ctx, nil, preference)
	if err != nil {
		f.efLogger.Error(err)
		return model.PreferenceModel{}, err
	}

	return res, nil
}
