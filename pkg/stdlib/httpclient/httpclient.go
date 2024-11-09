package httpclient

import (
	"context"
	"net/http"
	"time"

	"github.com/ssentinull/golang-boilerplate/pkg/stdlib/logger"
	"github.com/ssentinull/golang-boilerplate/pkg/stdlib/parser"

	"github.com/gojek/heimdall/v7"
	"github.com/gojek/heimdall/v7/httpclient"
)

const (
	DefaultBackOffInterval       = 2
	DefaultMaximumJitterInterval = 5
	DefaultTimeout               = 10
	DefaultMaxRetryCount         = 3
)

type HTTPClient interface {
	GetJSON(ctx context.Context, prop *RequestProp) (*http.Response, error)
	PostJSON(ctx context.Context, prop *RequestProp) (*http.Response, error)
	PutJSON(ctx context.Context, prop *RequestProp) (*http.Response, error)
	PatchJSON(ctx context.Context, prop *RequestProp) (*http.Response, error)
	DeleteJSON(ctx context.Context, prop *RequestProp) (*http.Response, error)
	GetJSONWithoutTelemetry(ctx context.Context, prop *RequestProp) (*http.Response, error)
	PostJSONWithoutTelemetry(ctx context.Context, prop *RequestProp) (*http.Response, error)
	PutJSONWithoutTelemetry(ctx context.Context, prop *RequestProp) (*http.Response, error)
	PatchJSONWithoutTelemetry(ctx context.Context, prop *RequestProp) (*http.Response, error)
	DeleteJSONWithoutTelemetry(ctx context.Context, prop *RequestProp) (*http.Response, error)
}

type httpClient struct {
	client   *httpclient.Client
	efLogger logger.Logger
	parser   parser.Parser
	opt      Options
}

type Options struct {
	BackOffInterval       int
	MaximumJitterInterval int
	Timeout               int
	MaxRetryCount         int
	DefaultClientID       string `validate:"required"`
}

func Init(efLogger logger.Logger, parser parser.Parser, opt Options) HTTPClient {
	if opt.BackOffInterval == 0 {
		opt.BackOffInterval = DefaultBackOffInterval
	}

	if opt.MaximumJitterInterval == 0 {
		opt.MaximumJitterInterval = DefaultMaximumJitterInterval
	}

	if opt.Timeout == 0 {
		opt.Timeout = DefaultTimeout
	}

	if opt.MaxRetryCount == 0 {
		opt.MaxRetryCount = DefaultMaxRetryCount
	}

	// First set a backoff mechanism. Constant backoff increases the backoff at a constant rate
	backoffInterval := time.Duration(opt.BackOffInterval) * time.Millisecond

	// Define a maximum jitter interval. It must be more than 1*time.Millisecond
	maximumJitterInterval := time.Duration(opt.MaximumJitterInterval) * time.Millisecond

	backoff := heimdall.NewConstantBackoff(backoffInterval, maximumJitterInterval)

	// Create a new retry mechanism with the backoff
	retrier := heimdall.NewRetrier(backoff)

	timeout := time.Duration(opt.Timeout) * time.Second // 10 seconds
	maxRetryCount := opt.MaxRetryCount

	return &httpClient{
		client: httpclient.NewClient(
			httpclient.WithHTTPTimeout(timeout),
			httpclient.WithRetrier(retrier),
			httpclient.WithRetryCount(maxRetryCount),
			httpclient.WithHTTPClient(&http.Client{
				Transport: &http.Transport{
					MaxResponseHeaderBytes: 10 * 1024 * 1024, // Set the maximum response header size
					IdleConnTimeout:        2 * time.Minute,
				},
			}),
		),
		efLogger: efLogger,
		parser:   parser,
		opt:      opt,
	}
}
