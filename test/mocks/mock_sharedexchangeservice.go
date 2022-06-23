// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ans-group/sdk-go/pkg/service/sharedexchange (interfaces: SharedExchangeService)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	connection "github.com/ans-group/sdk-go/pkg/connection"
	sharedexchange "github.com/ans-group/sdk-go/pkg/service/sharedexchange"
)

// MockSharedExchangeService is a mock of SharedExchangeService interface.
type MockSharedExchangeService struct {
	ctrl     *gomock.Controller
	recorder *MockSharedExchangeServiceMockRecorder
}

// MockSharedExchangeServiceMockRecorder is the mock recorder for MockSharedExchangeService.
type MockSharedExchangeServiceMockRecorder struct {
	mock *MockSharedExchangeService
}

// NewMockSharedExchangeService creates a new mock instance.
func NewMockSharedExchangeService(ctrl *gomock.Controller) *MockSharedExchangeService {
	mock := &MockSharedExchangeService{ctrl: ctrl}
	mock.recorder = &MockSharedExchangeServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSharedExchangeService) EXPECT() *MockSharedExchangeServiceMockRecorder {
	return m.recorder
}

// GetDomain mocks base method.
func (m *MockSharedExchangeService) GetDomain(arg0 int) (sharedexchange.Domain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDomain", arg0)
	ret0, _ := ret[0].(sharedexchange.Domain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDomain indicates an expected call of GetDomain.
func (mr *MockSharedExchangeServiceMockRecorder) GetDomain(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDomain", reflect.TypeOf((*MockSharedExchangeService)(nil).GetDomain), arg0)
}

// GetDomains mocks base method.
func (m *MockSharedExchangeService) GetDomains(arg0 connection.APIRequestParameters) ([]sharedexchange.Domain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDomains", arg0)
	ret0, _ := ret[0].([]sharedexchange.Domain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDomains indicates an expected call of GetDomains.
func (mr *MockSharedExchangeServiceMockRecorder) GetDomains(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDomains", reflect.TypeOf((*MockSharedExchangeService)(nil).GetDomains), arg0)
}

// GetDomainsPaginated mocks base method.
func (m *MockSharedExchangeService) GetDomainsPaginated(arg0 connection.APIRequestParameters) (*connection.Paginated[sharedexchange.Domain], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDomainsPaginated", arg0)
	ret0, _ := ret[0].(*connection.Paginated[sharedexchange.Domain])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDomainsPaginated indicates an expected call of GetDomainsPaginated.
func (mr *MockSharedExchangeServiceMockRecorder) GetDomainsPaginated(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDomainsPaginated", reflect.TypeOf((*MockSharedExchangeService)(nil).GetDomainsPaginated), arg0)
}
