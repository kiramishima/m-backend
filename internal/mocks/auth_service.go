// Code generated by MockGen. DO NOT EDIT.
// Source: .\internal\core\ports\services\auth_service.go
//
// Generated by this command:
//
//	mockgen -source .\internal\core\ports\services\auth_service.go -destination .\internal\mocks\auth_service.go -package mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	domain "kiramishima/m-backend/internal/core/domain"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockAuthService is a mock of AuthService interface.
type MockAuthService struct {
	ctrl     *gomock.Controller
	recorder *MockAuthServiceMockRecorder
}

// MockAuthServiceMockRecorder is the mock recorder for MockAuthService.
type MockAuthServiceMockRecorder struct {
	mock *MockAuthService
}

// NewMockAuthService creates a new mock instance.
func NewMockAuthService(ctrl *gomock.Controller) *MockAuthService {
	mock := &MockAuthService{ctrl: ctrl}
	mock.recorder = &MockAuthServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthService) EXPECT() *MockAuthServiceMockRecorder {
	return m.recorder
}

// FindByCredentials mocks base method.
func (m *MockAuthService) FindByCredentials(ctx context.Context, data *domain.AuthRequest) (*domain.AuthResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByCredentials", ctx, data)
	ret0, _ := ret[0].(*domain.AuthResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByCredentials indicates an expected call of FindByCredentials.
func (mr *MockAuthServiceMockRecorder) FindByCredentials(ctx, data any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByCredentials", reflect.TypeOf((*MockAuthService)(nil).FindByCredentials), ctx, data)
}

// Register mocks base method.
func (m *MockAuthService) Register(ctx context.Context, registerReq *domain.RegisterRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", ctx, registerReq)
	ret0, _ := ret[0].(error)
	return ret0
}

// Register indicates an expected call of Register.
func (mr *MockAuthServiceMockRecorder) Register(ctx, registerReq any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockAuthService)(nil).Register), ctx, registerReq)
}