package auth

import (
	"context"

	"github.com/ssentinull/dealls-dating-service/internal/business/domain/user"
	"github.com/ssentinull/dealls-dating-service/internal/business/model"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/auth"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/logger"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/parser"
)

type UsecaseItf interface {
	SignupUser(ctx context.Context, params model.SignupUserParams) (model.UserModel, error)
	LoginUser(ctx context.Context, params model.LoginUserParams) (model.JWTModel, error)
}

type authUc struct {
	userDom    user.DomainItf
	efLogger   logger.Logger
	jsonParser parser.JSONParser
	auth       auth.Auth
	opt        Options
}

type Options struct{}

func InitAuthUsecase(
	userDom user.DomainItf,
	efLogger logger.Logger,
	jsonParser parser.JSONParser,
	auth auth.Auth,
	opt Options,
) UsecaseItf {
	return &authUc{
		userDom:    userDom,
		efLogger:   efLogger,
		jsonParser: jsonParser,
		auth:       auth,
		opt:        opt,
	}
}
