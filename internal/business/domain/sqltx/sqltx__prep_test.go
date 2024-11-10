package sqltx_test

import (
	"testing"

	"github.com/ssentinull/dealls-dating-service/internal/business/domain/sqltx"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/libsql"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/logger"
)

type MockedDependency struct {
	efLogger logger.Logger
	sql      libsql.SQL
	opt      sqltx.Options
}

func NewMockedDependency(t *testing.T) *MockedDependency {
	log := logger.Init()
	return &MockedDependency{
		sql:      libsql.Init(log, libsql.Options{Leader: libsql.Config{Mock: true}, Follower: libsql.Config{Mock: true}}),
		efLogger: log,
		opt:      sqltx.Options{},
	}
}
