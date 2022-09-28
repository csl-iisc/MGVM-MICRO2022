// Code generated by MockGen. DO NOT EDIT.
// Source: gitlab.com/akita/mem/dram/internal/cmdq (interfaces: CommandQueue)

package trans

import (
	gomock "github.com/golang/mock/gomock"
	akita "gitlab.com/akita/akita"
	signal "gitlab.com/akita/mem/dram/internal/signal"
	reflect "reflect"
)

// MockCommandQueue is a mock of CommandQueue interface
type MockCommandQueue struct {
	ctrl     *gomock.Controller
	recorder *MockCommandQueueMockRecorder
}

// MockCommandQueueMockRecorder is the mock recorder for MockCommandQueue
type MockCommandQueueMockRecorder struct {
	mock *MockCommandQueue
}

// NewMockCommandQueue creates a new mock instance
func NewMockCommandQueue(ctrl *gomock.Controller) *MockCommandQueue {
	mock := &MockCommandQueue{ctrl: ctrl}
	mock.recorder = &MockCommandQueueMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCommandQueue) EXPECT() *MockCommandQueueMockRecorder {
	return m.recorder
}

// Accept mocks base method
func (m *MockCommandQueue) Accept(arg0 *signal.Command) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Accept", arg0)
}

// Accept indicates an expected call of Accept
func (mr *MockCommandQueueMockRecorder) Accept(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Accept", reflect.TypeOf((*MockCommandQueue)(nil).Accept), arg0)
}

// CanAccept mocks base method
func (m *MockCommandQueue) CanAccept(arg0 *signal.Command) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CanAccept", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// CanAccept indicates an expected call of CanAccept
func (mr *MockCommandQueueMockRecorder) CanAccept(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CanAccept", reflect.TypeOf((*MockCommandQueue)(nil).CanAccept), arg0)
}

// GetCommandToIssue mocks base method
func (m *MockCommandQueue) GetCommandToIssue(arg0 akita.VTimeInSec) *signal.Command {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCommandToIssue", arg0)
	ret0, _ := ret[0].(*signal.Command)
	return ret0
}

// GetCommandToIssue indicates an expected call of GetCommandToIssue
func (mr *MockCommandQueueMockRecorder) GetCommandToIssue(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCommandToIssue", reflect.TypeOf((*MockCommandQueue)(nil).GetCommandToIssue), arg0)
}
