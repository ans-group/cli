// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ukfast/sdk-go/pkg/service/managedcloudflare (interfaces: ManagedCloudflareService)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	connection "github.com/ukfast/sdk-go/pkg/connection"
	managedcloudflare "github.com/ukfast/sdk-go/pkg/service/managedcloudflare"
)

// MockManagedCloudflareService is a mock of ManagedCloudflareService interface.
type MockManagedCloudflareService struct {
	ctrl     *gomock.Controller
	recorder *MockManagedCloudflareServiceMockRecorder
}

// MockManagedCloudflareServiceMockRecorder is the mock recorder for MockManagedCloudflareService.
type MockManagedCloudflareServiceMockRecorder struct {
	mock *MockManagedCloudflareService
}

// NewMockManagedCloudflareService creates a new mock instance.
func NewMockManagedCloudflareService(ctrl *gomock.Controller) *MockManagedCloudflareService {
	mock := &MockManagedCloudflareService{ctrl: ctrl}
	mock.recorder = &MockManagedCloudflareServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockManagedCloudflareService) EXPECT() *MockManagedCloudflareServiceMockRecorder {
	return m.recorder
}

// CreateAccount mocks base method.
func (m *MockManagedCloudflareService) CreateAccount(arg0 managedcloudflare.CreateAccountRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAccount", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateAccount indicates an expected call of CreateAccount.
func (mr *MockManagedCloudflareServiceMockRecorder) CreateAccount(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccount", reflect.TypeOf((*MockManagedCloudflareService)(nil).CreateAccount), arg0)
}

// CreateAccountMember mocks base method.
func (m *MockManagedCloudflareService) CreateAccountMember(arg0 string, arg1 managedcloudflare.CreateAccountMemberRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAccountMember", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateAccountMember indicates an expected call of CreateAccountMember.
func (mr *MockManagedCloudflareServiceMockRecorder) CreateAccountMember(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccountMember", reflect.TypeOf((*MockManagedCloudflareService)(nil).CreateAccountMember), arg0, arg1)
}

// CreateOrchestration mocks base method.
func (m *MockManagedCloudflareService) CreateOrchestration(arg0 managedcloudflare.CreateOrchestrationRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrchestration", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateOrchestration indicates an expected call of CreateOrchestration.
func (mr *MockManagedCloudflareServiceMockRecorder) CreateOrchestration(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrchestration", reflect.TypeOf((*MockManagedCloudflareService)(nil).CreateOrchestration), arg0)
}

// CreateZone mocks base method.
func (m *MockManagedCloudflareService) CreateZone(arg0 managedcloudflare.CreateZoneRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateZone", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateZone indicates an expected call of CreateZone.
func (mr *MockManagedCloudflareServiceMockRecorder) CreateZone(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateZone", reflect.TypeOf((*MockManagedCloudflareService)(nil).CreateZone), arg0)
}

// DeleteZone mocks base method.
func (m *MockManagedCloudflareService) DeleteZone(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteZone", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteZone indicates an expected call of DeleteZone.
func (mr *MockManagedCloudflareServiceMockRecorder) DeleteZone(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteZone", reflect.TypeOf((*MockManagedCloudflareService)(nil).DeleteZone), arg0)
}

// GetAccount mocks base method.
func (m *MockManagedCloudflareService) GetAccount(arg0 string) (managedcloudflare.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccount", arg0)
	ret0, _ := ret[0].(managedcloudflare.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccount indicates an expected call of GetAccount.
func (mr *MockManagedCloudflareServiceMockRecorder) GetAccount(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccount", reflect.TypeOf((*MockManagedCloudflareService)(nil).GetAccount), arg0)
}

// GetAccounts mocks base method.
func (m *MockManagedCloudflareService) GetAccounts(arg0 connection.APIRequestParameters) ([]managedcloudflare.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccounts", arg0)
	ret0, _ := ret[0].([]managedcloudflare.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccounts indicates an expected call of GetAccounts.
func (mr *MockManagedCloudflareServiceMockRecorder) GetAccounts(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccounts", reflect.TypeOf((*MockManagedCloudflareService)(nil).GetAccounts), arg0)
}

// GetAccountsPaginated mocks base method.
func (m *MockManagedCloudflareService) GetAccountsPaginated(arg0 connection.APIRequestParameters) (*managedcloudflare.PaginatedAccount, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccountsPaginated", arg0)
	ret0, _ := ret[0].(*managedcloudflare.PaginatedAccount)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccountsPaginated indicates an expected call of GetAccountsPaginated.
func (mr *MockManagedCloudflareServiceMockRecorder) GetAccountsPaginated(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccountsPaginated", reflect.TypeOf((*MockManagedCloudflareService)(nil).GetAccountsPaginated), arg0)
}

// GetSpendPlans mocks base method.
func (m *MockManagedCloudflareService) GetSpendPlans(arg0 connection.APIRequestParameters) ([]managedcloudflare.SpendPlan, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSpendPlans", arg0)
	ret0, _ := ret[0].([]managedcloudflare.SpendPlan)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSpendPlans indicates an expected call of GetSpendPlans.
func (mr *MockManagedCloudflareServiceMockRecorder) GetSpendPlans(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSpendPlans", reflect.TypeOf((*MockManagedCloudflareService)(nil).GetSpendPlans), arg0)
}

// GetSpendPlansPaginated mocks base method.
func (m *MockManagedCloudflareService) GetSpendPlansPaginated(arg0 connection.APIRequestParameters) (*managedcloudflare.PaginatedSpendPlan, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSpendPlansPaginated", arg0)
	ret0, _ := ret[0].(*managedcloudflare.PaginatedSpendPlan)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSpendPlansPaginated indicates an expected call of GetSpendPlansPaginated.
func (mr *MockManagedCloudflareServiceMockRecorder) GetSpendPlansPaginated(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSpendPlansPaginated", reflect.TypeOf((*MockManagedCloudflareService)(nil).GetSpendPlansPaginated), arg0)
}

// GetSubscriptions mocks base method.
func (m *MockManagedCloudflareService) GetSubscriptions(arg0 connection.APIRequestParameters) ([]managedcloudflare.Subscription, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSubscriptions", arg0)
	ret0, _ := ret[0].([]managedcloudflare.Subscription)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSubscriptions indicates an expected call of GetSubscriptions.
func (mr *MockManagedCloudflareServiceMockRecorder) GetSubscriptions(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubscriptions", reflect.TypeOf((*MockManagedCloudflareService)(nil).GetSubscriptions), arg0)
}

// GetSubscriptionsPaginated mocks base method.
func (m *MockManagedCloudflareService) GetSubscriptionsPaginated(arg0 connection.APIRequestParameters) (*managedcloudflare.PaginatedSubscription, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSubscriptionsPaginated", arg0)
	ret0, _ := ret[0].(*managedcloudflare.PaginatedSubscription)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSubscriptionsPaginated indicates an expected call of GetSubscriptionsPaginated.
func (mr *MockManagedCloudflareServiceMockRecorder) GetSubscriptionsPaginated(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubscriptionsPaginated", reflect.TypeOf((*MockManagedCloudflareService)(nil).GetSubscriptionsPaginated), arg0)
}

// GetZone mocks base method.
func (m *MockManagedCloudflareService) GetZone(arg0 string) (managedcloudflare.Zone, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetZone", arg0)
	ret0, _ := ret[0].(managedcloudflare.Zone)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetZone indicates an expected call of GetZone.
func (mr *MockManagedCloudflareServiceMockRecorder) GetZone(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetZone", reflect.TypeOf((*MockManagedCloudflareService)(nil).GetZone), arg0)
}

// GetZones mocks base method.
func (m *MockManagedCloudflareService) GetZones(arg0 connection.APIRequestParameters) ([]managedcloudflare.Zone, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetZones", arg0)
	ret0, _ := ret[0].([]managedcloudflare.Zone)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetZones indicates an expected call of GetZones.
func (mr *MockManagedCloudflareServiceMockRecorder) GetZones(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetZones", reflect.TypeOf((*MockManagedCloudflareService)(nil).GetZones), arg0)
}

// GetZonesPaginated mocks base method.
func (m *MockManagedCloudflareService) GetZonesPaginated(arg0 connection.APIRequestParameters) (*managedcloudflare.PaginatedZone, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetZonesPaginated", arg0)
	ret0, _ := ret[0].(*managedcloudflare.PaginatedZone)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetZonesPaginated indicates an expected call of GetZonesPaginated.
func (mr *MockManagedCloudflareServiceMockRecorder) GetZonesPaginated(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetZonesPaginated", reflect.TypeOf((*MockManagedCloudflareService)(nil).GetZonesPaginated), arg0)
}
