package constants

type contextKey struct {
	name string
}

var (
	XAppId    = &contextKey{"X-App-ID"}
	XClientId = &contextKey{"X-Client-Id"}
	RequestId = &contextKey{"RequestId"}
)
