package book

import (
	"context"

	"github.com/ssentinull/dealls-dating-service/internal/business/domain/book"
	"github.com/ssentinull/dealls-dating-service/internal/business/model"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/logger"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/parser"
)

type UsecaseItf interface {
	CreateBook(ctx context.Context, params model.CreateBookParams) (model.BookModel, error)
}

type bookUc struct {
	bookDom    book.DomainItf
	efLogger   logger.Logger
	jsonParser parser.JSONParser
	opt        Options
}

type Options struct{}

func InitBookUsecase(bookDom book.DomainItf, efLogger logger.Logger, jsonParser parser.JSONParser, opt Options) UsecaseItf {
	return &bookUc{
		bookDom:    bookDom,
		efLogger:   efLogger,
		jsonParser: jsonParser,
		opt:        opt,
	}
}
