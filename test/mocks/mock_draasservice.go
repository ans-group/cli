// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ukfast/sdk-go/pkg/service/draas (interfaces: DRaaSService)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	connection "github.com/ukfast/sdk-go/pkg/connection"
	draas "github.com/ukfast/sdk-go/pkg/service/draas"
)

// MockDRaaSService is a mock of DRaaSService interface.
type MockDRaaSService struct {
	ctrl     *gomock.Controller
	recorder *MockDRaaSServiceMockRecorder
}

// MockDRaaSServiceMockRecorder is the mock recorder for MockDRaaSService.
type MockDRaaSServiceMockRecorder struct {
	mock *MockDRaaSService
}

// NewMockDRaaSService creates a new mock instance.
func NewMockDRaaSService(ctrl *gomock.Controller) *MockDRaaSService {
	mock := &MockDRaaSService{ctrl: ctrl}
	mock.recorder = &MockDRaaSServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDRaaSService) EXPECT() *MockDRaaSServiceMockRecorder {
	return m.recorder
}

// GetBillingType mocks base method.
func (m *MockDRaaSService) GetBillingType(arg0 string) (draas.BillingType, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBillingType", arg0)
	ret0, _ := ret[0].(draas.BillingType)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBillingType indicates an expected call of GetBillingType.
func (mr *MockDRaaSServiceMockRecorder) GetBillingType(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBillingType", reflect.TypeOf((*MockDRaaSService)(nil).GetBillingType), arg0)
}

// GetBillingTypes mocks base method.
func (m *MockDRaaSService) GetBillingTypes(arg0 connection.APIRequestParameters) ([]draas.BillingType, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBillingTypes", arg0)
	ret0, _ := ret[0].([]draas.BillingType)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBillingTypes indicates an expected call of GetBillingTypes.
func (mr *MockDRaaSServiceMockRecorder) GetBillingTypes(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBillingTypes", reflect.TypeOf((*MockDRaaSService)(nil).GetBillingTypes), arg0)
}

// GetBillingTypesPaginated mocks base method.
func (m *MockDRaaSService) GetBillingTypesPaginated(arg0 connection.APIRequestParameters) (*connection.Paginated[draas.BillingType], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBillingTypesPaginated", arg0)
	ret0, _ := ret[0].(*connection.Paginated[draas.BillingType])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBillingTypesPaginated indicates an expected call of GetBillingTypesPaginated.
func (mr *MockDRaaSServiceMockRecorder) GetBillingTypesPaginated(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBillingTypesPaginated", reflect.TypeOf((*MockDRaaSService)(nil).GetBillingTypesPaginated), arg0)
}

// GetIOPSTier mocks base method.
func (m *MockDRaaSService) GetIOPSTier(arg0 string) (draas.IOPSTier, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIOPSTier", arg0)
	ret0, _ := ret[0].(draas.IOPSTier)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetIOPSTier indicates an expected call of GetIOPSTier.
func (mr *MockDRaaSServiceMockRecorder) GetIOPSTier(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIOPSTier", reflect.TypeOf((*MockDRaaSService)(nil).GetIOPSTier), arg0)
}

// GetIOPSTiers mocks base method.
func (m *MockDRaaSService) GetIOPSTiers(arg0 connection.APIRequestParameters) ([]draas.IOPSTier, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIOPSTiers", arg0)
	ret0, _ := ret[0].([]draas.IOPSTier)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetIOPSTiers indicates an expected call of GetIOPSTiers.
func (mr *MockDRaaSServiceMockRecorder) GetIOPSTiers(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIOPSTiers", reflect.TypeOf((*MockDRaaSService)(nil).GetIOPSTiers), arg0)
}

// GetSolution mocks base method.
func (m *MockDRaaSService) GetSolution(arg0 string) (draas.Solution, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSolution", arg0)
	ret0, _ := ret[0].(draas.Solution)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSolution indicates an expected call of GetSolution.
func (mr *MockDRaaSServiceMockRecorder) GetSolution(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSolution", reflect.TypeOf((*MockDRaaSService)(nil).GetSolution), arg0)
}

// GetSolutionBackupResources mocks base method.
func (m *MockDRaaSService) GetSolutionBackupResources(arg0 string, arg1 connection.APIRequestParameters) ([]draas.BackupResource, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSolutionBackupResources", arg0, arg1)
	ret0, _ := ret[0].([]draas.BackupResource)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSolutionBackupResources indicates an expected call of GetSolutionBackupResources.
func (mr *MockDRaaSServiceMockRecorder) GetSolutionBackupResources(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSolutionBackupResources", reflect.TypeOf((*MockDRaaSService)(nil).GetSolutionBackupResources), arg0, arg1)
}

// GetSolutionBackupResourcesPaginated mocks base method.
func (m *MockDRaaSService) GetSolutionBackupResourcesPaginated(arg0 string, arg1 connection.APIRequestParameters) (*connection.Paginated[draas.BackupResource], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSolutionBackupResourcesPaginated", arg0, arg1)
	ret0, _ := ret[0].(*connection.Paginated[draas.BackupResource])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSolutionBackupResourcesPaginated indicates an expected call of GetSolutionBackupResourcesPaginated.
func (mr *MockDRaaSServiceMockRecorder) GetSolutionBackupResourcesPaginated(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSolutionBackupResourcesPaginated", reflect.TypeOf((*MockDRaaSService)(nil).GetSolutionBackupResourcesPaginated), arg0, arg1)
}

// GetSolutionBackupService mocks base method.
func (m *MockDRaaSService) GetSolutionBackupService(arg0 string) (draas.BackupService, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSolutionBackupService", arg0)
	ret0, _ := ret[0].(draas.BackupService)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSolutionBackupService indicates an expected call of GetSolutionBackupService.
func (mr *MockDRaaSServiceMockRecorder) GetSolutionBackupService(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSolutionBackupService", reflect.TypeOf((*MockDRaaSService)(nil).GetSolutionBackupService), arg0)
}

// GetSolutionComputeResource mocks base method.
func (m *MockDRaaSService) GetSolutionComputeResource(arg0, arg1 string) (draas.ComputeResource, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSolutionComputeResource", arg0, arg1)
	ret0, _ := ret[0].(draas.ComputeResource)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSolutionComputeResource indicates an expected call of GetSolutionComputeResource.
func (mr *MockDRaaSServiceMockRecorder) GetSolutionComputeResource(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSolutionComputeResource", reflect.TypeOf((*MockDRaaSService)(nil).GetSolutionComputeResource), arg0, arg1)
}

// GetSolutionComputeResources mocks base method.
func (m *MockDRaaSService) GetSolutionComputeResources(arg0 string, arg1 connection.APIRequestParameters) ([]draas.ComputeResource, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSolutionComputeResources", arg0, arg1)
	ret0, _ := ret[0].([]draas.ComputeResource)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSolutionComputeResources indicates an expected call of GetSolutionComputeResources.
func (mr *MockDRaaSServiceMockRecorder) GetSolutionComputeResources(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSolutionComputeResources", reflect.TypeOf((*MockDRaaSService)(nil).GetSolutionComputeResources), arg0, arg1)
}

// GetSolutionComputeResourcesPaginated mocks base method.
func (m *MockDRaaSService) GetSolutionComputeResourcesPaginated(arg0 string, arg1 connection.APIRequestParameters) (*connection.Paginated[draas.ComputeResource], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSolutionComputeResourcesPaginated", arg0, arg1)
	ret0, _ := ret[0].(*connection.Paginated[draas.ComputeResource])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSolutionComputeResourcesPaginated indicates an expected call of GetSolutionComputeResourcesPaginated.
func (mr *MockDRaaSServiceMockRecorder) GetSolutionComputeResourcesPaginated(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSolutionComputeResourcesPaginated", reflect.TypeOf((*MockDRaaSService)(nil).GetSolutionComputeResourcesPaginated), arg0, arg1)
}

// GetSolutionFailoverPlan mocks base method.
func (m *MockDRaaSService) GetSolutionFailoverPlan(arg0, arg1 string) (draas.FailoverPlan, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSolutionFailoverPlan", arg0, arg1)
	ret0, _ := ret[0].(draas.FailoverPlan)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSolutionFailoverPlan indicates an expected call of GetSolutionFailoverPlan.
func (mr *MockDRaaSServiceMockRecorder) GetSolutionFailoverPlan(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSolutionFailoverPlan", reflect.TypeOf((*MockDRaaSService)(nil).GetSolutionFailoverPlan), arg0, arg1)
}

// GetSolutionFailoverPlans mocks base method.
func (m *MockDRaaSService) GetSolutionFailoverPlans(arg0 string, arg1 connection.APIRequestParameters) ([]draas.FailoverPlan, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSolutionFailoverPlans", arg0, arg1)
	ret0, _ := ret[0].([]draas.FailoverPlan)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSolutionFailoverPlans indicates an expected call of GetSolutionFailoverPlans.
func (mr *MockDRaaSServiceMockRecorder) GetSolutionFailoverPlans(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSolutionFailoverPlans", reflect.TypeOf((*MockDRaaSService)(nil).GetSolutionFailoverPlans), arg0, arg1)
}

// GetSolutionFailoverPlansPaginated mocks base method.
func (m *MockDRaaSService) GetSolutionFailoverPlansPaginated(arg0 string, arg1 connection.APIRequestParameters) (*connection.Paginated[draas.FailoverPlan], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSolutionFailoverPlansPaginated", arg0, arg1)
	ret0, _ := ret[0].(*connection.Paginated[draas.FailoverPlan])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSolutionFailoverPlansPaginated indicates an expected call of GetSolutionFailoverPlansPaginated.
func (mr *MockDRaaSServiceMockRecorder) GetSolutionFailoverPlansPaginated(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSolutionFailoverPlansPaginated", reflect.TypeOf((*MockDRaaSService)(nil).GetSolutionFailoverPlansPaginated), arg0, arg1)
}

// GetSolutionHardwarePlan mocks base method.
func (m *MockDRaaSService) GetSolutionHardwarePlan(arg0, arg1 string) (draas.HardwarePlan, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSolutionHardwarePlan", arg0, arg1)
	ret0, _ := ret[0].(draas.HardwarePlan)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSolutionHardwarePlan indicates an expected call of GetSolutionHardwarePlan.
func (mr *MockDRaaSServiceMockRecorder) GetSolutionHardwarePlan(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSolutionHardwarePlan", reflect.TypeOf((*MockDRaaSService)(nil).GetSolutionHardwarePlan), arg0, arg1)
}

// GetSolutionHardwarePlanReplicas mocks base method.
func (m *MockDRaaSService) GetSolutionHardwarePlanReplicas(arg0, arg1 string, arg2 connection.APIRequestParameters) ([]draas.Replica, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSolutionHardwarePlanReplicas", arg0, arg1, arg2)
	ret0, _ := ret[0].([]draas.Replica)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSolutionHardwarePlanReplicas indicates an expected call of GetSolutionHardwarePlanReplicas.
func (mr *MockDRaaSServiceMockRecorder) GetSolutionHardwarePlanReplicas(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSolutionHardwarePlanReplicas", reflect.TypeOf((*MockDRaaSService)(nil).GetSolutionHardwarePlanReplicas), arg0, arg1, arg2)
}

// GetSolutionHardwarePlans mocks base method.
func (m *MockDRaaSService) GetSolutionHardwarePlans(arg0 string, arg1 connection.APIRequestParameters) ([]draas.HardwarePlan, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSolutionHardwarePlans", arg0, arg1)
	ret0, _ := ret[0].([]draas.HardwarePlan)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSolutionHardwarePlans indicates an expected call of GetSolutionHardwarePlans.
func (mr *MockDRaaSServiceMockRecorder) GetSolutionHardwarePlans(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSolutionHardwarePlans", reflect.TypeOf((*MockDRaaSService)(nil).GetSolutionHardwarePlans), arg0, arg1)
}

// GetSolutionHardwarePlansPaginated mocks base method.
func (m *MockDRaaSService) GetSolutionHardwarePlansPaginated(arg0 string, arg1 connection.APIRequestParameters) (*connection.Paginated[draas.HardwarePlan], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSolutionHardwarePlansPaginated", arg0, arg1)
	ret0, _ := ret[0].(*connection.Paginated[draas.HardwarePlan])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSolutionHardwarePlansPaginated indicates an expected call of GetSolutionHardwarePlansPaginated.
func (mr *MockDRaaSServiceMockRecorder) GetSolutionHardwarePlansPaginated(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSolutionHardwarePlansPaginated", reflect.TypeOf((*MockDRaaSService)(nil).GetSolutionHardwarePlansPaginated), arg0, arg1)
}

// GetSolutions mocks base method.
func (m *MockDRaaSService) GetSolutions(arg0 connection.APIRequestParameters) ([]draas.Solution, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSolutions", arg0)
	ret0, _ := ret[0].([]draas.Solution)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSolutions indicates an expected call of GetSolutions.
func (mr *MockDRaaSServiceMockRecorder) GetSolutions(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSolutions", reflect.TypeOf((*MockDRaaSService)(nil).GetSolutions), arg0)
}

// GetSolutionsPaginated mocks base method.
func (m *MockDRaaSService) GetSolutionsPaginated(arg0 connection.APIRequestParameters) (*connection.Paginated[draas.Solution], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSolutionsPaginated", arg0)
	ret0, _ := ret[0].(*connection.Paginated[draas.Solution])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSolutionsPaginated indicates an expected call of GetSolutionsPaginated.
func (mr *MockDRaaSServiceMockRecorder) GetSolutionsPaginated(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSolutionsPaginated", reflect.TypeOf((*MockDRaaSService)(nil).GetSolutionsPaginated), arg0)
}

// PatchSolution mocks base method.
func (m *MockDRaaSService) PatchSolution(arg0 string, arg1 draas.PatchSolutionRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PatchSolution", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// PatchSolution indicates an expected call of PatchSolution.
func (mr *MockDRaaSServiceMockRecorder) PatchSolution(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PatchSolution", reflect.TypeOf((*MockDRaaSService)(nil).PatchSolution), arg0, arg1)
}

// ResetSolutionBackupServiceCredentials mocks base method.
func (m *MockDRaaSService) ResetSolutionBackupServiceCredentials(arg0 string, arg1 draas.ResetBackupServiceCredentialsRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResetSolutionBackupServiceCredentials", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ResetSolutionBackupServiceCredentials indicates an expected call of ResetSolutionBackupServiceCredentials.
func (mr *MockDRaaSServiceMockRecorder) ResetSolutionBackupServiceCredentials(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetSolutionBackupServiceCredentials", reflect.TypeOf((*MockDRaaSService)(nil).ResetSolutionBackupServiceCredentials), arg0, arg1)
}

// StartSolutionFailoverPlan mocks base method.
func (m *MockDRaaSService) StartSolutionFailoverPlan(arg0, arg1 string, arg2 draas.StartFailoverPlanRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartSolutionFailoverPlan", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// StartSolutionFailoverPlan indicates an expected call of StartSolutionFailoverPlan.
func (mr *MockDRaaSServiceMockRecorder) StartSolutionFailoverPlan(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartSolutionFailoverPlan", reflect.TypeOf((*MockDRaaSService)(nil).StartSolutionFailoverPlan), arg0, arg1, arg2)
}

// StopSolutionFailoverPlan mocks base method.
func (m *MockDRaaSService) StopSolutionFailoverPlan(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StopSolutionFailoverPlan", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// StopSolutionFailoverPlan indicates an expected call of StopSolutionFailoverPlan.
func (mr *MockDRaaSServiceMockRecorder) StopSolutionFailoverPlan(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StopSolutionFailoverPlan", reflect.TypeOf((*MockDRaaSService)(nil).StopSolutionFailoverPlan), arg0, arg1)
}

// UpdateSolutionReplicaIOPS mocks base method.
func (m *MockDRaaSService) UpdateSolutionReplicaIOPS(arg0, arg1 string, arg2 draas.UpdateReplicaIOPSRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSolutionReplicaIOPS", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateSolutionReplicaIOPS indicates an expected call of UpdateSolutionReplicaIOPS.
func (mr *MockDRaaSServiceMockRecorder) UpdateSolutionReplicaIOPS(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSolutionReplicaIOPS", reflect.TypeOf((*MockDRaaSService)(nil).UpdateSolutionReplicaIOPS), arg0, arg1, arg2)
}
