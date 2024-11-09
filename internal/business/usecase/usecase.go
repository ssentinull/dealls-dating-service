package usecase

import (
	"github.com/ssentinull/dealls-dating-service/internal/business/domain"
	"github.com/ssentinull/dealls-dating-service/internal/business/usecase/auth"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/logger"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/parser"
)

type Usecase struct {
	Auth auth.UsecaseItf
}

type Options struct {
	Auth auth.Options
}

func Init(
	efLogger logger.Logger,
	parser parser.Parser,
	dom *domain.Domain,
	opt Options,
) *Usecase {
	return &Usecase{
		Auth: auth.InitAuthUsecase(
			dom.User,
			efLogger,
			parser.JSONParser(),
			opt.Auth,
		),
	}
}
