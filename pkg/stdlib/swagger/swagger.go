package swagger

import (
	"os"

	"github.com/ssentinull/golang-boilerplate/pkg/stdlib/logger"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag"
)

var (
	log     logger.Logger
	docPath string
)

type Swagger interface {
	GetPath() string
	GetName() string
	GetHandler() gin.HandlerFunc
}

type SWAG struct {
	handler gin.HandlerFunc
	opt     Options
}

type SwagDoc struct {
	doc string
}

type Options struct {
	Name    string
	Path    string
	DocPath string
}

func Init(efLogger logger.Logger, opt Options) Swagger {
	log = efLogger
	docPath = opt.DocPath
	sd := &SwagDoc{}
	swag.Register(swag.Name, sd)
	return &SWAG{
		opt:     opt,
		handler: ginSwagger.WrapHandler(swaggerFiles.Handler),
	}
}

func (s *SwagDoc) ReadDoc() string {
	if s.doc == "" {
		data, err := os.ReadFile(docPath)
		if err != nil {
			log.Fatal(err)
			return ""
		}
		s.doc = string(data)
	}
	return s.doc
}

func (s *SWAG) GetPath() string {
	return s.opt.Path
}

func (s *SWAG) GetName() string {
	return s.opt.Name
}

func (s *SWAG) GetHandler() gin.HandlerFunc {
	return s.handler
}
