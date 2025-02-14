// Code generated by MockGen. DO NOT EDIT.
// Source: ./service.go

// Package mock_distanceestimator is a generated GoMock package.
package mock_distanceestimator

import (
	context "context"
	cost "github.com/edusalguero/roteiro.git/internal/cost"
	point "github.com/edusalguero/roteiro.git/internal/point"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockService is a mock of Service interface
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// GetCost mocks base method
func (m *MockService) GetCost(ctx context.Context, from, to point.Point) (*cost.Cost, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCost", ctx, from, to)
	ret0, _ := ret[0].(*cost.Cost)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCost indicates an expected call of GetCost
func (mr *MockServiceMockRecorder) GetCost(ctx, from, to interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCost", reflect.TypeOf((*MockService)(nil).GetCost), ctx, from, to)
}
