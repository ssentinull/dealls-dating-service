package cache

import (
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/stacktrace"
)

const (
	EcodeUnknownMode = stacktrace.ErrorCode(iota)
	EcodeConnectTimeout
	EcodeCloseTimeout
	EcodeTelemetry
)

const (
	errRedis                  string = `%sRedis Error`
	errTelemetryRedisRecorder string = `Redis Telemetry Recorder Error`
)

var ErrUnknownMode error = stacktrace.New(`Unknown Redis Mode`)
