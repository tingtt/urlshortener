package server

import (
	"io"
	uiprovider "urlshortener/ui/provider"
	"urlshortener/usecase"

	"github.com/stretchr/testify/mock"
	"maragu.dev/gomponents"
)

var _ usecase.Handler = new(MockUsecase)

type MockUsecase struct {
	mock.Mock
}

// Delete implements usecase.Handler.
func (m *MockUsecase) Delete(shortURLs ...string) error {
	args := m.Called(shortURLs)
	return args.Error(0)
}

// Find implements usecase.Handler.
func (m *MockUsecase) Find(shortURL string) (redirectTarget string, err error) {
	args := m.Called(shortURL)
	return args.String(0), args.Error(1)
}

// FindAll implements usecase.Handler.
func (m *MockUsecase) FindAll() ([]usecase.ShortURL, error) {
	args := m.Called()
	return args.Get(0).([]usecase.ShortURL), args.Error(1)
}

// Save implements usecase.Handler.
func (m *MockUsecase) Save(shortURL string, redirectTarget string) error {
	args := m.Called(shortURL, redirectTarget)
	return args.Error(0)
}

var _ uiprovider.Provider = new(MockUI)

type MockUI struct {
	mock.Mock
}

// EditPage implements uiprovider.Provider.
func (m *MockUI) EditPage(reqPath string, redirectTargetURL string, shortURLs []usecase.ShortURL) gomponents.Node {
	args := m.Called(reqPath, redirectTargetURL, shortURLs)
	return args.Get(0).(gomponents.Node)
}

// RegisterPage implements uiprovider.Provider.
func (m *MockUI) RegisterPage(reqPath string, shortURLs []usecase.ShortURL) gomponents.Node {
	args := m.Called(reqPath, shortURLs)
	return args.Get(0).(gomponents.Node)
}

var _ gomponents.Node = new(MockNode)

type MockNode struct {
	mock.Mock
}

// Render implements gomponents.Node.
func (m *MockNode) Render(w io.Writer) error {
	args := m.Called(w)
	return args.Error(0)
}
