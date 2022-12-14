// Code generated by MockGen. DO NOT EDIT.
// Source: gitlab.com/akita/mem/device (interfaces: PageTable)

package mmu

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	device "gitlab.com/akita/mem/device"
	ca "gitlab.com/akita/util/ca"
)

// MockPageTable is a mock of PageTable interface.
type MockPageTable struct {
	ctrl     *gomock.Controller
	recorder *MockPageTableMockRecorder
}

// MockPageTableMockRecorder is the mock recorder for MockPageTable.
type MockPageTableMockRecorder struct {
	mock *MockPageTable
}

// NewMockPageTable creates a new mock instance.
func NewMockPageTable(ctrl *gomock.Controller) *MockPageTable {
	mock := &MockPageTable{ctrl: ctrl}
	mock.recorder = &MockPageTableMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPageTable) EXPECT() *MockPageTableMockRecorder {
	return m.recorder
}

// Find mocks base method.
func (m *MockPageTable) Find(arg0 ca.PID, arg1 uint64) (device.Page, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", arg0, arg1)
	ret0, _ := ret[0].(device.Page)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// Find indicates an expected call of Find.
func (mr *MockPageTableMockRecorder) Find(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockPageTable)(nil).Find), arg0, arg1)
}

// FindAddr mocks base method.
func (m *MockPageTable) FindAddr(arg0 ca.PID, arg1, arg2 uint64) uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAddr", arg0, arg1, arg2)
	ret0, _ := ret[0].(uint64)
	return ret0
}

// FindAddr indicates an expected call of FindAddr.
func (mr *MockPageTableMockRecorder) FindAddr(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAddr", reflect.TypeOf((*MockPageTable)(nil).FindAddr), arg0, arg1, arg2)
}

// Insert mocks base method.
func (m *MockPageTable) Insert(arg0 device.Page) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Insert", arg0)
}

// Insert indicates an expected call of Insert.
func (mr *MockPageTableMockRecorder) Insert(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockPageTable)(nil).Insert), arg0)
}

// PageTablePagesAsBytes mocks base method.
func (m *MockPageTable) PageTablePagesAsBytes(arg0 ca.PID) []uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PageTablePagesAsBytes", arg0)
	ret0, _ := ret[0].([]uint64)
	return ret0
}

// PageTablePagesAsBytes indicates an expected call of PageTablePagesAsBytes.
func (mr *MockPageTableMockRecorder) PageTablePagesAsBytes(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PageTablePagesAsBytes", reflect.TypeOf((*MockPageTable)(nil).PageTablePagesAsBytes), arg0)
}

// Remove mocks base method.
func (m *MockPageTable) Remove(arg0 ca.PID, arg1 uint64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Remove", arg0, arg1)
}

// Remove indicates an expected call of Remove.
func (mr *MockPageTableMockRecorder) Remove(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockPageTable)(nil).Remove), arg0, arg1)
}

// Update mocks base method.
func (m *MockPageTable) Update(arg0 device.Page) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Update", arg0)
}

// Update indicates an expected call of Update.
func (mr *MockPageTableMockRecorder) Update(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockPageTable)(nil).Update), arg0)
}
