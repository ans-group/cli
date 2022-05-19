// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ukfast/sdk-go/pkg/service/cloudflare (interfaces: CloudflareService)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	connection "github.com/ukfast/sdk-go/pkg/connection"
	cloudflare "github.com/ukfast/sdk-go/pkg/service/cloudflare"
)

// MockCloudflareService is a mock of CloudflareService interface.
type MockCloudflareService struct {
	ctrl     *gomock.Controller
	recorder *MockCloudflareServiceMockRecorder
}

// MockCloudflareServiceMockRecorder is the mock recorder for MockCloudflareService.
type MockCloudflareServiceMockRecorder struct {
	mock *MockCloudflareService
}

// NewMockCloudflareService creates a new mock instance.
func NewMockCloudflareService(ctrl *gomock.Controller) *MockCloudflareService {
	mock := &MockCloudflareService{ctrl: ctrl}
	mock.recorder = &MockCloudflareServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCloudflareService) EXPECT() *MockCloudflareServiceMockRecorder {
	return m.recorder
}

// CreateAccount mocks base method.
func (m *MockCloudflareService) CreateAccount(arg0 cloudflare.CreateAccountRequest) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAccount", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAccount indicates an expected call of CreateAccount.
func (mr *MockCloudflareServiceMockRecorder) CreateAccount(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccount", reflect.TypeOf((*MockCloudflareService)(nil).CreateAccount), arg0)
}

// CreateAccountMember mocks base method.
func (m *MockCloudflareService) CreateAccountMember(arg0 string, arg1 cloudflare.CreateAccountMemberRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAccountMember", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateAccountMember indicates an expected call of CreateAccountMember.
func (mr *MockCloudflareServiceMockRecorder) CreateAccountMember(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccountMember", reflect.TypeOf((*MockCloudflareService)(nil).CreateAccountMember), arg0, arg1)
}

// CreateOrchestration mocks base method.
func (m *MockCloudflareService) CreateOrchestration(arg0 cloudflare.CreateOrchestrationRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrchestration", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateOrchestration indicates an expected call of CreateOrchestration.
func (mr *MockCloudflareServiceMockRecorder) CreateOrchestration(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrchestration", reflect.TypeOf((*MockCloudflareService)(nil).CreateOrchestration), arg0)
}

// CreateZone mocks base method.
func (m *MockCloudflareService) CreateZone(arg0 cloudflare.CreateZoneRequest) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateZone", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateZone indicates an expected call of CreateZone.
func (mr *MockCloudflareServiceMockRecorder) CreateZone(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateZone", reflect.TypeOf((*MockCloudflareService)(nil).CreateZone), arg0)
}

// DeleteZone mocks base method.
func (m *MockCloudflareService) DeleteZone(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteZone", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteZone indicates an expected call of DeleteZone.
func (mr *MockCloudflareServiceMockRecorder) DeleteZone(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteZone", reflect.TypeOf((*MockCloudflareService)(nil).DeleteZone), arg0)
}

// GetAccount mocks base method.
func (m *MockCloudflareService) GetAccount(arg0 string) (cloudflare.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccount", arg0)
	ret0, _ := ret[0].(cloudflare.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccount indicates an expected call of GetAccount.
func (mr *MockCloudflareServiceMockRecorder) GetAccount(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccount", reflect.TypeOf((*MockCloudflareService)(nil).GetAccount), arg0)
}

// GetAccounts mocks base method.
func (m *MockCloudflareService) GetAccounts(arg0 connection.APIRequestParameters) ([]cloudflare.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccounts", arg0)
	ret0, _ := ret[0].([]cloudflare.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccounts indicates an expected call of GetAccounts.
func (mr *MockCloudflareServiceMockRecorder) GetAccounts(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccounts", reflect.TypeOf((*MockCloudflareService)(nil).GetAccounts), arg0)
}

// GetAccountsPaginated mocks base method.
func (m *MockCloudflareService) GetAccountsPaginated(arg0 connection.APIRequestParameters) (*connection.Paginated[cloudflare.Account], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccountsPaginated", arg0)
	ret0, _ := ret[0].(*connection.Paginated[cloudflare.Account])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccountsPaginated indicates an expected call of GetAccountsPaginated.
func (mr *MockCloudflareServiceMockRecorder) GetAccountsPaginated(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccountsPaginated", reflect.TypeOf((*MockCloudflareService)(nil).GetAccountsPaginated), arg0)
}

// GetSpendPlans mocks base method.
func (m *MockCloudflareService) GetSpendPlans(arg0 connection.APIRequestParameters) ([]cloudflare.SpendPlan, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSpendPlans", arg0)
	ret0, _ := ret[0].([]cloudflare.SpendPlan)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSpendPlans indicates an expected call of GetSpendPlans.
func (mr *MockCloudflareServiceMockRecorder) GetSpendPlans(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSpendPlans", reflect.TypeOf((*MockCloudflareService)(nil).GetSpendPlans), arg0)
}

// GetSpendPlansPaginated mocks base method.
func (m *MockCloudflareService) GetSpendPlansPaginated(arg0 connection.APIRequestParameters) (*connection.Paginated[cloudflare.SpendPlan], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSpendPlansPaginated", arg0)
	ret0, _ := ret[0].(*connection.Paginated[cloudflare.SpendPlan])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSpendPlansPaginated indicates an expected call of GetSpendPlansPaginated.
func (mr *MockCloudflareServiceMockRecorder) GetSpendPlansPaginated(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSpendPlansPaginated", reflect.TypeOf((*MockCloudflareService)(nil).GetSpendPlansPaginated), arg0)
}

// GetSubscriptions mocks base method.
func (m *MockCloudflareService) GetSubscriptions(arg0 connection.APIRequestParameters) ([]cloudflare.Subscription, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSubscriptions", arg0)
	ret0, _ := ret[0].([]cloudflare.Subscription)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSubscriptions indicates an expected call of GetSubscriptions.
func (mr *MockCloudflareServiceMockRecorder) GetSubscriptions(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubscriptions", reflect.TypeOf((*MockCloudflareService)(nil).GetSubscriptions), arg0)
}

// GetSubscriptionsPaginated mocks base method.
func (m *MockCloudflareService) GetSubscriptionsPaginated(arg0 connection.APIRequestParameters) (*connection.Paginated[cloudflare.Subscription], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSubscriptionsPaginated", arg0)
	ret0, _ := ret[0].(*connection.Paginated[cloudflare.Subscription])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSubscriptionsPaginated indicates an expected call of GetSubscriptionsPaginated.
func (mr *MockCloudflareServiceMockRecorder) GetSubscriptionsPaginated(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubscriptionsPaginated", reflect.TypeOf((*MockCloudflareService)(nil).GetSubscriptionsPaginated), arg0)
}

// GetTotalSpendMonthToDate mocks base method.
func (m *MockCloudflareService) GetTotalSpendMonthToDate() (cloudflare.TotalSpend, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTotalSpendMonthToDate")
	ret0, _ := ret[0].(cloudflare.TotalSpend)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTotalSpendMonthToDate indicates an expected call of GetTotalSpendMonthToDate.
func (mr *MockCloudflareServiceMockRecorder) GetTotalSpendMonthToDate() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTotalSpendMonthToDate", reflect.TypeOf((*MockCloudflareService)(nil).GetTotalSpendMonthToDate))
}

// GetZone mocks base method.
func (m *MockCloudflareService) GetZone(arg0 string) (cloudflare.Zone, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetZone", arg0)
	ret0, _ := ret[0].(cloudflare.Zone)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetZone indicates an expected call of GetZone.
func (mr *MockCloudflareServiceMockRecorder) GetZone(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetZone", reflect.TypeOf((*MockCloudflareService)(nil).GetZone), arg0)
}

// GetZones mocks base method.
func (m *MockCloudflareService) GetZones(arg0 connection.APIRequestParameters) ([]cloudflare.Zone, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetZones", arg0)
	ret0, _ := ret[0].([]cloudflare.Zone)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetZones indicates an expected call of GetZones.
func (mr *MockCloudflareServiceMockRecorder) GetZones(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetZones", reflect.TypeOf((*MockCloudflareService)(nil).GetZones), arg0)
}

// GetZonesPaginated mocks base method.
func (m *MockCloudflareService) GetZonesPaginated(arg0 connection.APIRequestParameters) (*connection.Paginated[cloudflare.Zone], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetZonesPaginated", arg0)
	ret0, _ := ret[0].(*connection.Paginated[cloudflare.Zone])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetZonesPaginated indicates an expected call of GetZonesPaginated.
func (mr *MockCloudflareServiceMockRecorder) GetZonesPaginated(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetZonesPaginated", reflect.TypeOf((*MockCloudflareService)(nil).GetZonesPaginated), arg0)
}

// PatchAccount mocks base method.
func (m *MockCloudflareService) PatchAccount(arg0 string, arg1 cloudflare.PatchAccountRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PatchAccount", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// PatchAccount indicates an expected call of PatchAccount.
func (mr *MockCloudflareServiceMockRecorder) PatchAccount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PatchAccount", reflect.TypeOf((*MockCloudflareService)(nil).PatchAccount), arg0, arg1)
}

// PatchZone mocks base method.
func (m *MockCloudflareService) PatchZone(arg0 string, arg1 cloudflare.PatchZoneRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PatchZone", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// PatchZone indicates an expected call of PatchZone.
func (mr *MockCloudflareServiceMockRecorder) PatchZone(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PatchZone", reflect.TypeOf((*MockCloudflareService)(nil).PatchZone), arg0, arg1)
}