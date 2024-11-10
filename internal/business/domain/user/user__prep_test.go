package user_test

import (
	"testing"

	"github.com/ssentinull/dealls-dating-service/internal/business/domain/sqltx"
	"github.com/ssentinull/dealls-dating-service/internal/business/domain/user"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/cache"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/libsql"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/logger"
	mockParser "github.com/ssentinull/dealls-dating-service/pkg/stdlib/tests/mock/parser"

	"go.uber.org/mock/gomock"
)

type MockedDependency struct {
	ctrl     *gomock.Controller
	efLogger logger.Logger
	json     *mockParser.MockJSONParser
	cache    cache.Redis
	sql      libsql.SQL
	opt      user.Options
	tx       sqltx.DomainItf
}

func NewMockedDependency(t *testing.T, ctrl *gomock.Controller) *MockedDependency {
	log := logger.Init()
	sql := libsql.Init(log, libsql.Options{Leader: libsql.Config{Mock: true}, Follower: libsql.Config{Mock: true}})
	jsonparser := mockParser.NewMockJSONParser(ctrl)

	return &MockedDependency{
		ctrl:     ctrl,
		json:     jsonparser,
		sql:      sql,
		efLogger: log,
		opt:      user.Options{},
		cache:    cache.Init(log, cache.Options{Mock: true}),
		tx:       sqltx.InitSQLTXDomain(log, sql, sqltx.Options{}),
	}
}
