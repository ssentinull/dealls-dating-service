package auth_test

import (
	"testing"

	authUsecase "github.com/ssentinull/dealls-dating-service/internal/business/usecase/auth"
	mockUserDomain "github.com/ssentinull/dealls-dating-service/internal/mocks/business/domain/user"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/logger"
	mockAuth "github.com/ssentinull/dealls-dating-service/pkg/stdlib/tests/mock/auth"
	mockParser "github.com/ssentinull/dealls-dating-service/pkg/stdlib/tests/mock/parser"

	"go.uber.org/mock/gomock"
)

type MockedDependency struct {
	ctrl       *gomock.Controller
	userDom    *mockUserDomain.MockDomainItf
	efLogger   logger.Logger
	jsonParser *mockParser.MockJSONParser
	auth       *mockAuth.MockAuth
	opt        authUsecase.Options
}

func NewMockedDependency(t *testing.T, ctrl *gomock.Controller) *MockedDependency {
	return &MockedDependency{
		ctrl:       ctrl,
		userDom:    mockUserDomain.NewMockDomainItf(ctrl),
		efLogger:   logger.Init(),
		jsonParser: mockParser.NewMockJSONParser(ctrl),
		auth:       mockAuth.NewMockAuth(ctrl),
		opt:        authUsecase.Options{},
	}
}
