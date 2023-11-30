// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/msmkdenis/yap-shortener/internal/handlers (interfaces: URLService)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	dto "github.com/msmkdenis/yap-shortener/internal/handlers/dto"
	model "github.com/msmkdenis/yap-shortener/internal/model"
)

// MockURLService is a mock of URLService interface.
type MockURLService struct {
	ctrl     *gomock.Controller
	recorder *MockURLServiceMockRecorder
}

// MockURLServiceMockRecorder is the mock recorder for MockURLService.
type MockURLServiceMockRecorder struct {
	mock *MockURLService
}

// NewMockURLService creates a new mock instance.
func NewMockURLService(ctrl *gomock.Controller) *MockURLService {
	mock := &MockURLService{ctrl: ctrl}
	mock.recorder = &MockURLServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockURLService) EXPECT() *MockURLServiceMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockURLService) Add(arg0 context.Context, arg1, arg2, arg3 string) (*model.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*model.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Add indicates an expected call of Add.
func (mr *MockURLServiceMockRecorder) Add(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockURLService)(nil).Add), arg0, arg1, arg2, arg3)
}

// AddAll mocks base method.
func (m *MockURLService) AddAll(arg0 context.Context, arg1 []dto.URLBatchRequest, arg2, arg3 string) ([]dto.URLBatchResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddAll", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].([]dto.URLBatchResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddAll indicates an expected call of AddAll.
func (mr *MockURLServiceMockRecorder) AddAll(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddAll", reflect.TypeOf((*MockURLService)(nil).AddAll), arg0, arg1, arg2, arg3)
}

// DeleteAll mocks base method.
func (m *MockURLService) DeleteAll(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAll", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAll indicates an expected call of DeleteAll.
func (mr *MockURLServiceMockRecorder) DeleteAll(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAll", reflect.TypeOf((*MockURLService)(nil).DeleteAll), arg0)
}

// DeleteURLByUserID mocks base method.
func (m *MockURLService) DeleteURLByUserID(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteURLByUserID", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteURLByUserID indicates an expected call of DeleteURLByUserID.
func (mr *MockURLServiceMockRecorder) DeleteURLByUserID(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteURLByUserID", reflect.TypeOf((*MockURLService)(nil).DeleteURLByUserID), arg0, arg1, arg2)
}

// GetAll mocks base method.
func (m *MockURLService) GetAll(arg0 context.Context) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", arg0)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockURLServiceMockRecorder) GetAll(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockURLService)(nil).GetAll), arg0)
}

// GetAllByUserID mocks base method.
func (m *MockURLService) GetAllByUserID(arg0 context.Context, arg1 string) ([]dto.URLBatchResponseByUserID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllByUserID", arg0, arg1)
	ret0, _ := ret[0].([]dto.URLBatchResponseByUserID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllByUserID indicates an expected call of GetAllByUserID.
func (mr *MockURLServiceMockRecorder) GetAllByUserID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllByUserID", reflect.TypeOf((*MockURLService)(nil).GetAllByUserID), arg0, arg1)
}

// GetByyID mocks base method.
func (m *MockURLService) GetByyID(arg0 context.Context, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByyID", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByyID indicates an expected call of GetByyID.
func (mr *MockURLServiceMockRecorder) GetByyID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByyID", reflect.TypeOf((*MockURLService)(nil).GetByyID), arg0, arg1)
}

// Ping mocks base method.
func (m *MockURLService) Ping(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockURLServiceMockRecorder) Ping(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockURLService)(nil).Ping), arg0)
}
