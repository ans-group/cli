// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ukfast/sdk-go/pkg/service/billing (interfaces: BillingService)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	connection "github.com/ukfast/sdk-go/pkg/connection"
	billing "github.com/ukfast/sdk-go/pkg/service/billing"
	reflect "reflect"
)

// MockBillingService is a mock of BillingService interface
type MockBillingService struct {
	ctrl     *gomock.Controller
	recorder *MockBillingServiceMockRecorder
}

// MockBillingServiceMockRecorder is the mock recorder for MockBillingService
type MockBillingServiceMockRecorder struct {
	mock *MockBillingService
}

// NewMockBillingService creates a new mock instance
func NewMockBillingService(ctrl *gomock.Controller) *MockBillingService {
	mock := &MockBillingService{ctrl: ctrl}
	mock.recorder = &MockBillingServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBillingService) EXPECT() *MockBillingServiceMockRecorder {
	return m.recorder
}

// CreateCard mocks base method
func (m *MockBillingService) CreateCard(arg0 billing.CreateCardRequest) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCard", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCard indicates an expected call of CreateCard
func (mr *MockBillingServiceMockRecorder) CreateCard(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCard", reflect.TypeOf((*MockBillingService)(nil).CreateCard), arg0)
}

// CreateInvoiceQuery mocks base method
func (m *MockBillingService) CreateInvoiceQuery(arg0 billing.CreateInvoiceQueryRequest) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateInvoiceQuery", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateInvoiceQuery indicates an expected call of CreateInvoiceQuery
func (mr *MockBillingServiceMockRecorder) CreateInvoiceQuery(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateInvoiceQuery", reflect.TypeOf((*MockBillingService)(nil).CreateInvoiceQuery), arg0)
}

// DeleteCard mocks base method
func (m *MockBillingService) DeleteCard(arg0 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCard", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCard indicates an expected call of DeleteCard
func (mr *MockBillingServiceMockRecorder) DeleteCard(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCard", reflect.TypeOf((*MockBillingService)(nil).DeleteCard), arg0)
}

// GetCard mocks base method
func (m *MockBillingService) GetCard(arg0 int) (billing.Card, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCard", arg0)
	ret0, _ := ret[0].(billing.Card)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCard indicates an expected call of GetCard
func (mr *MockBillingServiceMockRecorder) GetCard(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCard", reflect.TypeOf((*MockBillingService)(nil).GetCard), arg0)
}

// GetCards mocks base method
func (m *MockBillingService) GetCards(arg0 connection.APIRequestParameters) ([]billing.Card, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCards", arg0)
	ret0, _ := ret[0].([]billing.Card)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCards indicates an expected call of GetCards
func (mr *MockBillingServiceMockRecorder) GetCards(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCards", reflect.TypeOf((*MockBillingService)(nil).GetCards), arg0)
}

// GetCardsPaginated mocks base method
func (m *MockBillingService) GetCardsPaginated(arg0 connection.APIRequestParameters) (*billing.PaginatedCard, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCardsPaginated", arg0)
	ret0, _ := ret[0].(*billing.PaginatedCard)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCardsPaginated indicates an expected call of GetCardsPaginated
func (mr *MockBillingServiceMockRecorder) GetCardsPaginated(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCardsPaginated", reflect.TypeOf((*MockBillingService)(nil).GetCardsPaginated), arg0)
}

// GetCloudCost mocks base method
func (m *MockBillingService) GetCloudCost(arg0 int) (billing.CloudCost, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCloudCost", arg0)
	ret0, _ := ret[0].(billing.CloudCost)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCloudCost indicates an expected call of GetCloudCost
func (mr *MockBillingServiceMockRecorder) GetCloudCost(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCloudCost", reflect.TypeOf((*MockBillingService)(nil).GetCloudCost), arg0)
}

// GetCloudCosts mocks base method
func (m *MockBillingService) GetCloudCosts(arg0 connection.APIRequestParameters) ([]billing.CloudCost, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCloudCosts", arg0)
	ret0, _ := ret[0].([]billing.CloudCost)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCloudCosts indicates an expected call of GetCloudCosts
func (mr *MockBillingServiceMockRecorder) GetCloudCosts(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCloudCosts", reflect.TypeOf((*MockBillingService)(nil).GetCloudCosts), arg0)
}

// GetCloudCostsPaginated mocks base method
func (m *MockBillingService) GetCloudCostsPaginated(arg0 connection.APIRequestParameters) (*billing.PaginatedCloudCost, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCloudCostsPaginated", arg0)
	ret0, _ := ret[0].(*billing.PaginatedCloudCost)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCloudCostsPaginated indicates an expected call of GetCloudCostsPaginated
func (mr *MockBillingServiceMockRecorder) GetCloudCostsPaginated(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCloudCostsPaginated", reflect.TypeOf((*MockBillingService)(nil).GetCloudCostsPaginated), arg0)
}

// GetDirectDebit mocks base method
func (m *MockBillingService) GetDirectDebit() (billing.DirectDebit, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDirectDebit")
	ret0, _ := ret[0].(billing.DirectDebit)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDirectDebit indicates an expected call of GetDirectDebit
func (mr *MockBillingServiceMockRecorder) GetDirectDebit() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDirectDebit", reflect.TypeOf((*MockBillingService)(nil).GetDirectDebit))
}

// GetInvoice mocks base method
func (m *MockBillingService) GetInvoice(arg0 int) (billing.Invoice, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInvoice", arg0)
	ret0, _ := ret[0].(billing.Invoice)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInvoice indicates an expected call of GetInvoice
func (mr *MockBillingServiceMockRecorder) GetInvoice(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInvoice", reflect.TypeOf((*MockBillingService)(nil).GetInvoice), arg0)
}

// GetInvoiceQueries mocks base method
func (m *MockBillingService) GetInvoiceQueries(arg0 connection.APIRequestParameters) ([]billing.InvoiceQuery, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInvoiceQueries", arg0)
	ret0, _ := ret[0].([]billing.InvoiceQuery)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInvoiceQueries indicates an expected call of GetInvoiceQueries
func (mr *MockBillingServiceMockRecorder) GetInvoiceQueries(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInvoiceQueries", reflect.TypeOf((*MockBillingService)(nil).GetInvoiceQueries), arg0)
}

// GetInvoiceQueriesPaginated mocks base method
func (m *MockBillingService) GetInvoiceQueriesPaginated(arg0 connection.APIRequestParameters) (*billing.PaginatedInvoiceQuery, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInvoiceQueriesPaginated", arg0)
	ret0, _ := ret[0].(*billing.PaginatedInvoiceQuery)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInvoiceQueriesPaginated indicates an expected call of GetInvoiceQueriesPaginated
func (mr *MockBillingServiceMockRecorder) GetInvoiceQueriesPaginated(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInvoiceQueriesPaginated", reflect.TypeOf((*MockBillingService)(nil).GetInvoiceQueriesPaginated), arg0)
}

// GetInvoiceQuery mocks base method
func (m *MockBillingService) GetInvoiceQuery(arg0 int) (billing.InvoiceQuery, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInvoiceQuery", arg0)
	ret0, _ := ret[0].(billing.InvoiceQuery)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInvoiceQuery indicates an expected call of GetInvoiceQuery
func (mr *MockBillingServiceMockRecorder) GetInvoiceQuery(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInvoiceQuery", reflect.TypeOf((*MockBillingService)(nil).GetInvoiceQuery), arg0)
}

// GetInvoices mocks base method
func (m *MockBillingService) GetInvoices(arg0 connection.APIRequestParameters) ([]billing.Invoice, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInvoices", arg0)
	ret0, _ := ret[0].([]billing.Invoice)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInvoices indicates an expected call of GetInvoices
func (mr *MockBillingServiceMockRecorder) GetInvoices(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInvoices", reflect.TypeOf((*MockBillingService)(nil).GetInvoices), arg0)
}

// GetInvoicesPaginated mocks base method
func (m *MockBillingService) GetInvoicesPaginated(arg0 connection.APIRequestParameters) (*billing.PaginatedInvoice, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInvoicesPaginated", arg0)
	ret0, _ := ret[0].(*billing.PaginatedInvoice)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInvoicesPaginated indicates an expected call of GetInvoicesPaginated
func (mr *MockBillingServiceMockRecorder) GetInvoicesPaginated(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInvoicesPaginated", reflect.TypeOf((*MockBillingService)(nil).GetInvoicesPaginated), arg0)
}

// GetPayment mocks base method
func (m *MockBillingService) GetPayment(arg0 int) (billing.Payment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPayment", arg0)
	ret0, _ := ret[0].(billing.Payment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPayment indicates an expected call of GetPayment
func (mr *MockBillingServiceMockRecorder) GetPayment(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPayment", reflect.TypeOf((*MockBillingService)(nil).GetPayment), arg0)
}

// GetPayments mocks base method
func (m *MockBillingService) GetPayments(arg0 connection.APIRequestParameters) ([]billing.Payment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPayments", arg0)
	ret0, _ := ret[0].([]billing.Payment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPayments indicates an expected call of GetPayments
func (mr *MockBillingServiceMockRecorder) GetPayments(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPayments", reflect.TypeOf((*MockBillingService)(nil).GetPayments), arg0)
}

// GetPaymentsPaginated mocks base method
func (m *MockBillingService) GetPaymentsPaginated(arg0 connection.APIRequestParameters) (*billing.PaginatedPayment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPaymentsPaginated", arg0)
	ret0, _ := ret[0].(*billing.PaginatedPayment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPaymentsPaginated indicates an expected call of GetPaymentsPaginated
func (mr *MockBillingServiceMockRecorder) GetPaymentsPaginated(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPaymentsPaginated", reflect.TypeOf((*MockBillingService)(nil).GetPaymentsPaginated), arg0)
}

// GetRecurringCost mocks base method
func (m *MockBillingService) GetRecurringCost(arg0 int) (billing.RecurringCost, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRecurringCost", arg0)
	ret0, _ := ret[0].(billing.RecurringCost)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRecurringCost indicates an expected call of GetRecurringCost
func (mr *MockBillingServiceMockRecorder) GetRecurringCost(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRecurringCost", reflect.TypeOf((*MockBillingService)(nil).GetRecurringCost), arg0)
}

// GetRecurringCosts mocks base method
func (m *MockBillingService) GetRecurringCosts(arg0 connection.APIRequestParameters) ([]billing.RecurringCost, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRecurringCosts", arg0)
	ret0, _ := ret[0].([]billing.RecurringCost)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRecurringCosts indicates an expected call of GetRecurringCosts
func (mr *MockBillingServiceMockRecorder) GetRecurringCosts(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRecurringCosts", reflect.TypeOf((*MockBillingService)(nil).GetRecurringCosts), arg0)
}

// GetRecurringCostsPaginated mocks base method
func (m *MockBillingService) GetRecurringCostsPaginated(arg0 connection.APIRequestParameters) (*billing.PaginatedRecurringCost, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRecurringCostsPaginated", arg0)
	ret0, _ := ret[0].(*billing.PaginatedRecurringCost)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRecurringCostsPaginated indicates an expected call of GetRecurringCostsPaginated
func (mr *MockBillingServiceMockRecorder) GetRecurringCostsPaginated(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRecurringCostsPaginated", reflect.TypeOf((*MockBillingService)(nil).GetRecurringCostsPaginated), arg0)
}

// PatchCard mocks base method
func (m *MockBillingService) PatchCard(arg0 int, arg1 billing.PatchCardRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PatchCard", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// PatchCard indicates an expected call of PatchCard
func (mr *MockBillingServiceMockRecorder) PatchCard(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PatchCard", reflect.TypeOf((*MockBillingService)(nil).PatchCard), arg0, arg1)
}