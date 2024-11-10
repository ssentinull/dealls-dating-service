package restserver

import (
	"net/http"

	"github.com/ssentinull/dealls-dating-service/internal/business/model"
	feedTypes "github.com/ssentinull/dealls-dating-service/internal/types/feed"
	x "github.com/ssentinull/dealls-dating-service/pkg/stdlib/stacktrace"

	"github.com/gin-gonic/gin"
)

func (r *rest) CreatePreference(c *gin.Context) {
	params := feedTypes.NewCreateFeedPreferenceParams()
	if err := r.parser.Validator().BindAndValidateBody(c, &params); err != nil {
		r.ResponseError(c, err, http.StatusBadRequest, "params validation error")
		return
	}

	createPreferenceParam := model.CreatePreferenceParams{
		UserId:                     r.auth.GetID(c),
		CreateFeedPreferenceParams: params,
	}

	result, err := r.uc.Feed.CreatePreference(c.Request.Context(), createPreferenceParam)
	if err != nil {
		r.ResponseError(c, x.WrapWithCode(err, x.GetCode(err), "CreatePreference"), "Create Preference")
		return
	}

	r.responseSuccess(c, http.StatusCreated, result, "create preference successfully")
}

func (r *rest) GetFeed(c *gin.Context) {
	query := feedTypes.NewGetFeedParams()
	if err := r.parser.Validator().BindAndValidateBody(c, &query); err != nil {
		r.ResponseError(c, err, http.StatusBadRequest, "params validation error")
		return
	}

	param := model.GetFeedParams{
		UserId:        r.auth.GetID(c),
		GetFeedParams: query,
	}

	result, pagination, err := r.uc.Feed.GetFeed(c.Request.Context(), param)
	if err != nil {
		r.ResponseError(c, x.WrapWithCode(err, x.GetCode(err), "GetFeed"), "Get Feed")
		return
	}

	r.responseSuccess(c, http.StatusOK, result, pagination, "get feed successfully")
}
