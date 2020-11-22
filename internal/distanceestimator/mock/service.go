// Code generated by MockGen. DO NOT EDIT.
// Source: ./service.go

// Package mock_distanceestimator is a generated GoMock package.
package mock_distanceestimator

import (
	context "context"
	reflect "reflect"

	distanceestimator "github.com/edusalguero/roteiro.git/internal/distanceestimator"
	point "github.com/edusalguero/roteiro.git/internal/point"
	gomock "github.com/golang/mock/gomock"
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

// EstimateRouteDistances mocks base method
func (m *MockService) EstimateDistance(ctx context.Context, from, to point.Point) (*distanceestimator.RouteEstimation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EstimateDistance", ctx, from, to)
	ret0, _ := ret[0].(*distanceestimator.RouteEstimation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EstimateRouteDistances indicates an expected call of EstimateRouteDistances
func (mr *MockServiceMockRecorder) EstimateRouteDistances(ctx, from, to interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EstimateDistance", reflect.TypeOf((*MockService)(nil).EstimateDistance), ctx, from, to)
}
