package restserver

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/ssentinull/dealls-dating-service/internal/business/usecase"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/auth"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/httpmiddleware"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/httpmux"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/logger"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/parser"
)

var once = &sync.Once{}

type REST interface{}

type Options struct {
	DebugMode bool
}

type rest struct {
	efLogger logger.Logger
	parser   parser.Parser
	auth     auth.Auth
	mux      httpmux.HTTPMux
	mw       httpmiddleware.HTTPMWare
	uc       *usecase.Usecase
	opt      Options
}

func Init(
	efLogger logger.Logger,
	parser parser.Parser,
	auth auth.Auth,
	mux httpmux.HTTPMux,
	mw httpmiddleware.HTTPMWare,
	uc *usecase.Usecase,
	opt Options,
) REST {
	var r *rest
	once.Do(func() {
		r = &rest{
			efLogger: efLogger,
			parser:   parser,
			auth:     auth,
			mux:      mux,
			mw:       mw,
			uc:       uc,
			opt:      opt,
		}
		r.Serve()
	})

	return r
}

func (r *rest) Serve() {
	route := r.mux.Engine()
	route.HandleMethodNotAllowed = true

	route.NoMethod(func(c *gin.Context) {
		c.String(http.StatusMethodNotAllowed, "not allowed")
	})

	route.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "not found")
	})

	apiV1 := route.Group("/v1/api")
	apiV1.POST("/books", r.CreateBook)
}
