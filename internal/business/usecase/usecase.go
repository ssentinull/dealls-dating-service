package usecase

import (
	"github.com/ssentinull/dealls-dating-service/internal/business/domain"
	"github.com/ssentinull/dealls-dating-service/internal/business/usecase/book"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/logger"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/parser"
)

type Usecase struct {
	Book book.UsecaseItf
}

type Options struct {
	Book book.Options
}

func Init(
	efLogger logger.Logger,
	parser parser.Parser,
	dom *domain.Domain,
	opt Options,
) *Usecase {
	return &Usecase{
		Book: book.InitBookUsecase(
			dom.Book,
			efLogger,
			parser.JSONParser(),
			opt.Book,
		),
	}
}
