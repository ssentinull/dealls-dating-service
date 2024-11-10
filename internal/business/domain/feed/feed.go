package feed

import (
	"context"

	"github.com/ssentinull/dealls-dating-service/internal/business/model"
	"github.com/ssentinull/dealls-dating-service/internal/types"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/cache"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/libsql"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/logger"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/parser"

	"gorm.io/gorm"
)

type DomainItf interface {
	CreatePreference(ctx context.Context, tx *gorm.DB, p model.PreferenceModel) (model.PreferenceModel, error)
	GetPreferenceByParams(ctx context.Context, p model.GetPreferenceParams) (model.PreferenceModel, error)
	GetFeedByParams(ctx context.Context, p model.GetFeedParams) ([]model.FeedModel, *types.Pagination, error)
	CreateSwipe(ctx context.Context, tx *gorm.DB, p model.SwipeModel) (model.SwipeModel, error)
	GetSwipeByParams(ctx context.Context, p model.GetSwipeParams) (model.SwipeModel, error)
	CreateMatch(ctx context.Context, tx *gorm.DB, p model.MatchModel) (model.MatchModel, error)
}

type feedImpl struct {
	efLogger logger.Logger
	json     parser.JSONParser
	cache    cache.Redis
	sql      libsql.SQL
	opt      Options
}

type Options struct{}

func InitFeedDomain(efLogger logger.Logger, json parser.JSONParser, sql libsql.SQL, redisClient cache.Redis, opt Options) DomainItf {
	return &feedImpl{
		efLogger: efLogger,
		json:     json,
		cache:    redisClient,
		sql:      sql,
		opt:      opt,
	}
}
