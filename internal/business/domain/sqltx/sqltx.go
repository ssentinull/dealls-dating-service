package sqltx

import (
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/libsql"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/logger"
	"gorm.io/gorm"
)

type DomainItf interface {
	BeginTX() *gorm.DB
	CommitTX(tx *gorm.DB) error
	RollbackTX(tx *gorm.DB) error
}

type sqltxImpl struct {
	efLogger logger.Logger
	sql      libsql.SQL
	opt      Options
}

type Options struct{}

func InitSQLTXDomain(efLogger logger.Logger, sql libsql.SQL, opt Options) DomainItf {
	return &sqltxImpl{
		efLogger: efLogger,
		sql:      sql,
		opt:      opt,
	}
}
