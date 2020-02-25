// Code generated by MockGen. DO NOT EDIT.
// Source: manager.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	connection "github.com/stackrox/rox/central/sensor/service/connection"
	pipeline "github.com/stackrox/rox/central/sensor/service/pipeline"
	central "github.com/stackrox/rox/generated/internalapi/central"
	storage "github.com/stackrox/rox/generated/storage"
	concurrency "github.com/stackrox/rox/pkg/concurrency"
	reflect "reflect"
	time "time"
)

// MockClusterManager is a mock of ClusterManager interface
type MockClusterManager struct {
	ctrl     *gomock.Controller
	recorder *MockClusterManagerMockRecorder
}

// MockClusterManagerMockRecorder is the mock recorder for MockClusterManager
type MockClusterManagerMockRecorder struct {
	mock *MockClusterManager
}

// NewMockClusterManager creates a new mock instance
func NewMockClusterManager(ctrl *gomock.Controller) *MockClusterManager {
	mock := &MockClusterManager{ctrl: ctrl}
	mock.recorder = &MockClusterManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockClusterManager) EXPECT() *MockClusterManagerMockRecorder {
	return m.recorder
}

// UpdateClusterContactTimes mocks base method
func (m *MockClusterManager) UpdateClusterContactTimes(ctx context.Context, time time.Time, clusterID ...string) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, time}
	for _, a := range clusterID {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateClusterContactTimes", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateClusterContactTimes indicates an expected call of UpdateClusterContactTimes
func (mr *MockClusterManagerMockRecorder) UpdateClusterContactTimes(ctx, time interface{}, clusterID ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, time}, clusterID...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateClusterContactTimes", reflect.TypeOf((*MockClusterManager)(nil).UpdateClusterContactTimes), varargs...)
}

// UpdateClusterUpgradeStatus mocks base method
func (m *MockClusterManager) UpdateClusterUpgradeStatus(ctx context.Context, clusterID string, status *storage.ClusterUpgradeStatus) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateClusterUpgradeStatus", ctx, clusterID, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateClusterUpgradeStatus indicates an expected call of UpdateClusterUpgradeStatus
func (mr *MockClusterManagerMockRecorder) UpdateClusterUpgradeStatus(ctx, clusterID, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateClusterUpgradeStatus", reflect.TypeOf((*MockClusterManager)(nil).UpdateClusterUpgradeStatus), ctx, clusterID, status)
}

// GetCluster mocks base method
func (m *MockClusterManager) GetCluster(ctx context.Context, id string) (*storage.Cluster, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCluster", ctx, id)
	ret0, _ := ret[0].(*storage.Cluster)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetCluster indicates an expected call of GetCluster
func (mr *MockClusterManagerMockRecorder) GetCluster(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCluster", reflect.TypeOf((*MockClusterManager)(nil).GetCluster), ctx, id)
}

// GetClusters mocks base method
func (m *MockClusterManager) GetClusters(ctx context.Context) ([]*storage.Cluster, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetClusters", ctx)
	ret0, _ := ret[0].([]*storage.Cluster)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetClusters indicates an expected call of GetClusters
func (mr *MockClusterManagerMockRecorder) GetClusters(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClusters", reflect.TypeOf((*MockClusterManager)(nil).GetClusters), ctx)
}

// MockPolicyManager is a mock of PolicyManager interface
type MockPolicyManager struct {
	ctrl     *gomock.Controller
	recorder *MockPolicyManagerMockRecorder
}

// MockPolicyManagerMockRecorder is the mock recorder for MockPolicyManager
type MockPolicyManagerMockRecorder struct {
	mock *MockPolicyManager
}

// NewMockPolicyManager creates a new mock instance
func NewMockPolicyManager(ctrl *gomock.Controller) *MockPolicyManager {
	mock := &MockPolicyManager{ctrl: ctrl}
	mock.recorder = &MockPolicyManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPolicyManager) EXPECT() *MockPolicyManagerMockRecorder {
	return m.recorder
}

// GetPolicies mocks base method
func (m *MockPolicyManager) GetPolicies(ctx context.Context) ([]*storage.Policy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPolicies", ctx)
	ret0, _ := ret[0].([]*storage.Policy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPolicies indicates an expected call of GetPolicies
func (mr *MockPolicyManagerMockRecorder) GetPolicies(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPolicies", reflect.TypeOf((*MockPolicyManager)(nil).GetPolicies), ctx)
}

// MockWhitelistManager is a mock of WhitelistManager interface
type MockWhitelistManager struct {
	ctrl     *gomock.Controller
	recorder *MockWhitelistManagerMockRecorder
}

// MockWhitelistManagerMockRecorder is the mock recorder for MockWhitelistManager
type MockWhitelistManagerMockRecorder struct {
	mock *MockWhitelistManager
}

// NewMockWhitelistManager creates a new mock instance
func NewMockWhitelistManager(ctrl *gomock.Controller) *MockWhitelistManager {
	mock := &MockWhitelistManager{ctrl: ctrl}
	mock.recorder = &MockWhitelistManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockWhitelistManager) EXPECT() *MockWhitelistManagerMockRecorder {
	return m.recorder
}

// WalkAll mocks base method
func (m *MockWhitelistManager) WalkAll(ctx context.Context, fn func(*storage.ProcessWhitelist) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WalkAll", ctx, fn)
	ret0, _ := ret[0].(error)
	return ret0
}

// WalkAll indicates an expected call of WalkAll
func (mr *MockWhitelistManagerMockRecorder) WalkAll(ctx, fn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WalkAll", reflect.TypeOf((*MockWhitelistManager)(nil).WalkAll), ctx, fn)
}

// MockManager is a mock of Manager interface
type MockManager struct {
	ctrl     *gomock.Controller
	recorder *MockManagerMockRecorder
}

// MockManagerMockRecorder is the mock recorder for MockManager
type MockManagerMockRecorder struct {
	mock *MockManager
}

// NewMockManager creates a new mock instance
func NewMockManager(ctrl *gomock.Controller) *MockManager {
	mock := &MockManager{ctrl: ctrl}
	mock.recorder = &MockManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockManager) EXPECT() *MockManagerMockRecorder {
	return m.recorder
}

// Start mocks base method
func (m *MockManager) Start(mgr connection.ClusterManager, policyMgr connection.PolicyManager, whitelistMgr connection.WhitelistManager, autoTriggerUpgrades *concurrency.Flag) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start", mgr, policyMgr, whitelistMgr, autoTriggerUpgrades)
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start
func (mr *MockManagerMockRecorder) Start(mgr, policyMgr, whitelistMgr, autoTriggerUpgrades interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockManager)(nil).Start), mgr, policyMgr, whitelistMgr, autoTriggerUpgrades)
}

// HandleConnection mocks base method
func (m *MockManager) HandleConnection(ctx context.Context, clusterID string, eventPipeline pipeline.ClusterPipeline, server central.SensorService_CommunicateServer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleConnection", ctx, clusterID, eventPipeline, server)
	ret0, _ := ret[0].(error)
	return ret0
}

// HandleConnection indicates an expected call of HandleConnection
func (mr *MockManagerMockRecorder) HandleConnection(ctx, clusterID, eventPipeline, server interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleConnection", reflect.TypeOf((*MockManager)(nil).HandleConnection), ctx, clusterID, eventPipeline, server)
}

// GetConnection mocks base method
func (m *MockManager) GetConnection(clusterID string) connection.SensorConnection {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConnection", clusterID)
	ret0, _ := ret[0].(connection.SensorConnection)
	return ret0
}

// GetConnection indicates an expected call of GetConnection
func (mr *MockManagerMockRecorder) GetConnection(clusterID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConnection", reflect.TypeOf((*MockManager)(nil).GetConnection), clusterID)
}

// GetActiveConnections mocks base method
func (m *MockManager) GetActiveConnections() []connection.SensorConnection {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetActiveConnections")
	ret0, _ := ret[0].([]connection.SensorConnection)
	return ret0
}

// GetActiveConnections indicates an expected call of GetActiveConnections
func (mr *MockManagerMockRecorder) GetActiveConnections() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetActiveConnections", reflect.TypeOf((*MockManager)(nil).GetActiveConnections))
}

// BroadcastMessage mocks base method
func (m *MockManager) BroadcastMessage(msg *central.MsgToSensor) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "BroadcastMessage", msg)
}

// BroadcastMessage indicates an expected call of BroadcastMessage
func (mr *MockManagerMockRecorder) BroadcastMessage(msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BroadcastMessage", reflect.TypeOf((*MockManager)(nil).BroadcastMessage), msg)
}

// SendMessage mocks base method
func (m *MockManager) SendMessage(clusterID string, msg *central.MsgToSensor) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMessage", clusterID, msg)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMessage indicates an expected call of SendMessage
func (mr *MockManagerMockRecorder) SendMessage(clusterID, msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessage", reflect.TypeOf((*MockManager)(nil).SendMessage), clusterID, msg)
}

// TriggerUpgrade mocks base method
func (m *MockManager) TriggerUpgrade(ctx context.Context, clusterID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TriggerUpgrade", ctx, clusterID)
	ret0, _ := ret[0].(error)
	return ret0
}

// TriggerUpgrade indicates an expected call of TriggerUpgrade
func (mr *MockManagerMockRecorder) TriggerUpgrade(ctx, clusterID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TriggerUpgrade", reflect.TypeOf((*MockManager)(nil).TriggerUpgrade), ctx, clusterID)
}

// ProcessCheckInFromUpgrader mocks base method
func (m *MockManager) ProcessCheckInFromUpgrader(ctx context.Context, clusterID string, req *central.UpgradeCheckInFromUpgraderRequest) (*central.UpgradeCheckInFromUpgraderResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProcessCheckInFromUpgrader", ctx, clusterID, req)
	ret0, _ := ret[0].(*central.UpgradeCheckInFromUpgraderResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ProcessCheckInFromUpgrader indicates an expected call of ProcessCheckInFromUpgrader
func (mr *MockManagerMockRecorder) ProcessCheckInFromUpgrader(ctx, clusterID, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessCheckInFromUpgrader", reflect.TypeOf((*MockManager)(nil).ProcessCheckInFromUpgrader), ctx, clusterID, req)
}

// ProcessUpgradeCheckInFromSensor mocks base method
func (m *MockManager) ProcessUpgradeCheckInFromSensor(ctx context.Context, clusterID string, req *central.UpgradeCheckInFromSensorRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProcessUpgradeCheckInFromSensor", ctx, clusterID, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// ProcessUpgradeCheckInFromSensor indicates an expected call of ProcessUpgradeCheckInFromSensor
func (mr *MockManagerMockRecorder) ProcessUpgradeCheckInFromSensor(ctx, clusterID, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessUpgradeCheckInFromSensor", reflect.TypeOf((*MockManager)(nil).ProcessUpgradeCheckInFromSensor), ctx, clusterID, req)
}
