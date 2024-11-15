// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/stdlib/httpclient/httpclient.go
//
// Generated by this command:
//
//	mockgen -destination=pkg/stdlib/tests/mock/httpclient/httpclient.go -package=mocks -source=pkg/stdlib/httpclient/httpclient.go HttpClient
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	http "net/http"
	reflect "reflect"

	httpclient "github.com/ssentinull/dealls-dating-service/pkg/stdlib/httpclient"
	gomock "go.uber.org/mock/gomock"
)

// MockHTTPClient is a mock of HTTPClient interface.
type MockHTTPClient struct {
	ctrl     *gomock.Controller
	recorder *MockHTTPClientMockRecorder
	isgomock struct{}
}

// MockHTTPClientMockRecorder is the mock recorder for MockHTTPClient.
type MockHTTPClientMockRecorder struct {
	mock *MockHTTPClient
}

// NewMockHTTPClient creates a new mock instance.
func NewMockHTTPClient(ctrl *gomock.Controller) *MockHTTPClient {
	mock := &MockHTTPClient{ctrl: ctrl}
	mock.recorder = &MockHTTPClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHTTPClient) EXPECT() *MockHTTPClientMockRecorder {
	return m.recorder
}

// DeleteJSON mocks base method.
func (m *MockHTTPClient) DeleteJSON(ctx context.Context, prop *httpclient.RequestProp) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteJSON", ctx, prop)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteJSON indicates an expected call of DeleteJSON.
func (mr *MockHTTPClientMockRecorder) DeleteJSON(ctx, prop any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteJSON", reflect.TypeOf((*MockHTTPClient)(nil).DeleteJSON), ctx, prop)
}

// DeleteJSONWithoutTelemetry mocks base method.
func (m *MockHTTPClient) DeleteJSONWithoutTelemetry(ctx context.Context, prop *httpclient.RequestProp) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteJSONWithoutTelemetry", ctx, prop)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteJSONWithoutTelemetry indicates an expected call of DeleteJSONWithoutTelemetry.
func (mr *MockHTTPClientMockRecorder) DeleteJSONWithoutTelemetry(ctx, prop any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteJSONWithoutTelemetry", reflect.TypeOf((*MockHTTPClient)(nil).DeleteJSONWithoutTelemetry), ctx, prop)
}

// GetJSON mocks base method.
func (m *MockHTTPClient) GetJSON(ctx context.Context, prop *httpclient.RequestProp) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetJSON", ctx, prop)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetJSON indicates an expected call of GetJSON.
func (mr *MockHTTPClientMockRecorder) GetJSON(ctx, prop any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetJSON", reflect.TypeOf((*MockHTTPClient)(nil).GetJSON), ctx, prop)
}

// GetJSONWithoutTelemetry mocks base method.
func (m *MockHTTPClient) GetJSONWithoutTelemetry(ctx context.Context, prop *httpclient.RequestProp) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetJSONWithoutTelemetry", ctx, prop)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetJSONWithoutTelemetry indicates an expected call of GetJSONWithoutTelemetry.
func (mr *MockHTTPClientMockRecorder) GetJSONWithoutTelemetry(ctx, prop any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetJSONWithoutTelemetry", reflect.TypeOf((*MockHTTPClient)(nil).GetJSONWithoutTelemetry), ctx, prop)
}

// PatchJSON mocks base method.
func (m *MockHTTPClient) PatchJSON(ctx context.Context, prop *httpclient.RequestProp) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PatchJSON", ctx, prop)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PatchJSON indicates an expected call of PatchJSON.
func (mr *MockHTTPClientMockRecorder) PatchJSON(ctx, prop any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PatchJSON", reflect.TypeOf((*MockHTTPClient)(nil).PatchJSON), ctx, prop)
}

// PatchJSONWithoutTelemetry mocks base method.
func (m *MockHTTPClient) PatchJSONWithoutTelemetry(ctx context.Context, prop *httpclient.RequestProp) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PatchJSONWithoutTelemetry", ctx, prop)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PatchJSONWithoutTelemetry indicates an expected call of PatchJSONWithoutTelemetry.
func (mr *MockHTTPClientMockRecorder) PatchJSONWithoutTelemetry(ctx, prop any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PatchJSONWithoutTelemetry", reflect.TypeOf((*MockHTTPClient)(nil).PatchJSONWithoutTelemetry), ctx, prop)
}

// PostJSON mocks base method.
func (m *MockHTTPClient) PostJSON(ctx context.Context, prop *httpclient.RequestProp) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PostJSON", ctx, prop)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PostJSON indicates an expected call of PostJSON.
func (mr *MockHTTPClientMockRecorder) PostJSON(ctx, prop any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostJSON", reflect.TypeOf((*MockHTTPClient)(nil).PostJSON), ctx, prop)
}

// PostJSONWithoutTelemetry mocks base method.
func (m *MockHTTPClient) PostJSONWithoutTelemetry(ctx context.Context, prop *httpclient.RequestProp) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PostJSONWithoutTelemetry", ctx, prop)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PostJSONWithoutTelemetry indicates an expected call of PostJSONWithoutTelemetry.
func (mr *MockHTTPClientMockRecorder) PostJSONWithoutTelemetry(ctx, prop any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostJSONWithoutTelemetry", reflect.TypeOf((*MockHTTPClient)(nil).PostJSONWithoutTelemetry), ctx, prop)
}

// PutJSON mocks base method.
func (m *MockHTTPClient) PutJSON(ctx context.Context, prop *httpclient.RequestProp) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutJSON", ctx, prop)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PutJSON indicates an expected call of PutJSON.
func (mr *MockHTTPClientMockRecorder) PutJSON(ctx, prop any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutJSON", reflect.TypeOf((*MockHTTPClient)(nil).PutJSON), ctx, prop)
}

// PutJSONWithoutTelemetry mocks base method.
func (m *MockHTTPClient) PutJSONWithoutTelemetry(ctx context.Context, prop *httpclient.RequestProp) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutJSONWithoutTelemetry", ctx, prop)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PutJSONWithoutTelemetry indicates an expected call of PutJSONWithoutTelemetry.
func (mr *MockHTTPClientMockRecorder) PutJSONWithoutTelemetry(ctx, prop any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutJSONWithoutTelemetry", reflect.TypeOf((*MockHTTPClient)(nil).PutJSONWithoutTelemetry), ctx, prop)
}
