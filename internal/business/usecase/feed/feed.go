package feed

import (
	"context"

	"github.com/ssentinull/dealls-dating-service/internal/business/domain/feed"
	"github.com/ssentinull/dealls-dating-service/internal/business/domain/sqltx"
	"github.com/ssentinull/dealls-dating-service/internal/business/domain/user"
	"github.com/ssentinull/dealls-dating-service/internal/business/model"
	"github.com/ssentinull/dealls-dating-service/internal/types"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/auth"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/logger"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/parser"
)

type UsecaseItf interface {
	CreatePreference(ctx context.Context, params model.CreatePreferenceParams) (model.PreferenceModel, error)
	GetFeed(ctx context.Context, params model.GetFeedParams) ([]model.FeedModel, *types.Pagination, error)
	SwipeFeed(ctx context.Context, params model.SwipeFeedParams) (model.SwipeModel, error)
}

type feedUc struct {
	sqlTxDom   sqltx.DomainItf
	feedDom    feed.DomainItf
	userDom    user.DomainItf
	efLogger   logger.Logger
	jsonParser parser.JSONParser
	auth       auth.Auth
	opt        Options
}

type Options struct{}

func InitFeedUsecase(
	sqlTxDom sqltx.DomainItf,
	feedDom feed.DomainItf,
	userDom user.DomainItf,
	efLogger logger.Logger,
	jsonParser parser.JSONParser,
	auth auth.Auth,
	opt Options,
) UsecaseItf {
	return &feedUc{
		sqlTxDom:   sqlTxDom,
		feedDom:    feedDom,
		userDom:    userDom,
		efLogger:   efLogger,
		jsonParser: jsonParser,
		auth:       auth,
		opt:        opt,
	}
}
