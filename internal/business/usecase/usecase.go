package usecase

import (
	"github.com/ssentinull/dealls-dating-service/internal/business/domain"
	usecaseAuth "github.com/ssentinull/dealls-dating-service/internal/business/usecase/auth"
	"github.com/ssentinull/dealls-dating-service/internal/business/usecase/feed"
	stdLibAuth "github.com/ssentinull/dealls-dating-service/pkg/stdlib/auth"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/logger"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/parser"
)

type Usecase struct {
	Auth usecaseAuth.UsecaseItf
	Feed feed.UsecaseItf
}

type Options struct {
	Auth usecaseAuth.Options
	Feed feed.Options
}

func Init(
	efLogger logger.Logger,
	parser parser.Parser,
	stdAuth stdLibAuth.Auth,
	dom *domain.Domain,
	opt Options,
) *Usecase {
	return &Usecase{
		Auth: usecaseAuth.InitAuthUsecase(
			dom.User,
			efLogger,
			parser.JSONParser(),
			stdAuth,
			opt.Auth,
		),
		Feed: feed.InitFeedUsecase(
			dom.Feed,
			dom.User,
			efLogger,
			parser.JSONParser(),
			stdAuth,
			opt.Feed,
		),
	}
}
