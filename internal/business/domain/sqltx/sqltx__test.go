package sqltx_test

import (
	"errors"
	"testing"

	"github.com/ssentinull/dealls-dating-service/internal/business/domain/sqltx"

	"github.com/c2fo/testify/assert"
)

func TestSqlTxDomain_BeginTX(t *testing.T) {
	mockedDependency := NewMockedDependency(t)
	sqltxDomain := sqltx.InitSQLTXDomain(
		mockedDependency.efLogger,
		mockedDependency.sql,
		mockedDependency.opt,
	)

	t.Run("success", func(t *testing.T) {
		mockedDependency.sql.MockLeader().ExpectBegin()
		begin := sqltxDomain.BeginTX()
		assert.NotNil(t, begin)
		assert.Nil(t, mockedDependency.sql.MockLeader().ExpectationsWereMet())
	})
}

func TestSqlTxDomain_CommitTX(t *testing.T) {
	mockedDependency := NewMockedDependency(t)
	sqltxDomain := sqltx.InitSQLTXDomain(
		mockedDependency.efLogger,
		mockedDependency.sql,
		mockedDependency.opt,
	)

	begin := mockedDependency.sql.Leader()

	t.Run("failed", func(t *testing.T) {
		mockedDependency.sql.MockLeader().ExpectCommit().WillReturnError(errors.New("db error"))
		err := sqltxDomain.CommitTX(begin)
		assert.Error(t, err)
	})
}

func TestSqlTxDomain_RollbackTX(t *testing.T) {
	mockedDependency := NewMockedDependency(t)
	sqltxDomain := sqltx.InitSQLTXDomain(
		mockedDependency.efLogger,
		mockedDependency.sql,
		mockedDependency.opt,
	)

	begin := mockedDependency.sql.Leader()

	t.Run("failed", func(t *testing.T) {
		mockedDependency.sql.MockLeader().ExpectRollback().WillReturnError(errors.New("db error"))
		err := sqltxDomain.RollbackTX(begin)
		assert.Error(t, err)
	})
}
