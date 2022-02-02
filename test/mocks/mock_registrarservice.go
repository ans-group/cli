// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ukfast/sdk-go/pkg/service/registrar (interfaces: RegistrarService)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	connection "github.com/ukfast/sdk-go/pkg/connection"
	registrar "github.com/ukfast/sdk-go/pkg/service/registrar"
)

// MockRegistrarService is a mock of RegistrarService interface.
type MockRegistrarService struct {
	ctrl     *gomock.Controller
	recorder *MockRegistrarServiceMockRecorder
}

// MockRegistrarServiceMockRecorder is the mock recorder for MockRegistrarService.
type MockRegistrarServiceMockRecorder struct {
	mock *MockRegistrarService
}

// NewMockRegistrarService creates a new mock instance.
func NewMockRegistrarService(ctrl *gomock.Controller) *MockRegistrarService {
	mock := &MockRegistrarService{ctrl: ctrl}
	mock.recorder = &MockRegistrarServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRegistrarService) EXPECT() *MockRegistrarServiceMockRecorder {
	return m.recorder
}

// GetDomain mocks base method.
func (m *MockRegistrarService) GetDomain(arg0 string) (registrar.Domain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDomain", arg0)
	ret0, _ := ret[0].(registrar.Domain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDomain indicates an expected call of GetDomain.
func (mr *MockRegistrarServiceMockRecorder) GetDomain(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDomain", reflect.TypeOf((*MockRegistrarService)(nil).GetDomain), arg0)
}

// GetDomainNameservers mocks base method.
func (m *MockRegistrarService) GetDomainNameservers(arg0 string) ([]registrar.Nameserver, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDomainNameservers", arg0)
	ret0, _ := ret[0].([]registrar.Nameserver)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDomainNameservers indicates an expected call of GetDomainNameservers.
func (mr *MockRegistrarServiceMockRecorder) GetDomainNameservers(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDomainNameservers", reflect.TypeOf((*MockRegistrarService)(nil).GetDomainNameservers), arg0)
}

// GetDomains mocks base method.
func (m *MockRegistrarService) GetDomains(arg0 connection.APIRequestParameters) ([]registrar.Domain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDomains", arg0)
	ret0, _ := ret[0].([]registrar.Domain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDomains indicates an expected call of GetDomains.
func (mr *MockRegistrarServiceMockRecorder) GetDomains(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDomains", reflect.TypeOf((*MockRegistrarService)(nil).GetDomains), arg0)
}

// GetDomainsPaginated mocks base method.
func (m *MockRegistrarService) GetDomainsPaginated(arg0 connection.APIRequestParameters) (*registrar.PaginatedDomain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDomainsPaginated", arg0)
	ret0, _ := ret[0].(*registrar.PaginatedDomain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDomainsPaginated indicates an expected call of GetDomainsPaginated.
func (mr *MockRegistrarServiceMockRecorder) GetDomainsPaginated(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDomainsPaginated", reflect.TypeOf((*MockRegistrarService)(nil).GetDomainsPaginated), arg0)
}

// GetWhois mocks base method.
func (m *MockRegistrarService) GetWhois(arg0 string) (registrar.Whois, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWhois", arg0)
	ret0, _ := ret[0].(registrar.Whois)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWhois indicates an expected call of GetWhois.
func (mr *MockRegistrarServiceMockRecorder) GetWhois(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWhois", reflect.TypeOf((*MockRegistrarService)(nil).GetWhois), arg0)
}

// GetWhoisRaw mocks base method.
func (m *MockRegistrarService) GetWhoisRaw(arg0 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWhoisRaw", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWhoisRaw indicates an expected call of GetWhoisRaw.
func (mr *MockRegistrarServiceMockRecorder) GetWhoisRaw(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWhoisRaw", reflect.TypeOf((*MockRegistrarService)(nil).GetWhoisRaw), arg0)
}
