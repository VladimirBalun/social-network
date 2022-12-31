// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/services/profile_service.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"
	entities "social_network/internal/entities"

	gomock "github.com/golang/mock/gomock"
)

// MockProfileRepository is a mock of ProfileRepository interface.
type MockProfileRepository struct {
	ctrl     *gomock.Controller
	recorder *MockProfileRepositoryMockRecorder
}

// MockProfileRepositoryMockRecorder is the mock recorder for MockProfileRepository.
type MockProfileRepositoryMockRecorder struct {
	mock *MockProfileRepository
}

// NewMockProfileRepository creates a new mock instance.
func NewMockProfileRepository(ctrl *gomock.Controller) *MockProfileRepository {
	mock := &MockProfileRepository{ctrl: ctrl}
	mock.recorder = &MockProfileRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProfileRepository) EXPECT() *MockProfileRepositoryMockRecorder {
	return m.recorder
}

// GetProfiles mocks base method.
func (m *MockProfileRepository) GetProfiles(arg0 context.Context) ([]entities.Profile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProfiles", arg0)
	ret0, _ := ret[0].([]entities.Profile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProfiles indicates an expected call of GetProfiles.
func (mr *MockProfileRepositoryMockRecorder) GetProfiles(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProfiles", reflect.TypeOf((*MockProfileRepository)(nil).GetProfiles), arg0)
}

// SaveProfile mocks base method.
func (m *MockProfileRepository) SaveProfile(arg0 context.Context, arg1 *entities.Profile) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveProfile", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveProfile indicates an expected call of SaveProfile.
func (mr *MockProfileRepositoryMockRecorder) SaveProfile(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveProfile", reflect.TypeOf((*MockProfileRepository)(nil).SaveProfile), arg0, arg1)
}