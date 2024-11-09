package auth

import (
	"context"

	"github.com/ssentinull/dealls-dating-service/internal/business/domain/user"
	"github.com/ssentinull/dealls-dating-service/internal/business/model"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/logger"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/parser"
)

type UsecaseItf interface {
	SignupUser(ctx context.Context, params model.SignupUserParams) (model.UserModel, error)
}

type authUc struct {
	userDom    user.DomainItf
	efLogger   logger.Logger
	jsonParser parser.JSONParser
	opt        Options
}

type Options struct{}

func InitAuthUsecase(userDom user.DomainItf, efLogger logger.Logger, jsonParser parser.JSONParser, opt Options) UsecaseItf {
	return &authUc{
		userDom:    userDom,
		efLogger:   efLogger,
		jsonParser: jsonParser,
		opt:        opt,
	}
}
