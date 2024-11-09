package book

import (
	"context"

	"github.com/ssentinull/golang-boilerplate/internal/business/model"
	x "github.com/ssentinull/golang-boilerplate/pkg/stdlib/stacktrace"

	"gorm.io/gorm"
)

func (b *bookImpl) CreateBook(ctx context.Context, tx *gorm.DB, p model.BookModel) (model.BookModel, error) {
	result, err := b.createBookSQL(ctx, tx, p)
	if err != nil {
		err = x.WrapPassCode(err, "createBookSQL")
		return result, err
	}

	return result, nil
}
