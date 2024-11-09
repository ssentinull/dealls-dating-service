package cache

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	luaRefresh = redis.NewScript(`if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("pexpire", KEYS[1], ARGV[2]) else return 0 end`)
	luaRelease = redis.NewScript(`if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("del", KEYS[1]) else return 0 end`)
	luaPTTL    = redis.NewScript(`if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("pttl", KEYS[1]) else return -3 end`)
)

var (
	// ErrNotObtained is returned when a lock cannot be obtained.
	ErrNotObtained = errors.New("redislock: not obtained")

	// ErrLockNotHeld is returned when trying to release an inactive lock.
	ErrLockNotHeld = errors.New("redislock: lock not held")
)

// RedisLockClient is a minimal client interface.
type RedisLockClient interface {
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd
	Eval(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd
	EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) *redis.Cmd
	ScriptExists(ctx context.Context, scripts ...string) *redis.BoolSliceCmd
	ScriptLoad(ctx context.Context, script string) *redis.StringCmd
}

// Client wraps a redis client.
type LockClient struct {
	client RedisLockClient
	tmp    []byte
	tmpMu  sync.Mutex
}

// New creates a new Client instance with a custom namespace.
func NewLockClient(client RedisLockClient) *LockClient {
	return &LockClient{client: client}
}

// Obtain tries to obtain a new lock using a key with the given TTL.
// May return ErrNotObtained if not successful.
func (c *LockClient) Obtain(key string, ttl time.Duration, opt *LockOptions) (*Lock, error) {
	// Create a random token
	token, err := c.randomToken()
	if err != nil {
		return nil, err
	}

	value := token + opt.getMetadata()
	ctx := opt.getContext()

	var backoff *time.Timer
	for i, attempts := 0, opt.getRetryCount()+1; i < attempts; i++ {
		ok, err := c.obtain(ctx, key, value, ttl)
		if err != nil {
			return nil, err
		} else if ok {
			return &Lock{client: c, key: key, value: value}, nil
		}

		if backoff == nil {
			backoff = time.NewTimer(opt.getRetryBackoff())
			defer backoff.Stop()
		} else {
			backoff.Reset(opt.getRetryBackoff())
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-backoff.C:
		}
	}
	return nil, ErrNotObtained
}

func (c *LockClient) obtain(ctx context.Context, key, value string, ttl time.Duration) (bool, error) {
	ok, err := c.client.SetNX(ctx, key, value, ttl).Result()
	if errors.Is(err, redis.Nil) {
		err = nil
	}
	return ok, err
}

func (c *LockClient) randomToken() (string, error) {
	c.tmpMu.Lock()
	defer c.tmpMu.Unlock()

	if len(c.tmp) == 0 {
		c.tmp = make([]byte, 16)
	}

	if _, err := io.ReadFull(rand.Reader, c.tmp); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(c.tmp), nil
}

// --------------------------------------------------------------------

// Lock represents an obtained, distributed lock.
type Lock struct {
	client *LockClient
	key    string
	value  string
}

// Obtain is a short-cut for New(...).Obtain(...).
func Obtain(client RedisLockClient, key string, ttl time.Duration, opt *LockOptions) (*Lock, error) {
	return NewLockClient(client).Obtain(key, ttl, opt)
}

// Key returns the redis key used by the lock.
func (l *Lock) Key() string {
	return l.key
}

// Token returns the token value set by the lock.
func (l *Lock) Token() string {
	return l.value[:22]
}

// Metadata returns the metadata of the lock.
func (l *Lock) Metadata() string {
	return l.value[22:]
}

// TTL returns the remaining time-to-live. Returns 0 if the lock has expired.
func (l *Lock) TTL(ctx context.Context) (time.Duration, error) {
	res, err := luaPTTL.Eval(ctx, l.client.client, []string{l.key}, l.value).Result()
	if errors.Is(err, redis.Nil) {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	if num := res.(int64); num > 0 {
		return time.Duration(num) * time.Millisecond, nil
	}
	return 0, nil
}

// Refresh extends the lock with a new TTL.
// May return ErrNotObtained if refresh is unsuccessful.
func (l *Lock) Refresh(ctx context.Context, ttl time.Duration, opt *LockOptions) error {
	ttlVal := strconv.FormatInt(int64(ttl/time.Millisecond), 10)
	status, err := luaRefresh.Eval(ctx, l.client.client, []string{l.key}, l.value, ttlVal).Result()
	if err != nil {
		return err
	} else if status == int64(1) {
		return nil
	}
	return ErrNotObtained
}

// Release manually releases the lock.
// May return ErrLockNotHeld.
func (l *Lock) Release(ctx context.Context) error {
	res, err := luaRelease.Eval(ctx, l.client.client, []string{l.key}, l.value).Result()
	if errors.Is(err, redis.Nil) {
		return ErrLockNotHeld
	} else if err != nil {
		return err
	}

	if i, ok := res.(int64); !ok || i != 1 {
		return ErrLockNotHeld
	}
	return nil
}

// --------------------------------------------------------------------

// LockOptions describe the options for the lock
type LockOptions struct {
	// The number of time the acquisition of a lock will be retried.
	// Default: 0 = do not retry
	RetryCount int

	// RetryBackoff is the amount of time to wait between retries.
	// Default: 100ms
	RetryBackoff time.Duration

	// Metadata string is appended to the lock token.
	Metadata string

	// Optional context for Obtain timeout and cancellation control.
	Context context.Context
}

func (o *LockOptions) getRetryCount() int {
	if o != nil && o.RetryCount > 0 {
		return o.RetryCount
	}
	return 0
}

func (o *LockOptions) getRetryBackoff() time.Duration {
	if o != nil && o.RetryBackoff > 0 {
		return o.RetryBackoff
	}
	return 100 * time.Millisecond
}

func (o *LockOptions) getMetadata() string {
	if o != nil {
		return o.Metadata
	}
	return ""
}

func (o *LockOptions) getContext() context.Context {
	if o != nil && o.Context != nil {
		return o.Context
	}
	return context.Background()
}
