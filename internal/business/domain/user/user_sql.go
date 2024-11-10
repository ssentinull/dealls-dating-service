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
		u.efLogger.Error(err)
		return result, x.Wrap(err, libsql.SomethingWentWrongWithDB)
	}

	result = p

	return result, nil
}

func (u *userImpl) getUserSQL(ctx context.Context, p model.GetUserParams) (model.UserModel, error) {
	db := u.sql.Leader().WithContext(ctx)

	if p.Id > 0 {
		db = db.Where("id = ?", p.Id)
	}

	if p.Email != "" {
		db = db.Where("email = ?", p.Email)
	}

	var result model.UserModel
	if err := db.Take(&result).Error; err != nil {
		u.efLogger.Error(err)
		return result, err
	}

	return result, nil
}
