package httpmux

import (
	"net/http"

	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/healthcheck"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/httpmiddleware"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/logger"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/swagger"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

type HTTPMux interface {
	Engine() *gin.Engine
}

type mux struct {
	efLogger logger.Logger
	engine   *gin.Engine
	opt      Options
}

type Options struct {
	AppName string
	AppEnv  string
}

func Init(efLogger logger.Logger, swagg swagger.Swagger, middleware httpmiddleware.HTTPMWare, health healthcheck.HealthCheck, opt Options) HTTPMux {
	efLogger.Info("Initializing HTTP Mux")

	if opt.AppEnv != "local" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(middleware.AlertMiddleware())

	engine.GET("/", func(ctx *gin.Context) {
		message := `UP`
		ctx.String(http.StatusOK, message)
	})

	engine.GET(swagg.GetPath(), swagg.GetHandler())
	engine.GET("/health", gin.WrapH(health.Handler()))

	// NOTE: enable jwt middleware when needed
	// engine.Use(middleware.JWTAuthMiddleware())

	engine.Use(middleware.CORSMiddleware())
	engine.Use(middleware.LoggerMiddleware())
	engine.Use(middleware.MockApiMiddleware())
	pprof.Register(engine)

	return &mux{
		efLogger: efLogger,
		engine:   engine,
		opt:      opt,
	}
}

func (m *mux) Engine() *gin.Engine {
	return m.engine
}
