// Code generated by MockGen. DO NOT EDIT.
// Source: ./service.go

// Package mock_solver is a generated GoMock package.
package mock_solver

import (
	context "context"
	problem "github.com/edusalguero/roteiro.git/internal/problem"
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

// SolveProblem mocks base method
func (m *MockService) SolveProblem(ctx context.Context, p problem.Problem) (*problem.Solution, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SolveProblem", ctx, p)
	ret0, _ := ret[0].(*problem.Solution)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SolveProblem indicates an expected call of SolveProblem
func (mr *MockServiceMockRecorder) SolveProblem(ctx, p interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SolveProblem", reflect.TypeOf((*MockService)(nil).SolveProblem), ctx, p)
}