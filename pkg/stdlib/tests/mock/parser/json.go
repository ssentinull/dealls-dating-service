// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/stdlib/parser/json.go
//
// Generated by this command:
//
//	mockgen -destination=pkg/stdlib/tests/mock/parser/json.go -package=mocks -source=pkg/stdlib/parser/json.go JsonParser
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockJSONParser is a mock of JSONParser interface.
type MockJSONParser struct {
	ctrl     *gomock.Controller
	recorder *MockJSONParserMockRecorder
	isgomock struct{}
}

// MockJSONParserMockRecorder is the mock recorder for MockJSONParser.
type MockJSONParserMockRecorder struct {
	mock *MockJSONParser
}

// NewMockJSONParser creates a new mock instance.
func NewMockJSONParser(ctrl *gomock.Controller) *MockJSONParser {
	mock := &MockJSONParser{ctrl: ctrl}
	mock.recorder = &MockJSONParserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockJSONParser) EXPECT() *MockJSONParserMockRecorder {
	return m.recorder
}

// Marshal mocks base method.
func (m *MockJSONParser) Marshal(orig any) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Marshal", orig)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Marshal indicates an expected call of Marshal.
func (mr *MockJSONParserMockRecorder) Marshal(orig any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Marshal", reflect.TypeOf((*MockJSONParser)(nil).Marshal), orig)
}

// Unmarshal mocks base method.
func (m *MockJSONParser) Unmarshal(blob []byte, dest any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unmarshal", blob, dest)
	ret0, _ := ret[0].(error)
	return ret0
}

// Unmarshal indicates an expected call of Unmarshal.
func (mr *MockJSONParserMockRecorder) Unmarshal(blob, dest any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unmarshal", reflect.TypeOf((*MockJSONParser)(nil).Unmarshal), blob, dest)
}
