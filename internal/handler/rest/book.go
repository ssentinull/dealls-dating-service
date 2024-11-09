package restserver

import (
	"net/http"

	"github.com/ssentinull/dealls-dating-service/internal/business/model"
	bookTypes "github.com/ssentinull/dealls-dating-service/internal/types/book"
	x "github.com/ssentinull/dealls-dating-service/pkg/stdlib/stacktrace"

	"github.com/gin-gonic/gin"
)

func (r *rest) CreateBook(c *gin.Context) {
	params := bookTypes.NewCreateBookParams()
	if err := r.parser.Validator().BindAndValidateBody(c, &params); err != nil {
		r.ResponseError(c, err, http.StatusBadRequest, "params validation error")
		return
	}

	result, err := r.uc.Book.CreateBook(c.Request.Context(), model.CreateBookParams{CreateBookParams: params})
	if err != nil {
		r.ResponseError(c, x.WrapWithCode(err, x.GetCode(err), "CreateBook"), "Creat Book")
		return
	}

	r.responseSuccess(c, http.StatusOK, result, "create book successfully")
}
