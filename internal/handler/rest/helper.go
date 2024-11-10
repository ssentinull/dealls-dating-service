package restserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/ssentinull/dealls-dating-service/internal/business/model"
	"github.com/ssentinull/dealls-dating-service/internal/handler/rest/mapper"
	"github.com/ssentinull/dealls-dating-service/internal/types"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/stacktrace"

	"github.com/gin-gonic/gin"
	errorx "github.com/go-openapi/errors"
	"github.com/sirupsen/logrus"
)

const (
	// HTTPErrorTypeGeneric represents the generic error type returned as default for all HTTP errors without a type defined.
	HTTPErrorTypeGeneric    string = "generic"
	HTTPErrorTypeValidation string = "ValidationError"
)

// Payload in accordance with RFC 7807 (Problem Details for HTTP APIs) with the exception of the type
// value not being represented by a URI. https://tools.ietf.org/html/rfc7807 @ 2020-04-27T15:44:37Z

type HTTPError struct {
	types.PublicHTTPError
	Internal       error                  `json:"-"`
	AdditionalData map[string]interface{} `json:"-"`
}

type HTTPValidationError struct {
	types.PublicHTTPValidationError
	Internal       error                  `json:"-"`
	AdditionalData map[string]interface{} `json:"-"`
}

func NewHTTPValidationError(code int, errorType string, title string, validationErrors []*types.HTTPValidationErrorDetail) *HTTPValidationError {
	return &HTTPValidationError{
		PublicHTTPValidationError: types.PublicHTTPValidationError{
			Success: false,
			Message: title,
			Data: &types.PublicHTTPValidationErrorData{
				Code:             int64(code),
				Type:             errorType,
				ValidationErrors: validationErrors,
			},
		},
	}
}

func NewHTTPValidationErrorWithDetail(code int, errorType string, title string, validationErrors []*types.HTTPValidationErrorDetail, detail string) *HTTPValidationError {
	return &HTTPValidationError{
		PublicHTTPValidationError: types.PublicHTTPValidationError{
			Success: false,
			Message: title,
			Data: &types.PublicHTTPValidationErrorData{
				Code:             int64(code),
				Type:             errorType,
				Detail:           detail,
				ValidationErrors: validationErrors,
			},
		},
	}
}

func (e *HTTPValidationError) Error() string {
	var b strings.Builder

	fmt.Fprintf(&b, "HTTPValidationError %d (%s): %s", e.Data.Code, e.Data.Type, e.Message)

	if len(e.Data.Detail) > 0 {
		fmt.Fprintf(&b, " - %s", e.Data.Detail)
	}
	if e.Internal != nil {
		fmt.Fprintf(&b, ", %v", e.Internal)
	}
	if e.AdditionalData != nil && len(e.AdditionalData) > 0 {
		keys := make([]string, 0, len(e.AdditionalData))
		for k := range e.AdditionalData {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		b.WriteString(". Additional: ")
		for i, k := range keys {
			fmt.Fprintf(&b, "%s=%v", k, e.AdditionalData[k])
			if i < len(keys)-1 {
				b.WriteString(", ")
			}
		}
	}

	b.WriteString(" - Validation: ")
	for i, ve := range e.Data.ValidationErrors {
		fmt.Fprintf(&b, "%s (in %s): %s", ve.Key, ve.In, ve.Error)
		if i < len(e.Data.ValidationErrors)-1 {
			b.WriteString(", ")
		}
	}

	return b.String()
}

func utilFormatValidationErrors(ctx context.Context, err *errorx.CompositeError) []*types.HTTPValidationErrorDetail {
	valErrs := make([]*types.HTTPValidationErrorDetail, 0, len(err.Errors))
	for _, e := range err.Errors {
		switch ee := e.(type) {
		case *errorx.ParseError:
			valErrs = append(valErrs, &types.HTTPValidationErrorDetail{
				Key:   ee.Name,
				In:    ee.In,
				Error: ee.Error(),
			})
		case *errorx.Validation:
			valErrs = append(valErrs, &types.HTTPValidationErrorDetail{
				Key:   ee.Name,
				In:    ee.In,
				Error: ee.Error(),
			})
		case *errorx.CompositeError:
			valErrs = append(valErrs, utilFormatValidationErrors(ctx, ee)...)
		default:
			logrus.Error(err)
			valErrs = append(valErrs, &types.HTTPValidationErrorDetail{
				Error: ee.Error(),
			})
		}
	}

	return valErrs
}

func formatValidationErrors(c *gin.Context, unformatedErr error) error {
	switch ufErr := unformatedErr.(type) {
	case *errorx.CompositeError:
		valErrs := utilFormatValidationErrors(c.Request.Context(), ufErr)

		return NewHTTPValidationError(http.StatusBadRequest, HTTPErrorTypeValidation, http.StatusText(http.StatusBadRequest), valErrs)
	case *errorx.Validation:
		valErrs := []*types.HTTPValidationErrorDetail{
			{
				Key:   ufErr.Name,
				In:    ufErr.In,
				Error: ufErr.Error(),
			},
		}

		return NewHTTPValidationError(http.StatusBadRequest, HTTPErrorTypeValidation, http.StatusText(http.StatusBadRequest), valErrs)
	}

	return unformatedErr
}

func (r *rest) ResponseError(c *gin.Context, err error, opts ...interface{}) {
	err = formatValidationErrors(c, err)
	t := "InternalServerError"
	displayError := ""

	code := http.StatusInternalServerError

	for _, v := range opts {
		if tErr, ok := v.(string); ok {
			if strings.Contains(tErr, " ") {
				displayError = tErr
			} else {
				t = tErr
			}
		}
		if cErr, ok := v.(int); ok && cErr >= 100 && cErr <= 599 {
			code = cErr
		}
	}

	vErr := new(HTTPValidationError)
	if errors.As(err, &vErr) {
		c.AbortWithStatusJSON(http.StatusBadRequest, vErr.PublicHTTPValidationError)
		return
	}

	code, appErr := stacktrace.Compile(err, r.opt.DebugMode)
	displayError = appErr.Message
	detailError := ""
	if r.opt.DebugMode {
		detailError = *appErr.DebugError
	}

	t = http.StatusText(code)
	r.efLogger.Error(appErr.Error())

	c.AbortWithStatusJSON(code, types.PublicHTTPError{
		Success: false,
		Message: displayError,
		Data: &types.PublicHTTPErrorData{
			Type:   t,
			Code:   int64(code),
			Detail: detailError,
		},
	})
}

func (r *rest) responseSuccess(c *gin.Context, code int, resp interface{}, opts ...interface{}) {
	var (
		raw        []byte
		err        error
		message    string
		pagination *types.Pagination
	)

	for _, v := range opts {
		switch val := v.(type) {
		case int:
			code = val
		case string:
			message = val
		case *types.Pagination:
			baseURL := c.FullPath()
			if val.NextPage != nil {
				nextQuery := c.Request.URL.Query()
				nextQuery.Set("page", strconv.Itoa(int(*val.NextPage)))
				val.NextURL = fmt.Sprintf("%s?%s", baseURL, nextQuery.Encode())

				// purge NextPage Number to omit response the page number
				val.NextPage = nil
			}
			if val.PrevPage != nil {
				prevQuery := c.Request.URL.Query()
				prevQuery.Set("page", strconv.Itoa(int(*val.PrevPage)))
				val.PrevURL = fmt.Sprintf("%s?%s", baseURL, prevQuery.Encode())

				// purge PrevPage Number to omit response the page number
				val.PrevPage = nil
			}

			pagination = val
		}
	}

	switch data := resp.(type) {
	case model.UserModel:
		resultType := mapper.MapUserModelToUserType(data)
		res := types.SignupUserResponse{
			Success: true,
			Message: message,
			Data:    resultType,
		}

		raw, err = json.Marshal(res)
	case model.JWTModel:
		resultType := mapper.MapJWTModelToJWTType(data)
		res := types.LoginUserResponse{
			Success: true,
			Message: message,
			Data:    resultType,
		}

		raw, err = json.Marshal(res)
	case model.PreferenceModel:
		resultType := mapper.MapPreferenceModelToPreferenceType(data)
		res := types.CreateFeedPreferenceResponse{
			Success: true,
			Message: message,
			Data:    resultType,
		}

		raw, err = json.Marshal(res)
	case []model.FeedModel:
		resultTypes := make([]*types.Feed, 0)
		for _, d := range data {
			resultTypes = append(resultTypes, mapper.MapFeedModelToFeedType(d))
		}

		res := types.GetFeedResponse{
			Success:    true,
			Message:    message,
			Data:       resultTypes,
			Pagination: pagination,
		}

		raw, err = json.Marshal(res)
	default:
		r.ResponseError(c, errors.New("unknown response type"), fmt.Sprintf("cannot cast type of %+v", data))
		return
	}

	if err != nil {
		r.ResponseError(c, err, http.StatusInternalServerError, "Marshall HTTP Response")
		return
	}

	c.Data(code, "application/json", raw)
}
