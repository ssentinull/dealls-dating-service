package grace

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ssentinull/golang-boilerplate/pkg/stdlib/httpserver"
	"github.com/ssentinull/golang-boilerplate/pkg/stdlib/logger"

	"github.com/cloudflare/tableflip"
)

const (
	infoGrace string = `Grace Upgrader:`
	errGrace  string = `%s Grace Upgrader Error`
	infoTCP   string = `TCP Listen`
	errTCP    string = `%s TCP Listen Error`

	_UPGRADE string = `[UPGRADED]`
	_OK      string = `[OK]`
	_FAILED  string = `[FAILED]`
)

// App functions as graceful app upgrader
type App interface {
	// Start serving app
	Serve()
	// Stop stopping app
	Stop()
}

// app
type app struct {
	efLogger   logger.Logger
	httpServer httpserver.HTTPServer
	Upgrader   *tableflip.Upgrader
	Error      error
	Options    Options
	SigHUP     chan os.Signal
	done       chan struct{}
}

type Options struct {
	// Pidfile define custom Pidfile, default ""
	Pidfile string
	// UpgradeTimeout wait period for new app to be ready
	UpgradeTimeout time.Duration `validate:"required"`
	// ShutdownTimeout wait period to shutdown old app
	ShutdownTimeout time.Duration `validate:"required"`
	// Network define tcp or udp. (default: tcp)
	Network string
}

// Init initialize GraceApp Upgrader and http servers.
func Init(efLogger logger.Logger, httpserver httpserver.HTTPServer, opt Options) App {
	upg, err := tableflip.New(tableflip.Options{})
	if err != nil {
		efLogger.Fatal(err, errGrace, _FAILED)
	}
	efLogger.Info(_OK, infoGrace)

	if opt.Network == "" {
		opt.Network = `tcp`
	}
	gs := &app{
		efLogger:   efLogger,
		httpServer: httpserver,
		Upgrader:   upg,
		Error:      err,
		Options:    opt,
		SigHUP:     make(chan os.Signal, 1),
	}
	signal.Notify(gs.SigHUP, syscall.SIGHUP)
	go gs.sighup()
	return gs
}

// Done returns a channel that is closed when the server is fully shut down
func (g *app) Done() <-chan struct{} {
	return g.done
}

// Stop stops apps and its upgrader including all http(s) servers if any
func (g *app) Stop() {
	// Stop the HTTP server
	g.httpServer.Shutdown()

	// Stop the Upgrader
	g.efLogger.Info("Stopping Upgrader")
	g.Upgrader.Stop()
	os.Exit(0)
}

// sighup handle sighup signal
func (g *app) sighup() {
	for range g.SigHUP {
		g.efLogger.Info(_UPGRADE, infoGrace)

		// Create a context with a timeout for the upgrade process
		ctx, cancel := context.WithTimeout(context.Background(), g.Options.UpgradeTimeout)
		defer cancel()

		// Start the upgrade process in a goroutine
		go func() {
			err := g.Upgrader.Upgrade()
			if err != nil {
				g.efLogger.Error(err, errGrace, _FAILED)
				return
			}
			g.efLogger.Info(_OK, infoGrace)
		}()

		// Wait for the upgrade process to complete or the timeout to expire
		select {
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				g.efLogger.Info("Upgrade process timed out")
			}
		}
	}
}

// Server starts apps including all http(s) servers if any
func (g *app) Serve() {
	// get all http(s) servers
	for _, s := range g.httpServer.GetServers() {
		ln, err := g.Upgrader.Fds.Listen(g.Options.Network, s.Addr)
		if err != nil {
			g.efLogger.Fatal(err, fmt.Sprintf(errTCP, _FAILED))
		}
		g.efLogger.Info(_OK, infoTCP, fmt.Sprintf("@%s", s.Addr))
		go g.httpServer.Serve(s.Mode, ln)
	}

	if err := g.Upgrader.Ready(); err != nil {
		g.efLogger.Fatal(err, fmt.Sprintf(errGrace, _FAILED))
	}

	<-g.Upgrader.Exit()

	time.AfterFunc(g.Options.ShutdownTimeout, func() {
		err := fmt.Errorf(errGrace, _FAILED)
		g.efLogger.Fatal(err)
		os.Exit(1)
	})

	// Signal the completion of shutdown
	close(g.done)
	g.httpServer.Shutdown()
}
