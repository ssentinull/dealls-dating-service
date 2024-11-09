package httpmiddleware

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"

	"github.com/ssentinull/dealls-dating-service/pkg/common"
	"github.com/ssentinull/dealls-dating-service/pkg/constants"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/auth"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/logger"

	"github.com/gin-gonic/gin"
)

const (
	ContextKeyURLPath       string = "url.path"
	ContextKeyRequestMethod string = "request.method"

	// AuthorizationHeader const
	AuthorizationHeaderKey         string = "Authorization"
	AuthorizationHeaderPrefix      string = "Bearer"
	AuthorizationHeaderSeparator   string = " "
	AuthorizationHeaderPartsCount  int    = 2
	AuthorizationHeaderIndexPrefix int    = 0
	AuthorizationHeaderIndexToken  int    = 1
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

type HTTPMWare interface {
	CORSMiddleware() gin.HandlerFunc
	LoggerMiddleware() gin.HandlerFunc
	JWTAuthMiddleware() gin.HandlerFunc
	AlertMiddleware() gin.HandlerFunc
	MockApiMiddleware() gin.HandlerFunc
}

type httpmw struct {
	efLogger logger.Logger
	auth     auth.Auth
	opt      Options
}

type SignatureOptions struct {
	Secret        string
	TimeTolerance int // in minutes
}

type Options struct {
	AppName        string
	Env            string
	MockServerHost string
	MockServerPath string
	MockApiList    []string
	Signature      SignatureOptions
}

type simpleHTTPResponse struct {
	Message string
	Success bool
}

func Init(
	efLogger logger.Logger,
	auth auth.Auth,
	opt Options,
) HTTPMWare {
	return &httpmw{
		efLogger: efLogger,
		auth:     auth,
		opt:      opt,
	}
}

func (mw *httpmw) CORSMiddleware() gin.HandlerFunc {
	allowedHeaders := "Content-Type, Authorization, X-App-Id, X-Client-Id, X-Client-Version"
	if mw.opt.Env != "production" {
		allowedHeaders = "Content-Type, Authorization, X-App-Id, X-Client-Id, X-Client-Version, X-App-Debug"
	}

	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusOK)
		} else {
			c.Next()
		}
	}
}

func (mw *httpmw) LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(ContextKeyURLPath, c.Request.URL.Path)
		c.Set(ContextKeyRequestMethod, c.Request.Method)

		// set x-client-id on requests
		clientID := c.Request.Header.Get("X-Client-Id")
		if clientID == "" {
			clientID = mw.opt.AppName
		}

		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), constants.XClientId, clientID))
		ctype := c.Request.Header.Get(constants.HeaderContentType)

		if strings.HasPrefix(ctype, constants.HeaderContentTypeApplicationJSON) {
			// read request body content
			requestBody, err := io.ReadAll(c.Request.Body)

			if err != nil {
				c.Next()
			}

			// restore request body
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusOK)
		} else {
			c.Next()
		}
	}
}

func (mw *httpmw) JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get(AuthorizationHeaderKey)
		if len(authHeader) == 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, simpleHTTPResponse{
				Message: http.StatusText(http.StatusUnauthorized),
				Success: false,
			})
			return
		}

		parts := strings.SplitN(authHeader, AuthorizationHeaderSeparator, AuthorizationHeaderPartsCount)
		if len(parts) < AuthorizationHeaderPartsCount {
			c.AbortWithStatusJSON(http.StatusForbidden, simpleHTTPResponse{
				Message: http.StatusText(http.StatusUnauthorized),
				Success: false,
			})
			return
		}

		prefix := parts[AuthorizationHeaderIndexPrefix]
		tokenString := parts[AuthorizationHeaderIndexToken]

		if strings.Compare(prefix, AuthorizationHeaderPrefix) != 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, simpleHTTPResponse{
				Message: http.StatusText(http.StatusUnauthorized),
				Success: false,
			})
			return
		}

		claims, err := mw.auth.ParseJWTToken(tokenString)
		if err != nil {
			if tokenString == mw.auth.GetStaticToken() {
				c.Next()
			}
			c.AbortWithStatusJSON(http.StatusForbidden, simpleHTTPResponse{
				Message: err.Error(),
				Success: false,
			})
			return
		}

		c.Set(auth.ContextKeyJwtToken, tokenString)
		c.Set(auth.ContextKeyClaims, claims)
		c.Set(auth.ContextKeyUserID, claims.ID)
		c.Set(auth.ContextKeyUserEmail, claims.Email)
		c.Set(auth.ContextKeyUserName, claims.Name)

		c.Next()
	}
}

func (mw *httpmw) AlertMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		statusCode := c.Writer.Status()
		if statusCode >= 400 {
			ctx := c.Request.Context()

			if statusCode >= 500 {
				mw.efLogger.Error(ctx, blw.body.String())
			} else {
				mw.efLogger.Warn(ctx, blw.body.String())
			}
		}
	}
}

func (mw *httpmw) MockApiMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !common.SliceHas(mw.opt.MockApiList, c.Request.URL.Path) {
			c.Next()
			return
		}

		// remove base path
		// ref: https://github.com/stoplightio/prism#-faqs
		pathRegex := regexp.MustCompile(`^/v\d+/api(.*)$`)
		requestPath := pathRegex.ReplaceAllString(c.Request.URL.Path, mw.opt.MockServerPath+"$1")

		remote, err := url.Parse(mw.opt.MockServerHost)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, simpleHTTPResponse{
				Message: fmt.Sprintf("failed to parse mock server : %s", mw.opt.MockServerHost),
				Success: false,
			})
			return
		}

		// hardcoded to return dynamic response from prism
		if _, ok := c.Request.Header["Prefer"]; !ok {
			c.Request.Header["Prefer"] = []string{"dynamic=true"}
		}

		director := func(req *http.Request) {
			req.Header = c.Request.Header
			req.Host = remote.Host
			req.URL.Scheme = remote.Scheme
			req.URL.Host = remote.Host
			req.URL.Path = requestPath
		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(c.Writer, c.Request)
		c.Abort()
	}
}
