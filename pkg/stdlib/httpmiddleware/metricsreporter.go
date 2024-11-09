package httpmiddleware

import (
	"context"
	"regexp"

	"github.com/gin-gonic/gin"
)

type reporter struct {
	c *gin.Context
}

func (r *reporter) Method() string { return r.c.Request.Method }

func (r *reporter) Context() context.Context { return r.c.Request.Context() }

func (r *reporter) URLPath() string {
	reg := regexp.MustCompile(`[^\/]*\d[^\/]*`)
	fullPath := reg.ReplaceAllString(r.c.Request.URL.Path, ":id")
	return fullPath
}

func (r *reporter) StatusCode() int { return r.c.Writer.Status() }

func (r *reporter) BytesWritten() int64 { return int64(r.c.Writer.Size()) }
