package feed_test

import (
	"testing"

	feedUsecase "github.com/ssentinull/dealls-dating-service/internal/business/usecase/feed"
	mockFeedDomain "github.com/ssentinull/dealls-dating-service/internal/mocks/business/domain/feed"
	mockSqlTxDomain "github.com/ssentinull/dealls-dating-service/internal/mocks/business/domain/sqltx"
	mockUserDomain "github.com/ssentinull/dealls-dating-service/internal/mocks/business/domain/user"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/logger"
	mockAuth "github.com/ssentinull/dealls-dating-service/pkg/stdlib/tests/mock/auth"
	mockParser "github.com/ssentinull/dealls-dating-service/pkg/stdlib/tests/mock/parser"

	"go.uber.org/mock/gomock"
)

type MockedDependency struct {
	ctrl       *gomock.Controller
	sqlTxDom   *mockSqlTxDomain.MockDomainItf
	feedDom    *mockFeedDomain.MockDomainItf
	userDom    *mockUserDomain.MockDomainItf
	efLogger   logger.Logger
	jsonParser *mockParser.MockJSONParser
	auth       *mockAuth.MockAuth
	opt        feedUsecase.Options
}

func NewMockedDependency(t *testing.T, ctrl *gomock.Controller) *MockedDependency {
	return &MockedDependency{
		ctrl:       ctrl,
		sqlTxDom:   mockSqlTxDomain.NewMockDomainItf(ctrl),
		feedDom:    mockFeedDomain.NewMockDomainItf(ctrl),
		userDom:    mockUserDomain.NewMockDomainItf(ctrl),
		efLogger:   logger.Init(),
		jsonParser: mockParser.NewMockJSONParser(ctrl),
		auth:       mockAuth.NewMockAuth(ctrl),
		opt: feedUsecase.Options{
			DailySwipeThreshold: 2,
		},
	}
}
