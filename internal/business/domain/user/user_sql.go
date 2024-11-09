package user

import (
	"context"

	"github.com/ssentinull/dealls-dating-service/internal/business/model"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/libsql"
	x "github.com/ssentinull/dealls-dating-service/pkg/stdlib/stacktrace"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (u *userImpl) createUserSQL(ctx context.Context, tx *gorm.DB, p model.UserModel) (model.UserModel, error) {
	if tx == nil {
		tx = u.sql.Leader()
	}

	var result model.UserModel
	if err := tx.WithContext(ctx).Omit(clause.Associations).Create(&p).Error; err != nil {
		return result, x.Wrap(err, libsql.SomethingWentWrongWithDB)
	}

	result = p

	return result, nil
}
