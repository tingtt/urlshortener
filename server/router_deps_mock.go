package server

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

var _ Handler = new(MockHandler)

type MockHandler struct {
	mock.Mock
}

// HandleGet implements Handler.
func (m *MockHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

// HandlePost implements Handler.
func (m *MockHandler) HandlePost(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}
