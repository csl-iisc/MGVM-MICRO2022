// Code generated by MockGen. DO NOT EDIT.
// Source: coalescer.go

package cu

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	wavefront "gitlab.com/akita/mgpusim/timing/wavefront"
)

// Mockcoalescer is a mock of coalescer interface
type Mockcoalescer struct {
	ctrl     *gomock.Controller
	recorder *MockcoalescerMockRecorder
}

// MockcoalescerMockRecorder is the mock recorder for Mockcoalescer
type MockcoalescerMockRecorder struct {
	mock *Mockcoalescer
}

// NewMockcoalescer creates a new mock instance
func NewMockcoalescer(ctrl *gomock.Controller) *Mockcoalescer {
	mock := &Mockcoalescer{ctrl: ctrl}
	mock.recorder = &MockcoalescerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *Mockcoalescer) EXPECT() *MockcoalescerMockRecorder {
	return m.recorder
}

// generateMemTransactions mocks base method
func (m *Mockcoalescer) generateMemTransactions(wf *wavefront.Wavefront) []VectorMemAccessInfo {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "generateMemTransactions", wf)
	ret0, _ := ret[0].([]VectorMemAccessInfo)
	return ret0
}

// generateMemTransactions indicates an expected call of generateMemTransactions
func (mr *MockcoalescerMockRecorder) generateMemTransactions(wf interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "generateMemTransactions", reflect.TypeOf((*Mockcoalescer)(nil).generateMemTransactions), wf)
}
