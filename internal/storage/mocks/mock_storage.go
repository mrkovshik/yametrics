// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mrkovshik/yametrics/internal/service (interfaces: Storage)

// Package mock_storage is a generated GoMock package.
package mock_storage

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/mrkovshik/yametrics/internal/model"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// GetAllMetrics mocks base method.
func (m *MockStorage) GetAllMetrics(arg0 context.Context) (map[string]model.Metrics, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllMetrics", arg0)
	ret0, _ := ret[0].(map[string]model.Metrics)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllMetrics indicates an expected call of GetAllMetrics.
func (mr *MockStorageMockRecorder) GetAllMetrics(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllMetrics", reflect.TypeOf((*MockStorage)(nil).GetAllMetrics), arg0)
}

// GetMetricByModel mocks base method.
func (m *MockStorage) GetMetricByModel(arg0 context.Context, arg1 model.Metrics) (model.Metrics, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMetricByModel", arg0, arg1)
	ret0, _ := ret[0].(model.Metrics)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMetricByModel indicates an expected call of GetMetricByModel.
func (mr *MockStorageMockRecorder) GetMetricByModel(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMetricByModel", reflect.TypeOf((*MockStorage)(nil).GetMetricByModel), arg0, arg1)
}

// Ping mocks base method.
func (m *MockStorage) Ping(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockStorageMockRecorder) Ping(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockStorage)(nil).Ping), arg0)
}

// RestoreMetrics mocks base method.
func (m *MockStorage) RestoreMetrics(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RestoreMetrics", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RestoreMetrics indicates an expected call of RestoreMetrics.
func (mr *MockStorageMockRecorder) RestoreMetrics(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RestoreMetrics", reflect.TypeOf((*MockStorage)(nil).RestoreMetrics), arg0, arg1)
}

// StoreMetrics mocks base method.
func (m *MockStorage) StoreMetrics(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreMetrics", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// StoreMetrics indicates an expected call of StoreMetrics.
func (mr *MockStorageMockRecorder) StoreMetrics(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreMetrics", reflect.TypeOf((*MockStorage)(nil).StoreMetrics), arg0, arg1)
}

// UpdateMetricValue mocks base method.
func (m *MockStorage) UpdateMetricValue(arg0 context.Context, arg1 model.Metrics) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMetricValue", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateMetricValue indicates an expected call of UpdateMetricValue.
func (mr *MockStorageMockRecorder) UpdateMetricValue(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMetricValue", reflect.TypeOf((*MockStorage)(nil).UpdateMetricValue), arg0, arg1)
}

// UpdateMetrics mocks base method.
func (m *MockStorage) UpdateMetrics(arg0 context.Context, arg1 []model.Metrics) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMetrics", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateMetrics indicates an expected call of UpdateMetrics.
func (mr *MockStorageMockRecorder) UpdateMetrics(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMetrics", reflect.TypeOf((*MockStorage)(nil).UpdateMetrics), arg0, arg1)
}
