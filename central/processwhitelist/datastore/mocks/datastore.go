// Code generated by MockGen. DO NOT EDIT.
// Source: datastore.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	v1 "github.com/stackrox/rox/generated/api/v1"
	storage "github.com/stackrox/rox/generated/storage"
	search "github.com/stackrox/rox/pkg/search"
	reflect "reflect"
)

// MockDataStore is a mock of DataStore interface
type MockDataStore struct {
	ctrl     *gomock.Controller
	recorder *MockDataStoreMockRecorder
}

// MockDataStoreMockRecorder is the mock recorder for MockDataStore
type MockDataStoreMockRecorder struct {
	mock *MockDataStore
}

// NewMockDataStore creates a new mock instance
func NewMockDataStore(ctrl *gomock.Controller) *MockDataStore {
	mock := &MockDataStore{ctrl: ctrl}
	mock.recorder = &MockDataStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDataStore) EXPECT() *MockDataStoreMockRecorder {
	return m.recorder
}

// SearchRawProcessWhitelists mocks base method
func (m *MockDataStore) SearchRawProcessWhitelists(ctx context.Context, q *v1.Query) ([]*storage.ProcessWhitelist, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchRawProcessWhitelists", ctx, q)
	ret0, _ := ret[0].([]*storage.ProcessWhitelist)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchRawProcessWhitelists indicates an expected call of SearchRawProcessWhitelists
func (mr *MockDataStoreMockRecorder) SearchRawProcessWhitelists(ctx, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchRawProcessWhitelists", reflect.TypeOf((*MockDataStore)(nil).SearchRawProcessWhitelists), ctx, q)
}

// Search mocks base method
func (m *MockDataStore) Search(ctx context.Context, q *v1.Query) ([]search.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Search", ctx, q)
	ret0, _ := ret[0].([]search.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Search indicates an expected call of Search
func (mr *MockDataStoreMockRecorder) Search(ctx, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockDataStore)(nil).Search), ctx, q)
}

// GetProcessWhitelist mocks base method
func (m *MockDataStore) GetProcessWhitelist(ctx context.Context, key *storage.ProcessWhitelistKey) (*storage.ProcessWhitelist, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProcessWhitelist", ctx, key)
	ret0, _ := ret[0].(*storage.ProcessWhitelist)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetProcessWhitelist indicates an expected call of GetProcessWhitelist
func (mr *MockDataStoreMockRecorder) GetProcessWhitelist(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProcessWhitelist", reflect.TypeOf((*MockDataStore)(nil).GetProcessWhitelist), ctx, key)
}

// AddProcessWhitelist mocks base method
func (m *MockDataStore) AddProcessWhitelist(ctx context.Context, whitelist *storage.ProcessWhitelist) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddProcessWhitelist", ctx, whitelist)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddProcessWhitelist indicates an expected call of AddProcessWhitelist
func (mr *MockDataStoreMockRecorder) AddProcessWhitelist(ctx, whitelist interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddProcessWhitelist", reflect.TypeOf((*MockDataStore)(nil).AddProcessWhitelist), ctx, whitelist)
}

// RemoveProcessWhitelist mocks base method
func (m *MockDataStore) RemoveProcessWhitelist(ctx context.Context, key *storage.ProcessWhitelistKey) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveProcessWhitelist", ctx, key)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveProcessWhitelist indicates an expected call of RemoveProcessWhitelist
func (mr *MockDataStoreMockRecorder) RemoveProcessWhitelist(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveProcessWhitelist", reflect.TypeOf((*MockDataStore)(nil).RemoveProcessWhitelist), ctx, key)
}

// RemoveProcessWhitelistsByDeployment mocks base method
func (m *MockDataStore) RemoveProcessWhitelistsByDeployment(ctx context.Context, deploymentID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveProcessWhitelistsByDeployment", ctx, deploymentID)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveProcessWhitelistsByDeployment indicates an expected call of RemoveProcessWhitelistsByDeployment
func (mr *MockDataStoreMockRecorder) RemoveProcessWhitelistsByDeployment(ctx, deploymentID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveProcessWhitelistsByDeployment", reflect.TypeOf((*MockDataStore)(nil).RemoveProcessWhitelistsByDeployment), ctx, deploymentID)
}

// UpdateProcessWhitelistElements mocks base method
func (m *MockDataStore) UpdateProcessWhitelistElements(ctx context.Context, key *storage.ProcessWhitelistKey, addElements, removeElements []*storage.WhitelistItem, auto bool) (*storage.ProcessWhitelist, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProcessWhitelistElements", ctx, key, addElements, removeElements, auto)
	ret0, _ := ret[0].(*storage.ProcessWhitelist)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateProcessWhitelistElements indicates an expected call of UpdateProcessWhitelistElements
func (mr *MockDataStoreMockRecorder) UpdateProcessWhitelistElements(ctx, key, addElements, removeElements, auto interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProcessWhitelistElements", reflect.TypeOf((*MockDataStore)(nil).UpdateProcessWhitelistElements), ctx, key, addElements, removeElements, auto)
}

// UpsertProcessWhitelist mocks base method
func (m *MockDataStore) UpsertProcessWhitelist(ctx context.Context, key *storage.ProcessWhitelistKey, addElements []*storage.WhitelistItem, auto bool) (*storage.ProcessWhitelist, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertProcessWhitelist", ctx, key, addElements, auto)
	ret0, _ := ret[0].(*storage.ProcessWhitelist)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpsertProcessWhitelist indicates an expected call of UpsertProcessWhitelist
func (mr *MockDataStoreMockRecorder) UpsertProcessWhitelist(ctx, key, addElements, auto interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertProcessWhitelist", reflect.TypeOf((*MockDataStore)(nil).UpsertProcessWhitelist), ctx, key, addElements, auto)
}

// UserLockProcessWhitelist mocks base method
func (m *MockDataStore) UserLockProcessWhitelist(ctx context.Context, key *storage.ProcessWhitelistKey, locked bool) (*storage.ProcessWhitelist, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserLockProcessWhitelist", ctx, key, locked)
	ret0, _ := ret[0].(*storage.ProcessWhitelist)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserLockProcessWhitelist indicates an expected call of UserLockProcessWhitelist
func (mr *MockDataStoreMockRecorder) UserLockProcessWhitelist(ctx, key, locked interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserLockProcessWhitelist", reflect.TypeOf((*MockDataStore)(nil).UserLockProcessWhitelist), ctx, key, locked)
}

// WalkAll mocks base method
func (m *MockDataStore) WalkAll(ctx context.Context, fn func(*storage.ProcessWhitelist) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WalkAll", ctx, fn)
	ret0, _ := ret[0].(error)
	return ret0
}

// WalkAll indicates an expected call of WalkAll
func (mr *MockDataStoreMockRecorder) WalkAll(ctx, fn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WalkAll", reflect.TypeOf((*MockDataStore)(nil).WalkAll), ctx, fn)
}
