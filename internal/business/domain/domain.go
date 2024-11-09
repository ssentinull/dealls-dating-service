package domain

import (
	"github.com/ssentinull/golang-boilerplate/internal/business/domain/book"
	"github.com/ssentinull/golang-boilerplate/pkg/stdlib/cache"
	"github.com/ssentinull/golang-boilerplate/pkg/stdlib/libsql"
	"github.com/ssentinull/golang-boilerplate/pkg/stdlib/logger"
	"github.com/ssentinull/golang-boilerplate/pkg/stdlib/parser"
)

type Domain struct {
	Book book.DomainItf
}

type Options struct {
	Book book.Options
}

func Init(
	efLogger logger.Logger,
	parser parser.Parser,
	sqlClient libsql.SQL,
	redisClient cache.Redis,
	opt Options,
) *Domain {
	return &Domain{
		Book: book.InitBookDomain(
			efLogger,
			parser.JSONParser(),
			sqlClient,
			redisClient,
			opt.Book,
		),
	}
}
