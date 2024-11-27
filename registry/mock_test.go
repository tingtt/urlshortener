package registry

import "github.com/stretchr/testify/mock"

var _ fsregistry = new(MockFSRegistry)

type MockFSRegistry struct {
	mock.Mock
}

// loadToCache implements fsregistry.
func (m *MockFSRegistry) loadToCache() error {
	args := m.Called()
	return args.Error(0)
}

// savePersistently implements fsregistry.
func (m *MockFSRegistry) savePersistently() error {
	args := m.Called()
	return args.Error(0)
}
