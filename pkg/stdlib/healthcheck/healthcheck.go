package healthcheck

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/ssentinull/golang-boilerplate/pkg/build_util"
	"github.com/ssentinull/golang-boilerplate/pkg/stdlib/cache"
	"github.com/ssentinull/golang-boilerplate/pkg/stdlib/libsql"
	"github.com/ssentinull/golang-boilerplate/pkg/stdlib/logger"

	"gorm.io/gorm"
)

type HealthCheck interface {
	Handler() http.Handler
	Stop()
}

type healthcheck struct {
	efLogger  logger.Logger
	db        libsql.SQL
	cache     cache.Redis
	opt       Options
	status    *status
	termReady chan struct{}
}

type status struct {
	mu      *sync.RWMutex
	isReady bool
}

type healthCheckResponse struct {
	Database *string `json:"database,omitempty"`
	Queue    *string `json:"queue,omitempty"`
	Cache    *string `json:"cache,omitempty"`
}

type Options struct {
	WaitBeforeContinue bool
	MaxWaitingTime     time.Duration
	Readiness          ProbeOptions
}

type ProbeOptions struct {
	Enabled          bool
	SuccessThreshold int
	FailureThreshold int
	InitDelaySec     time.Duration
	PeriodSec        time.Duration
	CheckTimeout     time.Duration
	CheckF           func(ctx context.Context, cancel context.CancelFunc) error
}

func Init(efLogger logger.Logger, db libsql.SQL, cache cache.Redis, opt Options) HealthCheck {
	health := &healthcheck{
		efLogger: efLogger,
		db:       db,
		cache:    cache,
		opt:      opt,
		status: &status{
			mu:      &sync.RWMutex{},
			isReady: false,
		},
		termReady: make(chan struct{}, 1),
	}

	health.runChecker()
	if opt.WaitBeforeContinue {
		if err := health.waitUntilReady(); err != nil {
			efLogger.Panic(err)
			return nil
		}
	}

	return health
}

func (h *healthcheck) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dbHealth, cacheHealth := Checker(h.db, h.cache)

		resp := healthCheckResponse{
			Database: mapHealthCheckStatus(dbHealth),
			Cache:    mapHealthCheckStatus(cacheHealth),
		}

		httpStatus := http.StatusOK
		httpMessage := "OK"

		if (dbHealth != nil && !*dbHealth) || (cacheHealth != nil && !*cacheHealth) {
			httpStatus = http.StatusServiceUnavailable
			httpMessage = "failed"
		}

		w.WriteHeader(httpStatus)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": httpMessage,
			"status":  resp,
			"version": build_util.BuildVersion,
		})

	})
}

func Checker(db libsql.SQL, cache cache.Redis) (dbHealth, cacheHealth *bool) {
	// Start DB Check
	dbs_ := make([]*gorm.DB, 0)
	if db.Leader() != nil {
		dbs_ = append(dbs_, db.Leader())
	}

	if db.Follower() != nil {
		dbs_ = append(dbs_, db.Follower())
	}

	if len(dbs_) > 0 {
		dbHealth = new(bool)
		*dbHealth = true

		for _, db := range dbs_ {
			dbase, err := db.DB()
			if err != nil {
				*dbHealth = false
				break
			}

			if err := dbase.Ping(); err != nil {
				*dbHealth = false
				break
			}
		}
	}
	// End DB Check

	// Start Cache Check
	if cache != nil {
		cacheHealth = new(bool)
		*cacheHealth = true

		if err := cache.Ping(); err != nil {
			*cacheHealth = false
		}
	}
	// End Cache Check

	return
}

func (h *healthcheck) runChecker() {
	if h.opt.Readiness.Enabled {
		if h.opt.Readiness.CheckF == nil {
			h.efLogger.Panic("Readiness Check Function is Nil")
			return
		}
		go h.asyncChecker()
	}
}

func (h *healthcheck) waitUntilReady() error {
	ticker := time.NewTicker(1 * time.Second)
	maxAttempt := int64(h.opt.MaxWaitingTime)/int64(time.Second) + 1
	attempt := int64(0)
	for range ticker.C {
		if err := h.isReady(); err == nil || !h.opt.Readiness.Enabled {
			ticker.Stop()
			break
		}
		attempt++
		if attempt > maxAttempt {
			return errors.New("max waiting time is elapsed before App is ready")
		}
	}

	return nil
}

func (h *healthcheck) isReady() error {
	h.status.mu.RLock()
	defer h.status.mu.RUnlock()
	if h.status.isReady {
		return nil
	}

	return errors.New("app is not ready")
}

func (h *healthcheck) Stop() {
	if h.opt.Readiness.Enabled {
		h.setReadinessStatus(false)
		close(h.termReady)
	}
}

func mapHealthCheckStatus(status *bool) *string {
	if status == nil {
		return nil
	}

	res := "success"
	if !*status {
		res = "failed"
	}

	return &res
}
