package book

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
	CreateBook(ctx context.Context, tx *gorm.DB, p model.BookModel) (model.BookModel, error)
}

type bookImpl struct {
	efLogger logger.Logger
	json     parser.JSONParser
	cache    cache.Redis
	sql      libsql.SQL
	opt      Options
}

type Options struct{}

func InitBookDomain(efLogger logger.Logger, json parser.JSONParser, sql libsql.SQL, redisClient cache.Redis, opt Options) DomainItf {
	return &bookImpl{
		efLogger: efLogger,
		json:     json,
		cache:    redisClient,
		sql:      sql,
		opt:      opt,
	}
}
