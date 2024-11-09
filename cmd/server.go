package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	goValidator "github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"github.com/ssentinull/dealls-dating-service/internal/business/domain"
	"github.com/ssentinull/dealls-dating-service/internal/business/usecase"
	restserver "github.com/ssentinull/dealls-dating-service/internal/handler/rest"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/auth"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/cache"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/grace"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/healthcheck"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/httpclient"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/httpmiddleware"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/httpmux"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/httpserver"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/libsql"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/logger"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/parser"
	x "github.com/ssentinull/dealls-dating-service/pkg/stdlib/stacktrace"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/swagger"
)

var (
	conf Conf

	// Resource Storage
	sqlClient   libsql.SQL
	redisClient cache.Redis

	// Server Handler
	REST restserver.REST

	// Server Infrastructure
	efLogger       logger.Logger
	aut            auth.Auth
	httpMiddleware httpmiddleware.HTTPMWare
	httpMux        httpmux.HTTPMux
	httpServer     httpserver.HTTPServer
	parse          parser.Parser
	swagg          swagger.Swagger
	healthz        healthcheck.HealthCheck
	httpClient     httpclient.HTTPClient

	// Business Layer
	dom *domain.Domain
	uc  *usecase.Usecase

	// Application
	app grace.App
)

var runCmd = &cobra.Command{
	Use:   "server",
	Short: "run server",
	Long:  `This subcommand start the server`,
	Run:   run,
}

func init() {
	RootCmd.AddCommand(runCmd)
}

func run(cmd *cobra.Command, args []string) {
	// infrastructure initialization
	efLogger = logger.Init()
	swagg = swagger.Init(efLogger, conf.Swagger)
	goBindingValidator := goValidator.New()
	goBindingValidator.SetTagName("binding") // use to validate tag name binding

	parse = parser.Init(efLogger, goBindingValidator, conf.Parser)
	httpClient = httpclient.Init(efLogger, parse, conf.HTTPClient)
	aut = auth.Init(conf.Auth, httpClient, parse.JSONParser())
	httpMiddleware = httpmiddleware.Init(efLogger, aut, conf.HTTPMW)

	// storage initialization
	sqlClient = libsql.Init(efLogger, conf.SQL)
	redisClient = cache.Init(efLogger, conf.Redis)

	// register readiness health check function
	conf.HealthCheck.Readiness.CheckF = func(ctx context.Context, cancel context.CancelFunc) error {
		defer cancel()
		dbHealth, cacheHealth := healthcheck.Checker(sqlClient, redisClient)

		if dbHealth != nil && !*dbHealth {
			return x.New("Database is not ready")
		}

		if cacheHealth != nil && !*cacheHealth {
			return x.New("Cache is not ready")
		}

		return nil
	}

	// health check initialization
	healthz = healthcheck.Init(efLogger, sqlClient, redisClient, conf.HealthCheck)

	// mux initialization
	httpMux = httpmux.Init(efLogger, swagg, httpMiddleware, healthz, conf.HTTPMux)

	// business initialization
	dom = domain.Init(
		efLogger,
		parse,
		sqlClient,
		redisClient,
		conf.Business.Domain,
	)

	uc = usecase.Init(
		efLogger,
		parse,
		aut,
		dom,
		conf.Business.Usecase,
	)

	// handlers initialization
	REST = restserver.Init(efLogger, parse, aut, httpMux, httpMiddleware, uc, conf.Handler.REST)

	// HTTP Server Initialization
	httpServer = httpserver.Init(efLogger, httpMux, conf.HTTP)

	// App Initialization
	app = grace.Init(efLogger, httpServer, conf.Grace)

	defer func() {
		shutdown()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		for range quit {
			shutdown()
		}
	}()

	app.Serve()
}

func shutdown() {
	if healthz != nil {
		healthz.Stop()
	}

	if sqlClient != nil {
		sqlClient.Stop()
	}

	if redisClient != nil {
		redisClient.Stop()
	}

	app.Stop()
}
