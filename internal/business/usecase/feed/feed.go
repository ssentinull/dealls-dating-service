package feed

import (
	"context"

	"github.com/ssentinull/dealls-dating-service/internal/business/domain/feed"
	"github.com/ssentinull/dealls-dating-service/internal/business/model"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/auth"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/logger"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/parser"
)

type UsecaseItf interface {
	CreatePreference(ctx context.Context, params model.CreatePreferenceParams) (model.PreferenceModel, error)
}

type feedUc struct {
	feedDom    feed.DomainItf
	efLogger   logger.Logger
	jsonParser parser.JSONParser
	auth       auth.Auth
	opt        Options
}

type Options struct{}

func InitFeedUsecase(
	feedDom feed.DomainItf,
	efLogger logger.Logger,
	jsonParser parser.JSONParser,
	auth auth.Auth,
	opt Options,
) UsecaseItf {
	return &feedUc{
		feedDom:    feedDom,
		efLogger:   efLogger,
		jsonParser: jsonParser,
		auth:       auth,
		opt:        opt,
	}
}
