package book

import (
	"context"

	"github.com/ssentinull/dealls-dating-service/internal/business/model"
	"github.com/ssentinull/dealls-dating-service/pkg/common"
)

func (b *bookUc) CreateBook(ctx context.Context, params model.CreateBookParams) (model.BookModel, error) {
	book := model.BookModel{
		ID:     common.GenerateID(),
		Title:  params.Body.Title,
		Author: params.Body.Author,
	}

	res, err := b.bookDom.CreateBook(ctx, nil, book)
	if err != nil {
		b.efLogger.Error(err)
		return model.BookModel{}, err
	}

	return res, nil
}
