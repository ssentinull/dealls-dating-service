package libsql

import (
	"database/sql"
	"sync"
	"time"

	"github.com/ssentinull/golang-boilerplate/pkg/stdlib/logger"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gLogger "gorm.io/gorm/logger"
)

type SQL interface {
	Leader() *gorm.DB
	MockLeader() sqlmock.Sqlmock
	Follower() *gorm.DB
	MockFollower() sqlmock.Sqlmock
	Stop()
}

type sqlImpl struct {
	efLogger logger.Logger
	endOnce  *sync.Once
	leader   Connection
	follower Connection
	opt      Options
}

type Connection struct {
	ORM  *gorm.DB
	DB   *sql.DB
	mock sqlmock.Sqlmock
}

type Options struct {
	Leader   Config
	Follower Config
}

type Config struct {
	Enabled bool
	Mock    bool
	DSN     string
}

func Init(efLogger logger.Logger, opt Options) SQL {
	sqlx := &sqlImpl{
		efLogger: efLogger,
		endOnce:  &sync.Once{},
		opt:      opt,
	}

	if opt.Leader.Enabled && !opt.Leader.Mock {
		sqlDB, err := sql.Open("pgx", opt.Leader.DSN)
		if err != nil {
			return nil
		}

		// Create the connection pool
		sqlDB.SetConnMaxIdleTime(time.Minute * 5)

		// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
		sqlDB.SetMaxIdleConns(10)

		// SetMaxOpenConns sets the maximum number of open connections to the database.
		sqlDB.SetMaxOpenConns(100)

		// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
		sqlDB.SetConnMaxLifetime(time.Hour)

		gormDB, err := gorm.Open(
			postgres.New(postgres.Config{Conn: sqlDB}),
			&gorm.Config{
				Logger:                 gLogger.Default.LogMode(gLogger.Silent),
				SkipDefaultTransaction: true,
			},
		)

		if err != nil {
			efLogger.Fatal(err)
			return nil
		}

		sqlx.leader = Connection{
			ORM: gormDB,
			DB:  sqlDB,
		}
	}

	if opt.Follower.Enabled && !opt.Follower.Mock {
		sqlDB, err := sql.Open("pgx", opt.Follower.DSN)
		if err != nil {
			return nil
		}

		// Create the connection pool
		sqlDB.SetConnMaxIdleTime(time.Minute * 5)

		// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
		sqlDB.SetMaxIdleConns(10)

		// SetMaxOpenConns sets the maximum number of open connections to the database.
		sqlDB.SetMaxOpenConns(100)

		// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
		sqlDB.SetConnMaxLifetime(time.Hour)

		gormDB, err := gorm.Open(
			postgres.New(postgres.Config{Conn: sqlDB}),
			&gorm.Config{
				Logger:                 gLogger.Default.LogMode(gLogger.Silent),
				SkipDefaultTransaction: true,
			},
		)

		if err != nil {
			efLogger.Fatal(err)
			return nil
		}

		sqlx.follower = Connection{
			ORM: gormDB,
			DB:  sqlDB,
		}
	}

	if opt.Leader.Mock {
		sqlDB, mock, err := sqlmock.New()
		if err != nil {
			efLogger.Fatal("an error '%s' was not expected when opening a stub database connection", err)
		}

		gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
		if err != nil {
			efLogger.Fatal("an error '%s' was not expected when opening a stub database connection", err)
		}

		sqlx.leader = Connection{
			ORM:  gormDB,
			DB:   sqlDB,
			mock: mock,
		}
	}

	if opt.Follower.Mock {
		sqlDB, mock, err := sqlmock.New()
		if err != nil {
			efLogger.Fatal("an error '%s' was not expected when opening a stub database connection", err)
		}

		gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
		if err != nil {
			efLogger.Fatal("an error '%s' was not expected when opening a stub database connection", err)
		}

		sqlx.follower = Connection{
			ORM:  gormDB,
			DB:   sqlDB,
			mock: mock,
		}
	}

	return sqlx
}

func (s *sqlImpl) Leader() *gorm.DB {
	return s.leader.ORM
}

func (s *sqlImpl) MockLeader() sqlmock.Sqlmock {
	return s.leader.mock
}

func (s *sqlImpl) Follower() *gorm.DB {
	return s.follower.ORM
}

func (s *sqlImpl) MockFollower() sqlmock.Sqlmock {
	return s.follower.mock
}

func (s *sqlImpl) Stop() {
	s.endOnce.Do(func() {
		s.efLogger.Info("Shutting down database connection")
		if s.leader.DB != nil {
			s.leader.DB.Close()
		}

		if s.follower.DB != nil {
			s.follower.DB.Close()
		}

		s.efLogger.Info("[OK]: Shutdown Database connection")
	})
}
