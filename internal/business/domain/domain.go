package domain

import (
	"github.com/ssentinull/dealls-dating-service/internal/business/domain/feed"
	"github.com/ssentinull/dealls-dating-service/internal/business/domain/sqltx"
	"github.com/ssentinull/dealls-dating-service/internal/business/domain/user"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/cache"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/libsql"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/logger"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/parser"
)

type Domain struct {
	SqlTx sqltx.DomainItf
	User  user.DomainItf
	Feed  feed.DomainItf
}

type Options struct {
	SqlTx sqltx.Options
	User  user.Options
	Feed  feed.Options
}

func Init(
	efLogger logger.Logger,
	parser parser.Parser,
	sqlClient libsql.SQL,
	redisClient cache.Redis,
	opt Options,
) *Domain {
	return &Domain{
		SqlTx: sqltx.InitSQLTXDomain(
			efLogger,
			sqlClient,
			opt.SqlTx,
		),
		User: user.InitUserDomain(
			efLogger,
			parser.JSONParser(),
			sqlClient,
			redisClient,
			opt.User,
		),
		Feed: feed.InitFeedDomain(
			efLogger,
			parser.JSONParser(),
			sqlClient,
			redisClient,
			opt.Feed,
		),
	}
}
