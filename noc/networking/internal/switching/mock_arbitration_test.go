// Code generated by MockGen. DO NOT EDIT.
// Source: gitlab.com/akita/noc/networking/internal/arbitration (interfaces: Arbiter)

package switching

import (
	gomock "github.com/golang/mock/gomock"
	akita "gitlab.com/akita/akita"
	util "gitlab.com/akita/util"
	reflect "reflect"
)

// MockArbiter is a mock of Arbiter interface
type MockArbiter struct {
	ctrl     *gomock.Controller
	recorder *MockArbiterMockRecorder
}

// MockArbiterMockRecorder is the mock recorder for MockArbiter
type MockArbiterMockRecorder struct {
	mock *MockArbiter
}

// NewMockArbiter creates a new mock instance
func NewMockArbiter(ctrl *gomock.Controller) *MockArbiter {
	mock := &MockArbiter{ctrl: ctrl}
	mock.recorder = &MockArbiterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockArbiter) EXPECT() *MockArbiterMockRecorder {
	return m.recorder
}

// AddBuffer mocks base method
func (m *MockArbiter) AddBuffer(arg0 util.Buffer) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddBuffer", arg0)
}

// AddBuffer indicates an expected call of AddBuffer
func (mr *MockArbiterMockRecorder) AddBuffer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddBuffer", reflect.TypeOf((*MockArbiter)(nil).AddBuffer), arg0)
}

// Arbitrate mocks base method
func (m *MockArbiter) Arbitrate(arg0 akita.VTimeInSec) []util.Buffer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Arbitrate", arg0)
	ret0, _ := ret[0].([]util.Buffer)
	return ret0
}

// Arbitrate indicates an expected call of Arbitrate
func (mr *MockArbiterMockRecorder) Arbitrate(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Arbitrate", reflect.TypeOf((*MockArbiter)(nil).Arbitrate), arg0)
}
