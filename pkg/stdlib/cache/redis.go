package cache

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/logger"
	x "github.com/ssentinull/dealls-dating-service/pkg/stdlib/stacktrace"

	redis "github.com/go-redis/redis/v8"
	redismock "github.com/go-redis/redismock/v8"
)

const (
	cluster = iota
	single
	sentinel
)

const (
	defaultMaxConnectTimeout = 15 * time.Second

	infoRedis string = `Redis:`

	OK     string = "[OK]"
	FAILED string = "[FAILED]"
)

var Nil = redis.Nil

type Redis interface {
	Stop()
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Delete(ctx context.Context, key string) *redis.IntCmd
	Lock(ctx context.Context, key string, ttl time.Duration, opt *LockOptions) (*Lock, error)
	Mock() redismock.ClientMock
	Ping() error
}

type Cmd interface {
	redis.UniversalClient
}

type redisc struct {
	mu       *sync.Mutex
	mode     int
	logger   logger.Logger
	locker   *LockClient
	sentinel *redis.Client
	cluster  *redis.ClusterClient
	opt      Options
	driver   string
	termSig  chan struct{}
	mock     redismock.ClientMock
}

type Options struct {
	Enabled            bool
	Address            []string `validate:"required"`
	DB                 int      `validate:"required"`
	Password           string
	Mock               bool
	InsecureSkipVerify bool
	MaxRetries         int
	MinRetryBackoff    time.Duration
	MaxRetryBackoff    time.Duration
	DialTimeout        time.Duration
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	PoolSize           int
	MinIdleConns       int
	MaxConnAge         time.Duration
	PoolTimeout        time.Duration
	IdleTimeout        time.Duration
	IdleCheckFrequency time.Duration
}

func Init(efLogger logger.Logger, opt Options) Redis {
	if !opt.Enabled && !opt.Mock {
		return nil
	}

	var client *redisc

	if opt.Mock {
		rc, mock := redismock.NewClientMock()
		mode := single
		redisDriver := `single`

		client = &redisc{
			mu:       &sync.Mutex{},
			mode:     mode,
			logger:   efLogger,
			locker:   NewLockClient(rc),
			sentinel: rc,
			opt:      opt,
			driver:   redisDriver,
			termSig:  make(chan struct{}, 1),
			mock:     mock,
		}

		return client
	}

	univOptions := &redis.UniversalOptions{
		Addrs:              opt.Address,
		Password:           opt.Password,
		DB:                 opt.DB,
		MaxRetries:         opt.MaxRetries,
		MinRetryBackoff:    opt.MinRetryBackoff,
		MaxRetryBackoff:    opt.MaxRetryBackoff,
		DialTimeout:        opt.DialTimeout,
		ReadTimeout:        opt.ReadTimeout,
		WriteTimeout:       opt.WriteTimeout,
		PoolSize:           opt.PoolSize,
		MinIdleConns:       opt.MinIdleConns,
		MaxConnAge:         opt.MaxConnAge,
		PoolTimeout:        opt.PoolTimeout,
		IdleTimeout:        opt.IdleTimeout,
		IdleCheckFrequency: opt.IdleCheckFrequency,
	}

	univClient := redis.NewUniversalClient(univOptions)
	redisHost := strings.Join(univOptions.Addrs, ",")

	sentinelClient, ok := univClient.(*redis.Client)
	if ok {
		mode := single
		redisDriver := `single`

		client = &redisc{
			mu:       &sync.Mutex{},
			mode:     mode,
			logger:   efLogger,
			locker:   NewLockClient(sentinelClient),
			sentinel: sentinelClient,
			opt:      opt,
			driver:   redisDriver,
			termSig:  make(chan struct{}, 1),
		}

		if err := client.Connect(); err != nil {
			err = x.WrapWithCode(err, EcodeConnectTimeout, errRedis, FAILED)
			efLogger.Fatal(err)
		}

		efLogger.Info(OK, infoRedis, fmt.Sprintf("[%s] @%s", strings.ToUpper(redisDriver), redisHost))
		return client
	}

	clusterClient, ok := univClient.(*redis.ClusterClient)
	if ok {
		redisDriver := `cluster`

		client = &redisc{
			mu:      &sync.Mutex{},
			mode:    cluster,
			logger:  efLogger,
			locker:  NewLockClient(clusterClient),
			cluster: clusterClient,
			opt:     opt,
			driver:  redisDriver,
			termSig: make(chan struct{}, 1),
		}
		if err := client.Connect(); err != nil {
			err = x.WrapWithCode(err, EcodeConnectTimeout, errRedis, FAILED)
			efLogger.Fatal(err)
		}

		efLogger.Info(OK, infoRedis, fmt.Sprintf("[%s] host=%s", strings.ToUpper(redisDriver), redisHost))
		return client
	}

	err := x.WrapWithCode(ErrUnknownMode, EcodeUnknownMode, errRedis, FAILED)
	efLogger.Fatal(err)
	return nil
}

func (c *redisc) Stop() {
	c.logger.Info("Shutting Down redis connection")
	if c.mode == cluster {
		if err := c.cluster.Close(); err != nil {
			err = x.WrapWithCode(err, EcodeCloseTimeout, errRedis, FAILED)
			c.logger.Error(err)
		}
	} else {
		if err := c.sentinel.Close(); err != nil {
			err = x.WrapWithCode(err, EcodeCloseTimeout, errRedis, FAILED)
			c.logger.Error(err)
		}
	}
	close(c.termSig)
	c.logger.Info("[OK]: Shutdown Redis connection")
}

func (c *redisc) Connect() error {
	// no need to implement retry backoff
	// backoff retries is handled by redisc
	if c.mode == cluster {
		return c.cluster.Ping(context.Background()).Err()
	}
	return c.sentinel.Ping(context.Background()).Err()
}

func (c *redisc) Do(ctx context.Context) Cmd {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.mode != cluster {
		return c.sentinel.WithContext(ctx)
	}
	return c.cluster.WithContext(ctx)
}

func (c *redisc) Mock() redismock.ClientMock {
	return c.mock
}

func (c *redisc) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), (2 * time.Second))
	defer cancel()

	var pingErr error
	switch c.mode {
	case cluster:
		pingErr = c.cluster.Ping(ctx).Err()
	default:
		pingErr = c.sentinel.Ping(ctx).Err()
	}

	return pingErr
}
