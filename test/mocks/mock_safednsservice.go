// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ukfast/sdk-go/pkg/service/safedns (interfaces: SafeDNSService)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	connection "github.com/ukfast/sdk-go/pkg/connection"
	safedns "github.com/ukfast/sdk-go/pkg/service/safedns"
)

// MockSafeDNSService is a mock of SafeDNSService interface.
type MockSafeDNSService struct {
	ctrl     *gomock.Controller
	recorder *MockSafeDNSServiceMockRecorder
}

// MockSafeDNSServiceMockRecorder is the mock recorder for MockSafeDNSService.
type MockSafeDNSServiceMockRecorder struct {
	mock *MockSafeDNSService
}

// NewMockSafeDNSService creates a new mock instance.
func NewMockSafeDNSService(ctrl *gomock.Controller) *MockSafeDNSService {
	mock := &MockSafeDNSService{ctrl: ctrl}
	mock.recorder = &MockSafeDNSServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSafeDNSService) EXPECT() *MockSafeDNSServiceMockRecorder {
	return m.recorder
}

// CreateTemplate mocks base method.
func (m *MockSafeDNSService) CreateTemplate(arg0 safedns.CreateTemplateRequest) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTemplate", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTemplate indicates an expected call of CreateTemplate.
func (mr *MockSafeDNSServiceMockRecorder) CreateTemplate(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTemplate", reflect.TypeOf((*MockSafeDNSService)(nil).CreateTemplate), arg0)
}

// CreateTemplateRecord mocks base method.
func (m *MockSafeDNSService) CreateTemplateRecord(arg0 int, arg1 safedns.CreateRecordRequest) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTemplateRecord", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTemplateRecord indicates an expected call of CreateTemplateRecord.
func (mr *MockSafeDNSServiceMockRecorder) CreateTemplateRecord(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTemplateRecord", reflect.TypeOf((*MockSafeDNSService)(nil).CreateTemplateRecord), arg0, arg1)
}

// CreateZone mocks base method.
func (m *MockSafeDNSService) CreateZone(arg0 safedns.CreateZoneRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateZone", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateZone indicates an expected call of CreateZone.
func (mr *MockSafeDNSServiceMockRecorder) CreateZone(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateZone", reflect.TypeOf((*MockSafeDNSService)(nil).CreateZone), arg0)
}

// CreateZoneNote mocks base method.
func (m *MockSafeDNSService) CreateZoneNote(arg0 string, arg1 safedns.CreateNoteRequest) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateZoneNote", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateZoneNote indicates an expected call of CreateZoneNote.
func (mr *MockSafeDNSServiceMockRecorder) CreateZoneNote(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateZoneNote", reflect.TypeOf((*MockSafeDNSService)(nil).CreateZoneNote), arg0, arg1)
}

// CreateZoneRecord mocks base method.
func (m *MockSafeDNSService) CreateZoneRecord(arg0 string, arg1 safedns.CreateRecordRequest) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateZoneRecord", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateZoneRecord indicates an expected call of CreateZoneRecord.
func (mr *MockSafeDNSServiceMockRecorder) CreateZoneRecord(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateZoneRecord", reflect.TypeOf((*MockSafeDNSService)(nil).CreateZoneRecord), arg0, arg1)
}

// DeleteTemplate mocks base method.
func (m *MockSafeDNSService) DeleteTemplate(arg0 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTemplate", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTemplate indicates an expected call of DeleteTemplate.
func (mr *MockSafeDNSServiceMockRecorder) DeleteTemplate(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTemplate", reflect.TypeOf((*MockSafeDNSService)(nil).DeleteTemplate), arg0)
}

// DeleteTemplateRecord mocks base method.
func (m *MockSafeDNSService) DeleteTemplateRecord(arg0, arg1 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTemplateRecord", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTemplateRecord indicates an expected call of DeleteTemplateRecord.
func (mr *MockSafeDNSServiceMockRecorder) DeleteTemplateRecord(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTemplateRecord", reflect.TypeOf((*MockSafeDNSService)(nil).DeleteTemplateRecord), arg0, arg1)
}

// DeleteZone mocks base method.
func (m *MockSafeDNSService) DeleteZone(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteZone", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteZone indicates an expected call of DeleteZone.
func (mr *MockSafeDNSServiceMockRecorder) DeleteZone(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteZone", reflect.TypeOf((*MockSafeDNSService)(nil).DeleteZone), arg0)
}

// DeleteZoneRecord mocks base method.
func (m *MockSafeDNSService) DeleteZoneRecord(arg0 string, arg1 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteZoneRecord", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteZoneRecord indicates an expected call of DeleteZoneRecord.
func (mr *MockSafeDNSServiceMockRecorder) DeleteZoneRecord(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteZoneRecord", reflect.TypeOf((*MockSafeDNSService)(nil).DeleteZoneRecord), arg0, arg1)
}

// GetSettings mocks base method.
func (m *MockSafeDNSService) GetSettings() (safedns.Settings, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSettings")
	ret0, _ := ret[0].(safedns.Settings)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSettings indicates an expected call of GetSettings.
func (mr *MockSafeDNSServiceMockRecorder) GetSettings() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSettings", reflect.TypeOf((*MockSafeDNSService)(nil).GetSettings))
}

// GetTemplate mocks base method.
func (m *MockSafeDNSService) GetTemplate(arg0 int) (safedns.Template, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTemplate", arg0)
	ret0, _ := ret[0].(safedns.Template)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTemplate indicates an expected call of GetTemplate.
func (mr *MockSafeDNSServiceMockRecorder) GetTemplate(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTemplate", reflect.TypeOf((*MockSafeDNSService)(nil).GetTemplate), arg0)
}

// GetTemplateRecord mocks base method.
func (m *MockSafeDNSService) GetTemplateRecord(arg0, arg1 int) (safedns.Record, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTemplateRecord", arg0, arg1)
	ret0, _ := ret[0].(safedns.Record)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTemplateRecord indicates an expected call of GetTemplateRecord.
func (mr *MockSafeDNSServiceMockRecorder) GetTemplateRecord(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTemplateRecord", reflect.TypeOf((*MockSafeDNSService)(nil).GetTemplateRecord), arg0, arg1)
}

// GetTemplateRecords mocks base method.
func (m *MockSafeDNSService) GetTemplateRecords(arg0 int, arg1 connection.APIRequestParameters) ([]safedns.Record, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTemplateRecords", arg0, arg1)
	ret0, _ := ret[0].([]safedns.Record)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTemplateRecords indicates an expected call of GetTemplateRecords.
func (mr *MockSafeDNSServiceMockRecorder) GetTemplateRecords(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTemplateRecords", reflect.TypeOf((*MockSafeDNSService)(nil).GetTemplateRecords), arg0, arg1)
}

// GetTemplateRecordsPaginated mocks base method.
func (m *MockSafeDNSService) GetTemplateRecordsPaginated(arg0 int, arg1 connection.APIRequestParameters) (*connection.Paginated[safedns.Record], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTemplateRecordsPaginated", arg0, arg1)
	ret0, _ := ret[0].(*connection.Paginated[safedns.Record])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTemplateRecordsPaginated indicates an expected call of GetTemplateRecordsPaginated.
func (mr *MockSafeDNSServiceMockRecorder) GetTemplateRecordsPaginated(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTemplateRecordsPaginated", reflect.TypeOf((*MockSafeDNSService)(nil).GetTemplateRecordsPaginated), arg0, arg1)
}

// GetTemplates mocks base method.
func (m *MockSafeDNSService) GetTemplates(arg0 connection.APIRequestParameters) ([]safedns.Template, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTemplates", arg0)
	ret0, _ := ret[0].([]safedns.Template)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTemplates indicates an expected call of GetTemplates.
func (mr *MockSafeDNSServiceMockRecorder) GetTemplates(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTemplates", reflect.TypeOf((*MockSafeDNSService)(nil).GetTemplates), arg0)
}

// GetTemplatesPaginated mocks base method.
func (m *MockSafeDNSService) GetTemplatesPaginated(arg0 connection.APIRequestParameters) (*connection.Paginated[safedns.Template], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTemplatesPaginated", arg0)
	ret0, _ := ret[0].(*connection.Paginated[safedns.Template])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTemplatesPaginated indicates an expected call of GetTemplatesPaginated.
func (mr *MockSafeDNSServiceMockRecorder) GetTemplatesPaginated(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTemplatesPaginated", reflect.TypeOf((*MockSafeDNSService)(nil).GetTemplatesPaginated), arg0)
}

// GetZone mocks base method.
func (m *MockSafeDNSService) GetZone(arg0 string) (safedns.Zone, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetZone", arg0)
	ret0, _ := ret[0].(safedns.Zone)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetZone indicates an expected call of GetZone.
func (mr *MockSafeDNSServiceMockRecorder) GetZone(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetZone", reflect.TypeOf((*MockSafeDNSService)(nil).GetZone), arg0)
}

// GetZoneNote mocks base method.
func (m *MockSafeDNSService) GetZoneNote(arg0 string, arg1 int) (safedns.Note, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetZoneNote", arg0, arg1)
	ret0, _ := ret[0].(safedns.Note)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetZoneNote indicates an expected call of GetZoneNote.
func (mr *MockSafeDNSServiceMockRecorder) GetZoneNote(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetZoneNote", reflect.TypeOf((*MockSafeDNSService)(nil).GetZoneNote), arg0, arg1)
}

// GetZoneNotes mocks base method.
func (m *MockSafeDNSService) GetZoneNotes(arg0 string, arg1 connection.APIRequestParameters) ([]safedns.Note, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetZoneNotes", arg0, arg1)
	ret0, _ := ret[0].([]safedns.Note)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetZoneNotes indicates an expected call of GetZoneNotes.
func (mr *MockSafeDNSServiceMockRecorder) GetZoneNotes(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetZoneNotes", reflect.TypeOf((*MockSafeDNSService)(nil).GetZoneNotes), arg0, arg1)
}

// GetZoneNotesPaginated mocks base method.
func (m *MockSafeDNSService) GetZoneNotesPaginated(arg0 string, arg1 connection.APIRequestParameters) (*connection.Paginated[safedns.Note], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetZoneNotesPaginated", arg0, arg1)
	ret0, _ := ret[0].(*connection.Paginated[safedns.Note])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetZoneNotesPaginated indicates an expected call of GetZoneNotesPaginated.
func (mr *MockSafeDNSServiceMockRecorder) GetZoneNotesPaginated(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetZoneNotesPaginated", reflect.TypeOf((*MockSafeDNSService)(nil).GetZoneNotesPaginated), arg0, arg1)
}

// GetZoneRecord mocks base method.
func (m *MockSafeDNSService) GetZoneRecord(arg0 string, arg1 int) (safedns.Record, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetZoneRecord", arg0, arg1)
	ret0, _ := ret[0].(safedns.Record)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetZoneRecord indicates an expected call of GetZoneRecord.
func (mr *MockSafeDNSServiceMockRecorder) GetZoneRecord(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetZoneRecord", reflect.TypeOf((*MockSafeDNSService)(nil).GetZoneRecord), arg0, arg1)
}

// GetZoneRecords mocks base method.
func (m *MockSafeDNSService) GetZoneRecords(arg0 string, arg1 connection.APIRequestParameters) ([]safedns.Record, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetZoneRecords", arg0, arg1)
	ret0, _ := ret[0].([]safedns.Record)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetZoneRecords indicates an expected call of GetZoneRecords.
func (mr *MockSafeDNSServiceMockRecorder) GetZoneRecords(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetZoneRecords", reflect.TypeOf((*MockSafeDNSService)(nil).GetZoneRecords), arg0, arg1)
}

// GetZoneRecordsPaginated mocks base method.
func (m *MockSafeDNSService) GetZoneRecordsPaginated(arg0 string, arg1 connection.APIRequestParameters) (*connection.Paginated[safedns.Record], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetZoneRecordsPaginated", arg0, arg1)
	ret0, _ := ret[0].(*connection.Paginated[safedns.Record])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetZoneRecordsPaginated indicates an expected call of GetZoneRecordsPaginated.
func (mr *MockSafeDNSServiceMockRecorder) GetZoneRecordsPaginated(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetZoneRecordsPaginated", reflect.TypeOf((*MockSafeDNSService)(nil).GetZoneRecordsPaginated), arg0, arg1)
}

// GetZones mocks base method.
func (m *MockSafeDNSService) GetZones(arg0 connection.APIRequestParameters) ([]safedns.Zone, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetZones", arg0)
	ret0, _ := ret[0].([]safedns.Zone)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetZones indicates an expected call of GetZones.
func (mr *MockSafeDNSServiceMockRecorder) GetZones(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetZones", reflect.TypeOf((*MockSafeDNSService)(nil).GetZones), arg0)
}

// GetZonesPaginated mocks base method.
func (m *MockSafeDNSService) GetZonesPaginated(arg0 connection.APIRequestParameters) (*connection.Paginated[safedns.Zone], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetZonesPaginated", arg0)
	ret0, _ := ret[0].(*connection.Paginated[safedns.Zone])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetZonesPaginated indicates an expected call of GetZonesPaginated.
func (mr *MockSafeDNSServiceMockRecorder) GetZonesPaginated(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetZonesPaginated", reflect.TypeOf((*MockSafeDNSService)(nil).GetZonesPaginated), arg0)
}

// PatchTemplate mocks base method.
func (m *MockSafeDNSService) PatchTemplate(arg0 int, arg1 safedns.PatchTemplateRequest) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PatchTemplate", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PatchTemplate indicates an expected call of PatchTemplate.
func (mr *MockSafeDNSServiceMockRecorder) PatchTemplate(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PatchTemplate", reflect.TypeOf((*MockSafeDNSService)(nil).PatchTemplate), arg0, arg1)
}

// PatchTemplateRecord mocks base method.
func (m *MockSafeDNSService) PatchTemplateRecord(arg0, arg1 int, arg2 safedns.PatchRecordRequest) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PatchTemplateRecord", arg0, arg1, arg2)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PatchTemplateRecord indicates an expected call of PatchTemplateRecord.
func (mr *MockSafeDNSServiceMockRecorder) PatchTemplateRecord(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PatchTemplateRecord", reflect.TypeOf((*MockSafeDNSService)(nil).PatchTemplateRecord), arg0, arg1, arg2)
}

// PatchZone mocks base method.
func (m *MockSafeDNSService) PatchZone(arg0 string, arg1 safedns.PatchZoneRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PatchZone", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// PatchZone indicates an expected call of PatchZone.
func (mr *MockSafeDNSServiceMockRecorder) PatchZone(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PatchZone", reflect.TypeOf((*MockSafeDNSService)(nil).PatchZone), arg0, arg1)
}

// PatchZoneRecord mocks base method.
func (m *MockSafeDNSService) PatchZoneRecord(arg0 string, arg1 int, arg2 safedns.PatchRecordRequest) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PatchZoneRecord", arg0, arg1, arg2)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PatchZoneRecord indicates an expected call of PatchZoneRecord.
func (mr *MockSafeDNSServiceMockRecorder) PatchZoneRecord(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PatchZoneRecord", reflect.TypeOf((*MockSafeDNSService)(nil).PatchZoneRecord), arg0, arg1, arg2)
}

// UpdateZoneRecord mocks base method.
func (m *MockSafeDNSService) UpdateZoneRecord(arg0 string, arg1 safedns.Record) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateZoneRecord", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateZoneRecord indicates an expected call of UpdateZoneRecord.
func (mr *MockSafeDNSServiceMockRecorder) UpdateZoneRecord(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateZoneRecord", reflect.TypeOf((*MockSafeDNSService)(nil).UpdateZoneRecord), arg0, arg1)
}
