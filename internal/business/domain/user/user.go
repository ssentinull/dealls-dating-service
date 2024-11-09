package user

import (
	"context"

	"github.com/ssentinull/dealls-dating-service/internal/business/model"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/cache"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/libsql"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/logger"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/parser"
	"gorm.io/gorm"
)

type DomainItf interface {
	CreateUser(ctx context.Context, tx *gorm.DB, p model.UserModel) (model.UserModel, error)
	GetUserByParams(ctx context.Context, p model.GetUserParams) (model.UserModel, error)
}

type userImpl struct {
	efLogger logger.Logger
	json     parser.JSONParser
	cache    cache.Redis
	sql      libsql.SQL
	opt      Options
}

type Options struct{}

func InitUserDomain(efLogger logger.Logger, json parser.JSONParser, sql libsql.SQL, redisClient cache.Redis, opt Options) DomainItf {
	return &userImpl{
		efLogger: efLogger,
		json:     json,
		cache:    redisClient,
		sql:      sql,
		opt:      opt,
	}
}
