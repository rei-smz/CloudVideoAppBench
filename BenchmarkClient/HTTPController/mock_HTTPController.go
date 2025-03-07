// Code generated by MockGen. DO NOT EDIT.
// Source: D:\GolandProjects\BenchmarkClient\HTTPController\HTTPController.go
//
// Generated by this command:
//
//	mockgen -source=D:\GolandProjects\BenchmarkClient\HTTPController\HTTPController.go -destination=mock_HTTPController .go -package=HTTPController
//

// Package mock_HTTPController is a generated GoMock package.
package HTTPController

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockHTTPController is a mock of HTTPController interface.
type MockHTTPController struct {
	ctrl     *gomock.Controller
	recorder *MockHTTPControllerMockRecorder
}

// MockHTTPControllerMockRecorder is the mock recorder for MockHTTPController.
type MockHTTPControllerMockRecorder struct {
	mock *MockHTTPController
}

// NewMockHTTPController creates a new mock instance.
func NewMockHTTPController(ctrl *gomock.Controller) *MockHTTPController {
	mock := &MockHTTPController{ctrl: ctrl}
	mock.recorder = &MockHTTPControllerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHTTPController) EXPECT() *MockHTTPControllerMockRecorder {
	return m.recorder
}

// Post mocks base method.
func (m *MockHTTPController) Post(url string, body map[string]string) map[string]any {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Post", url, body)
	ret0, _ := ret[0].(map[string]any)
	return ret0
}

// Post indicates an expected call of Post.
func (mr *MockHTTPControllerMockRecorder) Post(url, body any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Post", reflect.TypeOf((*MockHTTPController)(nil).Post), url, body)
}
