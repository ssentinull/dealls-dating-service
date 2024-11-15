package restserver

import (
	"net/http"

	"github.com/ssentinull/dealls-dating-service/internal/business/model"
	authTypes "github.com/ssentinull/dealls-dating-service/internal/types/auth"
	x "github.com/ssentinull/dealls-dating-service/pkg/stdlib/stacktrace"

	"github.com/gin-gonic/gin"
)

func (r *rest) SignupUser(c *gin.Context) {
	params := authTypes.NewSignupUserParams()
	if err := r.parser.Validator().BindAndValidateBody(c, &params); err != nil {
		r.ResponseError(c, err, http.StatusBadRequest, "params validation error")
		return
	}

	result, err := r.uc.Auth.SignupUser(c.Request.Context(), model.SignupUserParams{SignupUserParams: params})
	if err != nil {
		r.ResponseError(c, x.WrapWithCode(err, x.GetCode(err), "SignupUser"), "Signup User")
		return
	}

	r.responseSuccess(c, http.StatusCreated, result, "signup user successfully")
}

func (r *rest) LoginUser(c *gin.Context) {
	params := authTypes.NewLoginUserParams()
	if err := r.parser.Validator().BindAndValidateBody(c, &params); err != nil {
		r.ResponseError(c, err, http.StatusBadRequest, "params validation error")
		return
	}

	result, err := r.uc.Auth.LoginUser(c.Request.Context(), model.LoginUserParams{LoginUserParams: params})
	if err != nil {
		r.ResponseError(c, x.WrapWithCode(err, x.GetCode(err), "LoginUser"), "Login User")
		return
	}

	r.responseSuccess(c, http.StatusOK, result, "login user successfully")
}
