package httpclient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/ssentinull/dealls-dating-service/pkg/constants"
)

type RequestProp struct {
	OperationName            string
	URI                      string
	Headers                  map[string]string
	QueryParams              map[string]string
	AuthToken                string
	Body                     interface{}
	IgnoreContentType        bool
	SeparateMultiValueParams bool
	TruncateRespBodySizeInMB int
}

func (h *httpClient) GetJSON(ctx context.Context, prop *RequestProp) (*http.Response, error) {
	return h.requestJSON(ctx, http.MethodGet, prop)
}

func (h *httpClient) PostJSON(ctx context.Context, prop *RequestProp) (*http.Response, error) {
	return h.requestJSON(ctx, http.MethodPost, prop)
}

func (h *httpClient) PutJSON(ctx context.Context, prop *RequestProp) (*http.Response, error) {
	return h.requestJSON(ctx, http.MethodPut, prop)
}

func (h *httpClient) PatchJSON(ctx context.Context, prop *RequestProp) (*http.Response, error) {
	return h.requestJSON(ctx, http.MethodPatch, prop)
}

func (h *httpClient) DeleteJSON(ctx context.Context, prop *RequestProp) (*http.Response, error) {
	return h.requestJSON(ctx, http.MethodDelete, prop)
}

func (h *httpClient) GetJSONWithoutTelemetry(ctx context.Context, prop *RequestProp) (*http.Response, error) {
	return h.requestJSONWithoutTelemetry(ctx, http.MethodGet, prop)
}

func (h *httpClient) PostJSONWithoutTelemetry(ctx context.Context, prop *RequestProp) (*http.Response, error) {
	return h.requestJSONWithoutTelemetry(ctx, http.MethodPost, prop)
}

func (h *httpClient) PutJSONWithoutTelemetry(ctx context.Context, prop *RequestProp) (*http.Response, error) {
	return h.requestJSONWithoutTelemetry(ctx, http.MethodPut, prop)
}

func (h *httpClient) PatchJSONWithoutTelemetry(ctx context.Context, prop *RequestProp) (*http.Response, error) {
	return h.requestJSONWithoutTelemetry(ctx, http.MethodPatch, prop)
}

func (h *httpClient) DeleteJSONWithoutTelemetry(ctx context.Context, prop *RequestProp) (*http.Response, error) {
	return h.requestJSONWithoutTelemetry(ctx, http.MethodDelete, prop)
}

func (h *httpClient) requestJSON(ctx context.Context, method string, prop *RequestProp) (*http.Response, error) {
	var buf io.ReadWriter
	if prop.Body != nil && !reflect.ValueOf(prop.Body).IsNil() {
		body, err := h.parser.JSONParser().Marshal(prop.Body)
		if err != nil {
			return nil, err
		}
		buf = bytes.NewBuffer(body)
	}

	u, err := url.Parse(prop.URI)
	if err != nil {
		return nil, err
	}

	if len(prop.QueryParams) > 0 {
		q := u.Query()
		for k, v := range prop.QueryParams {
			if prop.SeparateMultiValueParams {
				for _, split := range strings.Split(v, ",") {
					q.Add(k, split)
				}
			} else {
				q.Add(k, v)
			}
		}
		u.RawQuery = q.Encode()
	}

	urlx := u.String()
	req, err := http.NewRequest(method, urlx, buf)
	if err != nil {
		return nil, err
	}

	if prop.Headers == nil {
		prop.Headers = make(map[string]string)
	}

	for key := range prop.Headers {
		req.Header.Add(key, prop.Headers[key])
	}

	if prop.AuthToken != "" && req.Header.Get(Authorization) == "" {
		req.Header.Add(Authorization, fmt.Sprintf("Bearer %s", prop.AuthToken))
	}

	// default content type
	if !prop.IgnoreContentType {
		req.Header.Add(ContentType, "application/json")
	}

	// assign X-Client-Id
	clientID, ok := ctx.Value(constants.XClientId).(string)
	if !ok {
		clientID = h.opt.DefaultClientID
	}

	if clientID == "" {
		clientID = h.opt.DefaultClientID
	}

	if req.Header.Get(ClientId) == "" {
		req.Header.Add(ClientId, clientID)
	}

	req = req.WithContext(ctx)
	var rawBody []byte

	if req.Body == nil || req.Body == http.NoBody {
		rawBody = []byte{}
	} else {
		defer req.Body.Close()
		rawBody, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(rawBody))
	}

	res, err := h.client.Do(req)
	if err != nil {
		return res, err
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	res.Body = io.NopCloser(bytes.NewBuffer(resBody))

	if err != nil {
		return res, err
	}

	return res, err
}

func (h *httpClient) requestJSONWithoutTelemetry(ctx context.Context, method string, prop *RequestProp) (*http.Response, error) {
	var buf io.ReadWriter
	if prop.Body != nil {
		body, err := h.parser.JSONParser().Marshal(prop.Body)
		if err != nil {
			return nil, err
		}

		buf = bytes.NewBuffer(body)
	}

	u, err := url.Parse(prop.URI)
	if err != nil {
		return nil, err
	}

	if len(prop.QueryParams) > 0 {
		q := u.Query()
		for k, v := range prop.QueryParams {
			q.Add(k, v)
		}
		u.RawQuery = q.Encode()
	}

	urlx := u.String()
	req, err := http.NewRequest(method, urlx, buf)
	if err != nil {
		return nil, err
	}

	if prop.Headers == nil {
		prop.Headers = make(map[string]string)
	}

	for key := range prop.Headers {
		req.Header.Add(key, prop.Headers[key])
	}

	if prop.AuthToken != "" {
		req.Header.Add(Authorization, fmt.Sprintf("Bearer %s", prop.AuthToken))
	}

	// request-id
	requestID := ctx.Value(constants.RequestId)
	req.Header.Add(RequestId, fmt.Sprintf("%v", requestID))

	// default content type
	req.Header.Add(ContentType, "application/json")

	// assign X-Client-Id
	clientID, ok := ctx.Value(constants.XClientId).(string)
	if !ok {
		clientID = h.opt.DefaultClientID
	}

	if clientID == "" {
		clientID = h.opt.DefaultClientID
	}

	req.Header.Add(ClientId, clientID)
	req.WithContext(ctx)

	return h.client.Do(req)
}
