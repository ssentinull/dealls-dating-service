package domain

import (
	"github.com/ssentinull/dealls-dating-service/internal/business/domain/user"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/cache"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/libsql"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/logger"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/parser"
)

type Domain struct {
	User user.DomainItf
}

type Options struct {
	User user.Options
}

func Init(
	efLogger logger.Logger,
	parser parser.Parser,
	sqlClient libsql.SQL,
	redisClient cache.Redis,
	opt Options,
) *Domain {
	return &Domain{
		User: user.InitUserDomain(
			efLogger,
			parser.JSONParser(),
			sqlClient,
			redisClient,
			opt.User,
		),
	}
}
