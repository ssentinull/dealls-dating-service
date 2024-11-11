package user_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/c2fo/testify/assert"
	"github.com/ssentinull/dealls-dating-service/internal/business/domain/user"
	"github.com/ssentinull/dealls-dating-service/internal/business/model"

	"go.uber.org/mock/gomock"
)

func TestUserDomain_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	defer func() {
		ctx.Done()
		ctrl.Finish()
	}()

	mockedDependency := NewMockedDependency(t, ctrl)
	userDomain := user.InitUserDomain(
		mockedDependency.efLogger,
		mockedDependency.json,
		mockedDependency.sql,
		mockedDependency.cache,
		mockedDependency.opt,
	)

	ID := 1
	user := model.UserModel{
		Email:    "john.doe@mail.com",
		Password: "12345",
		Name:     "John Doe",
	}

	createUserSQL := `INSERT INTO "users" ("email","password","name","gender","birth_date","location","is_premium_user","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING "id"`

	t.Run("success", func(t *testing.T) {
		mockedDependency.sql.MockLeader().ExpectBegin()
		mockedDependency.sql.MockLeader().ExpectQuery(regexp.QuoteMeta(createUserSQL)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(ID))
		mockedDependency.sql.MockLeader().ExpectCommit()

		res, err := userDomain.CreateUser(ctx, nil, user)
		assert.NoError(t, err, err)
		assert.NotEmpty(t, res)
		assert.Nil(t, mockedDependency.sql.MockLeader().ExpectationsWereMet())
	})

	t.Run("failed", func(t *testing.T) {
		mockedDependency.sql.MockLeader().ExpectBegin()
		mockedDependency.sql.MockLeader().ExpectQuery(regexp.QuoteMeta(createUserSQL)).WillReturnError(errors.New("db error"))
		mockedDependency.sql.MockLeader().ExpectRollback()

		res, err := userDomain.CreateUser(ctx, nil, user)
		assert.NotNil(t, err, err)
		assert.NotEmpty(t, res)
		assert.Nil(t, mockedDependency.sql.MockLeader().ExpectationsWereMet())
	})
}

func TestUserDomain_GetUserByParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	defer func() {
		ctx.Done()
		ctrl.Finish()
	}()

	mockedDependency := NewMockedDependency(t, ctrl)
	userDomain := user.InitUserDomain(
		mockedDependency.efLogger,
		mockedDependency.json,
		mockedDependency.sql,
		mockedDependency.cache,
		mockedDependency.opt,
	)

	ID := int64(1)
	email := "john.doe@mail.com"

	param := model.GetUserParams{
		Id:    ID,
		Email: email,
	}

	paramRaw, err := json.Marshal(param)
	assert.NoError(t, err)

	cacheKey := fmt.Sprintf(user.UserCacheKey, string(paramRaw))
	getUserByParamSQL := `SELECT * FROM "users" WHERE id = $1 AND email = $2 AND "users"."deleted_at" IS NULL LIMIT 1`

	t.Run("success", func(t *testing.T) {
		mockedDependency.json.EXPECT().Marshal(param).Times(1).Return(paramRaw, nil)
		mockedDependency.cache.Mock().ExpectGet(cacheKey).SetErr(errors.New("redis error"))
		mockedDependency.sql.MockLeader().ExpectQuery(regexp.QuoteMeta(getUserByParamSQL)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).AddRow(ID, email))

		res, err := userDomain.GetUserByParams(ctx, param)
		assert.NoError(t, err, err)
		assert.NotEmpty(t, res)
		assert.Nil(t, mockedDependency.sql.MockLeader().ExpectationsWereMet())
	})

	t.Run("failed", func(t *testing.T) {
		mockedDependency.json.EXPECT().Marshal(param).Times(1).Return(paramRaw, nil)
		mockedDependency.cache.Mock().ExpectGet(cacheKey).SetErr(errors.New("redis error"))
		mockedDependency.sql.MockLeader().ExpectQuery(regexp.QuoteMeta(getUserByParamSQL)).WillReturnError(errors.New("db error"))

		activity, err := userDomain.GetUserByParams(ctx, param)
		assert.NotNil(t, err, err)
		assert.NotEmpty(t, activity)
		assert.Nil(t, mockedDependency.sql.MockLeader().ExpectationsWereMet())
	})
}
