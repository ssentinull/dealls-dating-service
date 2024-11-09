package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/httpmux"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/logger"
)

const (
	HTTP int = iota
	HTTPS
)

const (
	infoServe string = `Server:`

	OK     string = "[OK]"
	FAILED string = "[FAILED]"

	defaultTimeoutReadHeader = 60 * time.Second
	defaultTimeoutRead       = 5 * time.Minute
	defaultTimeoutWrite      = 5 * time.Minute
	defaultTimeoutIdle       = 1 * time.Hour
)

const (
	errServe          = `%s %s Server`
	errNetListener    = `failed to get %s net listener`
	errListenerFailed = `Cannot start listening %s @%s`
)

var (
	once   = &sync.Once{}
	server = []string{
		HTTP:  "[HTTP]",
		HTTPS: "[HTTPS]",
	}
)

type HTTPServer interface {
	Serve(mode int, ln net.Listener)
	Shutdown()
	GetServers() []*Server
}

type Server struct {
	Addr string
	Mode int
}

type httpserver struct {
	efLogger logger.Logger
	servers  []*http.Server
	opt      Options
}

type Options struct {
	Address           string
	Port              int `validate:"required"`
	TimeoutReadHeader time.Duration
	TimeoutRead       time.Duration
	TimeoutWrite      time.Duration
	TimeoutIdle       time.Duration
}

func Init(efLogger logger.Logger, mux httpmux.HTTPMux, opt Options) HTTPServer {
	var h *httpserver
	var servers []*http.Server

	// set default options value
	if opt.TimeoutReadHeader < 1 {
		opt.TimeoutReadHeader = defaultTimeoutReadHeader
	}

	if opt.TimeoutRead < 1 {
		opt.TimeoutRead = defaultTimeoutRead
	}

	if opt.TimeoutWrite < 1 {
		opt.TimeoutWrite = defaultTimeoutWrite
	}

	if opt.TimeoutIdle < 1 {
		opt.TimeoutIdle = defaultTimeoutIdle
	}

	servers = append(servers, &http.Server{
		Addr:              ":" + strconv.Itoa(opt.Port),
		Handler:           mux.Engine(),
		ReadHeaderTimeout: opt.TimeoutReadHeader,
		ReadTimeout:       opt.TimeoutRead,
		WriteTimeout:      opt.TimeoutWrite,
		IdleTimeout:       opt.TimeoutIdle,
	})

	once.Do(func() {
		// init http server
		h = &httpserver{
			efLogger: efLogger,
			opt:      opt,
			servers:  servers,
		}
	})

	return h
}

func (h *httpserver) Serve(mode int, ln net.Listener) {
	if ln == nil {
		err := fmt.Errorf(errNetListener, server[mode])
		h.efLogger.Panic(err)
	}

	h.efLogger.Info(OK, infoServe, fmt.Sprintf("%s @%s", server[mode], ln.Addr().String()))
	if h.servers[mode] != nil {
		if err := h.servers[mode].Serve(ln); !errors.Is(err, http.ErrServerClosed) {
			h.efLogger.Panic(err, fmt.Sprintf(errServe, FAILED, server[mode]))
		}
	} else {
		err := fmt.Errorf(errListenerFailed, server[mode], ln.Addr().String())
		h.efLogger.Panic(err)
	}
}

func (h *httpserver) Shutdown() {
	h.efLogger.Info("Shutting Down HTTP Server")
	for _, s := range h.servers {
		if err := s.Shutdown(context.Background()); err != nil {
			s.ErrorLog.Fatal(err, fmt.Sprintf(errServe, FAILED, s.Addr))
		}
	}

	h.efLogger.Info("[OK]: Shutdown HTTP Server")
}

func (h *httpserver) GetServers() []*Server {
	var servers []*Server
	for k, s := range h.servers {
		servers = append(servers, &Server{Mode: k, Addr: s.Addr})
	}

	return servers
}
