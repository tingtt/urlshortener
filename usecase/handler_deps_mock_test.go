package usecase

import (
	"urlshortener/registry"

	"github.com/stretchr/testify/mock"
)

var _ registry.Registry = new(MockRegistry)

type MockRegistry struct {
	mock.Mock
}

// Find implements registry.Registry.
func (m *MockRegistry) Find(path string) (redirectTarget string, err error) {
	args := m.Called(path)
	return args.String(0), args.Error(1)
}

// FindAll implements registry.Registry.
func (m *MockRegistry) FindAll() (map[string]string, error) {
	args := m.Called()
	return args.Get(0).(map[string]string), args.Error(1)
}

// Remove implements registry.Registry.
func (m *MockRegistry) Remove(path string) error {
	args := m.Called(path)
	return args.Error(0)
}

// Save implements registry.Registry.
func (m *MockRegistry) Save(path string, redirectTarget string) error {
	args := m.Called(path, redirectTarget)
	return args.Error(0)
}

type RegistryBehavior struct {
	find    *RegistryBehaviorFind
	findAll *RegistryBehaviorFindAll
	remove  *RegistryBehaviorRemove
	save    *RegistryBehaviorSave
}

type RegistryBehaviorFind struct {
	redirectTarget string
	err            error
}

type RegistryBehaviorFindAll struct {
	shortURLs map[string]string
	err       error
}

type RegistryBehaviorRemove struct {
	err error
}

type RegistryBehaviorSave struct {
	err error
}
