// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/pkg/archer/project.go

// Package mocks is a generated GoMock package.
package mocks

import (
	archer "github.com/aws/amazon-ecs-cli-v2/internal/pkg/archer"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockProjectStore is a mock of ProjectStore interface
type MockProjectStore struct {
	ctrl     *gomock.Controller
	recorder *MockProjectStoreMockRecorder
}

// MockProjectStoreMockRecorder is the mock recorder for MockProjectStore
type MockProjectStoreMockRecorder struct {
	mock *MockProjectStore
}

// NewMockProjectStore creates a new mock instance
func NewMockProjectStore(ctrl *gomock.Controller) *MockProjectStore {
	mock := &MockProjectStore{ctrl: ctrl}
	mock.recorder = &MockProjectStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockProjectStore) EXPECT() *MockProjectStoreMockRecorder {
	return m.recorder
}

// ListProjects mocks base method
func (m *MockProjectStore) ListProjects() ([]*archer.Project, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListProjects")
	ret0, _ := ret[0].([]*archer.Project)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListProjects indicates an expected call of ListProjects
func (mr *MockProjectStoreMockRecorder) ListProjects() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListProjects", reflect.TypeOf((*MockProjectStore)(nil).ListProjects))
}

// GetProject mocks base method
func (m *MockProjectStore) GetProject(projectName string) (*archer.Project, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProject", projectName)
	ret0, _ := ret[0].(*archer.Project)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProject indicates an expected call of GetProject
func (mr *MockProjectStoreMockRecorder) GetProject(projectName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProject", reflect.TypeOf((*MockProjectStore)(nil).GetProject), projectName)
}

// CreateProject mocks base method
func (m *MockProjectStore) CreateProject(project *archer.Project) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateProject", project)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateProject indicates an expected call of CreateProject
func (mr *MockProjectStoreMockRecorder) CreateProject(project interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProject", reflect.TypeOf((*MockProjectStore)(nil).CreateProject), project)
}

// DeleteProject mocks base method
func (m *MockProjectStore) DeleteProject(name string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteProject", name)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteProject indicates an expected call of DeleteProject
func (mr *MockProjectStoreMockRecorder) DeleteProject(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteProject", reflect.TypeOf((*MockProjectStore)(nil).DeleteProject), name)
}

// MockProjectLister is a mock of ProjectLister interface
type MockProjectLister struct {
	ctrl     *gomock.Controller
	recorder *MockProjectListerMockRecorder
}

// MockProjectListerMockRecorder is the mock recorder for MockProjectLister
type MockProjectListerMockRecorder struct {
	mock *MockProjectLister
}

// NewMockProjectLister creates a new mock instance
func NewMockProjectLister(ctrl *gomock.Controller) *MockProjectLister {
	mock := &MockProjectLister{ctrl: ctrl}
	mock.recorder = &MockProjectListerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockProjectLister) EXPECT() *MockProjectListerMockRecorder {
	return m.recorder
}

// ListProjects mocks base method
func (m *MockProjectLister) ListProjects() ([]*archer.Project, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListProjects")
	ret0, _ := ret[0].([]*archer.Project)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListProjects indicates an expected call of ListProjects
func (mr *MockProjectListerMockRecorder) ListProjects() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListProjects", reflect.TypeOf((*MockProjectLister)(nil).ListProjects))
}

// MockProjectCreator is a mock of ProjectCreator interface
type MockProjectCreator struct {
	ctrl     *gomock.Controller
	recorder *MockProjectCreatorMockRecorder
}

// MockProjectCreatorMockRecorder is the mock recorder for MockProjectCreator
type MockProjectCreatorMockRecorder struct {
	mock *MockProjectCreator
}

// NewMockProjectCreator creates a new mock instance
func NewMockProjectCreator(ctrl *gomock.Controller) *MockProjectCreator {
	mock := &MockProjectCreator{ctrl: ctrl}
	mock.recorder = &MockProjectCreatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockProjectCreator) EXPECT() *MockProjectCreatorMockRecorder {
	return m.recorder
}

// CreateProject mocks base method
func (m *MockProjectCreator) CreateProject(project *archer.Project) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateProject", project)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateProject indicates an expected call of CreateProject
func (mr *MockProjectCreatorMockRecorder) CreateProject(project interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProject", reflect.TypeOf((*MockProjectCreator)(nil).CreateProject), project)
}

// MockProjectGetter is a mock of ProjectGetter interface
type MockProjectGetter struct {
	ctrl     *gomock.Controller
	recorder *MockProjectGetterMockRecorder
}

// MockProjectGetterMockRecorder is the mock recorder for MockProjectGetter
type MockProjectGetterMockRecorder struct {
	mock *MockProjectGetter
}

// NewMockProjectGetter creates a new mock instance
func NewMockProjectGetter(ctrl *gomock.Controller) *MockProjectGetter {
	mock := &MockProjectGetter{ctrl: ctrl}
	mock.recorder = &MockProjectGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockProjectGetter) EXPECT() *MockProjectGetterMockRecorder {
	return m.recorder
}

// GetProject mocks base method
func (m *MockProjectGetter) GetProject(projectName string) (*archer.Project, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProject", projectName)
	ret0, _ := ret[0].(*archer.Project)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProject indicates an expected call of GetProject
func (mr *MockProjectGetterMockRecorder) GetProject(projectName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProject", reflect.TypeOf((*MockProjectGetter)(nil).GetProject), projectName)
}

// MockProjectDeleter is a mock of ProjectDeleter interface
type MockProjectDeleter struct {
	ctrl     *gomock.Controller
	recorder *MockProjectDeleterMockRecorder
}

// MockProjectDeleterMockRecorder is the mock recorder for MockProjectDeleter
type MockProjectDeleterMockRecorder struct {
	mock *MockProjectDeleter
}

// NewMockProjectDeleter creates a new mock instance
func NewMockProjectDeleter(ctrl *gomock.Controller) *MockProjectDeleter {
	mock := &MockProjectDeleter{ctrl: ctrl}
	mock.recorder = &MockProjectDeleterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockProjectDeleter) EXPECT() *MockProjectDeleterMockRecorder {
	return m.recorder
}

// DeleteProject mocks base method
func (m *MockProjectDeleter) DeleteProject(name string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteProject", name)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteProject indicates an expected call of DeleteProject
func (mr *MockProjectDeleterMockRecorder) DeleteProject(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteProject", reflect.TypeOf((*MockProjectDeleter)(nil).DeleteProject), name)
}

// MockProjectResourceStore is a mock of ProjectResourceStore interface
type MockProjectResourceStore struct {
	ctrl     *gomock.Controller
	recorder *MockProjectResourceStoreMockRecorder
}

// MockProjectResourceStoreMockRecorder is the mock recorder for MockProjectResourceStore
type MockProjectResourceStoreMockRecorder struct {
	mock *MockProjectResourceStore
}

// NewMockProjectResourceStore creates a new mock instance
func NewMockProjectResourceStore(ctrl *gomock.Controller) *MockProjectResourceStore {
	mock := &MockProjectResourceStore{ctrl: ctrl}
	mock.recorder = &MockProjectResourceStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockProjectResourceStore) EXPECT() *MockProjectResourceStoreMockRecorder {
	return m.recorder
}

// GetRegionalProjectResources mocks base method
func (m *MockProjectResourceStore) GetRegionalProjectResources(project *archer.Project) ([]*archer.ProjectRegionalResources, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRegionalProjectResources", project)
	ret0, _ := ret[0].([]*archer.ProjectRegionalResources)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRegionalProjectResources indicates an expected call of GetRegionalProjectResources
func (mr *MockProjectResourceStoreMockRecorder) GetRegionalProjectResources(project interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRegionalProjectResources", reflect.TypeOf((*MockProjectResourceStore)(nil).GetRegionalProjectResources), project)
}

// GetProjectResourcesByRegion mocks base method
func (m *MockProjectResourceStore) GetProjectResourcesByRegion(project *archer.Project, region string) (*archer.ProjectRegionalResources, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProjectResourcesByRegion", project, region)
	ret0, _ := ret[0].(*archer.ProjectRegionalResources)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProjectResourcesByRegion indicates an expected call of GetProjectResourcesByRegion
func (mr *MockProjectResourceStoreMockRecorder) GetProjectResourcesByRegion(project, region interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProjectResourcesByRegion", reflect.TypeOf((*MockProjectResourceStore)(nil).GetProjectResourcesByRegion), project, region)
}
