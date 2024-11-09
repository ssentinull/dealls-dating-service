package book

import (
	"context"

	"github.com/ssentinull/golang-boilerplate/internal/business/model"
	"github.com/ssentinull/golang-boilerplate/pkg/stdlib/libsql"
	x "github.com/ssentinull/golang-boilerplate/pkg/stdlib/stacktrace"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (b *bookImpl) createBookSQL(ctx context.Context, tx *gorm.DB, p model.BookModel) (model.BookModel, error) {
	if tx == nil {
		tx = b.sql.Leader()
	}

	var result model.BookModel
	if err := tx.WithContext(ctx).Omit(clause.Associations).Create(&p).Error; err != nil {
		return result, x.Wrap(err, libsql.SomethingWentWrongWithDB)
	}

	result = p

	return result, nil
}
